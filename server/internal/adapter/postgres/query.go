package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/dto"

	"github.com/doug-martin/goqu/v9"
)

type Query struct {
	*Postgres
}

type Row struct {
	Group    sqlc.SchoolGroup          `db:"group"`
	Curator  sql.Null[json.RawMessage] `db:"curator"`
	Students json.RawMessage           `db:"students"`
}

type SelectRow struct {
	Group sqlc.SchoolGroup `db:"group"`
}

func (q *Query) ListGroup(ctx context.Context) ([]dto.Group, error) {
	profile := q.tables.AccountProfile.As("curator")
	teacher := q.tables.AccountTeacher
	group := q.tables.SchoolGroup

	ds := q.goqu.From(group).
		Select(
			&SelectRow{},
			goqu.Func("row_to_json", profile.All()).As("curator"),
			goqu.COALESCE(q.withStudents(), goqu.L("'[]'::JSONB")).As("students"),
		).
		LeftJoin(teacher, goqu.On(
			teacher.Col("group_id").Eq(group.Col("id")),
		)).
		LeftJoin(profile, goqu.On(
			profile.Col("account_id").Eq(teacher.Col("account_id")),
		))

	var rows []Row
	if err := ds.Executor().ScanStructsContext(ctx, &rows); err != nil {
		return nil, err
	}

	groups := make([]dto.Group, 0, len(rows))
	for _, row := range rows {
		var jsonStudents []sqlc.AccountProfile
		if err := json.Unmarshal(row.Students, &jsonStudents); err != nil {
			return nil, err
		}
		students := make([]dto.User, 0, len(jsonStudents))
		for _, s := range jsonStudents {
			students = append(students, dto.User{
				ID:        s.AccountID,
				FirstName: s.FirstName,
				LastName:  s.LastName,
				Role:      string(s.Role),
			})
		}
		var curator *dto.User
		if row.Curator.Valid {
			var curatorRow sqlc.AccountProfile
			if err := json.Unmarshal(row.Curator.V, &curatorRow); err != nil {
				return nil, err
			}
			curator = &dto.User{
				ID:        curatorRow.AccountID,
				FirstName: curatorRow.FirstName,
				LastName:  curatorRow.LastName,
				Role:      string(curatorRow.Role),
			}
		}
		groups = append(groups, dto.Group{
			ID:       row.Group.ID,
			Name:     row.Group.Name,
			Curator:  curator,
			Students: students,
		})
	}

	return groups, nil
}

func (q *Query) withStudents() *goqu.SelectDataset {
	profile := q.tables.AccountProfile
	student := q.tables.AccountStudent
	group := q.tables.SchoolGroup

	return goqu.From(profile).
		Select(
			goqu.Func("jsonb_agg", profile.All()),
		).
		InnerJoin(student, goqu.On(
			student.Col("account_id").Eq(profile.Col("account_id")),
		)).
		Where(
			student.Col("group_id").Eq(group.Col("id")),
		)
}
