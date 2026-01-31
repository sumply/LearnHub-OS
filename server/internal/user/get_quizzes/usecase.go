package get_quizzes

import "context"

type UseCase struct {
}

func (u *UseCase) GetQuizzes(ctx context.Context) (Output, error) {
	return Output{}, nil
}
