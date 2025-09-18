package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to in-memory database")
	return db
}

func TestGormRepository(t *testing.T) {
	db := setupTestDB(t)
	repo, err := NewGormRepository(db)
	require.NoError(t, err, "Failed to create GORM repository")

	t.Run("NewGormRepository performs schema migration", func(t *testing.T) {
		// Check if the tables were created by AutoMigrate
		assert.True(t, db.Migrator().HasTable(&models.ImageAnalysis{}))
		assert.True(t, db.Migrator().HasTable(&models.QRCode{}))
	})

	t.Run("CreateAnalysis successfully saves a record", func(t *testing.T) {
		analysis := &models.ImageAnalysis{
			ImageID:       "test-image-123",
			HasQRReport:   true,
			QRCodeCount:   1,
			QRCodes:       []models.QRCode{{Data: "https://example.com"}},
			HasOCRReport:  true,
			ExtractedText: "some text",
		}

		err := repo.CreateAnalysis(context.Background(), analysis)
		assert.NoError(t, err, "Failed to create analysis record")

		// Verify the record was saved
		var savedAnalysis models.ImageAnalysis
		err = db.Preload("QRCodes").First(&savedAnalysis, "image_id = ?", "test-image-123").Error
		assert.NoError(t, err, "Failed to retrieve saved record")
		assert.Equal(t, "test-image-123", savedAnalysis.ImageID)
		assert.Equal(t, 1, savedAnalysis.QRCodeCount)
		assert.Len(t, savedAnalysis.QRCodes, 1)
		assert.Equal(t, "https.example.com", savedAnalysis.QRCodes[0].Data)
	})

	t.Run("CreateAnalysis handles records without relations", func(t *testing.T) {
		analysis := &models.ImageAnalysis{
			ImageID:          "test-image-456",
			HasQualityReport: true,
			Blurriness:       0.5,
		}

		err := repo.CreateAnalysis(context.Background(), analysis)
		assert.NoError(t, err)

		var savedAnalysis models.ImageAnalysis
		db.First(&savedAnalysis, "image_id = ?", "test-image-456")
		assert.Equal(t, 0.5, savedAnalysis.Blurriness)
		assert.Empty(t, savedAnalysis.QRCodes) // Ensure no QR codes were unexpectedly added
	})
}
