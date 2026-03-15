// Brief   :
// Author  : Nicolas Baeriswyl <nicolas@baeriswyl.dev>
// Created : 13.03.2026

package repository

import (
	"baeriswyl.dev/tournament-maker/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ─── Participant ──────────────────────────────────────────────────────────────

type ParticipantRepository struct{ db *gorm.DB }

func NewParticipantRepository(db *gorm.DB) *ParticipantRepository {
	return &ParticipantRepository{db}
}

func (r *ParticipantRepository) Add(tournamentID, userID uuid.UUID) error {
	p := models.TournamentParticipant{
		TournamentID: tournamentID,
		UserID:       userID,
	}
	return r.db.Create(&p).Error
}

func (r *ParticipantRepository) Remove(tournamentID, userID uuid.UUID) error {
	return r.db.Delete(&models.TournamentParticipant{},
		"tournament_id = ? AND user_id = ?", tournamentID, userID).Error
}

func (r *ParticipantRepository) ListWithStats(tournamentID uuid.UUID) ([]models.ParticipantResponse, error) {
	var participants []models.TournamentParticipant
	if err := r.db.Preload("User").
		Where("tournament_id = ?", tournamentID).
		Find(&participants).Error; err != nil {
		return nil, err
	}

	result := make([]models.ParticipantResponse, 0, len(participants))
	for _, p := range participants {
		var total, finished int64
		r.db.Model(&models.Encounter{}).
			Where("tournament_id = ? AND (party_a = ? OR party_b = ?)", tournamentID, p.UserID, p.UserID).
			Count(&total)
		r.db.Model(&models.Encounter{}).
			Where("tournament_id = ? AND (party_a = ? OR party_b = ?) AND state = ?",
				tournamentID, p.UserID, p.UserID, models.EncounterStateFinal).
			Count(&finished)

		// Calculate score: count wins (higher score)
		var score int64
		var encounters []models.Encounter
		r.db.Where("tournament_id = ? AND (party_a = ? OR party_b = ?) AND state = ?",
			tournamentID, p.UserID, p.UserID, models.EncounterStateFinal).
			Find(&encounters)
		for _, e := range encounters {
			if e.ScoreA != nil && e.ScoreB != nil {
				if e.PartyA == p.UserID && *e.ScoreA > *e.ScoreB {
					score++
				} else if e.PartyB == p.UserID && *e.ScoreB > *e.ScoreA {
					score++
				}
			}
		}

		result = append(result, models.ParticipantResponse{
			ID:          p.UserID,
			VisibleName: p.User.VisibleName,
			Encounters: models.EncounterSummary{
				Total:    int(total),
				Finished: int(finished),
				Score:    int(score),
			},
		})
	}
	return result, nil
}
