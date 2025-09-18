package analyzer

// ReportConfig holds the configuration for generating an analysis report.
// It is configured using functional options (ReportConfigurator).
type ReportConfig struct {
	WithQuality bool
	WithOCR     bool
}

// ReportConfigurator is a function that modifies the ReportConfig.
type ReportConfigurator func(*ReportConfig)

// WithQuality is a functional option to enable quality analysis.
func WithQuality() ReportConfigurator {
	return func(c *ReportConfig) {
		c.WithQuality = true
	}
}

// WithOCR is a functional option to enable OCR analysis.
func WithOCR() ReportConfigurator {
	return func(c *ReportConfig) {
		c.WithOCR = true
	}
}
