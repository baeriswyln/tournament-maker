package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Config struct {
	DatabasePath  string
	ServerAddress string

	// AccessSecret is the HMAC secret used to sign access tokens.
	AccessSecret []byte
	// RefreshSecret is the HMAC secret used to sign refresh tokens.
	// Use a different value from AccessSecret so tokens are not interchangeable.
	RefreshSecret []byte
	// AccessTokenTTL is how long an access token remains valid (e.g. 15 * time.Minute).
	AccessTokenTTL time.Duration
	// UserValidator is called during token refresh to confirm the user still exists
	// and is allowed to authenticate. Return a non-nil error to reject the refresh.
	UserValidator UserValidatorFunc
}

// UserValidatorFunc is implemented by callers to check whether a user is still
// active in the database. userID is the subject stored in the refresh token.
type UserValidatorFunc func(ctx context.Context, userID string, db *gorm.DB) error

func Load(userValidatorFunc UserValidatorFunc) *Config {
	jwtExpiry := 24 * time.Hour
	if v := os.Getenv("JWT_EXPIRY_MINUTES"); v != "" {
		if h, err := strconv.Atoi(v); err == nil {
			jwtExpiry = time.Duration(h) * time.Minute
		}
	}

	dbPath := os.Getenv("DB_CONNECTION_PATH")
	if dbPath == "" {
		dbPath = "tournament.db"
	}

	jwtAccessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if jwtAccessSecret == "" {
		jwtAccessSecret = "change-me-in-production"
	}

	jwtRefreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if jwtRefreshSecret == "" {
		jwtRefreshSecret = "change-me-in-production"
	}

	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	return &Config{
		DatabasePath:  dbPath,
		ServerAddress: addr,

		AccessSecret:   []byte(jwtAccessSecret),
		RefreshSecret:  []byte(jwtRefreshSecret),
		AccessTokenTTL: jwtExpiry,
		UserValidator:  userValidatorFunc,
	}
}
