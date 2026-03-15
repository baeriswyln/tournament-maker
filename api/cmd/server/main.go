package main

import (
	"log"
	"net/http"

	"baeriswyl.dev/tournament-maker/config"
	"baeriswyl.dev/tournament-maker/internal/handlers"
	"baeriswyl.dev/tournament-maker/internal/middleware"
	"baeriswyl.dev/tournament-maker/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	cfg := config.Load(repository.ValidateUser)

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := gorm.Open(postgres.Open(cfg.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := repository.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	tournamentRepo := repository.NewTournamentRepository(db)
	participantRepo := repository.NewParticipantRepository(db)
	encounterRepo := repository.NewEncounterRepository(db)

	// ── Init Middlewares───────────────────────────────────────────────────────
	auth := middleware.NewAuthenticator(cfg, db)

	// ── Handlers ──────────────────────────────────────────────────────────────
	authHandler := handlers.NewAuthHandler(userRepo, auth)
	userHandler := handlers.NewUserHandler(userRepo, auth)
	tourHandler := handlers.NewTournamentHandler(tournamentRepo)
	partHandler := handlers.NewParticipantHandler(participantRepo, tournamentRepo, userRepo)
	encHandler := handlers.NewEncounterHandler(encounterRepo, participantRepo, tournamentRepo, userRepo)

	r := handlers.NewRouter(auth, authHandler, userHandler, tourHandler, partHandler, encHandler)

	log.Printf("Tournament Maker API listening on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
