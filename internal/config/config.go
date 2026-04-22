package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// ProfileConfig represents the configuration for a profile
type ProfileConfig struct {
	AuthToken string
	BaseURL   string
	Model     string
}

// GetConfigDir returns the path to the configuration directory
func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home is not available
		return ".config/claudes"
	}
	return filepath.Join(home, ".config", "claudes")
}

// InitConfigDir creates the configuration directory if it doesn't exist
func InitConfigDir() error {
	configDir := GetConfigDir()
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return os.MkdirAll(configDir, 0755)
	}
	return nil
}

// ProfileExists checks if a profile already exists
func ProfileExists(profile string) bool {
	configDir := GetConfigDir()
	envPath := filepath.Join(configDir, fmt.Sprintf("%s.env", profile))
	_, err := os.Stat(envPath)
	return !os.IsNotExist(err)
}

// ListProfiles returns all available profile names (without .env extension)
func ListProfiles() ([]string, error) {
	configDir := GetConfigDir()
	files, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("error reading config directory: %w", err)
	}

	var profiles []string
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if strings.HasSuffix(name, ".env") {
			profileName := strings.TrimSuffix(name, ".env")
			profiles = append(profiles, profileName)
		}
	}

	return profiles, nil
}

// LoadProfile loads environment variables from a .env file for the given profile
func LoadProfile(profile string) ([]string, error) {
	configDir := GetConfigDir()
	envPath := filepath.Join(configDir, fmt.Sprintf("%s.env", profile))

	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("profile '%s' not found. Create %s first", profile, envPath)
	}

	envMap, err := godotenv.Read(envPath)
	if err != nil {
		return nil, fmt.Errorf("error loading profile %s: %w", profile, err)
	}

	var envVars []string
	for key, value := range envMap {
		envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
	}

	return envVars, nil
}

// SaveProfile creates or updates a profile .env file
func SaveProfile(profile string, config ProfileConfig) error {
	configDir := GetConfigDir()
	envPath := filepath.Join(configDir, fmt.Sprintf("%s.env", profile))

	content := fmt.Sprintf(
		"ANTHROPIC_AUTH_TOKEN=%s\nANTHROPIC_BASE_URL=%s\nANTHROPIC_MODEL=%s\n",
		config.AuthToken,
		config.BaseURL,
		config.Model,
	)

	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("error saving profile %s: %w", profile, err)
	}

	return nil
}

// DeleteProfile deletes a profile .env file
func DeleteProfile(profile string) error {
	if !ProfileExists(profile) {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	configDir := GetConfigDir()
	envPath := filepath.Join(configDir, fmt.Sprintf("%s.env", profile))

	if err := os.Remove(envPath); err != nil {
		return fmt.Errorf("error deleting profile %s: %w", profile, err)
	}

	return nil
}
