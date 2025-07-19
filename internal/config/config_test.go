package config

import (
	"encoding/json"
	"os"
	"path"
	"testing"
)

func TestGetConfigFilePath(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	configPath, err := getConfigFilePath()
	if err != nil {
		t.Fatalf("getConfigFilePath() returned error: %v", err)
	}

	expectedPath := path.Join(homeDir, CONFIG_FILE_NAME)

	if configPath != expectedPath {
		t.Errorf("getConfigFilePath() = %s, want %s", configPath, expectedPath)
	}

	if !path.IsAbs(configPath) {
		t.Errorf("getConfigFilePath() should return an absolute path, got %s", configPath)
	}

	if path.Base(configPath) != CONFIG_FILE_NAME {
		t.Errorf("getConfigFilePath() should end with %s, got %s", CONFIG_FILE_NAME, path.Base(configPath))
	}
}

func TestGetConfigFilePath_Error(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	os.Setenv("HOME", "/invalid/path/that/should/not/exist")

	os.Unsetenv("HOME")

	_, err := getConfigFilePath()
	if err == nil {
		t.Error("getConfigFilePath() should return error when home directory cannot be determined")
	}
}

// Helper function to create a temporary config file for testing
func createTempConfigFile(t *testing.T, config Config) string {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "gator-config-test-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Marshal config to JSON
	jsonData, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal test config: %v", err)
	}

	// Write to temp file
	err = os.WriteFile(tmpFile.Name(), jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	return tmpFile.Name()
}

func TestRead_ValidConfig(t *testing.T) {
	// Create a test config
	testConfig := Config{
		DbURL:           "postgresql://localhost:5432/testdb",
		CurrentUserName: "testuser",
	}

	// Create temporary config file
	tempPath := createTempConfigFile(t, testConfig)
	defer os.Remove(tempPath)

	// Temporarily backup the real config file if it exists
	realConfigPath, _ := getConfigFilePath()
	backupPath := realConfigPath + ".backup"
	if _, err := os.Stat(realConfigPath); err == nil {
		os.Rename(realConfigPath, backupPath)
		defer os.Rename(backupPath, realConfigPath)
	}

	// Copy our test file to the real config location
	testData, err := os.ReadFile(tempPath)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	err = os.WriteFile(realConfigPath, testData, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config to real location: %v", err)
	}

	// Test the Read function
	config, err := Read()
	if err != nil {
		t.Fatalf("Read() returned error: %v", err)
	}

	// Verify the config matches
	if config.DbURL != testConfig.DbURL {
		t.Errorf("Read() DbURL = %s, want %s", config.DbURL, testConfig.DbURL)
	}

	if config.CurrentUserName != testConfig.CurrentUserName {
		t.Errorf("Read() CurrentUserName = %s, want %s", config.CurrentUserName, testConfig.CurrentUserName)
	}
}

func TestRead_FileNotFound(t *testing.T) {
	// Temporarily backup the real config file if it exists
	realConfigPath, _ := getConfigFilePath()
	backupPath := realConfigPath + ".backup"
	if _, err := os.Stat(realConfigPath); err == nil {
		os.Rename(realConfigPath, backupPath)
		defer os.Rename(backupPath, realConfigPath)
	}

	// Remove the config file to simulate file not found
	os.Remove(realConfigPath)

	// Test the Read function
	config, err := Read()
	if err == nil {
		t.Error("Read() should return error when config file doesn't exist")
	}

	// Verify empty config is returned
	if config.DbURL != "" || config.CurrentUserName != "" {
		t.Error("Read() should return empty config when file doesn't exist")
	}
}

func TestRead_InvalidJSON(t *testing.T) {
	// Temporarily backup the real config file if it exists
	realConfigPath, _ := getConfigFilePath()
	backupPath := realConfigPath + ".backup"
	if _, err := os.Stat(realConfigPath); err == nil {
		os.Rename(realConfigPath, backupPath)
		defer os.Rename(backupPath, realConfigPath)
	}

	// Write invalid JSON to the real config location
	invalidJSON := `{"db_url": "postgresql://localhost:5432/testdb", "current_user_name": "testuser",}`
	err := os.WriteFile(realConfigPath, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid JSON config file: %v", err)
	}

	// Test the Read function
	config, err := Read()
	if err == nil {
		t.Error("Read() should return error when JSON is invalid")
	}

	// Verify empty config is returned
	if config.DbURL != "" || config.CurrentUserName != "" {
		t.Error("Read() should return empty config when JSON is invalid")
	}
}

func TestSetUser(t *testing.T) {
	// Create initial config
	initialConfig := Config{
		DbURL:           "postgresql://localhost:5432/testdb",
		CurrentUserName: "olduser",
	}

	// Temporarily backup the real config file if it exists
	realConfigPath, _ := getConfigFilePath()
	backupPath := realConfigPath + ".backup"
	if _, err := os.Stat(realConfigPath); err == nil {
		os.Rename(realConfigPath, backupPath)
		defer os.Rename(backupPath, realConfigPath)
	}

	// Write initial config to the real location
	jsonData, err := json.Marshal(initialConfig)
	if err != nil {
		t.Fatalf("Failed to marshal initial config: %v", err)
	}

	err = os.WriteFile(realConfigPath, jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write initial config file: %v", err)
	}

	// Test SetUser method
	config := &Config{}
	*config = initialConfig

	newUserName := "newuser"
	err = config.SetUser(newUserName)
	if err != nil {
		t.Fatalf("SetUser() returned error: %v", err)
	}

	// Verify the config was updated
	if config.CurrentUserName != newUserName {
		t.Errorf("SetUser() CurrentUserName = %s, want %s", config.CurrentUserName, newUserName)
	}

	// Verify the file was written correctly
	readConfig, err := Read()
	if err != nil {
		t.Fatalf("Failed to read config after SetUser: %v", err)
	}

	if readConfig.CurrentUserName != newUserName {
		t.Errorf("Config file CurrentUserName = %s, want %s", readConfig.CurrentUserName, newUserName)
	}

	if readConfig.DbURL != initialConfig.DbURL {
		t.Errorf("Config file DbURL = %s, want %s", readConfig.DbURL, initialConfig.DbURL)
	}
}
