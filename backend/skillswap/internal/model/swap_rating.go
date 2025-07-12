package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwapRating struct {
	RatingID  uuid.UUID `gorm:"type:uuid;primaryKey;column:rating_id;default:gen_random_uuid()"`
	SwapID    uuid.UUID `gorm:"type:uuid;column:swap_id;index"`
	RaterID   uuid.UUID `gorm:"type:uuid;column:rater_id"`
	RateeID   uuid.UUID `gorm:"type:uuid;column:ratee_id"`
	Score     int16     `gorm:"check:score >= 1 AND score <= 5"`
	Comment   *string
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relations - restored
	Swap  SwapRequest `gorm:"foreignKey:SwapID;references:SwapID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Rater User        `gorm:"foreignKey:RaterID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Ratee User        `gorm:"foreignKey:RateeID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// BeforeCreate is called by GORM before creating a SwapRating record
func (r *SwapRating) BeforeCreate(tx *gorm.DB) (err error) {
	if r.RatingID == uuid.Nil {
		r.RatingID = uuid.New()
	}
	return
}

func (SwapRating) TableName() string { return "swap_ratings" }