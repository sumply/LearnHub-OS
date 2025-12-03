package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type FakeUser struct {
	admins   []UserData
	teachers []UserData
	students []UserData
	roots    []UserData
}

func NewFakeUser() *FakeUser {
	u := FakeUser{}
	u.admins = initFakeSlices("admin")
	u.roots = initFakeSlices("root")
	u.students = initFakeSlices("student")
	u.teachers = initFakeSlices("teachers")
	return &u
}

func initFakeSlices(role string) []UserData {
	var users []UserData
	for i := range 20 {
		name := fmt.Sprintf("%s%d", role, i)
		d := UserData{
			ID:         int64(i),
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
	switch p.Email {
	case "admin", "teacher", "student", "root":
		tm := map[string]any{
			"Subject": 1,
			"Role":    p.Email,
		}
		data, _ := json.Marshal(&tm)
		return JWT{AccessToken: string(data), RefreshToken: string(data)}, nil
	}
	return JWT{}, fmt.Errorf("email is not role")
}

func (u *FakeUser) Get(ctx context.Context, auth AuthData, p UserGetParam) ([]UserData, error) {
	switch auth.Role {
	case "admin":
		return append(u.students, u.teachers...), nil
	case "teacher":
		return u.students, nil
	case "student":
		return nil, nil
	case "root":
		return append(append(u.admins, u.students...), u.teachers...), nil
	default:
		return nil, nil
	}
}

func (u *FakeUser) Create(ctx context.Context, auth AuthData, p UserCreateParam) error {
	return nil
}

func (u *FakeUser) GetMe(ctx context.Context, auth AuthData) (UserData, error) {
	switch auth.Role {
	case "admin":
		return u.admins[0], nil
	case "teacher":
		return u.teachers[0], nil
	case "student":
		return u.students[0], nil
	case "root":
		return u.roots[0], nil
	default:
		return UserData{}, nil
	}
}

func (u *FakeUser) GetByID(ctx context.Context, auth AuthData, id int64) (UserData, error) {
	return u.roots[0], nil
}

func (u *FakeUser) Put(ctx context.Context, auth AuthData, id int64, p UserPutParam) error {
	return nil
}

func (u *FakeUser) Delete(ctx context.Context, auth AuthData, id int64) error {
	return nil
}
