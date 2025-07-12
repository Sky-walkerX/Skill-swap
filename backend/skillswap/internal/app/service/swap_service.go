package service

import (
	"errors"

	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwapService interface {
	// Swap request CRUD operations
	CreateSwapRequest(req *CreateSwapRequestDTO) (*models.SwapRequest, error)
	GetSwapRequestByID(swapID uuid.UUID) (*models.SwapRequest, error)
	GetUserSwapRequests(userID uuid.UUID, filter SwapRequestFilter) ([]models.SwapRequest, error)
	UpdateSwapStatus(swapID uuid.UUID, userID uuid.UUID, status models.SwapStatus) (*models.SwapRequest, error)
	DeleteSwapRequest(swapID uuid.UUID, userID uuid.UUID) error

	// Swap request queries
	GetSwapRequestsForUser(userID uuid.UUID) (*SwapRequestsResponse, error)
	GetPendingSwapRequests(userID uuid.UUID) ([]models.SwapRequest, error)
	GetSwapHistory(userID uuid.UUID) ([]models.SwapRequest, error)

	// Matching and recommendations
	FindPotentialMatches(userID uuid.UUID) ([]SwapMatch, error)
}

// DTOs and Request structures
type CreateSwapRequestDTO struct {
	RequesterID    uuid.UUID `json:"requester_id"`
	ResponderID    uuid.UUID `json:"responder_id" binding:"required"`
	OfferedSkillID uuid.UUID `json:"offered_skill_id" binding:"required"`
	WantedSkillID  uuid.UUID `json:"wanted_skill_id" binding:"required"`
}

type SwapRequestFilter struct {
	Status   *models.SwapStatus `json:"status,omitempty"`
	Sent     bool               `json:"sent,omitempty"`     // Requests sent by user
	Received bool               `json:"received,omitempty"` // Requests received by user
	Limit    int                `json:"limit,omitempty"`
	Offset   int                `json:"offset,omitempty"`
}

type SwapRequestsResponse struct {
	Sent     []models.SwapRequest `json:"sent"`
	Received []models.SwapRequest `json:"received"`
}

type SwapMatch struct {
	User         models.User  `json:"user"`
	OfferedSkill models.Skill `json:"offered_skill"`
	WantedSkill  models.Skill `json:"wanted_skill"`
	MatchScore   int          `json:"match_score"` // 1-100 compatibility score
}

type swapService struct {
	db *gorm.DB
}

func NewSwapService(db *gorm.DB) SwapService {
	return &swapService{db: db}
}

// CreateSwapRequest creates a new swap request
func (s *swapService) CreateSwapRequest(req *CreateSwapRequestDTO) (*models.SwapRequest, error) {
	// Validate that requester and responder are different
	if req.RequesterID == req.ResponderID {
		return nil, errors.New("cannot create swap request with yourself")
	}

	// Validate that requester offers the offered skill
	var offeredCount int64
	s.db.Model(&models.UserSkillOffered{}).
		Where("user_id = ? AND skill_id = ?", req.RequesterID, req.OfferedSkillID).
		Count(&offeredCount)
	if offeredCount == 0 {
		return nil, errors.New("you don't offer the specified skill")
	}

	// Validate that responder wants the offered skill
	var wantedCount int64
	s.db.Model(&models.UserSkillWanted{}).
		Where("user_id = ? AND skill_id = ?", req.ResponderID, req.OfferedSkillID).
		Count(&wantedCount)
	if wantedCount == 0 {
		return nil, errors.New("responder doesn't want the offered skill")
	}

	// Validate that responder offers the wanted skill
	var responderOffersCount int64
	s.db.Model(&models.UserSkillOffered{}).
		Where("user_id = ? AND skill_id = ?", req.ResponderID, req.WantedSkillID).
		Count(&responderOffersCount)
	if responderOffersCount == 0 {
		return nil, errors.New("responder doesn't offer the requested skill")
	}

	// Check for existing pending request between same users and skills
	var existingCount int64
	s.db.Model(&models.SwapRequest{}).
		Where("requester_id = ? AND responder_id = ? AND offered_skill_id = ? AND wanted_skill_id = ? AND status = ?",
			req.RequesterID, req.ResponderID, req.OfferedSkillID, req.WantedSkillID, models.StatusPending).
		Count(&existingCount)
	if existingCount > 0 {
		return nil, errors.New("pending swap request already exists")
	}

	swapRequest := &models.SwapRequest{
		RequesterID:    req.RequesterID,
		ResponderID:    req.ResponderID,
		OfferedSkillID: req.OfferedSkillID,
		WantedSkillID:  req.WantedSkillID,
		Status:         models.StatusPending,
	}

	err := s.db.Create(swapRequest).Error
	if err != nil {
		return nil, err
	}

	// Load relationships
	err = s.db.Preload("Requester").Preload("Responder").
		Preload("OfferedSkill").Preload("WantedSkill").
		First(swapRequest, swapRequest.SwapID).Error

	return swapRequest, err
}

// GetSwapRequestByID retrieves a swap request by ID
func (s *swapService) GetSwapRequestByID(swapID uuid.UUID) (*models.SwapRequest, error) {
	var swapRequest models.SwapRequest
	err := s.db.Preload("Requester").Preload("Responder").
		Preload("OfferedSkill").Preload("WantedSkill").
		First(&swapRequest, "swap_id = ?", swapID).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("swap request not found")
		}
		return nil, err
	}

	return &swapRequest, nil
}

