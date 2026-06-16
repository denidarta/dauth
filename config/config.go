package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	DatabaseURL     string
	PrivateKeyPath  string
	PublicKeyPath   string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	InvitationTTL   time.Duration
	Port            string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		PrivateKeyPath: os.Getenv("PRIVATE_KEY_PATH"),
		PublicKeyPath:  os.Getenv("PUBLIC_KEY_PATH"),
		Port:           os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.PrivateKeyPath == "" {
		return nil, fmt.Errorf("PRIVATE_KEY_PATH is required")
	}
	if cfg.PublicKeyPath == "" {
		return nil, fmt.Errorf("PUBLIC_KEY_PATH is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	var err error

	cfg.AccessTokenTTL, err = parseDuration(os.Getenv("ACCESS_TOKEN_TTL"), 15*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("ACCESS_TOKEN_TTL: %w", err)
	}

	cfg.RefreshTokenTTL, err = parseDuration(os.Getenv("REFRESH_TOKEN_TTL"), 720*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("REFRESH_TOKEN_TTL: %w", err)
	}

	cfg.InvitationTTL, err = parseDuration(os.Getenv("INVITATION_TTL"), 168*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("INVITATION_TTL: %w", err)
	}

	return cfg, nil
}

func parseDuration(val string, fallback time.Duration) (time.Duration, error) {
	if val == "" {
		return fallback, nil
	}
	return time.ParseDuration(val)
}
