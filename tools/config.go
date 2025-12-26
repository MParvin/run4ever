package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the run4ever configuration structure
type Config struct {
	TelegramToken     string `yaml:"telegram_token"`
	TelegramChatID    string `yaml:"telegram_chat_id"`
	TelegramCustomAPI string `yaml:"telegram_custom_api"`
	NotifyMethod      string `yaml:"notify_method"`
	NotifyOn          string `yaml:"notify_on"`
}

// LoadConfig loads configuration from files and environment variables
// Priority: CLI flags > env vars > user config > system config
func LoadConfig(verbose bool) (*Config, error) {
	config := &Config{}

	// Load from system config first (lowest priority)
	systemConfigPath := "/etc/run4ever/config.yaml"
	if err := loadConfigFile(systemConfigPath, config, verbose); err != nil && !os.IsNotExist(err) {
		if verbose {
			fmt.Printf("Warning: failed to load system config: %v\n", err)
		}
	}

	// Load from user config (higher priority)
	homeDir := os.Getenv("HOME")
	if homeDir != "" {
		userConfigPath := filepath.Join(homeDir, ".config", "run4ever", "config.yaml")
		if err := loadConfigFile(userConfigPath, config, verbose); err != nil && !os.IsNotExist(err) {
			if verbose {
				fmt.Printf("Warning: failed to load user config: %v\n", err)
			}
		}
	}

	// Override with environment variables (higher priority)
	if token := os.Getenv("RUN4EVER_TELEGRAM_TOKEN"); token != "" {
		config.TelegramToken = token
	}
	if chatID := os.Getenv("RUN4EVER_TELEGRAM_CHAT_ID"); chatID != "" {
		config.TelegramChatID = chatID
	}
	if customAPI := os.Getenv("RUN4EVER_TELEGRAM_CUSTOM_API"); customAPI != "" {
		config.TelegramCustomAPI = customAPI
	}
	if notifyMethod := os.Getenv("RUN4EVER_NOTIFY_METHOD"); notifyMethod != "" {
		config.NotifyMethod = notifyMethod
	}
	if notifyOn := os.Getenv("RUN4EVER_NOTIFY_ON"); notifyOn != "" {
		config.NotifyOn = notifyOn
	}

	return config, nil
}

// loadConfigFile loads configuration from a YAML file
func loadConfigFile(path string, config *Config, verbose bool) error {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	// Check file permissions
	if err := checkConfigPermissions(path, verbose); err != nil {
		if verbose {
			fmt.Printf("Warning: config file %s has insecure permissions: %v\n", path, err)
		}
		// Continue loading despite permission warning
	}

	// Read and parse config file
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return fmt.Errorf("failed to parse config file: %w", err)
	}

	if verbose {
		fmt.Printf("Loaded config from: %s\n", path)
	}

	return nil
}

// checkConfigPermissions verifies that config file has 0600 permissions
func checkConfigPermissions(path string, verbose bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	mode := info.Mode().Perm()
	// Check if permissions are more permissive than 0600 (owner read/write only)
	if mode&0077 != 0 {
		return fmt.Errorf("config file %s has permissions %s, should be 0600 (owner read/write only)", path, mode.String())
	}

	return nil
}

