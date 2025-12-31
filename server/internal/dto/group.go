package dto

import (
	"server/internal/common"
	"server/internal/domain"
)

type GroupResp struct {
	ID      common.ID      `json:"id"`
	Name    string         `json:"name"`
	Curator *UserShortResp `json:"curator"`
}

func NewGroupResp(d *domain.Group) *GroupResp {
	return &GroupResp{
		ID:      d.ID,
		Name:    string(d.Name),
		Curator: NewUserShortResp(d.Curator),
	}
}

type GroupCreateReq struct {
	Name         string    `json:"name"`
	CuratorID    common.ID `json:"curator_id"`
	SpecialityID common.ID `json:"speciality_id"`
}
