package postgres

import (
	"context"
	"fmt"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/dto"
	"server/internal/repository/filter"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exec"
)

type quizItemRowScanner struct {
	scanner exec.Scanner
}

func (q quizItemRowScanner) scanToDTOQuizItems() ([]dto.QuizItem, error) {
	var items []dto.QuizItem

	for q.scanner.Next() {
		var row quizItemRow

		err := q.scanner.ScanStruct(&row)
		if err != nil {
			return nil, err
		}

		items = append(items, row.ToDTOQuizItem())
	}

	return items, nil
}

type quizItemRow struct {
	Quiz    sqlc.QuizInfo       `db:"info"`
	Subject sqlc.SchoolSubject  `db:"subject"`
	Owner   sqlc.AccountProfile `db:"owner"`
}

func (q *quizItemRow) ToDTOQuizItem() dto.QuizItem {
	owner := &q.Owner
	subject := &q.Subject
	quiz := &q.Quiz

	var deadline *time.Time
	if quiz.Deadline.Valid {
		deadline = &q.Quiz.Deadline.Time
	}

	return dto.QuizItem{
		ID:          quiz.QuizID,
		Title:       quiz.Title,
		Summary:     quiz.Summary,
		TotalScore:  int(quiz.TotalScore.(int64)),
		Deadline:    deadline,
		MaxAttempts: int(quiz.MaxAttempts),
		CreatedAt:   quiz.CreatedAt,

		Owner: dto.User{
			ID:        owner.AccountID,
			FirstName: owner.FirstName,
			LastName:  owner.LastName,
			Role:      owner.AccountID.String(),
		},

		Subject: dto.Subject{
			ID:   subject.ID,
			Name: subject.Name,
		},
	}
}

func (f *Postgres) FindQuizItems(ctx context.Context, filter *filter.Quiz) ([]dto.QuizItem, error) {
	quiz := goqu.S("quiz")

	info := quiz.Table("info")
	subject := goqu.S("school").Table("subject")
	owner := goqu.S("account").Table("profile").As("owner")

	subjectOn := goqu.On(subject.Col("id").Eq(info.Col("subject_id")))
	ownerOn := goqu.On(owner.Col("account_id").Eq(info.Col("owner_id")))

	ds := f.goqu.From(info).
		SelectDistinct(&quizItemRow{}).
		InnerJoin(subject, subjectOn).
		InnerJoin(owner, ownerOn)

	if quiz != nil {
		ds = f.joinAttemptToQuiz(ds, filter.Attempt)
	}

	fmt.Println(ds.ToSQL())

	scanner, err := ds.Executor().Scanner()
	if err != nil {
		return nil, err
	}

	defer scanner.Close()

	return quizItemRowScanner{
		scanner: scanner,
	}.scanToDTOQuizItems()
}

func (f *Postgres) joinAttemptToQuiz(ds *goqu.SelectDataset, filter *filter.Attempt) *goqu.SelectDataset {
	if filter == nil {
		return ds
	}

	quiz := goqu.S("quiz")

	attempt := quiz.Table("attempt")
	info := quiz.Table("info")

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
