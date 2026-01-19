package domain

import "server/archive/internal/common"

type GroupName string

type GroupID common.ID

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

func (g Group) Copy() *Group {
	if g.Curator != nil {
		g.Curator = g.Curator.Copy()
	}
	for i, student := range g.Students {
		if student != nil {
			g.Students[i] = student.Copy()
		}
	}
	return &g
}
