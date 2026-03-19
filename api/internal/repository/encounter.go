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

// ─── Encounter ────────────────────────────────────────────────────────────────

type EncounterRepository struct{ db *gorm.DB }

func NewEncounterRepository(db *gorm.DB) *EncounterRepository {
	return &EncounterRepository{db}
}

func (r *EncounterRepository) CreateBatch(encounters []models.Encounter) error {
	return r.db.Create(&encounters).Error
}

func (r *EncounterRepository) List(tournamentID uuid.UUID, userID *uuid.UUID) ([]models.Encounter, error) {
	q := r.db.Where("tournament_id = ?", tournamentID)
	if userID != nil {
		q = q.Where("party_a = ? OR party_b = ?", *userID, *userID)
	}
	var encounters []models.Encounter
	return encounters, q.Find(&encounters).Error
}

func (r *EncounterRepository) Get(tournamentID, partyA, partyB uuid.UUID, round string) (*models.Encounter, error) {
	var e models.Encounter
	err := r.db.First(&e,
		"tournament_id = ? AND party_a = ? AND party_b = ? AND round = ?",
		tournamentID, partyA, partyB, round).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &e, nil
}

func (r *EncounterRepository) Update(e *models.Encounter) error {
	return r.db.Save(e).Error
}

func (r *EncounterRepository) ClearScore(tournamentID, partyA, partyB uuid.UUID, round string) error {
	return r.db.Model(&models.Encounter{}).
		Where("tournament_id = ? AND party_a = ? AND party_b = ? AND round = ?",
			tournamentID, partyA, partyB, round).
		Updates(map[string]interface{}{
			"score_a": nil,
			"score_b": nil,
			"state":   models.EncounterStatePlanned,
		}).Error
}

func (r *EncounterRepository) DeleteEncounter(tournamentID, partyA, partyB uuid.UUID, round string) error {
	return r.db.
		Where("tournament_id = ? AND party_a = ? AND party_b = ? AND round = ?",
			tournamentID, partyA, partyB, round).
		Delete(&models.Encounter{}).Error
}
