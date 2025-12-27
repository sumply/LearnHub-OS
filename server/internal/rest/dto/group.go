package dto

import "server/internal/domain"

type GroupResp struct {
	ID         ID              `json:"id"`
	Name       string          `json:"name"`
	Curator    *UserShortResp  `json:"curator"`
	Speciality *SpecialityResp `json:"speciality"`
}

func NewGroupResp(d *domain.Group) *GroupResp {
	return &GroupResp{
		ID:         ID(d.ID),
		Name:       string(d.Name),
		Curator:    NewUserShortResp(d.Curator),
		Speciality: NewSpecialityResp(d.Speciality),
	}
}

type GroupCreateReq struct {
	Name         string `json:"name"`
	CuratorID    ID     `json:"curator_id"`
	SpecialityID ID     `json:"speciality_id"`
}
