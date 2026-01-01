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

func NewSliceGroupResp(domains []*domain.Group) []*GroupResp {
	resp := make([]*GroupResp, len(domains))
	for i, d := range domains {
		resp[i] = NewGroupResp(d)
	}
	return resp
}

type GroupCreateReq struct {
	Name         string    `json:"name"`
	CuratorID    common.ID `json:"curator_id"`
	SpecialityID common.ID `json:"speciality_id"`
}

type GroupAddStudentsReq struct {
	StudentIDs []common.ID `json:"student_ids"`
}
