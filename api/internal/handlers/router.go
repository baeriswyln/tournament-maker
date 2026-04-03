package handlers

import (
	"baeriswyl.dev/tournament-maker/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	auth *middleware.Authenticator,
	hAuth *AuthHandler,
	hUser *UserHandler,
	hTour *TournamentHandler,
	hPart *ParticipantHandler,
	hEnc *EncounterHandler,
) chi.Router {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {

		// Global Middlewares
		r.Use(chimiddleware.Logger)
		r.Use(chimiddleware.Recoverer)
		r.Use(chimiddleware.RealIP)
		r.Use(chimiddleware.RequestID)
		r.Use(cors.Handler(cors.Options{
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "access_token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300,
		}))

		r.Post("/auth", hAuth.Authenticate)
		r.Put("/auth", hAuth.Refresh)

		// Registrations without a login get a refresh token that is valid indefinitely.
		// On token refresh with such a token, the database is consulted to check if the user still has no login.
		// If that is not the case, the refresh fails.
		r.Post("/register", hUser.RegisterUser)

		r.Group(func(r chi.Router) {
			r.Use(auth.Middleware)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", hUser.GetUsers)
				r.Post("/", hUser.CreateUser)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", hUser.GetUserSingle)
					r.Put("/", hUser.UpdateUser)
					r.Delete("/", hUser.DeleteUser)
				})
			})

			r.Route("/tournaments", func(r chi.Router) {
				r.Get("/", hTour.GetTournaments)
				r.Post("/", hTour.CreateTournament)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", hTour.GetTournamentSingle)
					r.Put("/", hTour.UpdateTournament)
					r.Delete("/", hTour.DeleteTournament)
					r.Delete("/reset", hTour.ResetTournament)

					r.Route("/participants", func(r chi.Router) {
						r.Get("/", hPart.GetParticipants)
						r.Post("/", hPart.AddParticipants)
						r.Delete("/", hPart.DeleteParticipant)
					})

					r.Get("/scoreboard", hTour.GetScoreboard)

					r.Route("/encounters", func(r chi.Router) {
						r.Get("/", hEnc.GetEncounters)
						r.Post("/", hEnc.CreateEncounters)
						r.Put("/", hEnc.UpdateEncounterScore)
						r.Delete("/", hEnc.DeleteEncounterScore)
					})
				})
			})
		})
	})

	return r
}
