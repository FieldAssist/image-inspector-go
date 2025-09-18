
package models

// QRCode represents a single QR code found in an image.
type QRCode struct {
	ID              uint   `gorm:"primaryKey"`
	ImageAnalysisID string `gorm:"not null"` // Foreign key to ImageAnalysis
	Data            string `json:"data"`
}

// ImageAnalysis holds the consolidated results of all analyses performed on an image.
// It serves as the central data model for analysis reports.
type ImageAnalysis struct {
	ImageID string `gorm:"primaryKey" json:"image_id"`

	// Quality Analysis Report
	HasQualityReport bool    `json:"has_quality_report"`
	Blurriness       float64 `json:"blurriness,omitempty"`
	Overexposure     float64 `json:"overexposure,omitempty"`
	AvgLuminance     float64 `json:"avg_luminance,omitempty"`
	AvgSaturation    float64 `json:"avg_saturation,omitempty"`

	// OCR Analysis Report
	HasOCRReport  bool   `json:"has_ocr_report"`
	ExtractedText string `json:"extracted_text,omitempty"`

	// QR Code Analysis Report
	HasQRReport bool     `json:"has_qr_report"`
	QRCodeCount int      `json:"qr_code_count,omitempty"`
	QRCodes     []QRCode `gorm:"foreignKey:ImageAnalysisID" json:"qr_codes,omitempty"`
}
