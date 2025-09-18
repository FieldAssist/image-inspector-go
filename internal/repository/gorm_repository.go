package repository

import (
	"context"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"gorm.io/gorm"
)

// GormRepository is a GORM-based implementation of the Repository interface.
type GormRepository struct {
	db *gorm.DB
}

// NewGormRepository creates a new GormRepository.
func NewGormRepository(db *gorm.DB) (*GormRepository, error) {
	// Auto-migrate the schema to keep it up to date.
	if err := db.AutoMigrate(&models.ImageAnalysis{}, &models.QRCode{}); err != nil {
		return nil, err
	}
	return &GormRepository{db: db}, nil
}

// CreateAnalysis saves a new ImageAnalysis record to the database.
func (r *GormRepository) CreateAnalysis(ctx context.Context, analysis *models.ImageAnalysis) error {
	// Use GORM's context-aware methods.
	return r.db.WithContext(ctx).Create(analysis).Error
}
