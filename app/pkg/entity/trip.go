package entity

import (
	"github.com/haandol/devops-academy-eda-demo/pkg/dto"
)

type Trip struct {
	PK        string `json:"PK" binding:"required" validate:"required"`
	SK        string `json:"SK" binding:"required" validate:"required"`
	GS1PK     string `json:"GS1PK" binding:"required" validate:"required"`
	GS1SK     string `json:"GS1SK" binding:"required" validate:"required"`
	ID        string `json:"id" binding:"required" validate:"required"`
	Status    string `json:"status" binding:"required" validate:"required"`
	CreatedAt string `json:"createdAt" binding:"required" validate:"required"`
	UpdatedAt string `json:"updatedAt"`
}

func (t *Trip) DTO() dto.Trip {
	return dto.Trip{
		ID:        t.ID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
