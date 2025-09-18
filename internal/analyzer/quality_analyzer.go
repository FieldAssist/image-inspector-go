package analyzer

import (
	"context"
	"image"

	"github.com/anime-shed/image-inspector-go/pkg/models"
	"golang.org/x/image/draw"
)

// QualityAnalyzer performs image quality checks, such as blurriness and overexposure.
// It implements the FeatureAnalyzer interface.
type QualityAnalyzer struct {
	metricsCalc MetricsCalculator
}

// NewQualityAnalyzer creates a new analyzer for image quality.
func NewQualityAnalyzer() *QualityAnalyzer {
	return &QualityAnalyzer{
		metricsCalc: NewMetricsCalculator(),
	}
}

// Analyze assesses the image for quality metrics like blurriness and overexposure
// and populates the results in the provided report.
func (qa *QualityAnalyzer) Analyze(_ context.Context, img image.Image, report *models.ImageAnalysis) error {
	// Create a grayscale version of the image for certain calculations.
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Src)

	report.HasQualityReport = true
	report.Blurriness = qa.metricsCalc.CalculateLaplacianVariance(gray)
	report.Overexposure = qa.metricsCalc.CalculateBrightness(gray)

	// Also, get the basic metrics.
	basicMetrics := qa.metricsCalc.CalculateBasicMetrics(img)
	report.AvgLuminance = basicMetrics.avgLuminance
	report.AvgSaturation = basicMetrics.avgSaturation

	return nil
}
