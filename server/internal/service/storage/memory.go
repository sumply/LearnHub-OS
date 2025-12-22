package storage

type Storage interface {
	User() User
}

type User interface {
	Save(UserSaveParam) error
}

type UserSaveParam struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Login      string
	HashedPwd  string
}

func NewStub() *Stub {
	return &Stub{}
}

type Stub struct{}

func (s *Stub) User() User {
	return &StubUser{}
}

type StubUser struct{}

func (s *StubUser) Save(UserSaveParam) error {
	return nil
}
