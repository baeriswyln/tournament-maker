// Package auth provides JWT-based authentication middleware for Chi routers.
// It supports short-lived access tokens, week-long refresh tokens, and
// permanent refresh tokens, with DB-backed user validation on refresh.
package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"baeriswyl.dev/tournament-maker/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContextKey is the type used for context keys in this package.
type ContextKey string

const (
	// ContextKeyUserID is the context key under which the authenticated user ID is stored.
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyClaims is the context key under which the full JWT claims are stored.
	ContextKeyClaims ContextKey = "claims"
)

// RefreshTokenKind distinguishes the two supported refresh token lifetimes.
type RefreshTokenKind string

const (
	// RefreshTokenWeek expires after 7 days.
	RefreshTokenWeek RefreshTokenKind = "week"
	// RefreshTokenPermanent never expires.
	RefreshTokenPermanent RefreshTokenKind = "permanent"
)

// Claims is the JWT payload stored inside every access token.
type Claims struct {
	UserID   string `json:"uid"`
	HasLogin bool   `json:"has_login"`
	jwt.RegisteredClaims
}

// RefreshClaims is the JWT payload stored inside every refresh token.
type RefreshClaims struct {
	UserID string           `json:"uid"`
	Kind   RefreshTokenKind `json:"kind"`
	jwt.RegisteredClaims
}

// TokenPair is returned after a successful login or token refresh.
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Authenticator is the main object that issues tokens and provides
// the Chi middleware handler.
type Authenticator struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewAuthenticator creates a new Authenticator with the given Config.
// It panics if AccessSecret or RefreshSecret are empty, or if UserValidator
// is nil (it is required for safe token refresh).
func NewAuthenticator(cfg *config.Config, db *gorm.DB) *Authenticator {
	if len(cfg.AccessSecret) == 0 {
		panic("auth: AccessSecret must not be empty")
	}
	if len(cfg.RefreshSecret) == 0 {
		panic("auth: RefreshSecret must not be empty")
	}
	if cfg.UserValidator == nil {
		panic("auth: UserValidator must not be nil")
	}
	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = 15 * time.Minute
	}
	return &Authenticator{cfg: cfg, db: db}
}

// -------------------------------------------------------------------
// Token issuance
// -------------------------------------------------------------------

// IssueTokenPair creates a fresh access token and a refresh token for userID.
// kind controls whether the refresh token lives for a week or forever.
func (a *Authenticator) IssueTokenPair(userID string, kind RefreshTokenKind) (TokenPair, error) {
	accessToken, expiresAt, err := a.newAccessToken(userID)
	if err != nil {
		return TokenPair{}, err
	}

	refreshToken, err := a.newRefreshToken(userID, kind)
	if err != nil {
		return TokenPair{}, err
	}

	return TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// newAccessToken signs and returns a new access token for userID.
func (a *Authenticator) newAccessToken(userID string) (string, time.Time, error) {
	expiresAt := time.Now().Add(a.cfg.AccessTokenTTL)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(a.cfg.AccessSecret)
	return signed, expiresAt, err
}

// newRefreshToken signs and returns a new refresh token for userID.
// Permanent tokens intentionally omit an ExpiresAt claim.
func (a *Authenticator) newRefreshToken(userID string, kind RefreshTokenKind) (string, error) {
	rc := jwt.RegisteredClaims{
		Subject:  userID,
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	if kind == RefreshTokenWeek {
		rc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour))
	}
	// For RefreshTokenPermanent we leave ExpiresAt nil — no expiry.

	claims := RefreshClaims{
		UserID:           userID,
		Kind:             kind,
		RegisteredClaims: rc,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.cfg.RefreshSecret)
}

// -------------------------------------------------------------------
// Token refresh
// -------------------------------------------------------------------

// Refresh validates the provided refresh token, checks the user in the DB,
// and issues a new TokenPair.  The new refresh token preserves the same Kind
// as the original one.
func (a *Authenticator) Refresh(ctx context.Context, rawRefreshToken string) (TokenPair, error) {
	rc, err := a.parseRefreshToken(rawRefreshToken)
	if err != nil {
		return TokenPair{}, err
	}

	// DB check — confirm the user is still active.
	if err = a.cfg.UserValidator(ctx, rc.UserID, a.db); err != nil {
		return TokenPair{}, ErrUserInvalid
	}

	return a.IssueTokenPair(rc.UserID, rc.Kind)
}

// parseRefreshToken verifies the signature and (if present) expiry of a
// refresh token, then returns its claims.
func (a *Authenticator) parseRefreshToken(raw string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(raw, &RefreshClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return a.cfg.RefreshSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// -------------------------------------------------------------------
// Chi middleware
// -------------------------------------------------------------------

// Middleware returns a Chi-compatible http.Handler that enforces access token
// authentication on every request.  On success the user ID and full claims are
// stored in the request context so downstream handlers can retrieve them with
// UserIDFromContext and ClaimsFromContext.
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, err := extractBearerToken(r)
		if err != nil {
			http.Error(w, "missing or malformed Authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := a.parseAccessToken(raw)
		if err != nil {
			switch {
			case errors.Is(err, ErrTokenExpired):
				http.Error(w, "access token expired", http.StatusUnauthorized)
			default:
				http.Error(w, "invalid access token", http.StatusUnauthorized)
			}
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ContextKeyClaims, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// parseAccessToken verifies the signature and expiry of an access token.
func (a *Authenticator) parseAccessToken(raw string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(raw, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return a.cfg.AccessSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// -------------------------------------------------------------------
// Context helpers
// -------------------------------------------------------------------

// UserIDFromContext retrieves the authenticated user ID from a request context.
// Returns an empty string if not present.
func GetIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(ContextKeyUserID).(string)
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(v)
	if err != nil {
		return uuid.Nil, false
	}

	return id, ok
}

// ClaimsFromContext retrieves the full access token claims from a request context.
// Returns nil if not present.
func ClaimsFromContext(ctx context.Context) *Claims {
	v, _ := ctx.Value(ContextKeyClaims).(*Claims)
	return v
}

func HasLogin(ctx context.Context) bool {
	v, ok := ctx.Value(ContextKeyClaims).(*Claims)
	return ok && v.HasLogin
}

// -------------------------------------------------------------------
// Internal helpers
// -------------------------------------------------------------------

// extractBearerToken pulls the token string from "Authorization: Bearer <token>".
func extractBearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", ErrMissingToken
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return "", ErrMalformedToken
	}
	if parts[1] == "" {
		return "", ErrMissingToken
	}
	return parts[1], nil
}
