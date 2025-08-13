package config

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

//go:embed default.env
var defaultEnvContent string

type Config struct {
	EmailUsername string
	EmailPassword string
	EmailImapHost string
	EmailImapPort string
	EmailSmtpHost string
	EmailSmtpPort string
}

// GetConfigDir returns the user's config directory for the app
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(homeDir, ".config", "GoMail")
	return configDir, nil
}

// EnsureConfigExists creates the config directory and default .env file if they don't exist
func EnsureConfigExists() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	envPath := filepath.Join(configDir, ".env")

	// Check if .env file already exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Create .env file with default content
		if err := os.WriteFile(envPath, []byte(defaultEnvContent), 0600); err != nil {
			return fmt.Errorf("failed to create .env file: %w", err)
		}

		fmt.Printf("Created default configuration file at: %s\n", envPath)
		fmt.Println("Please edit this file with your email settings before running the application.")
		return fmt.Errorf("configuration file created, please edit it with your settings")
	}

	return nil
}

// InitConfig creates a new .env file, even if one already exists
func InitConfig() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	envPath := filepath.Join(configDir, "GoMail.env")

	// Check if .env file already exists
	if _, err := os.Stat(envPath); err == nil {
		fmt.Printf("Configuration file already exists at: %s\n", envPath)
		fmt.Print("Do you want to overwrite it? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Configuration initialization cancelled.")
			return nil
		}
	}

	// Create .env file with default content
	if err := os.WriteFile(envPath, []byte(defaultEnvContent), 0600); err != nil {
		return fmt.Errorf("failed to create .env file: %w", err)
	}

	fmt.Printf("Configuration file created at: %s\n", envPath)
	fmt.Println("Please edit this file with your email settings.")

	return nil
}

// LoadConfig loads configuration from the user's config directory
func LoadConfig() (*Config, error) {
	// Ensure config exists first
	if err := EnsureConfigExists(); err != nil {
		return nil, err
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	envPath := filepath.Join(configDir, ".env")

	// Load the .env file
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("failed to load .env file from %s: %w", envPath, err)
	}

	config := &Config{
		EmailUsername: os.Getenv("EMAIL_USERNAME"),
		EmailPassword: os.Getenv("EMAIL_PASSWORD"),
		EmailImapHost: os.Getenv("EMAIL_IMAP_HOST"),
		EmailImapPort: os.Getenv("EMAIL_IMAP_PORT"),
		EmailSmtpHost: os.Getenv("EMAIL_SMTP_HOST"),
		EmailSmtpPort: os.Getenv("EMAIL_SMTP_PORT"),
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Validate checks if all required configuration fields are set
func (c *Config) Validate() error {
	var missing []string

	if c.EmailUsername == "" {
		missing = append(missing, "EMAIL_USERNAME")
	}
	if c.EmailPassword == "" {
		missing = append(missing, "EMAIL_PASSWORD")
	}
	if c.EmailImapHost == "" {
		missing = append(missing, "EMAIL_IMAP_HOST")
	}
	if c.EmailImapPort == "" {
		missing = append(missing, "EMAIL_IMAP_PORT")
	}
	if c.EmailSmtpHost == "" {
		missing = append(missing, "EMAIL_SMTP_HOST")
	}
	if c.EmailSmtpPort == "" {
		missing = append(missing, "EMAIL_SMTP_PORT")
	}

	if len(missing) > 0 {
		configDir, _ := GetConfigDir()
		envPath := filepath.Join(configDir, ".env")
		return fmt.Errorf("missing required configuration values: %s. Please edit %s",
			strings.Join(missing, ", "), envPath)
	}

	return nil
}

// GetConfigPath returns the path to the user's .env file
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, ".env"), nil
}
