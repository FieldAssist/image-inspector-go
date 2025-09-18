package analyzer

import (
	"context"
	"image"
	"testing"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestQualityAnalyzer(t *testing.T) {
	t.Run("Analyze populates quality report", func(t *testing.T) {
		qa := NewQualityAnalyzer()
		img := image.NewRGBA(image.Rect(0, 0, 100, 100)) // 100x100 = 10000 pixels
		report := &models.ImageAnalysis{}

		err := qa.Analyze(context.Background(), img, report)

		assert.NoError(t, err)
		assert.True(t, report.HasQualityReport, "Expected HasQualityReport to be true")

		// Based on our placeholder logic:
		// Blurriness = 1.0 - (100*100)/10000.0 = 0.0
		assert.InDelta(t, 0.0, report.Blurriness, 0.01, "Unexpected blurriness score")

		// Overexposure depends on the color of the generated image (default is all black/zeroes)
		assert.InDelta(t, 0.0, report.Overexposure, 0.01, "Unexpected overexposure score")
	})

	t.Run("Calculate blurriness placeholder logic", func(t *testing.T) {
		// Small image -> high blurriness score
		smallImg := image.NewRGBA(image.Rect(0, 0, 50, 50))
		blurrinessSmall := calculateBlurriness(smallImg) // 1.0 - (2500/10000) = 0.75
		assert.InDelta(t, 0.75, blurrinessSmall, 0.01)

		// Large image -> low blurriness score
		largeImg := image.NewRGBA(image.Rect(0, 0, 200, 200))
		blurrinessLarge := calculateBlurriness(largeImg) // 1.0 - (40000/10000) capped at 1.0 -> 0.0
		assert.InDelta(t, 0.0, blurrinessLarge, 0.01)
	})
}
