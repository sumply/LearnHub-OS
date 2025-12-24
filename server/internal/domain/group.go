package domain

type GroupID uint64

type GroupName string

func NewGroupName(s string) (GroupName, error) {
	return GroupName(s), nil
}

type Group struct {
	ID   GroupID
	Name GroupName

	Curator    *User
	Speciality *Speciality
}

func NewGroup(name string, curator UserID, speciality SpecialityID) (*Group, error) {
	n, err := NewGroupName(name)
	if err != nil {
		return nil, err
	}
	return &Group{
		Name:       n,
		Curator:    &User{ID: curator},
		Speciality: &Speciality{ID: speciality},
	}, nil
}
