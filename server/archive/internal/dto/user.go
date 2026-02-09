package dto

import (
	"fmt"
	"server/archive/internal/common"
	"server/archive/internal/domain"
)

type Identity struct {
	ID   common.ID       `json:"id"`
	Role domain.UserRole `json:"role"`
}

type UserShortResp struct {
	ID        common.ID       `json:"id"`
	ShortName string          `json:"short_name"`
	Role      domain.UserRole `json:"role"`
}

func NewUserShortResp(d *domain.User) *UserShortResp {
	if d == nil {
		return nil
	}
	resp := &UserShortResp{
		ID:   d.ID,
		Role: d.Role,
	}
	resp.ShortName = resp.formatShortName(
		string(d.FirstName),
		string(d.LastName),
		string(d.MiddleName),
	)
	return resp
}

func NewSliceUserShortResp(d []*domain.User) []*UserShortResp {
	resp := make([]*UserShortResp, len(d))
	for i := range resp {
		resp[i] = NewUserShortResp(d[i])
	}
	return resp
}

func (u *UserShortResp) formatShortName(f string, l string, m string) string {
	if f == "" || l == "" {
		return ""
	}
	name := fmt.Sprintf("%s %s.", f, string(l[0]))
	if m != "" {
		name = fmt.Sprintf("%s%s.", name, string(m[0]))
	}
	return name
}

type UserFullResp struct {
	ID         common.ID       `json:"id"`
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	MiddleName string          `json:"middle_name"`
	Group      *GroupShortResp `json:"group"`
	Role       domain.UserRole `json:"role"`
}

func NewUserFullResp(d *domain.User, group *domain.Group) *UserFullResp {
	return &UserFullResp{
		ID:         d.ID,
		FirstName:  string(d.FirstName),
		LastName:   string(d.LastName),
		MiddleName: string(d.MiddleName),
		Role:       d.Role,
		Group:      NewGroupShortResp(group),
	}
}

type LoginResp struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func NewLoginResp(d *domain.TokenPair) *LoginResp {
	return &LoginResp{
		RefreshToken: d.Refresh,
		AccessToken:  d.Access,
	}
}

type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserCreateReq struct {
	FirstName  string          `json:"first_name"`
	LastName   string          `json:"last_name"`
	MiddleName string          `json:"middle_name,omitempty"`
	Email      string          `json:"email"`
	Role       domain.UserRole `json:"role"`
}
