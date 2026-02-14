package repository

type NotFoundError struct{}

func NewNotFoundError() *NotFoundError {
	return &NotFoundError{}
}

func (n NotFoundError) Error() string {
	return "not found"
}
