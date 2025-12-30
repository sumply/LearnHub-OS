package usecase

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/logger"
)

type Quiz struct {
}

func (q *Quiz) Create(ctx context.Context, identity *dto.Identity, req *dto.QuizCreateReq) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(req),
	)
	log.Debug("Called a create quiz method")

	if !identity.Role.IsHigherOrEqual(domain.UserAdmin) {
		log.Warn("User role is less than admin")
		return ErrAccess
	}

	quiz, err := q.createQuiz(identity, req)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	log.With(
		logger.TraceFieldFromAny(quiz),
	).Debug("Created quiz")

	return nil
}

func (q *Quiz) Get(ctx context.Context, identity *dto.Identity) ([]*domain.Quiz, error) {
	return nil, nil
}

func (q *Quiz) createQuiz(identity *dto.Identity, req *dto.QuizCreateReq) (*domain.Quiz, error) {
	questions, err := q.createQuestionSlice(req)
	if err != nil {
		return nil, err
	}

	new, err := domain.NewQuiz(
		req.Title,
		req.Summary,
		questions,
		identity.ID,
		req.SubjectID,
		req.GroupIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidField, err)
	}
	return new, nil
}

func (q *Quiz) createQuestionSlice(req *dto.QuizCreateReq) ([]*domain.Question, error) {
	questions := make([]*domain.Question, len(req.Questions))
	for i, question := range req.Questions {
		options, err := q.createOptionSlice(question)
		if err != nil {
			return nil, err
		}
		new, err := domain.NewQuestion(question.Text, options)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidField, err)
		}
		questions[i] = new
	}
	return questions, nil
}

func (q *Quiz) createOptionSlice(req *dto.QuestionCreateReq) ([]*domain.Option, error) {
	options := make([]*domain.Option, len(req.Options))
	for i, option := range req.Options {
		new, err := domain.NewOption(option.Text, option.IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidField, err)
		}
		options[i] = new
	}
	return options, nil
}
