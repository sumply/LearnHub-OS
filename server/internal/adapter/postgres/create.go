package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"server/internal/adapter/postgres/sqlc"
	"server/internal/domain"
)

func (p *Postgres) CreateUser(ctx context.Context, user *domain.User) error {
	err := p.createAccount(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (p *Postgres) CreateGroup(ctx context.Context, group *domain.Group) error {
	err := p.sqlc.InsertSchoolGroup(ctx, sqlc.InsertSchoolGroupParams{
		ID:   group.ID,
		Name: group.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateSubject(ctx context.Context, subject *domain.Subject) error {
	err := p.sqlc.InsertSchoolSubject(ctx, sqlc.InsertSchoolSubjectParams{
		ID:   subject.ID,
		Name: subject.Name,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateStudent(ctx context.Context, student *domain.Student) error {
	err := p.createAccount(ctx, &student.User)
	if err != nil {
		return err
	}

	err = p.sqlc.InsertStudent(ctx, sqlc.InsertStudentParams{
		AccountID: student.User.ID,
		GroupID:   student.Group,
	})

	return nil
}

func (p *Postgres) CreateTeacher(ctx context.Context, teacher *domain.Teacher) error {
	err := p.createAccount(ctx, &teacher.User)
	if err != nil {
		return err
	}

	err = p.sqlc.InsertTeacher(ctx, sqlc.InsertTeacherParams{
		AccountID: teacher.ID,
		GroupID:   teacher.Groups[0],
		SubjectID: teacher.Subjects[0],
	})

	return nil
}

func (p *Postgres) createAccount(ctx context.Context, user *domain.User) error {
	err := p.sqlc.InsertAccountCredential(ctx, sqlc.InsertAccountCredentialParams{
		AccountID: user.ID,
		Email:     user.Email,
		PwdHash:   user.PwdHash,
	})
	if err != nil {
		return err
	}

	err = p.sqlc.InsertAccountProfile(ctx, sqlc.InsertAccountProfileParams{
		AccountID: user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      sqlc.AccountUserRole(user.Role),
		CreatedAt: user.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}

func (p *Postgres) CreateQuiz(ctx context.Context, quiz *domain.Quiz) error {
	deadline := sql.NullTime{}
	deadline.Scan(quiz.Deadline)
	err := p.sqlc.InsertQuizInfo(ctx, sqlc.InsertQuizInfoParams{
		QuizID:      quiz.ID,
		Title:       quiz.Title,
		Summary:     quiz.Summary,
		SubjectID:   quiz.SubjectID,
		OwnerID:     quiz.OwnerID,
		MaxAttempts: int32(quiz.MaxAttempts),
		TotalScore:  quiz.TotalScore,
		Deadline:    deadline,
		CreatedAt:   quiz.CreatedAt,
	})
	if err != nil {
		return err
	}

	for _, question := range quiz.Content {
		details, err := json.Marshal(question.Details)
		if err != nil {
			return err
		}

		err = p.sqlc.InsertQuizQuestion(ctx, sqlc.InsertQuizQuestionParams{
			ID:      question.ID,
			QuizID:  question.QuizID,
			Title:   question.Text,
			Variant: sqlc.QuizQuestionType(question.Type),
			Score:   question.Score,
			Details: details,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) CreateAttempt(ctx context.Context, attempt *domain.Attempt) error {
	err := p.sqlc.InsertQuizAttempt(ctx, sqlc.InsertQuizAttemptParams{
		ID:        attempt.ID,
		QuizID:    attempt.QuizID,
		UserID:    attempt.UserID,
		Score:     attempt.Score,
		StartedAt: attempt.StartedAt,
	})
	if err != nil {
		return err
	}

	for _, answer := range attempt.Answers {
		details, err := json.Marshal(map[string]any{
			"answer": answer.Answer,
		})
		if err != nil {
			return err
		}

		err = p.sqlc.InsertQuizAnswer(ctx, sqlc.InsertQuizAnswerParams{
			ID:         answer.ID,
			AttemptID:  attempt.ID,
			QuestionID: answer.QuestionID,
			Details:    details,
			Score:      answer.Score,
			IsCorrect:  answer.IsCorrect,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
