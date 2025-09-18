package analyzer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportConfig(t *testing.T) {
	t.Run("Default config is empty", func(t *testing.T) {
		config := &ReportConfig{}
		assert.False(t, config.WithQuality, "Expected WithQuality to be false by default")
		assert.False(t, config.WithOCR, "Expected WithOCR to be false by default")
	})

	t.Run("Applies WithQuality option", func(t *testing.T) {
		config := &ReportConfig{}
		WithQuality()(config)
		assert.True(t, config.WithQuality, "Expected WithQuality to be true")
	})

	t.Run("Applies WithOCR option", func(t *testing.T) {
		config := &ReportConfig{}
		WithOCR()(config)
		assert.True(t, config.WithOCR, "Expected WithOCR to be true")
	})

	t.Run("Applies multiple options", func(t *testing.T) {
		config := &ReportConfig{}
		configs := []ReportConfigurator{WithQuality(), WithOCR()}
		for _, cfg := range configs {
			cfg(config)
		}
		assert.True(t, config.WithQuality, "Expected WithQuality to be true")
		assert.True(t, config.WithOCR, "Expected WithOCR to be true")
	})
}
