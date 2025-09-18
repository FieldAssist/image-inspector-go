package analyzer

import "image"

// metrics holds the basic image metrics.
type metrics struct {
	avgLuminance  float64
	avgSaturation float64
	avgR          float64
	avgG          float64
	avgB          float64
}

// MetricsCalculator defines the interface for calculating various image metrics.
// This allows for different implementations and simplifies testing.
type MetricsCalculator interface {
	// CalculateBasicMetrics computes fundamental metrics like average color and luminance.
	CalculateBasicMetrics(img image.Image) metrics

	// CalculateLaplacianVariance measures blurriness by detecting sharp edges.
	CalculateLaplacianVariance(gray *image.Gray) float64

	// CalculateBrightness computes the average brightness of a grayscale image.
	CalculateBrightness(gray *image.Gray) float64

	// DetectSkew determines the skew angle of the image, if any.
	DetectSkew(gray *image.Gray) *float64

	// DetectContours finds and counts the number of distinct shapes or objects.
	DetectContours(gray *image.Gray) int
}
