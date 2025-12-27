package dto

import "server/internal/domain"

type SpecialityCreateReq struct {
	Name string `json:"name"`
}

type SpecialityResp struct {
	ID   ID     `json:"id"`
	Name string `json:"name"`
}

func NewSpecialityResp(d *domain.Speciality) *SpecialityResp {
	return &SpecialityResp{
		ID:   ID(d.ID),
		Name: string(d.Name),
	}
}
