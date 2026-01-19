package repository

import (
	"context"
	"errors"
	"server/archive/internal/domain"
	"server/archive/internal/logger"
	"testing"
)

func TestQuizOptionSave(t *testing.T) {
	t.Log("creating quiz options")
	repo := QuizMemory{
		storage: NewStorage(),
	}
	options := []*domain.Option{
		{
			Text:      "true",
			IsCorrect: true,
		},
		{
			Text:      "false",
			IsCorrect: false,
		},
	}
	newOptions := repo.saveOptions(options)
	if newOptions[0].ID == options[0].ID {
		t.Error("method 'saveOptions' is mutable")
	}

	for _, opt := range newOptions {
		if _, ok := repo.storage.option[opt.ID]; !ok {
			t.Error("a option has not been saved")
		}
	}
}

func TestQuizQuestionSave(t *testing.T) {
	t.Log("creating quiz questions")
	repo := QuizMemory{
		storage: NewStorage(),
	}
	questions := []*domain.Question{
		{
			Text:      "true",
			IsWritten: true,
		},
		{
			Text:      "false",
			IsWritten: false,
		},
	}
	newQuestions := repo.saveQuestions(questions)
	if newQuestions[0].ID == questions[0].ID {
		t.Error("method 'saveQuestions' is mutable")
	}

	for _, q := range newQuestions {
		if _, ok := repo.storage.question[q.ID]; !ok {
			t.Error("a question has not been saved")
		}
	}
}

func TestQuizSave(t *testing.T) {
	logger.SetNewFunc(func() logger.Logger {
		return &logger.Stub{}
	})
	t.Log("creating a quiz")
	repo := QuizMemory{
		storage: NewStorage(),
	}
	repo.storage.user[1] = &domain.User{
		ID:        1,
		FirstName: "test",
	}
	repo.storage.group[1] = &domain.Group{
		ID:   1,
		Name: "test",
	}
	repo.storage.subject[1] = &domain.Subject{
		ID:   1,
		Name: "test",
	}
	quiz := &domain.Quiz{
		Title:         "title",
		Summary:       "summary",
		IsForEveryone: true,
		Owner:         repo.storage.user[1],
		Subject:       repo.storage.subject[1],
		Questions: []*domain.Question{
			{
				Text:      "true",
				IsWritten: true,
			},
			{
				Text:      "false",
				IsWritten: true,
			},
		},
	}
	err := repo.Save(context.Background(), quiz)
	if err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
	if quiz.ID != 0 {
		t.Errorf("a method 'Save' is mutable")
	}
	quiz.Owner.ID = 0
	if err := repo.Save(context.Background(), quiz); !errors.Is(err, ErrDependence) {
		t.Errorf("a quiz can be created without owner_id")
	}
	quiz.Owner.ID = 1
	quiz.Subject.ID = 0
	if err := repo.Save(context.Background(), quiz); !errors.Is(err, ErrDependence) {
		t.Errorf("a quiz can be created without subject_id")
	}
	quiz.Subject.ID = 1
	quiz.Groups = []*domain.Group{
		{
			ID: 0,
		},
	}
	if err := repo.Save(context.Background(), quiz); !errors.Is(err, ErrDependence) {
		t.Errorf("a quiz can be created without group_id")
	}
	quiz.Groups = []*domain.Group{
		{
			ID: 1,
		},
	}
	if err := repo.Save(context.Background(), quiz); err != nil {
		t.Errorf("unexcepted error: %v", err)
	}
}
