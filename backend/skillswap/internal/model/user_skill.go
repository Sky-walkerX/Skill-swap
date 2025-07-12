package models

import "github.com/google/uuid"

// UserSkillOffered represents table user_skills_offered
type UserSkillOffered struct {
	UserID  uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	SkillID uuid.UUID `gorm:"type:uuid;primaryKey;column:skill_id"`
	// FK constraints
	User  User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Skill Skill `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}

func (UserSkillOffered) TableName() string { return "user_skills_offered" }

// UserSkillWanted represents table user_skills_wanted
type UserSkillWanted struct {
	UserID  uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	SkillID uuid.UUID `gorm:"type:uuid;primaryKey;column:skill_id"`
	User    User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Skill   Skill     `gorm:"foreignKey:SkillID;constraint:OnDelete:CASCADE"`
}

func (UserSkillWanted) TableName() string { return "user_skills_wanted" }