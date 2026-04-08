package create_quiz

import (
	"context"
	"server/internal/domain"
	"server/internal/pkg/repository"
)

type QuizRepository interface {
	repository.Saver[*domain.Quiz]
}

type UseCase struct {
	quizRepo QuizRepository
}

func New(quiz QuizRepository) *UseCase {
	return &UseCase{
		quizRepo: quiz,
	}
}

func (u *UseCase) CreateQuiz(ctx context.Context, req Request) (Response, error) {
	questions := make([]domain.IQuestion, 0, len(req.Questions))
	for _, q := range req.Questions {
		d, err := q.Unpack()
		if err != nil {
			return Response{}, err
		}
		questions = append(questions, d)
	}

	quiz, err := domain.NewQuiz(
		req.OwnerID,
		req.SubjectID,
		req.Title,
		req.Summary,
		questions,
		req.GroupIDs,
		req.MaxAttempts,
		req.Deadline,
	)
	if err != nil {
		return Response{}, err
	}

	if err := u.quizRepo.Save(ctx, quiz); err != nil {
		return Response{}, err
	}

	return Response{
		ID: quiz.ID,
	}, nil
}
