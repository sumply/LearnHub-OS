package domain

import "server/internal/common"

type GroupName string

func NewGroupName(s string) (GroupName, error) {
	return GroupName(s), nil
}

type Group struct {
	ID   common.ID
	Name GroupName

	Curator  *User
	Students []*User
}

func NewGroup(name string, curator common.ID) (*Group, error) {
	n, err := NewGroupName(name)
	if err != nil {
		return nil, err
	}
	return &Group{
		Name:    n,
		Curator: &User{ID: curator},
	}, nil
}
