package dto

import "server/internal/domain"

type SubjectResp struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

func NewSubjectResp(d *domain.Subject) *SubjectResp {
	return &SubjectResp{
		ID:   ID(d.ID),
		Name: string(d.Name),
	}
}

type SubjectCreateReq struct {
	Name          string `json:"name"`
	SpecialityIds []ID   `json:"speciality_ids"`
}
