package postgres

import (
	"context"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
	"server/internal/pkg/repository"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type Group struct {
	*Postgres
	goqu   *goqu.Database
	tables goquTableNames
}

func NewGroup(p *Postgres) *Group {
	return &Group{
		Postgres: p,
		goqu:     p.goqu,
		tables:   p.tables,
	}
}

func (g *Group) Save(ctx context.Context, group domain.Group) error {
	err := g.sqlc.InsertSchoolGroup(ctx, sqlc.InsertSchoolGroupParams{
		ID:   group.ID,
		Name: group.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (g *Group) GetByStudent(ctx context.Context, studentID uuid.UUID) (domain.Group, error) {
	group := g.tables.SchoolGroup
	student := g.tables.AccountStudent

	ds := g.goqu.From(group).
		Select(group.All()).
		InnerJoin(student, goqu.On(
			student.Col("group_id").Eq(group.Col("id")),
		)).
		Where(
			student.Col("account_id").Eq(studentID),
		)

	var row sqlc.SchoolGroup
	ok, err := ds.Executor().ScanStructContext(ctx, &row)
	if err != nil {
		return domain.Group{}, err
	}

	if !ok {
		return domain.Group{}, repository.NewNotFoundError()
	}

	return domain.Group{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func (g *Group) AddTeacher(ctx context.Context, groupID uuid.UUID, teacherID uuid.UUID) error {
	return g.sqlc.InsertAccountTeacher(ctx, sqlc.InsertAccountTeacherParams{
		AccountID: teacherID,
		GroupID:   groupID,
	})
}

func (g *Group) AddStudents(ctx context.Context, groupID uuid.UUID, studentIDs uuid.UUIDs) error {
	return g.withQueries(ctx, func(q *sqlc.Queries) error {
		for _, id := range studentIDs {
			err := g.sqlc.InsertAccountStudent(ctx, sqlc.InsertAccountStudentParams{
				AccountID: id,
				GroupID:   groupID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (g *Group) Remove(ctx context.Context, id uuid.UUID) error {
	return g.sqlc.DeleteSchoolGroup(ctx, id)
}

func (g *Group) Get(ctx context.Context, id uuid.UUID) (domain.Group, error) {
	row, err := g.sqlc.GetDomainGroup(ctx, id)
	if err != nil {
		return domain.Group{}, err
	}
	return domain.Group{
		ID:         row.ID,
		Name:       row.Name,
		CuratorID:  row.CuratorID.UUID,
		StudentIDs: row.StudentIds,
	}, nil
}

func (g *Group) Update(ctx context.Context, group domain.Group) error {
	err := g.sqlc.UpdateSchoolGroup(ctx, sqlc.UpdateSchoolGroupParams{
		ID:   group.ID,
		Name: group.Name,
	})
	if err != nil {
		return err
	}
	return nil
}
