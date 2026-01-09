package storage

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// ImageStore handles image upload and storage
type ImageStore struct {
	baseDir string
}

// NewImageStore creates a new image store
func NewImageStore(baseDir string) *ImageStore {
	return &ImageStore{
		baseDir: baseDir,
	}
}

// SaveImage saves an uploaded image and returns its ID and path
func (is *ImageStore) SaveImage(filename string, data io.Reader) (string, string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(is.baseDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Generate unique ID
	imageID := uuid.New().String()

	// Determine file extension
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".png" // Default extension
	}

	// Create unique filename
	storedFilename := fmt.Sprintf("%s%s", imageID, ext)
	imagePath := filepath.Join(is.baseDir, storedFilename)

	// Create file
	file, err := os.Create(imagePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy data to file
	if _, err := io.Copy(file, data); err != nil {
		os.Remove(imagePath)
		return "", "", fmt.Errorf("failed to write image data: %w", err)
	}

	log.Printf("ImageStore: Saved image %s to %s", imageID, imagePath)
	return imageID, imagePath, nil
}

// GetImagePath returns the filesystem path for an image ID
func (is *ImageStore) GetImagePath(imageID string) (string, error) {
	// Find file with matching ID (any extension)
	pattern := filepath.Join(is.baseDir, imageID+".*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to search for image: %w", err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("image not found: %s", imageID)
	}

	return matches[0], nil
}

// DeleteImage deletes an image by ID
func (is *ImageStore) DeleteImage(imageID string) error {
	imagePath, err := is.GetImagePath(imageID)
	if err != nil {
		return err
	}

	if err := os.Remove(imagePath); err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	log.Printf("ImageStore: Deleted image %s", imageID)
	return nil
}

// CleanupOldImages removes images older than the specified duration
func (is *ImageStore) CleanupOldImages(maxAge time.Duration) error {
	entries, err := os.ReadDir(is.baseDir)
	if err != nil {
		return fmt.Errorf("failed to read storage directory: %w", err)
	}

	now := time.Now()
	deletedCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.Printf("Warning: Failed to get file info for %s: %v", entry.Name(), err)
			continue
		}

		age := now.Sub(info.ModTime())
		if age > maxAge {
			imagePath := filepath.Join(is.baseDir, entry.Name())
			if err := os.Remove(imagePath); err != nil {
				log.Printf("Warning: Failed to delete old image %s: %v", imagePath, err)
			} else {
				deletedCount++
			}
		}
	}

	if deletedCount > 0 {
		log.Printf("ImageStore: Cleaned up %d old images", deletedCount)
	}

	return nil
}

// ImageMetadata holds information about an uploaded image
type ImageMetadata struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	Path       string    `json:"path"`
	UploadedAt time.Time `json:"uploaded_at"`
	Size       int64     `json:"size"`
}
