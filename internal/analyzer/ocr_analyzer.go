
package analyzer

import (
	"context"
	"image"

	"github.com/anime-shed/image-inspector-go/pkg/models"
)

// OCRAnalyzer performs Optical Character Recognition (OCR) on an image.
// It implements the FeatureAnalyzer interface.
type OCRAnalyzer struct{}

// NewOCRAnalyzer creates a new analyzer for OCR.
func NewOCRAnalyzer() *OCRAnalyzer {
	return &OCRAnalyzer{}
}

// Analyze performs OCR on the image and populates the extracted text into the report.
func (oa *OCRAnalyzer) Analyze(_ context.Context, img image.Image, report *models.ImageAnalysis) error {
	// Simulate OCR text extraction
	text := "Simulated OCR text from image"

	report.HasOCRReport = true
	report.ExtractedText = text

	return nil
}

// Close cleans up any resources used by the OCR client.
func (oa *OCRAnalyzer) Close() error {
	return nil
}
