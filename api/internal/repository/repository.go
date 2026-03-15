// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package repository

import (
	"encoding/json"
	"errors"

	"baeriswyl.dev/tournament-maker/internal/models"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("already exists")

// ─── Database setup ───────────────────────────────────────────────────────────

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.Tournament{},
		&models.TournamentParticipant{},
		&models.Encounter{},
	)
}

// ─── helpers ─────────────────────────────────────────────────────────────────

func marshalTournament(t *models.Tournament) error {
	v, err := json.Marshal(t.Visibility)
	if err != nil {
		return err
	}
	t.VisibilityJSON = string(v)

	f, err := json.Marshal(t.Format)
	if err != nil {
		return err
	}
	t.FormatJSON = string(f)
	return nil
}

func unmarshalTournament(t *models.Tournament) {
	if t.VisibilityJSON != "" {
		_ = json.Unmarshal([]byte(t.VisibilityJSON), &t.Visibility)
	}
	if t.FormatJSON != "" {
		_ = json.Unmarshal([]byte(t.FormatJSON), &t.Format)
	}
}

func unmarshalTournaments(ts []models.Tournament) {
	for i := range ts {
		unmarshalTournament(&ts[i])
	}
}
