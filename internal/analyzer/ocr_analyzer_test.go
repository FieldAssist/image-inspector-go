package analyzer

import (
	"context"
	"image"
	"testing"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"github.com/stretchr/testify/assert"
)

func TestOCRAnalyzer(t *testing.T) {
	// This test uses the real OCRAnalyzer, but since the Analyze method is a simulation,
	// it doesn't actually call the Tesseract binary.
	t.Run("Analyze populates OCR report with simulated text", func(t *testing.T) {
		ocrAnalyzer := NewOCRAnalyzer()
		defer ocrAnalyzer.Close() // Ensure the client is closed after the test

		img := image.NewRGBA(image.Rect(0, 0, 100, 100))
		report := &models.ImageAnalysis{}

		err := ocrAnalyzer.Analyze(context.Background(), img, report)

		assert.NoError(t, err)
		assert.True(t, report.HasOCRReport)
		assert.Equal(t, "Simulated OCR text from image", report.ExtractedText)
	})
}
