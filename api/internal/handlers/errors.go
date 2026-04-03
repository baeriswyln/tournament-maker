package handlers

import (
	"net/http"

	"baeriswyl.dev/tournament-maker/internal/apierror"
)

var (
	// ── Auth ──────────────────────────────────────────────────────────────────
	errAuthTokenMissing = &apierror.APIError{
		Status:    http.StatusUnauthorized,
		ErrorCode: "AUTH_TOKEN_MISSING",
		Title:     "Authentication required",
		Details:   "No access token was provided.",
	}
	errAuthTokenExpired = &apierror.APIError{
		Status:    http.StatusUnauthorized,
		ErrorCode: "AUTH_TOKEN_EXPIRED",
		Title:     "Token expired",
		Details:   "The access token has expired. Please refresh your session.",
	}
	errAuthTokenInvalid = &apierror.APIError{
		Status:    http.StatusUnauthorized,
		ErrorCode: "AUTH_TOKEN_INVALID",
		Title:     "Invalid token",
		Details:   "The provided access token is invalid.",
	}
	errAuthForbidden = &apierror.APIError{
		Status:    http.StatusForbidden,
		ErrorCode: "AUTH_FORBIDDEN",
		Title:     "Forbidden",
		Details:   "You do not have permission to perform this action.",
	}
	errAuthInvalidCredentials = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "AUTH_INVALID_CREDENTIALS",
		Title:     "Invalid credentials",
		Details:   "The username or password is incorrect.",
	}
	errAuthRefreshInvalid = &apierror.APIError{
		Status:    http.StatusForbidden,
		ErrorCode: "AUTH_REFRESH_INVALID",
		Title:     "Invalid refresh token",
		Details:   "The provided refresh token is invalid or has expired.",
	}

	// ── Request ───────────────────────────────────────────────────────────────
	errRequestInvalidBody = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "REQUEST_INVALID_BODY",
		Title:     "Invalid request body",
		Details:   "The request body could not be parsed or is missing required fields.",
	}
	errRequestInvalidID = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "REQUEST_INVALID_ID",
		Title:     "Invalid identifier",
		Details:   "The provided ID is not a valid UUID.",
	}

	// ── User ──────────────────────────────────────────────────────────────────
	errUserNotFound = &apierror.APIError{
		Status:    http.StatusNotFound,
		ErrorCode: "USER_NOT_FOUND",
		Title:     "User not found",
		Details:   "No user was found with the given identifier.",
	}
	errUserUsernameTaken = &apierror.APIError{
		Status:    http.StatusConflict,
		ErrorCode: "USER_USERNAME_TAKEN",
		Title:     "Username taken",
		Details:   "This username is already in use. Please choose another.",
	}
	errUserCredentialsIncomplete = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "USER_CREDENTIALS_INCOMPLETE",
		Title:     "Incomplete credentials",
		Details:   "Username and password must both be provided or both be omitted.",
	}
	errUserPasswordTooShort = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "USER_PASSWORD_TOO_SHORT",
		Title:     "Password too short",
		Details:   "Password must be at least 8 characters long.",
	}
	errUserNameRequired = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "USER_NAME_REQUIRED",
		Title:     "Name required",
		Details:   "A visible name is required to create this user.",
	}

	// ── Tournament ────────────────────────────────────────────────────────────
	errTournamentNotFound = &apierror.APIError{
		Status:    http.StatusNotFound,
		ErrorCode: "TOURNAMENT_NOT_FOUND",
		Title:     "Tournament not found",
		Details:   "No tournament was found with the given identifier.",
	}
	errTournamentNotOpen = &apierror.APIError{
		Status:    http.StatusConflict,
		ErrorCode: "TOURNAMENT_NOT_OPEN",
		Title:     "Tournament not open",
		Details:   "Participants can only join tournaments that are in the 'open' state.",
	}

	// ── Encounter ─────────────────────────────────────────────────────────────
	errEncounterNotFound = &apierror.APIError{
		Status:    http.StatusNotFound,
		ErrorCode: "ENCOUNTER_NOT_FOUND",
		Title:     "Encounter not found",
		Details:   "No encounter was found matching the given participants and round.",
	}
	errEncounterConfirmForbidden = &apierror.APIError{
		Status:    http.StatusForbidden,
		ErrorCode: "ENCOUNTER_CONFIRM_FORBIDDEN",
		Title:     "Score confirmation not allowed",
		Details:   "You cannot confirm this score: either you submitted the original vote, or the submitted score does not match.",
	}
	errEncounterAlreadyFinal = &apierror.APIError{
		Status:    http.StatusBadRequest,
		ErrorCode: "ENCOUNTER_ALREADY_FINAL",
		Title:     "Encounter already finalized",
		Details:   "The score for this encounter has already been finalized and cannot be changed.",
	}

	// ── Internal ──────────────────────────────────────────────────────────────
	errInternal = &apierror.APIError{
		Status:    http.StatusInternalServerError,
		ErrorCode: "INTERNAL_ERROR",
		Title:     "Internal server error",
		Details:   "An unexpected error occurred. Please try again later.",
	}
)