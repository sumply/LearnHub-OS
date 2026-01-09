package dto

import (
	"server/internal/common"
	"server/internal/domain"
)

type SubjectResp struct {
	ID   common.ID `json:"id"`
	Name string    `json:"name"`
}

func NewSubjectResp(d *domain.Subject) *SubjectResp {
	if d == nil {
		return nil
	}
	return &SubjectResp{
		ID:   d.ID,
		Name: string(d.Name),
	}
}

type SubjectCreateReq struct {
	Name string `json:"name"`
}
