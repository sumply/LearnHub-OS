package postgres

import (
	"context"
	"server/internal/adapter/postgres/row"
	"server/internal/adapter/postgres/utils"
	"server/internal/dto"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type QuizItem struct {
	goqu   *goqu.Database
	tables goquTableNames
}

func (p *Postgres) FindUserLastAttempt(ctx context.Context, quizID uuid.UUID) ([]dto.UserLastAttempt, error) {
	ds := p.buildFindUserLastAttemptQuery(quizID)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}
	defer scanner.Close()

	row := new(row.UserLastAttempt)
	return utils.ScanDTO(scanner, row)
}

func (p *Postgres) FindStudents(ctx context.Context, groupID uuid.UUID) ([]dto.User, error) {
	ds := p.buildFindStudentsQuery(groupID)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}
	defer scanner.Close()

	row := new(row.User)
	return utils.ScanDTO(scanner, row)
}

func (p *Postgres) buildQuizLastAttemptQuery(userID uuid.UUID) *goqu.SelectDataset {
	info := p.tables.QuizInfo
	attempt := p.tables.QuizAttempt

	owner := p.tables.AccountProfile.As("owner")
	subject := p.tables.SchoolSubject

	attemptOn := goqu.On(
		attempt.Col("user_id").Eq(userID),
		attempt.Col("quiz_id").Eq(info.Col("quiz_id")),
	)
	ownerOn := goqu.On(
		info.Col("owner_id").Eq(owner.Col("account_id")),
	)
	subjectOn := goqu.On(
		info.Col("subject_id").Eq(subject.Col("id")),
	)

	return p.goqu.From(info).
		Distinct(info.Col("quiz_id")).
		//Select(&row.QuizLastAttempt{}).
		LeftJoin(attempt, attemptOn).
		InnerJoin(owner, ownerOn).
		InnerJoin(subject, subjectOn).
		Order(
			info.Col("quiz_id").Asc(),
			attempt.Col("started_at").Desc(),
		)
}

func (p *Postgres) buildFindUserLastAttemptQuery(quizID uuid.UUID) *goqu.SelectDataset {
	attempt := p.tables.QuizAttempt
	user := p.tables.AccountProfile.As("user")

	attemptOn := goqu.On(attempt.Col("user_id").Eq(user.Col("account_id")))

	return p.goqu.From(user).
		Distinct(user.Col("account_id")).
		Select(&row.UserLastAttempt{}).
		LeftJoin(attempt, attemptOn).
		Where(
			goqu.Or(
				attempt.Col("quiz_id").Eq(quizID),
				attempt.Col("quiz_id").IsNull(),
			),
		).
		Order(
			user.Col("account_id").Asc(),
			attempt.Col("started_at").Desc(),
		)
}

func (p *Postgres) buildFindStudentsQuery(groupID uuid.UUID) *goqu.SelectDataset {
	student := p.tables.AccountStudent
	profile := p.tables.AccountProfile

	studentOn := goqu.On(
		student.Col("account_id").Eq(profile.Col("account_id")),
		student.Col("group_id").Eq(groupID),
	)

	return p.goqu.From(profile).
		Select(profile.All()).
		InnerJoin(student, studentOn)
}
