package review

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oxynote/oxynote/server/core/internal/document"
	"github.com/oxynote/oxynote/server/core/internal/server/internal/auth"
	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func Test_NewHandler(t *testing.T) {
	t.Parallel()

	db := &DBMock{}
	h := NewHandler(slog.New(slog.DiscardHandler), db)
	require.NotNil(t, h)
	assert.NotNil(t, h.log)
	assert.Same(t, db, h.db)
}

func Test_Handler_FetchRequestedReviews(t *testing.T) {
	t.Parallel()

	id := xid.New()
	branchID := xid.New()
	cc := map[string]struct {
		NoSession bool
		Result    []document.RequestedReview
		Err       error
		Code      int
		Body      string
		Calls     int
	}{
		"No session": {
			NoSession: true,
			Code:      http.StatusUnauthorized,
			Body:      `{"code":"account.not_authenticated","message":"not authenticated"}`,
		},
		"Error returned by DB.FetchRequestedReviews": {
			Err:   assert.AnError,
			Code:  http.StatusInternalServerError,
			Body:  `{"code":"general","message":"internal server error"}`,
			Calls: 1,
		},
		"No pending reviews": {
			Result: []document.RequestedReview{},
			Code:   http.StatusOK,
			Body:   `[]`,
			Calls:  1,
		},
		"Pending branch": {
			Result: []document.RequestedReview{{DocumentID: id, BranchID: branchID, DocumentName: "Runbook", BranchName: "draft", Icon: "lucide:file"}},
			Code:   http.StatusOK,
			Body:   `[{"documentId":"` + id.String() + `","branchId":"` + branchID.String() + `","documentName":"Runbook","branchName":"draft","icon":"lucide:file"}]`,
			Calls:  1,
		},
	}

	for cn, c := range cc {
		t.Run(cn, func(t *testing.T) {
			t.Parallel()

			db := &DBMock{
				FetchRequestedReviewsFunc: func(_ context.Context, organizationID, userID string) ([]document.RequestedReview, error) {
					assert.Equal(t, "org", organizationID)
					assert.Equal(t, "user", userID)

					return c.Result, c.Err
				},
			}
			req := httptest.NewRequest(http.MethodGet, "http://test.com/?userId=other&organizationId=other", nil)
			if !c.NoSession {
				req = req.WithContext(auth.AddSessionToContext(req.Context(), auth.Session{UserID: "user", ActiveOrganizationID: "org"}))
			}

			rec := httptest.NewRecorder()
			NewHandler(slog.New(slog.DiscardHandler), db).FetchRequestedReviews(rec, req)

			assert.Equal(t, c.Code, rec.Code)
			assert.JSONEq(t, c.Body, rec.Body.String())
			assert.Len(t, db.FetchRequestedReviewsCalls(), c.Calls)
		})
	}
}
