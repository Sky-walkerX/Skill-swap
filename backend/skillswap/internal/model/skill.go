package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Skill struct {
	SkillID   uuid.UUID `gorm:"type:uuid;primaryKey;column:skill_id;default:gen_random_uuid()"`
	Name      string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

// BeforeCreate is called by GORM before creating a Skill record
func (s *Skill) BeforeCreate(tx *gorm.DB) (err error) {
	if s.SkillID == uuid.Nil {
		s.SkillID = uuid.New()
	}
	return
}

func (Skill) TableName() string { return "skills" }