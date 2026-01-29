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
		return Attempt{}, fmt.Errorf("quiz is nil: %w", ErrValidate)
	}

	if userID == uuid.Nil {
		return Attempt{}, fmt.Errorf("user id is empty: %w", ErrValidate)
	}

	answers := make([]Answer, len(quiz.Content))
	for i, question := range quiz.Content {
		answers[i] = NewAnswer(question.ID)
	}

	return Attempt{
		ID:        uuid.New(),
		UserID:    userID,
		QuizID:    quiz.ID,
		Answers:   answers,
		StartedAt: time.Now().UTC(),
	}, nil
}

func RestoreAttempt(
	id, userID, quizID uuid.UUID,
	answers []Answer,
	score int,
	startedAt time.Time,
	endedAt *time.Time,
) (Attempt, error) {
	if id == uuid.Nil {
		return Attempt{}, fmt.Errorf("id is empty: %w", ErrInvalid)
	}

	if userID == uuid.Nil {
		return Attempt{}, fmt.Errorf("user id is empty: %w", ErrInvalid)
	}

	if quizID == uuid.Nil {
		return Attempt{}, fmt.Errorf("quiz id is empty: %w", ErrInvalid)
	}

	if len(answers) == 0 {
		return Attempt{}, fmt.Errorf("answers are empty: %w", ErrInvalid)
	}

	return Attempt{
		ID:        id,
		UserID:    userID,
		QuizID:    quizID,
		Answers:   answers,
		Score:     score,
		StartedAt: startedAt,
		EndedAt:   endedAt,
	}, nil
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

func RestoreAnswer(
	id, questionID uuid.UUID,
	answer any,
	isCorrect bool,
	score int,
) (Answer, error) {
	if id == uuid.Nil {
		return Answer{}, fmt.Errorf("id is empty: %w", ErrInvalid)
	}

	if questionID == uuid.Nil {
		return Answer{}, fmt.Errorf("question id is empty: %w", ErrInvalid)
	}

	if score < 0 {
		return Answer{}, fmt.Errorf("score less 0: %w", ErrInvalid)
	}

	return Answer{
		ID:         id,
		QuestionID: questionID,
		Answer:     answer,
		IsCorrect:  isCorrect,
		Score:      score,
	}, nil
}
