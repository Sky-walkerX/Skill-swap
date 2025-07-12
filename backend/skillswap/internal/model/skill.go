package models

import (
	"time"

	"github.com/google/uuid"
)

type Skill struct {
	SkillID   uuid.UUID `gorm:"type:uuid;primaryKey;column:skill_id"`
	Name      string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Skill) TableName() string { return "skills" }