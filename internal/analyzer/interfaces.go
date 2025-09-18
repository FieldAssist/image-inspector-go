package analyzer

import (
	"context"
	"image"

	"github.com/anime-shed/image-inspector-go/pkg/models"
)

// Analyzer is the primary, public-facing interface for the analyzer package.
// It provides a single, clear entry point for all image analysis operations.
type Analyzer interface {
	// Analyze performs an analysis of the given image based on the provided
	// configuration options. It returns a consolidated analysis report.
	Analyze(ctx context.Context, image image.Image, cfgs ...ReportConfigurator) (*models.ImageAnalysis, error)

	// Close releases any resources held by the analyzer.
	Close() error
}

// FeatureAnalyzer is an internal interface implemented by specialized, single-purpose analyzers
// (e.g., for quality, QR codes, or OCR). These are orchestrated by the CoreAnalyzer.
type FeatureAnalyzer interface {
	// Analyze performs a specific analysis feature on the image and contributes
	// its findings to the provided report.
	Analyze(ctx context.Context, image image.Image, report *models.ImageAnalysis) error
}
