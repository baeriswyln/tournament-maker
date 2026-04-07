// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Hooks ────────────────────────────────────────────────────────────────────

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type User struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"     json:"id"`
	ParentID    *uuid.UUID `gorm:"type:uuid"                json:"parent_id"`
	VisibleName *string    `gorm:"not null"                 json:"visible_name"`
	Username    *string    `gorm:"uniqueIndex"              json:"username,omitempty"`
	Password    *string    `                                json:"-"` // hashed, never serialised
	IsActive    bool       `gorm:"not null;default:true"    json:"-"`
	CreatedAt   time.Time  `                                json:"created_at"`
	UpdatedAt   time.Time  `                                json:"updated_at"`

	// relations
	Parent *User `gorm:"foreignKey:ParentID"                 json:"-"`
}

// ─── Request / Response DTOs ─────────────────────────────────────────────────

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	RefreshToken string    `json:"refresh_token"`
	AccessToken  string    `json:"access_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	UserID       uuid.UUID `json:"user_id"`
}

type RefreshRequest struct {
	RefreshToken *string `json:"refresh_token"`
}

type CreateUserRequest struct {
	ID          *uuid.UUID `json:"id,omitempty"`
	VisibleName string     `json:"visible_name"`
	Username    *string    `json:"username,omitempty"`
	Password    *string    `json:"password,omitempty"`
}

type CreateUserResponse struct {
	ID           uuid.UUID  `json:"id"`
	AccessToken  *string    `json:"access_token,omitempty"`
	RefreshToken *string    `json:"refresh_token,omitempty"`
	ExpiresAt    *time.Time `json:"expires_at,omitempty"`
}

type UpdateUserRequest struct {
	VisibleName *string `json:"visible_name,omitempty"`
	Password    *string `json:"password,omitempty"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	VisibleName string    `json:"visible_name"`
	Username    *string   `json:"username,omitempty"`
}

// ─── JWT Claims ─────────────────────────────────────────────────

type Claims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
}
