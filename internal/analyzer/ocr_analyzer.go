
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
	// This analyzer's purpose is to confirm that the OCR analysis path was triggered.
	// In a real-world scenario, this is where you would check for text-related quality
	// issues (e.g., text contrast, font size, etc.).
	// Actual text extraction is out of scope.

	report.HasOCRReport = true

	return nil
}

// Close cleans up any resources used by the OCR client.
func (oa *OCRAnalyzer) Close() error {
	// No resources to clean up in the current simulated implementation.
	return nil
}
