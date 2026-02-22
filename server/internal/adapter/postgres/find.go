package postgres

import (
	"context"
	"fmt"
	"server/internal/adapter/postgres/row"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/adapter/postgres/utils"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/pkg/repository"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type QuizItem struct {
	goqu   *goqu.Database
	tables goquTableNames
}

func (p *QuizItem) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]dto.QuizItem, error) {
	quiz := p.tables.QuizInfo

	ds := p.goqu.From(quiz).
		Select(quiz.All()).
		Where(quiz.Col("owner_id").Eq(ownerID))

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}

	return utils.ScanDTO(scanner, new(row.QuizItem))
}

func (p *QuizItem) ListByGroup(ctx context.Context, groupID uuid.UUID) ([]dto.QuizItem, error) {
	quiz := p.tables.QuizInfo
	assignment := p.tables.QuizAssignment
	subject := p.tables.SchoolSubject
	owner := p.tables.AccountProfile.As("owner")

	ds := p.goqu.From(quiz).
		Select(&row.QuizItem{}).
		InnerJoin(owner, goqu.On(
			owner.Col("account_id").Eq(quiz.Col("owner_id")),
		)).
		InnerJoin(subject, goqu.On(
			subject.Col("id").Eq(quiz.Col("subject_id")),
		)).
		InnerJoin(assignment, goqu.On(
			assignment.Col("quiz_id").Eq(quiz.Col("quiz_id")),
		)).
		Where(
			assignment.Col("group_id").Eq(groupID),
		)

	fmt.Println(ds.ToSQL())

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}

	return utils.ScanDTO(scanner, new(row.QuizItem))
}

type Group struct {
	goqu   *goqu.Database
	tables goquTableNames
}

func (p *Group) GetByStudent(ctx context.Context, studentID uuid.UUID) (domain.Group, error) {
	group := p.tables.SchoolGroup
	student := p.tables.AccountStudent

	ds := p.goqu.From(group).
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

func (p *Postgres) FindQuizLastAttempt(ctx context.Context, userID uuid.UUID) ([]dto.QuizLastAttempt, error) {
	ds := p.buildQuizLastAttemptQuery(userID)

	scanner, err := ds.Executor().ScannerContext(ctx)
	if err != nil {
		return nil, err
	}

	row := new(row.QuizLastAttempt)
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
		Select(&row.QuizLastAttempt{}).
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
