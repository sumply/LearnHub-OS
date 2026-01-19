package usecase

type usecase struct {
}

func (u *usecase) mapStorageError(err error) error {
	return err
}
