package get_quiz

import (
	"context"
	"os"
	"server/internal/domain"
	"server/internal/pkg/encoder"
	"server/internal/pkg/repository"

	"github.com/google/uuid"
)

type QuizRepository interface {
	repository.Geter[*domain.Quiz]
}

type UseCase struct {
	quizRepo QuizRepository
}

func New(quiz QuizRepository) *UseCase {
	return &UseCase{
		quizRepo: quiz,
	}
}

func (uc *UseCase) GetQuiz(ctx context.Context, id uuid.UUID) (Response, error) {
	quiz, err := uc.quizRepo.Get(ctx, id)
	if err != nil {
		return Response{}, err
	}

	questions := make([]*ResponseQuestion, 0, len(quiz.Questions))
	for _, q := range quiz.Questions {
		resp, err := NewResponseQuestion(q)
		if err != nil {
			return Response{}, err
		}
		questions = append(questions, resp)
	}

	resp := Response{
		ID:          quiz.ID,
		Title:       quiz.Title,
		Summary:     quiz.Summary,
		OwnerID:     quiz.OwnerID,
		SubjectID:   quiz.SubjectID,
		GroupIDs:    quiz.GroupIDs,
		Deadline:    quiz.Deadline,
		MaxAttempts: quiz.MaxAttempts,
		Questions:   questions,
		TotalScore:  quiz.TotalScore,
		CreatedAt:   quiz.CreatedAt,
	}

	encoder.JSON(os.Stdout, resp)

	return resp, nil
}
