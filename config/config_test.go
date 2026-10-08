package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReadsEnvFileWithoutOverridingProcessEnvironment(t *testing.T) {
	keys := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "APP_PORT"}
	previous := make(map[string]string, len(keys))
	wasSet := make(map[string]bool, len(keys))
	for _, key := range keys {
		previous[key], wasSet[key] = os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() {
		for _, key := range keys {
			if wasSet[key] {
				_ = os.Setenv(key, previous[key])
			} else {
				_ = os.Unsetenv(key)
			}
		}
	})
	if err := os.Setenv("DB_HOST", "process-host"); err != nil {
		t.Fatal(err)
	}

	filename := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(filename, []byte("DB_HOST=file-host\nDB_PORT=5433\nDB_NAME=pokemon_test\nAPP_PORT=9090\n"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := load(filename)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.DBHost != "process-host" {
		t.Errorf("DBHost = %q, want process-host", cfg.DBHost)
	}
	if cfg.DBPort != "5433" || cfg.DBName != "pokemon_test" || cfg.AppPort != "9090" {
		t.Errorf("file configuration not loaded: DBPort=%q DBName=%q AppPort=%q", cfg.DBPort, cfg.DBName, cfg.AppPort)
	}
}