// GetUserSwapRequests retrieves swap requests for a user with filtering
func (s *swapService) GetUserSwapRequests(userID uuid.UUID, filter SwapRequestFilter) ([]models.SwapRequest, error) {
	query := s.db.Model(&models.SwapRequest{}).
		Preload("Requester").Preload("Responder").
		Preload("OfferedSkill").Preload("WantedSkill")

	// Apply filters
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.Sent && filter.Received {
		query = query.Where("requester_id = ? OR responder_id = ?", userID, userID)
	} else if filter.Sent {
		query = query.Where("requester_id = ?", userID)
	} else if filter.Received {
		query = query.Where("responder_id = ?", userID)
	} else {
		// Default: show both sent and received
		query = query.Where("requester_id = ? OR responder_id = ?", userID, userID)
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Order by creation date (newest first)
	query = query.Order("created_at DESC")

	var swapRequests []models.SwapRequest
	err := query.Find(&swapRequests).Error

	return swapRequests, err
}

// UpdateSwapStatus updates the status of a swap request
func (s *swapService) UpdateSwapStatus(swapID uuid.UUID, userID uuid.UUID, status models.SwapStatus) (*models.SwapRequest, error) {
	// Get the swap request
	swapRequest, err := s.GetSwapRequestByID(swapID)
	if err != nil {
		return nil, err
	}

	// Only responder can accept/reject requests
	if status == models.StatusAccepted || status == models.StatusRejected {
		if swapRequest.ResponderID != userID {
			return nil, errors.New("only responder can accept or reject requests")
		}
	}

	// Only requester or responder can cancel
	if status == models.StatusCancelled {
		if swapRequest.RequesterID != userID && swapRequest.ResponderID != userID {
			return nil, errors.New("only requester or responder can cancel requests")
		}
	}

	// Validate status transitions
	if swapRequest.Status != models.StatusPending && status != models.StatusCancelled {
		return nil, errors.New("can only modify pending requests")
	}

	// Update status
	swapRequest.Status = status
	err = s.db.Save(swapRequest).Error
	if err != nil {
		return nil, err
	}

	return swapRequest, nil
}

// DeleteSwapRequest deletes a swap request (only requester can delete)
func (s *swapService) DeleteSwapRequest(swapID uuid.UUID, userID uuid.UUID) error {
	swapRequest, err := s.GetSwapRequestByID(swapID)
	if err != nil {
		return err
	}

	// Only requester can delete
	if swapRequest.RequesterID != userID {
		return errors.New("only requester can delete swap requests")
	}

	// Can only delete pending requests
	if swapRequest.Status != models.StatusPending {
		return errors.New("can only delete pending requests")
	}

	return s.db.Delete(&models.SwapRequest{}, "swap_id = ?", swapID).Error
}

// GetSwapRequestsForUser retrieves organized swap requests for a user
func (s *swapService) GetSwapRequestsForUser(userID uuid.UUID) (*SwapRequestsResponse, error) {
	// Get sent requests
	sentFilter := SwapRequestFilter{Sent: true, Limit: 50}
	sent, err := s.GetUserSwapRequests(userID, sentFilter)
	if err != nil {
		return nil, err
	}

	// Get received requests
	receivedFilter := SwapRequestFilter{Received: true, Limit: 50}
	received, err := s.GetUserSwapRequests(userID, receivedFilter)
	if err != nil {
		return nil, err
	}

	return &SwapRequestsResponse{
		Sent:     sent,
		Received: received,
	}, nil
}

// GetPendingSwapRequests retrieves pending swap requests for a user
func (s *swapService) GetPendingSwapRequests(userID uuid.UUID) ([]models.SwapRequest, error) {
	status := models.StatusPending
	filter := SwapRequestFilter{Status: &status, Limit: 100}
	return s.GetUserSwapRequests(userID, filter)
}

// GetSwapHistory retrieves completed swap requests for a user
func (s *swapService) GetSwapHistory(userID uuid.UUID) ([]models.SwapRequest, error) {
	var swapRequests []models.SwapRequest
	err := s.db.Model(&models.SwapRequest{}).
		Preload("Requester").Preload("Responder").
		Preload("OfferedSkill").Preload("WantedSkill").
		Where("(requester_id = ? OR responder_id = ?) AND status IN ?",
			userID, userID, []models.SwapStatus{models.StatusAccepted, models.StatusRejected, models.StatusCancelled}).
		Order("updated_at DESC").
		Limit(50).
		Find(&swapRequests).Error

	return swapRequests, err
}

// FindPotentialMatches finds potential swap matches for a user
func (s *swapService) FindPotentialMatches(userID uuid.UUID) ([]SwapMatch, error) {
	var matches []SwapMatch

	// Get user's offered skills
	var userOfferedSkills []models.Skill
	err := s.db.Table("skills").
		Joins("JOIN user_skills_offered ON skills.skill_id = user_skills_offered.skill_id").
		Where("user_skills_offered.user_id = ?", userID).
		Find(&userOfferedSkills).Error
	if err != nil {
		return nil, err
	}

	// Get user's wanted skills
	var userWantedSkills []models.Skill
	err = s.db.Table("skills").
		Joins("JOIN user_skills_wanted ON skills.skill_id = user_skills_wanted.skill_id").
		Where("user_skills_wanted.user_id = ?", userID).
		Find(&userWantedSkills).Error
	if err != nil {
		return nil, err
	}

	// Find users who want what we offer and offer what we want
	for _, offeredSkill := range userOfferedSkills {
		for _, wantedSkill := range userWantedSkills {
			// Find users who want our offered skill AND offer our wanted skill
			var potentialUsers []models.User
			err := s.db.Table("users").
				Joins("JOIN user_skills_wanted ON users.user_id = user_skills_wanted.user_id").
				Joins("JOIN user_skills_offered ON users.user_id = user_skills_offered.user_id").
				Where("user_skills_wanted.skill_id = ? AND user_skills_offered.skill_id = ? AND users.user_id != ? AND users.is_public = true AND users.deleted_at IS NULL",
					offeredSkill.SkillID, wantedSkill.SkillID, userID).
				Find(&potentialUsers).Error

			if err != nil {
				continue
			}

			for _, user := range potentialUsers {
				// Calculate match score (simple algorithm for now)
				matchScore := 80 // Base score for mutual skill match

				matches = append(matches, SwapMatch{
					User:         user,
					OfferedSkill: offeredSkill,
					WantedSkill:  wantedSkill,
					MatchScore:   matchScore,
				})
			}
		}
	}

	// Limit results
	if len(matches) > 20 {
		matches = matches[:20]
	}

	return matches, nil
}
