package service

import (
	"errors"

	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillService interface {
	// Skill CRUD operations
	GetAllSkills() ([]models.Skill, error)
	GetSkillByID(skillID uuid.UUID) (*models.Skill, error)
	CreateSkill(name string) (*models.Skill, error)
	UpdateSkill(skillID uuid.UUID, name string) (*models.Skill, error)
	DeleteSkill(skillID uuid.UUID) error

	// User skill management
	AddOfferedSkill(userID, skillID uuid.UUID) error
	RemoveOfferedSkill(userID, skillID uuid.UUID) error
	AddWantedSkill(userID, skillID uuid.UUID) error
	RemoveWantedSkill(userID, skillID uuid.UUID) error

	// User skill queries
	GetUserOfferedSkills(userID uuid.UUID) ([]models.Skill, error)
	GetUserWantedSkills(userID uuid.UUID) ([]models.Skill, error)
	GetUsersWithOfferedSkill(skillID uuid.UUID) ([]models.User, error)
	GetUsersWithWantedSkill(skillID uuid.UUID) ([]models.User, error)
}

type skillService struct {
	db *gorm.DB
}

func NewSkillService(db *gorm.DB) SkillService {
	return &skillService{db: db}
}

// GetAllSkills retrieves all available skills
func (s *skillService) GetAllSkills() ([]models.Skill, error) {
	var skills []models.Skill
	err := s.db.Find(&skills).Error
	return skills, err
}

// GetSkillByID retrieves a skill by its ID
func (s *skillService) GetSkillByID(skillID uuid.UUID) (*models.Skill, error) {
	var skill models.Skill
	err := s.db.First(&skill, "skill_id = ?", skillID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("skill not found")
		}
		return nil, err
	}
	return &skill, nil
}

// CreateSkill creates a new skill
func (s *skillService) CreateSkill(name string) (*models.Skill, error) {
	skill := &models.Skill{
		Name: name,
	}

	err := s.db.Create(skill).Error
	if err != nil {
		return nil, err
	}

	return skill, nil
}

// UpdateSkill updates an existing skill
func (s *skillService) UpdateSkill(skillID uuid.UUID, name string) (*models.Skill, error) {
	skill, err := s.GetSkillByID(skillID)
	if err != nil {
		return nil, err
	}

	skill.Name = name
	err = s.db.Save(skill).Error
	if err != nil {
		return nil, err
	}

	return skill, nil
}

// DeleteSkill deletes a skill (admin only)
func (s *skillService) DeleteSkill(skillID uuid.UUID) error {
	// Check if skill exists
	_, err := s.GetSkillByID(skillID)
	if err != nil {
		return err
	}

	// Check if skill is referenced in user skills or swap requests
	var offeredCount, wantedCount int64
	s.db.Model(&models.UserSkillOffered{}).Where("skill_id = ?", skillID).Count(&offeredCount)
	s.db.Model(&models.UserSkillWanted{}).Where("skill_id = ?", skillID).Count(&wantedCount)

	if offeredCount > 0 || wantedCount > 0 {
		return errors.New("skill is in use and cannot be deleted")
	}

	return s.db.Delete(&models.Skill{}, "skill_id = ?", skillID).Error
}

// AddOfferedSkill adds a skill to user's offered skills
func (s *skillService) AddOfferedSkill(userID, skillID uuid.UUID) error {
	// Check if skill exists
	_, err := s.GetSkillByID(skillID)
	if err != nil {
		return err
	}

	// Check if already exists
	var count int64
	s.db.Model(&models.UserSkillOffered{}).Where("user_id = ? AND skill_id = ?", userID, skillID).Count(&count)
	if count > 0 {
		return errors.New("skill already in offered skills")
	}

	userSkill := &models.UserSkillOffered{
		UserID:  userID,
		SkillID: skillID,
	}

	return s.db.Create(userSkill).Error
}

// RemoveOfferedSkill removes a skill from user's offered skills
func (s *skillService) RemoveOfferedSkill(userID, skillID uuid.UUID) error {
	result := s.db.Delete(&models.UserSkillOffered{}, "user_id = ? AND skill_id = ?", userID, skillID)
	if result.RowsAffected == 0 {
		return errors.New("offered skill not found")
	}
	return result.Error
}

// AddWantedSkill adds a skill to user's wanted skills
func (s *skillService) AddWantedSkill(userID, skillID uuid.UUID) error {
	// Check if skill exists
	_, err := s.GetSkillByID(skillID)
	if err != nil {
		return err
	}

	// Check if already exists
	var count int64
	s.db.Model(&models.UserSkillWanted{}).Where("user_id = ? AND skill_id = ?", userID, skillID).Count(&count)
	if count > 0 {
		return errors.New("skill already in wanted skills")
	}

	userSkill := &models.UserSkillWanted{
		UserID:  userID,
		SkillID: skillID,
	}

	return s.db.Create(userSkill).Error
}

// RemoveWantedSkill removes a skill from user's wanted skills
func (s *skillService) RemoveWantedSkill(userID, skillID uuid.UUID) error {
	result := s.db.Delete(&models.UserSkillWanted{}, "user_id = ? AND skill_id = ?", userID, skillID)
	if result.RowsAffected == 0 {
		return errors.New("wanted skill not found")
	}
	return result.Error
}

// GetUserOfferedSkills retrieves all skills offered by a user
func (s *skillService) GetUserOfferedSkills(userID uuid.UUID) ([]models.Skill, error) {
	var skills []models.Skill
	err := s.db.Table("skills").
		Joins("JOIN user_skills_offered ON skills.skill_id = user_skills_offered.skill_id").
		Where("user_skills_offered.user_id = ?", userID).
		Find(&skills).Error

	return skills, err
}

// GetUserWantedSkills retrieves all skills wanted by a user
func (s *skillService) GetUserWantedSkills(userID uuid.UUID) ([]models.Skill, error) {
	var skills []models.Skill
	err := s.db.Table("skills").
		Joins("JOIN user_skills_wanted ON skills.skill_id = user_skills_wanted.skill_id").
		Where("user_skills_wanted.user_id = ?", userID).
		Find(&skills).Error

	return skills, err
}

// GetUsersWithOfferedSkill retrieves all users who offer a specific skill
func (s *skillService) GetUsersWithOfferedSkill(skillID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := s.db.Table("users").
		Joins("JOIN user_skills_offered ON users.user_id = user_skills_offered.user_id").
		Where("user_skills_offered.skill_id = ? AND users.is_public = true AND users.deleted_at IS NULL", skillID).
		Find(&users).Error

	return users, err
}

// GetUsersWithWantedSkill retrieves all users who want a specific skill
func (s *skillService) GetUsersWithWantedSkill(skillID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := s.db.Table("users").
		Joins("JOIN user_skills_wanted ON users.user_id = user_skills_wanted.user_id").
		Where("user_skills_wanted.skill_id = ? AND users.is_public = true AND users.deleted_at IS NULL", skillID).
		Find(&users).Error

	return users, err
}
