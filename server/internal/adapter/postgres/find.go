package postgres

import (
	"context"
	"server/internal/adapter/postgres/row"
	"server/internal/adapter/postgres/utils"
	"server/internal/dto"
	"server/internal/repository/filter"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

func (p *Postgres) FindQuizItems(ctx context.Context, filter *filter.Quiz) ([]dto.QuizItem, error) {
	ds := p.buildFindQuizItemsQuery(filter)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}

	defer scanner.Close()

	return utils.ScanDTO[row.QuizItem, dto.QuizItem](scanner)
}

func (p *Postgres) joinAttemptToQuiz(ds *goqu.SelectDataset, filter *filter.Attempt) *goqu.SelectDataset {
	if filter == nil {
		return ds
	}

	attempt := p.tables.QuizAttempt
	info := p.tables.QuizInfo

	attemptOn := goqu.On(attempt.Col("quiz_id").Eq(info.Col("quiz_id")))

	if filter.IsCompleted != nil {
		if *filter.IsCompleted {
			ds = ds.InnerJoin(attempt, attemptOn).
				Where(
					attempt.Col("ended_at").IsNotNull(),
				)
			if filter.UserID != nil {
				ds = ds.Where(
					attempt.Col("user_id").Eq(filter.UserID),
				)
			}
		}
	}

	return ds
}

func (p *Postgres) FindUserLastAttempt(ctx context.Context, quizID uuid.UUID) ([]dto.UserLastAttempt, error) {
	ds := p.buildFindUserLastAttemptQuery(quizID)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}
	defer scanner.Close()

	return utils.ScanDTO[row.UserLastAttempt, dto.UserLastAttempt](scanner)
}

func (p *Postgres) FindStudents(ctx context.Context, groupID uuid.UUID) ([]dto.User, error) {
	ds := p.buildFindStudentsQuery(groupID)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}
	defer scanner.Close()

	return utils.ScanDTO[row.User, dto.User](scanner)
}

func (p *Postgres) buildFindUserLastAttemptQuery(quizID uuid.UUID) *goqu.SelectDataset {
	attempt := p.tables.QuizAttempt
	user := p.tables.AccountProfile

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

func (p *Postgres) buildFindQuizItemsQuery(filter *filter.Quiz) *goqu.SelectDataset {
	info := p.tables.QuizInfo
	subject := p.tables.SchoolSubject
	owner := p.tables.AccountProfile.As("owner")

	subjectOn := goqu.On(subject.Col("id").Eq(info.Col("subject_id")))
	ownerOn := goqu.On(owner.Col("account_id").Eq(info.Col("owner_id")))

	ds := p.goqu.From(info).
		SelectDistinct(&row.QuizItem{}).
		InnerJoin(subject, subjectOn).
		InnerJoin(owner, ownerOn)

	if filter != nil {
		ds = p.joinAttemptToQuiz(ds, filter.Attempt)
	}

	return ds
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
