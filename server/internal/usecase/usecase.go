package usecase

import "server/internal/logger"

type usecase struct {
}

func (u *usecase) mapStorageError(err error) error {
	return err
}

func (u *usecase) tracedFieldWithUsecase(data map[string]any) logger.TraceField {
	return logger.TraceField{
		Key:   "usecase",
		Value: data,
	}
}
