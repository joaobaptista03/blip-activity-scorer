package app

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("LoadingNonExistentConfigReturnsDefaults", func(t *testing.T) {
		cfg, err := LoadConfig("non-existent-config-file.yaml")
		if err != nil {
			t.Fatalf("unexpected error loading non-existent config: %v", err)
		}
		if cfg.Weights.Commits != 0.30 || cfg.Weights.Contributors != 0.20 {
			t.Errorf("expected default config, got %+v", cfg)
		}
	})

	t.Run("ValidTemporaryYAMLConfigParsing", func(t *testing.T) {
		tempFile, err := os.CreateTemp("", "config-test-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFile.Name())

		yamlContent := `
weights:
  commits: 0.10
  contributors: 0.40
  churn: 0.30
  consistency: 0.20
`
		if _, err := tempFile.WriteString(yamlContent); err != nil {
			t.Fatalf("failed to write temp file content: %v", err)
		}
		_ = tempFile.Close()

		cfg, err := LoadConfig(tempFile.Name())
		if err != nil {
			t.Fatalf("unexpected error loading valid config: %v", err)
		}
		if cfg.Weights.Commits != 0.10 || cfg.Weights.Contributors != 0.40 || cfg.Weights.Churn != 0.30 || cfg.Weights.Consistency != 0.20 {
			t.Errorf("incorrect weights loaded: %+v", cfg.Weights)
		}
	})

	t.Run("InvalidYAMLConfigValidation", func(t *testing.T) {
		tempFileInvalid, err := os.CreateTemp("", "config-test-invalid-*.yaml")
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		defer os.Remove(tempFileInvalid.Name())

		yamlContentInvalid := `
weights:
  commits: 0.10
  contributors: 0.10
  churn: 0.10
  consistency: 0.10
`
		if _, err := tempFileInvalid.WriteString(yamlContentInvalid); err != nil {
			t.Fatalf("failed to write temp file content: %v", err)
		}
		_ = tempFileInvalid.Close()

		_, err = LoadConfig(tempFileInvalid.Name())
		if err == nil {
			t.Error("expected validation error for weights not summing to 1.0, got nil")
		}
	})
}
