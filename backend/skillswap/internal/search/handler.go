package search

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	searchService service.SearchService
}

func NewHandler(searchService service.SearchService) *Handler {
	return &Handler{
		searchService: searchService,
	}
}

// SearchUsers performs advanced user search
// @Summary Advanced user search
// @Description Search users with advanced filtering options
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query (name, location)"
// @Param location query string false "Filter by location"
// @Param skills_offered query string false "Comma-separated skill IDs offered"
// @Param skills_wanted query string false "Comma-separated skill IDs wanted"
// @Param min_rating query number false "Minimum average rating"
// @Param is_public query bool false "Public profiles only"
// @Param sort_by query string false "Sort by (created_at, name, rating)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search/users [get]
func (h *Handler) SearchUsers(c *gin.Context) {
	filter := service.UserSearchFilter{
		Query:     c.Query("q"),
		Location:  c.Query("location"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	// Parse skills offered
	if skillsOfferedStr := c.Query("skills_offered"); skillsOfferedStr != "" {
		skillIDs := strings.Split(skillsOfferedStr, ",")
		for _, idStr := range skillIDs {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filter.SkillsOffered = append(filter.SkillsOffered, id)
			}
		}
	}

	// Parse skills wanted
	if skillsWantedStr := c.Query("skills_wanted"); skillsWantedStr != "" {
		skillIDs := strings.Split(skillsWantedStr, ",")
		for _, idStr := range skillIDs {
			if id, err := uuid.Parse(strings.TrimSpace(idStr)); err == nil {
				filter.SkillsWanted = append(filter.SkillsWanted, id)
			}
		}
	}

	// Parse min rating
	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			filter.MinRating = &minRating
		}
	}

	// Parse is_public
	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		if isPublic, err := strconv.ParseBool(isPublicStr); err == nil {
			filter.IsPublic = &isPublic
		}
	}

	// Parse pagination
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

	users, total, err := h.searchService.SearchUsers(filter)
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

// SearchSwaps performs advanced swap search
// @Summary Advanced swap search
// @Description Search swap requests with advanced filtering options
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query (description)"
// @Param status query string false "Filter by status"
// @Param offered_skill_id query string false "Filter by offered skill ID"
// @Param wanted_skill_id query string false "Filter by wanted skill ID"
// @Param requester_id query string false "Filter by requester ID"
// @Param responder_id query string false "Filter by responder ID"
// @Param created_after query string false "Created after date (ISO format)"
// @Param created_before query string false "Created before date (ISO format)"
// @Param sort_by query string false "Sort by (created_at, updated_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search/swaps [get]
func (h *Handler) SearchSwaps(c *gin.Context) {
	filter := service.SwapSearchFilter{
		Query:     c.Query("q"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	// Parse status
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}

	// Parse skill IDs
	if offeredSkillIDStr := c.Query("offered_skill_id"); offeredSkillIDStr != "" {
		if id, err := uuid.Parse(offeredSkillIDStr); err == nil {
			filter.OfferedSkillID = &id
		}
	}
	if wantedSkillIDStr := c.Query("wanted_skill_id"); wantedSkillIDStr != "" {
		if id, err := uuid.Parse(wantedSkillIDStr); err == nil {
			filter.WantedSkillID = &id
		}
	}

	// Parse user IDs
	if requesterIDStr := c.Query("requester_id"); requesterIDStr != "" {
		if id, err := uuid.Parse(requesterIDStr); err == nil {
			filter.RequesterID = &id
		}
	}
	if responderIDStr := c.Query("responder_id"); responderIDStr != "" {
		if id, err := uuid.Parse(responderIDStr); err == nil {
			filter.ResponderID = &id
		}
	}

	// Parse date filters
	if createdAfter := c.Query("created_after"); createdAfter != "" {
		filter.CreatedAfter = &createdAfter
	}
	if createdBefore := c.Query("created_before"); createdBefore != "" {
		filter.CreatedBefore = &createdBefore
	}

	// Parse pagination
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

	swaps, total, err := h.searchService.SearchSwaps(filter)
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

// SearchSkills performs advanced skill search
// @Summary Advanced skill search
// @Description Search skills with advanced filtering options
// @Tags search
// @Accept json
// @Produce json
// @Param q query string false "Search query (name, description)"
// @Param category query string false "Filter by category"
// @Param sort_by query string false "Sort by (name, created_at, popularity)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param limit query int false "Limit results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search/skills [get]
func (h *Handler) SearchSkills(c *gin.Context) {
	filter := service.SkillSearchFilter{
		Query:     c.Query("q"),
		Category:  c.Query("category"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
	}

	// Parse pagination
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

	skills, total, err := h.searchService.SearchSkills(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skills": skills,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// GlobalSearch performs search across all entities
// @Summary Global search
// @Description Search across users, skills, and swaps
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param types query string false "Comma-separated entity types (users,skills,swaps)"
// @Param limit query int false "Limit results per entity type"
// @Success 200 {object} service.GlobalSearchResults
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search/global [get]
func (h *Handler) GlobalSearch(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	// Parse entity types
	entityTypes := []string{"users", "skills", "swaps"} // Default: search all
	if typesStr := c.Query("types"); typesStr != "" {
		entityTypes = strings.Split(typesStr, ",")
		for i, t := range entityTypes {
			entityTypes[i] = strings.TrimSpace(t)
		}
	}

	// Parse limit
	limit := 5 // Default limit per entity type
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	results, err := h.searchService.GlobalSearch(query, entityTypes, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// SearchSuggestions provides search suggestions/autocomplete
// @Summary Search suggestions
// @Description Get search suggestions for autocomplete
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param type query string false "Entity type (users,skills,swaps)"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/v1/search/suggestions [get]
func (h *Handler) SearchSuggestions(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	entityType := c.Query("type")
	if entityType == "" {
		entityType = "skills" // Default to skills for suggestions
	}

	var suggestions []string

	switch entityType {
	case "skills":
		// Get skill name suggestions
		results, _, err := h.searchService.SearchSkills(service.SkillSearchFilter{
			Query: query,
			Limit: 10,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, skill := range results {
			suggestions = append(suggestions, skill.Name)
		}

	case "users":
		// Get user name suggestions (public users only)
		isPublic := true
		results, _, err := h.searchService.SearchUsers(service.UserSearchFilter{
			Query:    query,
			IsPublic: &isPublic,
			Limit:    10,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, user := range results {
			suggestions = append(suggestions, user.Name)
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid entity type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"suggestions": suggestions,
		"query":       query,
		"type":        entityType,
	})
}
