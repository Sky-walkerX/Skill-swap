package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UserID       uuid.UUID      `gorm:"type:uuid;primaryKey;column:user_id;default:gen_random_uuid()"`
	Name         string         `gorm:"not null"`
	Email        string         `gorm:"uniqueIndex;not null"`
	PasswordHash string         `gorm:"column:password_hash;not null"`
	Location     *string
	PhotoURL     *string        `gorm:"column:photo_url"`
	IsPublic     bool           `gorm:"column:is_public;default:true"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
	
	// Relations - restored
	SkillsOffered     []UserSkillOffered `gorm:"foreignKey:UserID;references:UserID"`
	SkillsWanted      []UserSkillWanted  `gorm:"foreignKey:UserID;references:UserID"`
	AvailabilitySlots []AvailabilitySlot `gorm:"foreignKey:UserID;references:UserID"`
}

// BeforeCreate is called by GORM before creating a User record
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return
}

func (User) TableName() string { return "users" }