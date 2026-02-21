package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CredentialsPath string
	SourceSheetID   string
	DestSheetID     string
	Timezone        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		CredentialsPath: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		SourceSheetID:   os.Getenv("SOURCE_SHEET_ID"),
		DestSheetID:     os.Getenv("DEST_SHEET_ID"),
		Timezone:        os.Getenv("TIMEZONE"),
	}

	if cfg.CredentialsPath == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is required")
	}
	if cfg.SourceSheetID == "" {
		return nil, fmt.Errorf("SOURCE_SHEET_ID is required")
	}
	if cfg.DestSheetID == "" {
		return nil, fmt.Errorf("DEST_SHEET_ID is required")
	}
	if cfg.Timezone == "" {
		cfg.Timezone = "Asia/Jakarta"
	}

	return cfg, nil
}
