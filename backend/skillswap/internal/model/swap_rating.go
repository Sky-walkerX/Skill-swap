package models

import (
	"time"

	"github.com/google/uuid"
)

type SwapRating struct {
	RatingID  uuid.UUID `gorm:"type:uuid;primaryKey;column:rating_id"`
	SwapID    uuid.UUID `gorm:"type:uuid;column:swap_id;index"`
	RaterID   uuid.UUID `gorm:"type:uuid;column:rater_id"`
	RateeID   uuid.UUID `gorm:"type:uuid;column:ratee_id"`
	Score     int16     `gorm:"check:score >= 1 AND score <= 5"`
	Comment   *string
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`

	// Relations
	Swap  SwapRequest `gorm:"foreignKey:SwapID;constraint:OnDelete:CASCADE"`
	Rater User        `gorm:"foreignKey:RaterID"`
	Ratee User        `gorm:"foreignKey:RateeID"`
}

func (SwapRating) TableName() string { return "swap_ratings" }