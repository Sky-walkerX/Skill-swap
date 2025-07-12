package models

import (
	"time"
)

// FileUploadResponse represents a successful file upload response
type FileUploadResponse struct {
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
}

// FileInfo represents file information
type FileInfo struct {
	Filename string    `json:"filename"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mod_time"`
	URL      string    `json:"url"`
}

// FileUploadRequest represents a file upload request
type FileUploadRequest struct {
	File string `form:"file" binding:"required"` // This will be handled by multipart form
}

// FileDeleteResponse represents a file deletion response
type FileDeleteResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
