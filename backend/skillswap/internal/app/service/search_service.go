package service

import (
	"fmt"
	"strings"

	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SearchService interface {
	// Advanced user search
	SearchUsers(filter UserSearchFilter) ([]models.User, int64, error)

	// Advanced swap search
	SearchSwaps(filter SwapSearchFilter) ([]models.SwapRequest, int64, error)

	// Advanced skill search
	SearchSkills(filter SkillSearchFilter) ([]models.Skill, int64, error)

	// Global search across all entities
	GlobalSearch(query string, entityTypes []string, limit int) (*GlobalSearchResults, error)
}

// Search filters and DTOs
type UserSearchFilter struct {
	Query         string      `json:"query,omitempty"`          // Search in name, location
	Location      string      `json:"location,omitempty"`       // Filter by location
	SkillsOffered []uuid.UUID `json:"skills_offered,omitempty"` // Users offering these skills
	SkillsWanted  []uuid.UUID `json:"skills_wanted,omitempty"`  // Users wanting these skills
	MinRating     *float64    `json:"min_rating,omitempty"`     // Minimum average rating
	IsPublic      *bool       `json:"is_public,omitempty"`      // Public profiles only
	SortBy        string      `json:"sort_by,omitempty"`        // "created_at", "name", "rating"
	SortOrder     string      `json:"sort_order,omitempty"`     // "asc", "desc"
	Limit         int         `json:"limit,omitempty"`
	Offset        int         `json:"offset,omitempty"`
}

type SwapSearchFilter struct {
	Query          string     `json:"query,omitempty"`            // Search in description
	Status         *string    `json:"status,omitempty"`           // Filter by status
	OfferedSkillID *uuid.UUID `json:"offered_skill_id,omitempty"` // Filter by offered skill
	WantedSkillID  *uuid.UUID `json:"wanted_skill_id,omitempty"`  // Filter by wanted skill
	RequesterID    *uuid.UUID `json:"requester_id,omitempty"`     // Filter by requester
	ResponderID    *uuid.UUID `json:"responder_id,omitempty"`     // Filter by responder
	LocationRadius *float64   `json:"location_radius,omitempty"`  // Search within radius (future)
	CreatedAfter   *string    `json:"created_after,omitempty"`    // Created after date
	CreatedBefore  *string    `json:"created_before,omitempty"`   // Created before date
	SortBy         string     `json:"sort_by,omitempty"`          // "created_at", "updated_at"
	SortOrder      string     `json:"sort_order,omitempty"`       // "asc", "desc"
	Limit          int        `json:"limit,omitempty"`
	Offset         int        `json:"offset,omitempty"`
}

type SkillSearchFilter struct {
	Query     string `json:"query,omitempty"`      // Search in name, description
	Category  string `json:"category,omitempty"`   // Filter by category
	SortBy    string `json:"sort_by,omitempty"`    // "name", "created_at", "popularity"
	SortOrder string `json:"sort_order,omitempty"` // "asc", "desc"
	Limit     int    `json:"limit,omitempty"`
	Offset    int    `json:"offset,omitempty"`
}

type GlobalSearchResults struct {
	Users        []models.User        `json:"users,omitempty"`
	Skills       []models.Skill       `json:"skills,omitempty"`
	SwapRequests []models.SwapRequest `json:"swap_requests,omitempty"`
	Total        int                  `json:"total"`
}

type searchService struct {
	db *gorm.DB
}

func NewSearchService(db *gorm.DB) SearchService {
	return &searchService{db: db}
}

// SearchUsers performs advanced user search with filtering
func (s *searchService) SearchUsers(filter UserSearchFilter) ([]models.User, int64, error) {
	query := s.db.Model(&models.User{}).Where("deleted_at IS NULL")

	// Apply filters
	if filter.Query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query))
		query = query.Where("LOWER(name) LIKE ? OR LOWER(location) LIKE ?", searchTerm, searchTerm)
	}

	if filter.Location != "" {
		locationTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Location))
		query = query.Where("LOWER(location) LIKE ?", locationTerm)
	}

	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}

	// Filter by skills offered
	if len(filter.SkillsOffered) > 0 {
		query = query.Joins("JOIN user_skills_offered uso ON users.user_id = uso.user_id").
			Where("uso.skill_id IN ?", filter.SkillsOffered).
			Distinct("users.user_id")
	}

	// Filter by skills wanted
	if len(filter.SkillsWanted) > 0 {
		query = query.Joins("JOIN user_skills_wanted usw ON users.user_id = usw.user_id").
			Where("usw.skill_id IN ?", filter.SkillsWanted).
			Distinct("users.user_id")
	}

	// Filter by minimum rating (requires calculating average rating)
	if filter.MinRating != nil {
		query = query.Joins("LEFT JOIN swap_ratings sr ON users.user_id = sr.ratee_id").
			Group("users.user_id").
			Having("COALESCE(AVG(sr.score), 0) >= ?", *filter.MinRating)
	}

	// Get total count
	var total int64
	countQuery := s.db.Model(&models.User{}).Where("deleted_at IS NULL")
	if err := applyUserFiltersForCount(countQuery, filter).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name", "created_at":
			sortBy = filter.SortBy
		case "rating":
			sortBy = "AVG(sr.score)"
		}
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20) // Default limit
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var users []models.User
	err := query.Preload("SkillsOffered.Skill").Preload("SkillsWanted.Skill").Find(&users).Error
	return users, total, err
}

