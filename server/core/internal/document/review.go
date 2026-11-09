package document

import "github.com/rs/xid"

// RequestedReview identifies a branch awaiting the caller's review.
type RequestedReview struct {
	// DocumentID specifies the document that owns the branch.
	DocumentID xid.ID `json:"documentId" db:"document_id"`

	// BranchID specifies the branch to review.
	BranchID xid.ID `json:"branchId" db:"branch_id"`

	// DocumentName specifies the default branch's title used in document URLs.
	DocumentName string `json:"documentName" db:"document_name"`

	// BranchName specifies the name of the branch to review.
	BranchName string `json:"branchName" db:"branch_name"`

	// Icon specifies the document's default branch icon.
	Icon string `json:"icon" db:"icon"`
}
