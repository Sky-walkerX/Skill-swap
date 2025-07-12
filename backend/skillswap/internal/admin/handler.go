package admin

import (
	"net/http"
	"strconv"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	adminService service.AdminService
}

func NewHandler(adminService service.AdminService) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

// GetAllUsers retrieves all users with filtering
// @Summary Get all users (admin only)
// @Description Get all users with filtering and pagination
// @Tags admin
// @Accept json
// @Produce json
// @Param search query string false "Search by name or email"
// @Param is_banned query bool false "Filter by banned status"
// @Param is_admin query bool false "Filter by admin status"
// @Param sort_by query string false "Sort by field (created_at, name, email)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	// Parse query parameters
	filter := service.AdminUserFilter{
		Search:    c.Query("search"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	if isBanned := c.Query("is_banned"); isBanned != "" {
		if val, err := strconv.ParseBool(isBanned); err == nil {
			filter.IsBanned = &val
		}
	}

	if isAdmin := c.Query("is_admin"); isAdmin != "" {
		if val, err := strconv.ParseBool(isAdmin); err == nil {
			filter.IsAdmin = &val
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			filter.Offset = val
		}
	}

	users, total, err := h.adminService.GetAllUsers(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":  users,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// BanUser bans a user
// @Summary Ban a user (admin only)
// @Description Ban a user from the platform
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users/{id}/ban [put]
func (h *Handler) BanUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get admin ID from JWT token
	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.BanUser(adminID.(uuid.UUID), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err.Error() == "cannot ban an admin user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot ban an admin user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// UnbanUser unbans a user
// @Summary Unban a user (admin only)
// @Description Unban a user from the platform
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users/{id}/unban [put]
func (h *Handler) UnbanUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.UnbanUser(adminID.(uuid.UUID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteUser deletes a user
// @Summary Delete a user (admin only)
// @Description Soft delete a user from the platform
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 404 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.DeleteUser(adminID.(uuid.UUID), userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if err.Error() == "cannot delete an admin user" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete an admin user"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// MakeUserAdmin grants admin privileges to a user
// @Summary Make user admin (admin only)
// @Description Grant admin privileges to a user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users/{id}/make-admin [put]
func (h *Handler) MakeUserAdmin(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.MakeUserAdmin(adminID.(uuid.UUID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// RemoveUserAdmin removes admin privileges from a user
// @Summary Remove admin privileges (admin only)
// @Description Remove admin privileges from a user
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/users/{id}/remove-admin [put]
func (h *Handler) RemoveUserAdmin(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.RemoveUserAdmin(adminID.(uuid.UUID), userID)
	if err != nil {
		if err.Error() == "cannot remove admin privileges from yourself" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot remove admin privileges from yourself"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetAllSwaps retrieves all swaps with filtering
// @Summary Get all swaps (admin only)
// @Description Get all swaps with filtering and pagination
// @Tags admin
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param requester_id query string false "Filter by requester ID"
// @Param responder_id query string false "Filter by responder ID"
// @Param sort_by query string false "Sort by field (created_at, updated_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/swaps [get]
func (h *Handler) GetAllSwaps(c *gin.Context) {
	filter := service.AdminSwapFilter{
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	if requesterID := c.Query("requester_id"); requesterID != "" {
		if id, err := uuid.Parse(requesterID); err == nil {
			filter.RequesterID = &id
		}
	}

	if responderID := c.Query("responder_id"); responderID != "" {
		if id, err := uuid.Parse(responderID); err == nil {
			filter.ResponderID = &id
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			filter.Limit = val
		}
	}

	if offset := c.Query("offset"); offset != "" {
		if val, err := strconv.Atoi(offset); err == nil {
			filter.Offset = val
		}
	}

	swaps, total, err := h.adminService.GetAllSwaps(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"swaps":  swaps,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// CancelSwap cancels a swap with admin intervention
// @Summary Cancel a swap (admin only)
// @Description Cancel a swap as admin intervention
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "Swap ID"
// @Param reason body gin.H true "Reason for cancellation"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/swaps/{id}/cancel [put]
func (h *Handler) CancelSwap(c *gin.Context) {
	swapIDStr := c.Param("id")
	swapID, err := uuid.Parse(swapIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid swap ID"})
		return
	}

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err = h.adminService.CancelSwap(adminID.(uuid.UUID), swapID, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetPlatformStats retrieves platform statistics
// @Summary Get platform statistics (admin only)
// @Description Get platform-wide statistics and metrics
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} service.PlatformStats
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/stats [get]
func (h *Handler) GetPlatformStats(c *gin.Context) {
	stats, err := h.adminService.GetPlatformStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetReportedContent retrieves reported content
// @Summary Get reported content (admin only)
// @Description Get list of reported content for moderation
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {array} service.ReportedContent
// @Failure 401 {object} gin.H
// @Failure 403 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/admin/reports [get]
func (h *Handler) GetReportedContent(c *gin.Context) {
	reports, err := h.adminService.GetReportedContent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}