// SearchSwaps performs advanced swap search with filtering
func (s *searchService) SearchSwaps(filter SwapSearchFilter) ([]models.SwapRequest, int64, error) {
	query := s.db.Model(&models.SwapRequest{}).
		Preload("Requester").
		Preload("Responder").
		Preload("OfferedSkill").
		Preload("WantedSkill")

	// Apply filters
	if filter.Query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query))
		query = query.Where("LOWER(description) LIKE ?", searchTerm)
	}

	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}

	if filter.OfferedSkillID != nil {
		query = query.Where("offered_skill_id = ?", *filter.OfferedSkillID)
	}

	if filter.WantedSkillID != nil {
		query = query.Where("wanted_skill_id = ?", *filter.WantedSkillID)
	}

	if filter.RequesterID != nil {
		query = query.Where("requester_id = ?", *filter.RequesterID)
	}

	if filter.ResponderID != nil {
		query = query.Where("responder_id = ?", *filter.ResponderID)
	}

	if filter.CreatedAfter != nil {
		query = query.Where("created_at >= ?", *filter.CreatedAfter)
	}

	if filter.CreatedBefore != nil {
		query = query.Where("created_at <= ?", *filter.CreatedBefore)
	}

	// Get total count - rebuild query for count
	var total int64
	countQuery := s.db.Model(&models.SwapRequest{})
	if filter.Query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query))
		countQuery = countQuery.Where("LOWER(description) LIKE ?", searchTerm)
	}
	if filter.Status != nil {
		countQuery = countQuery.Where("status = ?", *filter.Status)
	}
	if filter.OfferedSkillID != nil {
		countQuery = countQuery.Where("offered_skill_id = ?", *filter.OfferedSkillID)
	}
	if filter.WantedSkillID != nil {
		countQuery = countQuery.Where("wanted_skill_id = ?", *filter.WantedSkillID)
	}
	if filter.RequesterID != nil {
		countQuery = countQuery.Where("requester_id = ?", *filter.RequesterID)
	}
	if filter.ResponderID != nil {
		countQuery = countQuery.Where("responder_id = ?", *filter.ResponderID)
	}
	if filter.CreatedAfter != nil {
		countQuery = countQuery.Where("created_at >= ?", *filter.CreatedAfter)
	}
	if filter.CreatedBefore != nil {
		countQuery = countQuery.Where("created_at <= ?", *filter.CreatedBefore)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "created_at"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "created_at", "updated_at":
			sortBy = filter.SortBy
		}
	}
	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20) // Default limit
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var swaps []models.SwapRequest
	err := query.Find(&swaps).Error
	return swaps, total, err
}

