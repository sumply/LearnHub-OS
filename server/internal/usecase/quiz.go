package usecase

import (
	"context"
	"fmt"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/logger"
	"server/internal/repository"
)

type QuizReal struct {
	repo *repository.Repository
}

func NewQuiz(repo *repository.Repository) *QuizReal {
	return &QuizReal{
		repo: repo,
	}
}

func (q *QuizReal) Create(ctx context.Context, identity *dto.Identity, req *dto.QuizCreateReq) error {
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

	if err := q.repo.Quiz().Save(ctx, quiz); err != nil {
		log.Debug(err.Error())
		return err
	}

	return nil
}

func (q *QuizReal) Get(ctx context.Context, identity *dto.Identity) ([]*domain.Quiz, error) {
	return nil, nil
}

func (q *QuizReal) createQuiz(identity *dto.Identity, req *dto.QuizCreateReq) (*domain.Quiz, error) {
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

func (q *QuizReal) createQuestionSlice(req *dto.QuizCreateReq) ([]*domain.Question, error) {
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

func (q *QuizReal) createOptionSlice(req *dto.QuestionCreateReq) ([]*domain.Option, error) {
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
