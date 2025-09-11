package analyzer

import "image"

// ImageAnalyzer defines the main interface for image analysis
type ImageAnalyzer interface {
	// Analyze performs the analysis on the provided image.
	Analyze(img image.Image) (*AnalysisResult, error)
	// AnalyzeWithOCR performs OCR-specific image analysis (legacy method for backward compatibility).
	AnalyzeWithOCR(img image.Image, expectedText string) (*AnalysisResult, error)
	// AnalyzeWithOptions performs image analysis with enhanced parallel processing and memory optimization.
	AnalyzeWithOptions(img image.Image, options AnalysisOptions) (*AnalysisResult, error)

	// Cleanup releases resources held by the analyzer (e.g., stops its worker pool).
	Cleanup() error
}

// MetricsCalculator handles image metrics computation
type MetricsCalculator interface {
	CalculateBasicMetrics(img image.Image) metrics
	CalculateLaplacianVariance(gray *image.Gray) float64
	CalculateBrightness(gray *image.Gray) float64
	DetectSkew(gray *image.Gray) *float64
	DetectContours(gray *image.Gray) int
}

// QRDetector handles QR code detection
type QRDetector interface {
	DetectQRCode(img image.Image) bool
}

// OCRAnalyzer handles OCR-specific analysis
type OCRAnalyzer interface {
	PerformOCRAnalysis(img image.Image, expectedText string) AnalysisResult
}
