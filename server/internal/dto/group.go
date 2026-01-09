package dto

import (
	"server/internal/common"
	"server/internal/domain"
)

type GroupShortResp struct {
	ID      common.ID      `json:"id"`
	Name    string         `json:"name"`
	Curator *UserShortResp `json:"curator"`
}

func NewGroupShortResp(d *domain.Group) *GroupShortResp {
	if d == nil {
		return nil
	}
	return &GroupShortResp{
		ID:      d.ID,
		Name:    string(d.Name),
		Curator: NewUserShortResp(d.Curator),
	}
}

func NewSliceGroupShortResp(domains []*domain.Group) []*GroupShortResp {
	resp := make([]*GroupShortResp, len(domains))
	for i, d := range domains {
		resp[i] = NewGroupShortResp(d)
	}
	return resp
}

type GroupCreateReq struct {
	Name      string    `json:"name"`
	CuratorID common.ID `json:"curator_id"`
}

type GroupAddStudentsReq struct {
	StudentIDs []common.ID `json:"student_ids"`
}

type GroupFullResp struct {
	ID       common.ID        `json:"id"`
	Name     string           `json:"name"`
	Curator  *UserShortResp   `json:"curator"`
	Students []*UserShortResp `json:"students"`
}

func NewGroupFullResp(d *domain.Group) *GroupFullResp {
	if d == nil {
		return nil
	}
	return &GroupFullResp{
		ID:       d.ID,
		Name:     string(d.Name),
		Curator:  NewUserShortResp(d.Curator),
		Students: NewSliceUserShortResp(d.Students),
	}
}
