package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	ID        uuid.UUID
	QuizID    uuid.UUID
	UserID    uuid.UUID
	Answers   []Answer
	Score     int
	StartedAt time.Time
	EndedAt   *time.Time
}

func NewAttempt(quiz *Quiz, userID uuid.UUID) (Attempt, error) {
	if quiz == nil {
		panic("quiz is panic")
	}

	err := NewError("attempt")

	if userID == uuid.Nil {
		err.add("userID", fmt.Errorf("userID is empty"))
	}

	if !err.Empty() {
		return Attempt{}, err
	}

	attemptID := uuid.New()

	answers := make([]Answer, len(quiz.Questions))
	for i := range answers {
		answers[i] = Answer{
			ID:         uuid.New(),
			AttemptID:  attemptID,
			QuestionID: quiz.Questions[i].ID,
		}
	}

	return Attempt{
		ID:        attemptID,
		QuizID:    quiz.ID,
		UserID:    userID,
		Answers:   answers,
		StartedAt: time.Now().UTC(),
	}, nil
}

func (a *Attempt) Finish() {
	endedAt := time.Now().UTC()
	a.EndedAt = &endedAt
}

type Answer struct {
	ID         uuid.UUID
	AttemptID  uuid.UUID
	QuestionID uuid.UUID
	Answer     any
	Score      int
	IsCorrect  bool
}
