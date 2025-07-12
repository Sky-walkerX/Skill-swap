package availability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	availabilityService service.AvailabilityService
}

func NewHandler(availabilityService service.AvailabilityService) *Handler {
	return &Handler{
		availabilityService: availabilityService,
	}
}

// CreateAvailabilitySlot creates a new availability slot
// @Summary Create a new availability slot
// @Description Create a new availability slot for the authenticated user
// @Tags availability
// @Accept json
// @Produce json
// @Param availability body service.CreateAvailabilitySlotDTO true "Availability slot data"
// @Success 201 {object} models.AvailabilitySlot
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability [post]
func (h *Handler) CreateAvailabilitySlot(c *gin.Context) {
	var req service.CreateAvailabilitySlotDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	req.UserID = userID.(uuid.UUID)

	slot, err := h.availabilityService.CreateAvailabilitySlot(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, slot)
}

// GetUserAvailabilitySlots retrieves all availability slots for the authenticated user
// @Summary Get user's availability slots
// @Description Get all availability slots for the authenticated user
// @Tags availability
// @Accept json
// @Produce json
// @Success 200 {array} models.AvailabilitySlot
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability [get]
func (h *Handler) GetUserAvailabilitySlots(c *gin.Context) {
	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	slots, err := h.availabilityService.GetUserAvailabilitySlots(userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availability_slots": slots})
}

// GetAvailabilitySlot retrieves a specific availability slot
// @Summary Get an availability slot
// @Description Get a specific availability slot by ID
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability slot ID"
// @Success 200 {object} models.AvailabilitySlot
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability/{id} [get]
func (h *Handler) GetAvailabilitySlot(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot ID"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	slot, err := h.availabilityService.GetAvailabilitySlot(slotID, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "availability slot not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Availability slot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slot)
}

// UpdateAvailabilitySlot updates an existing availability slot
// @Summary Update an availability slot
// @Description Update an existing availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability slot ID"
// @Param availability body service.UpdateAvailabilitySlotDTO true "Updated availability slot data"
// @Success 200 {object} models.AvailabilitySlot
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability/{id} [put]
func (h *Handler) UpdateAvailabilitySlot(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot ID"})
		return
	}

	var req service.UpdateAvailabilitySlotDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	slot, err := h.availabilityService.UpdateAvailabilitySlot(slotID, userID.(uuid.UUID), &req)
	if err != nil {
		if err.Error() == "availability slot not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Availability slot not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slot)
}

// DeleteAvailabilitySlot deletes an availability slot
// @Summary Delete an availability slot
// @Description Delete an availability slot
// @Tags availability
// @Accept json
// @Produce json
// @Param id path string true "Availability slot ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability/{id} [delete]
func (h *Handler) DeleteAvailabilitySlot(c *gin.Context) {
	slotIDStr := c.Param("id")
	slotID, err := uuid.Parse(slotIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slot ID"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.availabilityService.DeleteAvailabilitySlot(slotID, userID.(uuid.UUID))
	if err != nil {
		if err.Error() == "availability slot not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Availability slot not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// FindCommonAvailability finds overlapping availability between two users
// @Summary Find common availability
// @Description Find overlapping availability between authenticated user and another user
// @Tags availability
// @Accept json
// @Produce json
// @Param user_id path string true "Other user ID"
// @Success 200 {array} service.CommonAvailabilitySlot
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability/common/{user_id} [get]
func (h *Handler) FindCommonAvailability(c *gin.Context) {
	otherUserIDStr := c.Param("user_id")
	otherUserID, err := uuid.Parse(otherUserIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	commonSlots, err := h.availabilityService.FindCommonAvailability(userID.(uuid.UUID), otherUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"common_availability": commonSlots})
}

// GetAvailabilityByDayAndTime finds availability slots for specific day and time
// @Summary Get availability by day and time
// @Description Get availability slots for a specific day and time range
// @Tags availability
// @Accept json
// @Produce json
// @Param day query int true "Day of week (1=Monday, 7=Sunday)"
// @Param start_time query string true "Start time (HH:MM)"
// @Param end_time query string true "End time (HH:MM)"
// @Success 200 {array} models.AvailabilitySlot
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/availability/search [get]
func (h *Handler) GetAvailabilityByDayAndTime(c *gin.Context) {
	// Parse query parameters
	dayStr := c.Query("day")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if dayStr == "" || startTimeStr == "" || endTimeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "day, start_time, and end_time are required"})
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "day must be between 1 and 7"})
		return
	}

	startTime, err := time.Parse("15:04", startTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time format, use HH:MM"})
		return
	}

	endTime, err := time.Parse("15:04", endTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time format, use HH:MM"})
		return
	}

	// Get user ID from JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	slots, err := h.availabilityService.GetAvailabilityByDayAndTime(userID.(uuid.UUID), day, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"availability_slots": slots})
}
