package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	models "github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileUploadService struct {
	db           *gorm.DB
	uploadDir    string
	maxFileSize  int64
	allowedTypes map[string]bool
	baseURL      string
}

func NewFileUploadService(db *gorm.DB, cfg config.Config) *FileUploadService {
	// Set default upload directory
	uploadDir := cfg.UploadDir
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	// Set allowed image types
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	return &FileUploadService{
		db:           db,
		uploadDir:    uploadDir,
		maxFileSize:  5 * 1024 * 1024, // 5MB default
		allowedTypes: allowedTypes,
		baseURL:      cfg.BaseURL,
	}
}

// UploadUserPhoto uploads a user's profile photo
func (s *FileUploadService) UploadUserPhoto(userID uuid.UUID, file *multipart.FileHeader) (*models.FileUploadResponse, error) {
	// Validate file
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	// Generate unique filename
	filename, err := s.generateFilename(file.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	// Create user-specific directory
	userDir := filepath.Join(s.uploadDir, "users", userID.String())
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create user directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(userDir, filename)

	// Save file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Generate URL
	fileURL := s.generateFileURL(userID, filename)

	// Update user's photo URL in database
	if err := s.updateUserPhotoURL(userID, fileURL); err != nil {
		// Cleanup uploaded file on database error
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to update user photo URL: %w", err)
	}

	return &models.FileUploadResponse{
		Filename: filename,
		URL:      fileURL,
		Size:     file.Size,
		MimeType: file.Header.Get("Content-Type"),
	}, nil
}

// DeleteUserPhoto deletes a user's profile photo
func (s *FileUploadService) DeleteUserPhoto(userID uuid.UUID) error {
	// Get current photo URL from database
	var user models.User
	if err := s.db.Select("photo_url").First(&user, "user_id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// If no photo URL, nothing to delete
	if user.PhotoURL == nil || *user.PhotoURL == "" {
		return fmt.Errorf("user has no photo to delete")
	}

	// Extract filename from URL
	filename := s.extractFilenameFromURL(*user.PhotoURL)
	if filename == "" {
		return fmt.Errorf("invalid photo URL")
	}

	// Delete file from filesystem
	filePath := filepath.Join(s.uploadDir, "users", userID.String(), filename)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Clear photo URL in database
	if err := s.db.Model(&user).Where("user_id = ?", userID).Update("photo_url", nil).Error; err != nil {
		return fmt.Errorf("failed to clear photo URL: %w", err)
	}

	return nil
}

// GetFileInfo returns information about a file
func (s *FileUploadService) GetFileInfo(userID uuid.UUID, filename string) (*models.FileInfo, error) {
	filePath := filepath.Join(s.uploadDir, "users", userID.String(), filename)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found")
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	return &models.FileInfo{
		Filename: filename,
		Size:     info.Size(),
		ModTime:  info.ModTime(),
		URL:      s.generateFileURL(userID, filename),
	}, nil
}

// validateFile validates the uploaded file
func (s *FileUploadService) validateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > s.maxFileSize {
		return fmt.Errorf("file too large: %d bytes (max: %d bytes)", file.Size, s.maxFileSize)
	}

	// Check file type
	contentType := file.Header.Get("Content-Type")
	if !s.allowedTypes[contentType] {
		return fmt.Errorf("invalid file type: %s", contentType)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	validExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExts[ext] {
		return fmt.Errorf("invalid file extension: %s", ext)
	}

	return nil
}

// generateFilename generates a unique filename
func (s *FileUploadService) generateFilename(originalName string) (string, error) {
	ext := filepath.Ext(originalName)

	// Generate random bytes for unique filename
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", hex.EncodeToString(bytes), ext), nil
}

// generateFileURL generates the full URL for a file
func (s *FileUploadService) generateFileURL(userID uuid.UUID, filename string) string {
	return fmt.Sprintf("%s/api/v1/files/users/%s/%s", s.baseURL, userID.String(), filename)
}

// extractFilenameFromURL extracts filename from file URL
func (s *FileUploadService) extractFilenameFromURL(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// updateUserPhotoURL updates the user's photo URL in the database
func (s *FileUploadService) updateUserPhotoURL(userID uuid.UUID, photoURL string) error {
	return s.db.Model(&models.User{}).
		Where("user_id = ?", userID).
		Update("photo_url", photoURL).Error
}

// CleanupOrphanedFiles removes files that are not referenced in the database
func (s *FileUploadService) CleanupOrphanedFiles() error {
	userDir := filepath.Join(s.uploadDir, "users")

	return filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Extract user ID and filename from path
		rel, err := filepath.Rel(userDir, path)
		if err != nil {
			return err
		}

		parts := strings.Split(rel, string(os.PathSeparator))
		if len(parts) != 2 {
			return nil
		}

		userIDStr := parts[0]
		filename := parts[1]

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return nil
		}

		// Check if file is referenced in database
		var count int64
		expectedURL := s.generateFileURL(userID, filename)

		err = s.db.Model(&models.User{}).
			Where("user_id = ? AND photo_url = ?", userID, expectedURL).
			Count(&count).Error

		if err != nil {
			return err
		}

		// If file is not referenced, delete it
		if count == 0 {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("failed to remove orphaned file %s: %w", path, err)
			}
		}

		return nil
	})
}
