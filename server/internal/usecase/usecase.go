package usecase

import (
	"context"
	"time"
)

type UserLoginParam struct {
	Email    string
	Password string
}

type JWT struct {
	AccessToken  string
	RefreshToken string
}

type UserGetParam struct {
	ID         *int64
	FirstName  *string
	LastName   *string
	MiddleName *string
}

type UserData struct {
	ID         int64
	FirstName  string
	LastName   string
	MiddleName string
	Role       string
	CreatedAt  time.Time
}

type UserPutParam struct {
	FirstName  string
	LastName   string
	MiddleName string
}

type UserCreateParam struct {
	FirstName  string
	LastName   string
	MiddleName string
}

type AuthData struct {
	Subject int64
	Role    string
}

type User interface {
	Login(ctx context.Context, p UserLoginParam) (JWT, error)
	Get(ctx context.Context, auth AuthData, p UserGetParam) ([]UserData, error)
	Create(ctx context.Context, auth AuthData, p UserCreateParam) error
	GetMe(ctx context.Context, auth AuthData) (UserData, error)
	GetByID(ctx context.Context, auth AuthData, id int64) (UserData, error)
	Put(ctx context.Context, auth AuthData, id int64, p UserPutParam) error
	Delete(ctx context.Context, auth AuthData, id int64) error
}