// SearchSkills performs advanced skill search with filtering
func (s *searchService) SearchSkills(filter SkillSearchFilter) ([]models.Skill, int64, error) {
	query := s.db.Model(&models.Skill{})

	// Apply filters
	if filter.Query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query))
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	if filter.Category != "" {
		query = query.Where("LOWER(category) = LOWER(?)", filter.Category)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortBy := "name"
	if filter.SortBy != "" {
		switch filter.SortBy {
		case "name", "created_at":
			sortBy = filter.SortBy
		case "popularity":
			// Sort by number of times this skill is offered/wanted
			query = query.Select("skills.*, (SELECT COUNT(*) FROM user_skills_offered WHERE skill_id = skills.skill_id) + (SELECT COUNT(*) FROM user_skills_wanted WHERE skill_id = skills.skill_id) as popularity").
				Order("popularity")
		}
	}

	if filter.SortBy != "popularity" {
		sortOrder := "ASC"
		if filter.SortOrder == "desc" {
			sortOrder = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		sortOrder := "DESC"
		if filter.SortOrder == "asc" {
			sortOrder = "ASC"
		}
		query = query.Order(fmt.Sprintf("popularity %s", sortOrder))
	}

	// Apply pagination
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	} else {
		query = query.Limit(20) // Default limit
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	var skills []models.Skill
	err := query.Find(&skills).Error
	return skills, total, err
}

// GlobalSearch performs search across all entities
func (s *searchService) GlobalSearch(query string, entityTypes []string, limit int) (*GlobalSearchResults, error) {
	results := &GlobalSearchResults{}
	searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(query))

	if limit <= 0 {
		limit = 10 // Default limit per entity type
	}

	for _, entityType := range entityTypes {
		switch entityType {
		case "users":
			var users []models.User
			err := s.db.Model(&models.User{}).
				Where("deleted_at IS NULL AND is_public = true AND (LOWER(name) LIKE ? OR LOWER(location) LIKE ?)", searchTerm, searchTerm).
				Limit(limit).
				Preload("SkillsOffered.Skill").
				Preload("SkillsWanted.Skill").
				Find(&users).Error
			if err != nil {
				return nil, err
			}
			results.Users = users
			results.Total += len(users)

		case "skills":
			var skills []models.Skill
			err := s.db.Model(&models.Skill{}).
				Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm).
				Limit(limit).
				Find(&skills).Error
			if err != nil {
				return nil, err
			}
			results.Skills = skills
			results.Total += len(skills)

		case "swaps":
			var swaps []models.SwapRequest
			err := s.db.Model(&models.SwapRequest{}).
				Where("LOWER(description) LIKE ?", searchTerm).
				Limit(limit).
				Preload("Requester").
				Preload("Responder").
				Preload("OfferedSkill").
				Preload("WantedSkill").
				Find(&swaps).Error
			if err != nil {
				return nil, err
			}
			results.SwapRequests = swaps
			results.Total += len(swaps)
		}
	}

	return results, nil
}

// Helper function to apply user filters for count query (simplified)
func applyUserFiltersForCount(query *gorm.DB, filter UserSearchFilter) *gorm.DB {
	if filter.Query != "" {
		searchTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Query))
		query = query.Where("LOWER(name) LIKE ? OR LOWER(location) LIKE ?", searchTerm, searchTerm)
	}

	if filter.Location != "" {
		locationTerm := fmt.Sprintf("%%%s%%", strings.ToLower(filter.Location))
		query = query.Where("LOWER(location) LIKE ?", locationTerm)
	}

	if filter.IsPublic != nil {
		query = query.Where("is_public = ?", *filter.IsPublic)
	}

	return query
}
