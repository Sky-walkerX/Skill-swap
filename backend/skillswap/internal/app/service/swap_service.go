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
	GetSwapRequestByID(swapID uuid.UUID, userID uuid.UUID) (*models.SwapRequest, error)
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
	// Use a transaction to ensure atomicity
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	// Validate that requester and responder are different
	if req.RequesterID == req.ResponderID {
		tx.Rollback()
		return nil, errors.New("cannot create swap request with yourself")
	}

	// Validate that requester offers the offered skill
	var offeredCount int64
	tx.Model(&models.UserSkillOffered{}).
		Where("user_id = ? AND skill_id = ?", req.RequesterID, req.OfferedSkillID).
		Count(&offeredCount)
	if offeredCount == 0 {
		tx.Rollback()
		return nil, errors.New("you don't offer the specified skill")
	}

	// Validate that responder wants the offered skill
	var wantedCount int64
	tx.Model(&models.UserSkillWanted{}).
		Where("user_id = ? AND skill_id = ?", req.ResponderID, req.OfferedSkillID).
		Count(&wantedCount)
	if wantedCount == 0 {
		tx.Rollback()
		return nil, errors.New("responder doesn't want the offered skill")
	}

	// Validate that responder offers the wanted skill
	var responderOffersCount int64
	tx.Model(&models.UserSkillOffered{}).
		Where("user_id = ? AND skill_id = ?", req.ResponderID, req.WantedSkillID).
		Count(&responderOffersCount)
	if responderOffersCount == 0 {
		tx.Rollback()
		return nil, errors.New("responder doesn't offer the requested skill")
	}

	// Check for existing pending request between same users and skills
	var existingCount int64
	tx.Model(&models.SwapRequest{}).
		Where("requester_id = ? AND responder_id = ? AND offered_skill_id = ? AND wanted_skill_id = ? AND status = ?",
			req.RequesterID, req.ResponderID, req.OfferedSkillID, req.WantedSkillID, models.StatusPending).
		Count(&existingCount)
	if existingCount > 0 {
		tx.Rollback()
		return nil, errors.New("pending swap request already exists")
	}

	swapRequest := &models.SwapRequest{
		RequesterID:    req.RequesterID,
		ResponderID:    req.ResponderID,
		OfferedSkillID: req.OfferedSkillID,
		WantedSkillID:  req.WantedSkillID,
		Status:         models.StatusPending,
	}

	if err := tx.Create(swapRequest).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Load relationships
	if err := tx.Preload("Requester").Preload("Responder").
		Preload("OfferedSkill").Preload("WantedSkill").
		First(swapRequest, swapRequest.SwapID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return swapRequest, tx.Commit().Error
}

// GetSwapRequestByID retrieves a swap request by ID
func (s *swapService) GetSwapRequestByID(swapID uuid.UUID, userID uuid.UUID) (*models.SwapRequest, error) {
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

	// Check if user is involved in the swap
	if swapRequest.RequesterID != userID && swapRequest.ResponderID != userID {
		return nil, errors.New("access denied")
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
	swapRequest, err := s.GetSwapRequestByID(swapID, userID)
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
			swapRequest, err := s.GetSwapRequestByID(swapID, userID)
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

const (
	DefaultSwapRequestLimit = 50
)

// GetSwapRequestsForUser retrieves organized swap requests for a user
func (s *swapService) GetSwapRequestsForUser(userID uuid.UUID) (*SwapRequestsResponse, error) {
	// Get sent requests
	sentFilter := SwapRequestFilter{Sent: true, Limit: DefaultSwapRequestLimit}
	sent, err := s.GetUserSwapRequests(userID, sentFilter)
	if err != nil {
		return nil, err
	}

	// Get received requests
	receivedFilter := SwapRequestFilter{Received: true, Limit: DefaultSwapRequestLimit}
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

	// This query finds users who want a skill that the current user offers,
	// and who offer a skill that the current user wants.
	// It joins the user's offered and wanted skills with other users' wanted and offered skills.
	// It's a more efficient way to find matches than iterating through all skills in Go.
	err := s.db.Raw(`
		SELECT
			u.*,
			so.skill_id AS offered_skill_id,
			sw.skill_id AS wanted_skill_id,
			80 AS match_score
		FROM
			users u
		JOIN
			user_skills_offered uso ON u.user_id = uso.user_id
		JOIN
			user_skills_wanted usw ON u.user_id = usw.user_id
		JOIN
			skills so ON uso.skill_id = so.skill_id
		JOIN
			skills sw ON usw.skill_id = sw.skill_id
		WHERE
			uso.skill_id IN (SELECT skill_id FROM user_skills_wanted WHERE user_id = ?)
			AND usw.skill_id IN (SELECT skill_id FROM user_skills_offered WHERE user_id = ?)
			AND u.user_id != ?
			AND u.is_public = true
			AND u.deleted_at IS NULL
		LIMIT 20
	`, userID, userID, userID).Scan(&matches).Error

	return matches, err
}
