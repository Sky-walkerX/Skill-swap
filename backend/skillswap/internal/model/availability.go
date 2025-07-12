package models

import (
	"time"

	"github.com/google/uuid"
)

type AvailabilitySlot struct {
	SlotID     uuid.UUID `gorm:"type:uuid;primaryKey;column:slot_id"`
	UserID     uuid.UUID `gorm:"type:uuid;column:user_id;index"`
	Label      string    `gorm:"not null"`
	DayBitmask int32     `gorm:"column:day_bitmask;not null"`
	StartTime  time.Time `gorm:"column:start_time;type:time;not null"`
	EndTime    time.Time `gorm:"column:end_time;type:time;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`
	// FK
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (AvailabilitySlot) TableName() string { return "availability_slots" }