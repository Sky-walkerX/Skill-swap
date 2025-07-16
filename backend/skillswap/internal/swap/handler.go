package swap

import (
	"net/http"
	"strconv"

	appservice "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	swapService appservice.SwapService
}

func NewHandler(swapService appservice.SwapService) *Handler {
	return &Handler{
		swapService: swapService,
	}
}

// Request and response structures
type CreateSwapRequestRequest struct {
	ResponderID    string `json:"responder_id" binding:"required,uuid"`
	OfferedSkillID string `json:"offered_skill_id" binding:"required,uuid"`
	WantedSkillID  string `json:"wanted_skill_id" binding:"required,uuid"`
}

type UpdateSwapStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected cancelled"`
}

type SwapRequestResponse struct {
	SwapID         string         `json:"swap_id"`
	RequesterID    string         `json:"requester_id"`
	ResponderID    string         `json:"responder_id"`
	OfferedSkillID string         `json:"offered_skill_id"`
	WantedSkillID  string         `json:"wanted_skill_id"`
	Status         string         `json:"status"`
	CreatedAt      string         `json:"created_at"`
	UpdatedAt      string         `json:"updated_at"`
	Requester      *UserResponse  `json:"requester,omitempty"`
	Responder      *UserResponse  `json:"responder,omitempty"`
	OfferedSkill   *SkillResponse `json:"offered_skill,omitempty"`
	WantedSkill    *SkillResponse `json:"wanted_skill,omitempty"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email,omitempty"`
	Location string `json:"location,omitempty"`
	PhotoURL string `json:"photo_url,omitempty"`
}

type SkillResponse struct {
	SkillID string `json:"skill_id"`
	Name    string `json:"name"`
}

type SwapRequestsResponse struct {
	Sent     []SwapRequestResponse `json:"sent"`
	Received []SwapRequestResponse `json:"received"`
}

