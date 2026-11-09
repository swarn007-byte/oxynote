// Package review provides the current user's requested reviews.
package review

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/oxynote/oxynote/server/core/internal/document"
	"github.com/oxynote/oxynote/server/core/internal/server/internal/auth"
	"github.com/oxynote/oxynote/server/core/pkg/httpserver"
)

// Handler holds dependencies for requested reviews.
type Handler struct {
	log *slog.Logger
	db  DB
}

// NewHandler creates a requested reviews handler.
func NewHandler(log *slog.Logger, db DB) *Handler {
	return &Handler{
		log: log.With("component", "review-handler"),
		db:  db,
	}
}

// FetchRequestedReviews fetches the caller's pending reviews in their workspace.
func (h *Handler) FetchRequestedReviews(w http.ResponseWriter, r *http.Request) {
	session, ok := auth.RequireSession(h.log, w, r)
	if !ok {
		return
	}

	reviews, err := h.db.FetchRequestedReviews(r.Context(), session.ActiveOrganizationID, session.UserID)
	if err != nil {
		httpserver.RespondError(h.log, w, err)
		return
	}

	httpserver.Respond(h.log, w, reviews, http.StatusOK)
}

// DB specifies the queries required by requested reviews.
//
//go:generate ../../../../scripts/codegen/mock -t internal DB db
type DB interface {
	// FetchRequestedReviews should fetch the user's pending reviews in the organization.
	FetchRequestedReviews(ctx context.Context, organizationID, userID string) ([]document.RequestedReview, error)
}
