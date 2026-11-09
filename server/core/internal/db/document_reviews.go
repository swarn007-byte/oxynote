package db

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/oxynote/oxynote/server/core/internal/document"
)

// FetchRequestedReviews retrieves pending non-default branch reviews for a
// user in one organization. The default branch supplies the document URL title.
func (a *agent) FetchRequestedReviews(ctx context.Context, organizationID, userID string) ([]document.RequestedReview, error) {
	q, args := a.builder.Select(
		"documents.id AS document_id",
		"db.id AS branch_id",
		"main.document_name",
		"db.branch_name",
		"main.icon",
	).From("branch_reviewers br").
		Join("document_branches db ON db.id = br.fk_branch_id").
		Join("documents ON documents.id = db.fk_document_id AND documents.fk_organization_id = br.fk_organization_id").
		Join(`document_branches main ON main.fk_document_id = documents.id AND main."default" = true`).
		Where(sq.Eq{
			"br.fk_organization_id": organizationID,
			"br.fk_user_id":         userID,
			"br.currently_approved": false,
			`db."default"`:          false,
		}).OrderBy("main.document_name ASC", "documents.id ASC", "db.id ASC").MustSql()

	reviews := []document.RequestedReview{}

	if err := sqlx.SelectContext(ctx, a.sql, &reviews, q, args...); err != nil {
		return nil, err
	}

	return reviews, nil
}