type MatchResponse struct {
	User         UserResponse  `json:"user"`
	OfferedSkill SkillResponse `json:"offered_skill"`
	WantedSkill  SkillResponse `json:"wanted_skill"`
	MatchScore   int           `json:"match_score"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Helper function to convert models to responses
func (h *Handler) convertToSwapResponse(swap *models.SwapRequest, includeDetails bool) SwapRequestResponse {
	response := SwapRequestResponse{
		SwapID:         swap.SwapID.String(),
		RequesterID:    swap.RequesterID.String(),
		ResponderID:    swap.ResponderID.String(),
		OfferedSkillID: swap.OfferedSkillID.String(),
		WantedSkillID:  swap.WantedSkillID.String(),
		Status:         string(swap.Status),
		CreatedAt:      swap.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      swap.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if includeDetails {
		if swap.Requester.UserID != uuid.Nil {
			response.Requester = &UserResponse{
				UserID:   swap.Requester.UserID.String(),
				Name:     swap.Requester.Name,
				Location: h.getStringValue(swap.Requester.Location),
				PhotoURL: h.getStringValue(swap.Requester.PhotoURL),
			}
		}

		if swap.Responder.UserID != uuid.Nil {
			response.Responder = &UserResponse{
				UserID:   swap.Responder.UserID.String(),
				Name:     swap.Responder.Name,
				Location: h.getStringValue(swap.Responder.Location),
				PhotoURL: h.getStringValue(swap.Responder.PhotoURL),
			}
		}

		if swap.OfferedSkill.SkillID != uuid.Nil {
			response.OfferedSkill = &SkillResponse{
				SkillID: swap.OfferedSkill.SkillID.String(),
				Name:    swap.OfferedSkill.Name,
			}
		}

		if swap.WantedSkill.SkillID != uuid.Nil {
			response.WantedSkill = &SkillResponse{
				SkillID: swap.WantedSkill.SkillID.String(),
				Name:    swap.WantedSkill.Name,
			}
		}
	}

	return response
}

func (h *Handler) getStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// CreateSwapRequest godoc
// @Summary Create swap request
// @Description Create a new skill swap request
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param swap body CreateSwapRequestRequest true "Swap request data"
// @Success 201 {object} SwapRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/swaps [post]
func (h *Handler) CreateSwapRequest(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	var req CreateSwapRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Parse UUIDs
	responderID, err := uuid.Parse(req.ResponderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid responder ID"})
		return
	}

	offeredSkillID, err := uuid.Parse(req.OfferedSkillID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid offered skill ID"})
		return
	}

	wantedSkillID, err := uuid.Parse(req.WantedSkillID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid wanted skill ID"})
		return
	}

	// Create swap request
	swapDTO := &appservice.CreateSwapRequestDTO{
		RequesterID:    userID,
		ResponderID:    responderID,
		OfferedSkillID: offeredSkillID,
		WantedSkillID:  wantedSkillID,
	}

	swap, err := h.swapService.CreateSwapRequest(swapDTO)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	response := h.convertToSwapResponse(swap, true)
	c.JSON(http.StatusCreated, response)
}

// GetSwapRequest godoc
// @Summary Get swap request
// @Description Get a specific swap request by ID
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Swap ID"
// @Success 200 {object} SwapRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/swaps/{id} [get]
func (h *Handler) GetSwapRequest(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	swapIDStr := c.Param("id")
	swapID, err := uuid.Parse(swapIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid swap ID"})
		return
	}

	swap, err := h.swapService.GetSwapRequestByID(swapID, userID)
	if err != nil {
		if err.Error() == "swap request not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Swap request not found"})
			return
		}
		if err.Error() == "access denied" {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch swap request"})
		return
	}

	response := h.convertToSwapResponse(swap, true)
	c.JSON(http.StatusOK, response)
}

// GetUserSwapRequests godoc
// @Summary Get user's swap requests
// @Description Get all swap requests for the authenticated user
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(pending, accepted, rejected, cancelled)
// @Param sent query bool false "Include sent requests"
// @Param received query bool false "Include received requests"
// @Param limit query int false "Limit number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} SwapRequestsResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/swaps [get]
func (h *Handler) GetUserSwapRequests(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	// Parse query parameters
	filter := appservice.SwapRequestFilter{}

	if statusStr := c.Query("status"); statusStr != "" {
		status := models.SwapStatus(statusStr)
		filter.Status = &status
	}

	if sentStr := c.Query("sent"); sentStr != "" {
		filter.Sent = sentStr == "true"
	}

	if receivedStr := c.Query("received"); receivedStr != "" {
		filter.Received = receivedStr == "true"
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			filter.Offset = offset
		}
	}

	// If neither sent nor received specified, show organized view
	if !filter.Sent && !filter.Received {
		swapRequests, err := h.swapService.GetSwapRequestsForUser(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch swap requests"})
			return
		}

		response := SwapRequestsResponse{
			Sent:     make([]SwapRequestResponse, len(swapRequests.Sent)),
			Received: make([]SwapRequestResponse, len(swapRequests.Received)),
		}

		for i, swap := range swapRequests.Sent {
			response.Sent[i] = h.convertToSwapResponse(&swap, true)
		}

		for i, swap := range swapRequests.Received {
			response.Received[i] = h.convertToSwapResponse(&swap, true)
		}

		c.JSON(http.StatusOK, response)
		return
	}

	// Get filtered requests
	swaps, err := h.swapService.GetUserSwapRequests(userID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to fetch swap requests"})
		return
	}

	var response []SwapRequestResponse
	for _, swap := range swaps {
		response = append(response, h.convertToSwapResponse(&swap, true))
	}

	c.JSON(http.StatusOK, response)
}

// UpdateSwapStatus godoc
// @Summary Update swap status
// @Description Accept, reject, or cancel a swap request
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Swap ID"
// @Param status body UpdateSwapStatusRequest true "New status"
// @Success 200 {object} SwapRequestResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/swaps/{id}/status [put]
func (h *Handler) UpdateSwapStatus(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	swapIDStr := c.Param("id")
	swapID, err := uuid.Parse(swapIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid swap ID"})
		return
	}

	var req UpdateSwapStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	status := models.SwapStatus(req.Status)
	swap, err := h.swapService.UpdateSwapStatus(swapID, userID, status)
	if err != nil {
		if err.Error() == "swap request not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Swap request not found"})
			return
		}
		if err.Error() == "only responder can accept or reject requests" ||
			err.Error() == "only requester or responder can cancel requests" {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	response := h.convertToSwapResponse(swap, true)
	c.JSON(http.StatusOK, response)
}

// DeleteSwapRequest godoc
// @Summary Delete swap request
// @Description Delete a swap request (requester only)
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Swap ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/swaps/{id} [delete]
func (h *Handler) DeleteSwapRequest(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	swapIDStr := c.Param("id")
	swapID, err := uuid.Parse(swapIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid swap ID"})
		return
	}

	err = h.swapService.DeleteSwapRequest(swapID, userID)
	if err != nil {
		if err.Error() == "swap request not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Swap request not found"})
			return
		}
		if err.Error() == "only requester can delete swap requests" {
			c.JSON(http.StatusForbidden, ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPotentialMatches godoc
// @Summary Get potential matches
// @Description Find potential swap matches for the authenticated user
// @Tags swaps
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} MatchResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/swaps/matches [get]
func (h *Handler) GetPotentialMatches(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID"})
		return
	}

	matches, err := h.swapService.FindPotentialMatches(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to find matches"})
		return
	}

	var response []MatchResponse
	for _, match := range matches {
		response = append(response, MatchResponse{
			User: UserResponse{
				UserID:   match.User.UserID.String(),
				Name:     match.User.Name,
				Location: h.getStringValue(match.User.Location),
				PhotoURL: h.getStringValue(match.User.PhotoURL),
			},
			OfferedSkill: SkillResponse{
				SkillID: match.OfferedSkill.SkillID.String(),
				Name:    match.OfferedSkill.Name,
			},
			WantedSkill: SkillResponse{
				SkillID: match.WantedSkill.SkillID.String(),
				Name:    match.WantedSkill.Name,
			},
			MatchScore: match.MatchScore,
		})
	}

	c.JSON(http.StatusOK, response)
}
