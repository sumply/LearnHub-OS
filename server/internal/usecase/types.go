package usecase

import (
	"context"
	"errors"
	"server/internal/common"
	"server/internal/domain"
	"server/internal/dto"
)

var (
	ErrAccess       = errors.New("access permission")
	ErrNotFound     = errors.New("not found")
	ErrCollision    = errors.New("collision")
	ErrInvalidField = errors.New("invalid field")
)

type ProgressInterface interface {
	UpdateAnswer(
		ctx context.Context,
		identity *dto.Identity,
		req *dto.AnswerPatchReq,
		progressID, answerID common.ID,
	) error
	ReviewAnswer(
		ctx context.Context,
		identity *dto.Identity,
		isCorrect bool,
		progressID, answerID common.ID,
	) error
	Get(context.Context, *dto.Identity) ([]*domain.QuizProgress, error)
	Start(ctx context.Context, identity *dto.Identity, progressID common.ID) error
	Finish(ctx context.Context, idenity *dto.Identity, progressID common.ID) error
}

type QuizInterface interface {
	Create(context.Context, *dto.Identity, *dto.QuizCreateReq) error
	Delete(ctx context.Context, identity *dto.Identity, quizID common.ID) error
	Get(context.Context, *dto.Identity) ([]*domain.Quiz, error)
}

type UserInterface interface {
	Login(context.Context, *dto.LoginReq) (*domain.TokenPair, error)
	Get(context.Context, *dto.Identity) ([]*domain.User, error)
	Create(context.Context, *dto.Identity, *dto.UserCreateReq) error
	GetMe(context.Context, *dto.Identity) (*domain.User, error)
	GetByID(context.Context, *dto.Identity, common.ID) (*domain.User, error)
}

type GroupInterface interface {
	Create(context.Context, *dto.Identity, *dto.GroupCreateReq) error
	Get(context.Context) ([]*domain.Group, error)
	GetByID(context.Context, *dto.Identity, common.ID) (*domain.Group, error)
	AddStudents(
		ctx context.Context,
		identity *dto.Identity,
		groupID common.ID,
		req *dto.GroupAddStudentsReq,
	) error
}

type SubjectInterface interface {
	Create(context.Context, *dto.Identity, *dto.SubjectCreateReq) error
	Get(context.Context, *dto.Identity) ([]*domain.Subject, error)
}
