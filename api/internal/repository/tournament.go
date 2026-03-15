// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package repository

import (
	"errors"

	"baeriswyl.dev/tournament-maker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Tournament ───────────────────────────────────────────────────────────────

type TournamentRepository struct{ db *gorm.DB }

func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	return &TournamentRepository{db}
}

func (r *TournamentRepository) Create(t *models.Tournament) error {
	if err := marshalTournament(t); err != nil {
		return err
	}
	return r.db.Create(t).Error
}

func (r *TournamentRepository) GetByID(id uuid.UUID) (*models.Tournament, error) {
	var t models.Tournament
	if err := r.db.First(&t, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	unmarshalTournament(&t)
	return &t, nil
}

func (r *TournamentRepository) ListOwned(ownerID uuid.UUID) ([]models.Tournament, error) {
	var ts []models.Tournament
	err := r.db.Where("owner_id = ?", ownerID).Find(&ts).Error
	unmarshalTournaments(ts)
	return ts, err
}

func (r *TournamentRepository) ListParticipants(tournamentID uuid.UUID) ([]models.User, error) {
	var us []models.User
	err := r.db.
		Joins("JOIN tournament_participants tp ON tp.user_id = users.id").
		Where("tp.tournament_id = ?", tournamentID).
		Find(&us).Error

	return us, err
}

func (r *TournamentRepository) ListParticipating(userID uuid.UUID) ([]models.Tournament, error) {
	var ts []models.Tournament
	err := r.db.
		Joins("JOIN tournament_participants tp ON tp.tournament_id = tournaments.id").
		Where("tp.user_id = ? AND tournaments.owner_id != ?", userID, userID).
		Find(&ts).Error
	unmarshalTournaments(ts)
	return ts, err
}

func (r *TournamentRepository) Update(t *models.Tournament) error {
	if err := marshalTournament(t); err != nil {
		return err
	}
	return r.db.Save(t).Error
}

func (r *TournamentRepository) Delete(id uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Encounter{}, "tournament_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.TournamentParticipant{}, "tournament_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Tournament{}, "id = ?", id).Error
	})
}

func (r *TournamentRepository) Reset(id uuid.UUID) error {
	return r.db.Delete(&models.Encounter{}, "tournament_id = ?", id).Error
}

// IsRelated returns true when userID is owner or participant of the tournament.
func (r *TournamentRepository) IsRelated(tournamentID, userID uuid.UUID) bool {
	var t models.Tournament
	if err := r.db.First(&t, "id = ?", tournamentID).Error; err != nil {
		return false
	}
	if t.OwnerID == userID {
		return true
	}
	var count int64
	r.db.Model(&models.TournamentParticipant{}).
		Where("tournament_id = ? AND user_id = ?", tournamentID, userID).
		Count(&count)
	return count > 0
}

// IsOwner returns true when userID owns the tournament.
func (r *TournamentRepository) IsOwner(tournamentID, userID uuid.UUID) bool {
	var count int64
	r.db.Model(&models.Tournament{}).
		Where("id = ? AND owner_id = ?", tournamentID, userID).
		Count(&count)
	return count > 0
}
