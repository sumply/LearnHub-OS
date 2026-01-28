package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Attempt struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	QuizID    uuid.UUID
	Answers   []Answer
	Score     int
	StartedAt time.Time
	EndedAt   *time.Time
}

func NewAttempt(quiz *Quiz, userID uuid.UUID) (Attempt, error) {
	if quiz == nil {
		return Attempt{}, fmt.Errorf("quiz is nil: %w", ErrInvalid)
	}
	answers := make([]Answer, len(quiz.Content))
	for i, question := range quiz.Content {
		answers[i] = NewAnswer(question.ID)
	}
	attempt := Attempt{
		ID:        uuid.New(),
		UserID:    userID,
		QuizID:    quiz.ID,
		Answers:   answers,
		StartedAt: time.Now().UTC(),
	}
	return attempt, nil
}

type Answer struct {
	ID         uuid.UUID
	QuestionID uuid.UUID
	Answer     any
	IsCorrect  bool
	Score      int
}

func NewAnswer(questionID uuid.UUID) Answer {
	return Answer{
		ID:         uuid.New(),
		QuestionID: questionID,
	}
}
