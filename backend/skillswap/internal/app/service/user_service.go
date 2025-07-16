package service

import (
	"errors"
	"time"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/repository"
	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/google/uuid"
)

type UserService interface {
	GetProfile(userID uuid.UUID) (*UserProfileResponse, error)
	UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) error
	SearchUsers(req *SearchUsersRequest) (*SearchUsersResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// DTOs for API responses
type UserProfileResponse struct {
	UserID        uuid.UUID       `json:"user_id"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Location      *string         `json:"location"`
	PhotoURL      *string         `json:"photo_url"`
	IsPublic      bool            `json:"is_public"`
	SkillsOffered []SkillResponse `json:"skills_offered"`
	SkillsWanted  []SkillResponse `json:"skills_wanted"`
	CreatedAt     time.Time       `json:"created_at"`
}

type SkillResponse struct {
	SkillID uuid.UUID `json:"skill_id"`
	Name    string    `json:"name"`
}

type UpdateProfileRequest struct {
	Name     *string `json:"name,omitempty"`
	Location *string `json:"location,omitempty"`
	PhotoURL *string `json:"photo_url,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

type SearchUsersRequest struct {
	Location   string `json:"location,omitempty"`
	SearchTerm string `json:"search_term,omitempty"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type SearchUsersResponse struct {
	Users      []UserProfileResponse `json:"users"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
}

func (s *userService) GetProfile(userID uuid.UUID) (*UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return s.toUserProfileResponse(user), nil
}

func (s *userService) UpdateProfile(userID uuid.UUID, req *UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if req.Name != nil {
		if len(*req.Name) < 2 || len(*req.Name) > 100 {
			return errors.New("name must be between 2 and 100 characters")
		}
		user.Name = *req.Name
	}
	if req.Location != nil {
		if len(*req.Location) > 100 {
			return errors.New("location must be less than 100 characters")
		}
		user.Location = req.Location
	}
	if req.PhotoURL != nil {
		if len(*req.PhotoURL) > 255 {
			return errors.New("photo URL must be less than 255 characters")
		}
		user.PhotoURL = req.PhotoURL
	}
	if req.IsPublic != nil {
		user.IsPublic = *req.IsPublic
	}

	return s.userRepo.Update(user)
}

func (s *userService) SearchUsers(req *SearchUsersRequest) (*SearchUsersResponse, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Page < 1 {
		req.Page = 1
	}

	offset := (req.Page - 1) * req.Limit

	filters := repository.UserFilters{
		IsPublic:   boolPtr(true), // Only show public profiles
		Location:   req.Location,
		SearchTerm: req.SearchTerm,
	}

	users, total, err := s.userRepo.List(req.Limit, offset, filters)
	if err != nil {
		return nil, err
	}

	userResponses := make([]UserProfileResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserProfileResponse(user)
	}

	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return &SearchUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}, nil
}

// Helper functions
func (s *userService) toUserProfileResponse(user *models.User) *UserProfileResponse {
	skillsOffered := make([]SkillResponse, len(user.SkillsOffered))
	for i, skill := range user.SkillsOffered {
		skillsOffered[i] = SkillResponse{
			SkillID: skill.Skill.SkillID,
			Name:    skill.Skill.Name,
		}
	}

	skillsWanted := make([]SkillResponse, len(user.SkillsWanted))
	for i, skill := range user.SkillsWanted {
		skillsWanted[i] = SkillResponse{
			SkillID: skill.Skill.SkillID,
			Name:    skill.Skill.Name,
		}
	}

	return &UserProfileResponse{
		UserID:        user.UserID,
		Name:          user.Name,
		Email:         user.Email,
		Location:      user.Location,
		PhotoURL:      user.PhotoURL,
		IsPublic:      user.IsPublic,
		SkillsOffered: skillsOffered,
		SkillsWanted:  skillsWanted,
		CreatedAt:     user.CreatedAt,
	}
}

func boolPtr(b bool) *bool {
	return &b
}
