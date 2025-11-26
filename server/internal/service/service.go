package service

import (
	"server/internal/dto"
	"time"
)

type Interface interface {
	Authz(*dto.AuthzRequest) (*dto.AuthzResponse, error)
	Registration(*dto.RegistrationRequest) error
}

func New() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Authz(req *dto.AuthzRequest) (*dto.AuthzResponse, error) {
	return nil, nil
}

func (s *Service) Registration(req *dto.RegistrationRequest) error {
	return nil
}

func NewMock() *Mock {
	return &Mock{}
}

type Mock struct{}

func (m *Mock) Authz(req *dto.AuthzRequest) (*dto.AuthzResponse, error) {
	return &dto.AuthzResponse{
		JWT: dto.JWTResponse{
			Refresh: "{id: 10}",
			Access:  "{id: 10}",
		},
		User: dto.UserResponse{
			ID:         10,
			FirstName:  "Name",
			LastName:   "Name",
			MiddleName: "Name",
			IconRef:    "",
			Role:       "user",
			CreatedAt:  time.Now(),
		},
	}, nil
}

func (m *Mock) Registration(req *dto.RegistrationRequest) error {
	return nil
}
