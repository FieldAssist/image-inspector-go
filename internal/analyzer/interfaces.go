package analyzer

import (
	"image"
	"github.com/anime-shed/image-inspector-go/internal/loader"
)

// ImageAnalyzer defines the main interface for image analysis
type ImageAnalyzer interface {
	// Analyze performs the analysis on the provided image.
	Analyze(img *loader.FastImage) (*AnalysisResult, error)
	// AnalyzeWithOCR performs OCR-specific image analysis (legacy method for backward compatibility).
	AnalyzeWithOCR(img *loader.FastImage, expectedText string) (*AnalysisResult, error)
	// AnalyzeWithOptions performs image analysis with enhanced parallel processing and memory optimization.
	AnalyzeWithOptions(img *loader.FastImage, options AnalysisOptions) (*AnalysisResult, error)

	// Cleanup releases resources held by the analyzer (e.g., stops its worker pool).
	Cleanup() error
}

// MetricsCalculator handles image metrics computation
type MetricsCalculator interface {
	CalculateBasicMetrics(img *loader.FastImage) (metrics, error)
	CalculateLaplacianVariance(img *loader.FastImage) (float64, error)
	CalculateBrightness(img *loader.FastImage) (float64, error)
	DetectSkew(img *loader.FastImage) (*float64, error)
	DetectContours(img *loader.FastImage) (int, error)
}

// QRDetector handles QR code detection
type QRDetector interface {
	DetectQRCode(img *loader.FastImage) bool
}

// OCRAnalyzer handles OCR-specific analysis
type OCRAnalyzer interface {
	PerformOCRAnalysis(img image.Image, expectedText string) AnalysisResult
}
