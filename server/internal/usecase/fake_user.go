package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type FakeUser struct {
	admins   []UserDomain
	teachers []UserDomain
	students []UserDomain
	roots    []UserDomain
}

func NewFakeUser() *FakeUser {
	u := FakeUser{}
	u.admins = initFakeSlices(Admin)
	u.roots = initFakeSlices(Root)
	u.students = initFakeSlices(Student)
	u.teachers = initFakeSlices(Teacher)
	return &u
}

func initFakeSlices(role UserRole) []UserDomain {
	var users []UserDomain
	for i := range 20 {
		name := fmt.Sprintf("%d%d", role, i)
		d := UserDomain{
			ID:         ID(i),
			FirstName:  name,
			LastName:   name,
			MiddleName: name,
			Role:       role,
			CreatedAt:  time.Now(),
		}
		users = append(users, d)
	}
	return users
}

func (u *FakeUser) Login(ctx context.Context, p UserLoginParam) (JWT, error) {
	switch p.Login {
	case "admin", "teacher", "student", "root":
		tm := map[string]any{
			"Subject": 1,
			"Role":    p.Login,
		}
		data, _ := json.Marshal(&tm)
		return JWT{AccessToken: string(data), RefreshToken: string(data)}, nil
	}
	return JWT{}, fmt.Errorf("email is not role")
}

func (u *FakeUser) Get(ctx context.Context, auth Identity) ([]UserDomain, error) {
	switch auth.Role {
	case Admin:
		return append(u.students, u.teachers...), nil
	case Teacher:
		return u.students, nil
	case Student:
		return nil, nil
	case Root:
		return append(append(u.admins, u.students...), u.teachers...), nil
	default:
		return nil, nil
	}
}

func (u *FakeUser) Create(ctx context.Context, auth Identity, p UserCreateParam) error {
	return nil
}

func (u *FakeUser) GetMe(ctx context.Context, auth Identity) (UserDomain, error) {
	switch auth.Role {
	case Admin:
		return u.admins[0], nil
	case Teacher:
		return u.teachers[0], nil
	case Student:
		return u.students[0], nil
	case Root:
		return u.roots[0], nil
	default:
		return UserDomain{}, nil
	}
}

func (u *FakeUser) GetByID(ctx context.Context, auth Identity, id ID) (UserDomain, error) {
	return u.roots[0], nil
}

func (u *FakeUser) Put(ctx context.Context, auth Identity, id ID, p UserPutParam) error {
	return nil
}

func (u *FakeUser) Delete(ctx context.Context, auth Identity, id ID) error {
	return nil
}
