package mainpkg

import (
	"os"
	"testing"
)

func TestServiceConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config_test.ini")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	configContent := `
[database]
host = localhost
port = 3306
username = user
password = pass
`
	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}
	tmpFile.Close()

	cfg, err := ServiceConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("ServiceConfig error: %v", err)
	}

	expectedValue := "localhost"
	actualValue := cfg.GetString("database.host")
	if expectedValue != actualValue {
		t.Errorf("Expected: %s, Got: %s", expectedValue, actualValue)
	}
}
