// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package repository

import (
	"context"
	"errors"
	"fmt"

	"baeriswyl.dev/tournament-maker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── User ─────────────────────────────────────────────────────────────────────

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db} }

func ValidateUser(ctx context.Context, userID string, db *gorm.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	var user models.User
	db.First(&user, "id = ? AND is_active AND username is NULL", id) //todo more conditions required?

	return nil
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("%w: %v", ErrConflict, err)
	}
	return nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "id = ? AND is_active", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, "username = ? AND is_active", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ListByOwnerTournaments returns users participating in any tournament owned by ownerID.
func (r *UserRepository) ListByOwnerTournaments(ownerID uuid.UUID) ([]models.User, error) {
	var users []models.User
	err := r.db.
		Joins("JOIN tournament_participants tp ON tp.user_id = users.id").
		Joins("JOIN tournaments t ON t.id = tp.tournament_id").
		Where("t.owner_id = ? AND users.is_active", ownerID).
		Distinct("users.*").
		Find(&users).Error
	return users, err
}

func (r *UserRepository) IsOwner(userID uuid.UUID, ownerID uuid.UUID) bool {
	var count int64
	r.db.Model(&models.User{}).
		Where("id = ? AND owner_id = ?", userID, ownerID).
		Count(&count)
	return count > 0
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) SoftDelete(id uuid.UUID) error {
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}
	user.OwnerID = nil
	user.VisibleName = nil
	user.Username = nil
	user.Password = nil
	user.IsActive = false

	return r.Update(user)
}
