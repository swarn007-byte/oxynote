package db

import (
	"context"
	"testing"

	"github.com/oxynote/oxynote/server/core/internal/document"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_agent_FetchRequestedReviews(t *testing.T) {
	db := prepTempDB(t)
	users := prepUsers(t, db, 2)
	docs := prepDocuments(t, db, 2, nil)
	branches := prepDocumentBranches(t, db, 3, func(i int, d *document.Document) {
		parent := docs[0]
		if i == 2 {
			parent = docs[1]
		}

		d.ID = parent.ID
		d.OrganizationID = parent.OrganizationID
	})
	for i, b := range branches {
		prepBranchReviewers(t, db, 1, func(_ int, r *document.BranchReviewer) {
			r.BranchID = b.BranchID
			r.OrganizationID = b.OrganizationID
			r.UserID = users[0]
			r.CurrentlyApproved = i == 1
		})
	}

	prepBranchReviewers(t, db, 1, func(_ int, r *document.BranchReviewer) {
		r.BranchID = docs[0].BranchID
		r.OrganizationID = docs[0].OrganizationID
		r.UserID = users[0]
	})
	want := document.RequestedReview{
		DocumentID:   docs[0].ID,
		BranchID:     branches[0].BranchID,
		DocumentName: docs[0].DocumentName,
		BranchName:   branches[0].BranchName,
		Icon:         docs[0].Icon,
	}
	cc := map[string]struct {
		OrganizationID string
		UserID         string
		Canceled       bool
		Result         []document.RequestedReview
	}{
		"Only pending drafts in the workspace": {OrganizationID: docs[0].OrganizationID, UserID: users[0], Result: []document.RequestedReview{want}},
		"Another user":                         {OrganizationID: docs[0].OrganizationID, UserID: users[1], Result: []document.RequestedReview{}},
		"Another organization":                 {OrganizationID: "other", UserID: users[0], Result: []document.RequestedReview{}},
		"Canceled context":                     {OrganizationID: docs[0].OrganizationID, UserID: users[0], Canceled: true},
	}

	for cn, c := range cc {
		t.Run(cn, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			if c.Canceled {
				cancel()
			}

			got, err := db.FetchRequestedReviews(ctx, c.OrganizationID, c.UserID)
			if c.Canceled {
				require.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, c.Result, got)
		})
	}
}
