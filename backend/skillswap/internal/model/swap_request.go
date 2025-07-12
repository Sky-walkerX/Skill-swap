package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwapStatus string

const (
	StatusPending   SwapStatus = "pending"
	StatusAccepted  SwapStatus = "accepted"
	StatusRejected  SwapStatus = "rejected"
	StatusCancelled SwapStatus = "cancelled"
)

type SwapRequest struct {
	SwapID         uuid.UUID   `gorm:"type:uuid;primaryKey;column:swap_id"`
	RequesterID    uuid.UUID   `gorm:"type:uuid;column:requester_id;index"`
	ResponderID    uuid.UUID   `gorm:"type:uuid;column:responder_id;index"`
	OfferedSkillID uuid.UUID   `gorm:"type:uuid;column:offered_skill_id"`
	WantedSkillID  uuid.UUID   `gorm:"type:uuid;column:wanted_skill_id"`
	Status         SwapStatus  `gorm:"type:swap_status;default:'pending'"`
	CreatedAt      time.Time   `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time   `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Relations
	Requester    User  `gorm:"foreignKey:RequesterID;constraint:OnDelete:CASCADE"`
	Responder    User  `gorm:"foreignKey:ResponderID;constraint:OnDelete:CASCADE"`
	OfferedSkill Skill `gorm:"foreignKey:OfferedSkillID"`
	WantedSkill  Skill `gorm:"foreignKey:WantedSkillID"`
	Ratings      []SwapRating `gorm:"foreignKey:SwapID"`
}

func (SwapRequest) TableName() string { return "swap_requests" }