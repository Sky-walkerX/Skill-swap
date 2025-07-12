package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	SlotID     uuid.UUID `gorm:"type:uuid;primaryKey;column:slot_id;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;column:user_id;index"`
	Label      string    `gorm:"not null"`
	DayBitmask int32     `gorm:"column:day_bitmask;not null"`
	StartTime  time.Time `gorm:"column:start_time;type:time;not null"`
	EndTime    time.Time `gorm:"column:end_time;type:time;not null"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relations - restored
	User User `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// BeforeCreate is called by GORM before creating an AvailabilitySlot record
func (a *AvailabilitySlot) BeforeCreate(tx *gorm.DB) (err error) {
	if a.SlotID == uuid.Nil {
		a.SlotID = uuid.New()
	}
	return
}

func (AvailabilitySlot) TableName() string { return "availability_slots" }
