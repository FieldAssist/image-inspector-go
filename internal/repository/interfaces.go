package repository

import (
	"context"
	"image"

	"github.com/anime-shed/image-inspector-go/pkg/models"
)

// Repository defines the interface for data access operations related to image analysis.
// It abstracts the underlying database technology (e.g., GORM, pgx) from the service layer.
type Repository interface {
	// CreateAnalysis records a new image analysis report in the data store.
	CreateAnalysis(ctx context.Context, analysis *models.ImageAnalysis) error
}

// ImageRepository defines the interface for fetching images.
type ImageRepository interface {
	Fetch(ctx context.Context, url string) (image.Image, error)
}
