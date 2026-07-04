package app

import (
	"fmt"
	"math"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration.
type Config struct {
	// Weights holds the scoring weights configuration.
	Weights WeightsConfig `yaml:"weights"`
}

// WeightsConfig holds the individual scoring weights.
type WeightsConfig struct {
	// Commits is the weight for the commit frequency metric.
	Commits float64 `yaml:"commits"`
	// Contributors is the weight for the contributor diversity metric.
	Contributors float64 `yaml:"contributors"`
	// Churn is the weight for the code churn intensity metric.
	Churn float64 `yaml:"churn"`
	// Consistency is the weight for the consistency metric.
	Consistency float64 `yaml:"consistency"`
}

// DefaultConfig returns the default configuration values.
func DefaultConfig() Config {
	return Config{
		Weights: WeightsConfig{
			Commits:      0.30,
			Contributors: 0.20,
			Churn:        0.25,
			Consistency:  0.25,
		},
	}
}

// LoadConfig loads the configuration from a YAML file.
// If the file does not exist, it falls back to the default configuration.
func LoadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to decode YAML config: %w", err)
	}

	// Validate weights sum to 1.0 (with a small epsilon check)
	sum := cfg.Weights.Commits + cfg.Weights.Contributors + cfg.Weights.Churn + cfg.Weights.Consistency
	if math.Abs(sum-1.0) > 1e-6 {
		return Config{}, fmt.Errorf("invalid weights sum: expected 1.0, got %f", sum)
	}

	return cfg, nil
}
