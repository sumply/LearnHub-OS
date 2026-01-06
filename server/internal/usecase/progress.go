package usecase

import (
	"context"
	"fmt"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/dto"
	"server/internal/logger"
	"server/internal/repository"
)

type ProgressUsecase struct {
	repo *repository.Repository
}

func NewProgressUsecase(repo *repository.Repository) *ProgressUsecase {
	return &ProgressUsecase{repo: repo}
}

func (p *ProgressUsecase) UpdateAnswer(
	ctx context.Context,
	identity *dto.Identity,
	req *dto.AnswerPatchReq,
	progressID, answerID common.ID,
) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.TraceFieldFromAny(req),
		logger.NewTracedField("progress_id", progressID),
		logger.NewTracedField("answer_id", answerID),
	)
	log.Debug("Called a patch answer usecase")

	progress, err := p.repo.Progress().GetByID(ctx, progressID)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if progress.User.ID != identity.ID {
		return fmt.Errorf("%w: user (id=%d) is not progress owner", ErrAccess, identity.ID)
	}

	if !progress.InProgress() {
		return fmt.Errorf("%w: quiz progress (id=%d) status is not progress", ErrAccess, progressID)
	}

	answer, ok := progress.GetAnswer(answerID)
	if !ok {
		return fmt.Errorf("%w: answer (id=%d) is not exists", ErrNotFound, answerID)
	}

	if err := answer.Submit(req.Text); err != nil {
		log.Warn(err.Error())
		return err
	}

	if err := p.repo.Progress().UpdateAnswer(ctx, answer); err != nil {
		log.Warn(err.Error())
		return err
	}

	return nil
}

func (p *ProgressUsecase) ReviewAnswer(
	ctx context.Context,
	identity *dto.Identity,
	isCorrect bool,
	progressID, answerID common.ID,
) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.NewTracedField("answer_correct", isCorrect),
		logger.NewTracedField("progress_id", progressID),
		logger.NewTracedField("answer_id", answerID),
	)
	log.Debug("Called a review answer usecase")

	ctx = logger.WithLoggerCtx(ctx, log)
	progress, err := p.repo.Progress().GetByID(ctx, progressID)
	if err != nil {
		log.Warn(err.Error())
		return err
	}

	if !progress.Quiz.IsOwner(identity.ID) {
		return fmt.Errorf("%w: user (id=%d) is not owner", ErrAccess, identity.ID)
	}

	if !progress.IsPendingReview() {
		return fmt.Errorf("%w: progress (id=%d) is not awaiting review", ErrAccess, progressID)
	}

	answer, ok := progress.GetAnswer(answerID)
	if !ok {
		return fmt.Errorf("%w: answer (id=%d) is not exists", ErrNotFound, answerID)
	}

	if err := answer.Review(isCorrect); err != nil {
		return err
	}

	if err := p.repo.Progress().UpdateAnswer(ctx, answer); err != nil {
		return err
	}
	return nil
}

func (p *ProgressUsecase) Get(ctx context.Context, identity *dto.Identity) ([]*domain.QuizProgress, error) {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
	)
	log.Debug("Called a get usecase")

	var filter *repository.ProgressFilter

	switch identity.Role {
	case domain.UserAdmin, domain.UserRoot:
		filter = nil
	case domain.UserTeacher:
		filter = &repository.ProgressFilter{
			Quiz: &repository.QuizFilter{
				OwnerID: &identity.ID,
			},
		}
	case domain.UserStudent:
		filter = &repository.ProgressFilter{
			UserID: &identity.ID,
		}
	default:
		return nil, fmt.Errorf("%w: invalid role (%d)", ErrAccess, identity.Role)
	}

	progresses, err := p.repo.Progress().Get(ctx, filter)
	if err != nil {
		return nil, err
	}
	return progresses, nil
}

func (p *ProgressUsecase) Start(ctx context.Context, identity *dto.Identity, progressID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.NewTracedField("progress_id", progressID),
	)
	log.Debug("Called a start usecase")

	ctx = logger.WithLoggerCtx(ctx, log)
	progress, err := p.repo.Progress().GetByID(ctx, progressID)
	if err != nil {
		return err
	}

	if progress.User.ID != identity.ID {
		return fmt.Errorf("%w: user (id=%d) is not progress owner", ErrAccess, identity.ID)
	}

	if err := progress.Start(); err != nil {
		return err
	}

	if err := p.repo.Progress().Update(progress); err != nil {
		return err
	}

	return nil
}

func (p *ProgressUsecase) Finish(ctx context.Context, identity *dto.Identity, progressID common.ID) error {
	log := logger.FromCtx(ctx).With(
		logger.TraceFieldFromAny(identity),
		logger.NewTracedField("progress_id", progressID),
	)
	log.Debug("Called a finish usecase")

	ctx = logger.WithLoggerCtx(ctx, log)
	progress, err := p.repo.Progress().GetByID(ctx, progressID)
	if err != nil {
		return err
	}

	if progress.User.ID == identity.ID {
		if err := progress.Complete(); err != nil {
			return err
		}
	} else if progress.Quiz.IsOwner(identity.ID) {
		if err := progress.CompleteReview(); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("%w: user (id=%d) does not have acces for progress (id=%d)", ErrAccess, identity.ID, progressID)
	}

	if err := p.repo.Progress().Update(progress); err != nil {
		return err
	}

	return nil
}
