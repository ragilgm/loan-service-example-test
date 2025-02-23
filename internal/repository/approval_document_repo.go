package repo

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=$PROJ/internal/generated/mock/mock_$GOPACKAGE/$GOFILE

// ApprovalDocument represents the structure of approval documents
type (
	ApprovalDocument struct {
		ID             int64      `db:"id"`
		LoanApprovalID int64      `db:"loan_approval_id"`
		DocumentType   string     `db:"document_type"`
		FileURL        string     `db:"file_url"`
		Description    *string    `db:"description"`
		CreatedAt      time.Time  `db:"created_at"`
		UpdatedAt      time.Time  `db:"updated_at"`
		DeletedAt      *time.Time `db:"deleted_at"`
	}
)

type ApprovalDocumentRepo interface {
	Create(ctx context.Context, approvalDocument *ApprovalDocument) (int64, error)
	GetByApprovalID(ctx context.Context, documentID int64) ([]ApprovalDocument, error)
}

// ApprovalDocumentRepoImpl is the implementation of ApprovalDocumentRepo
type ApprovalDocumentRepoImpl struct {
	dig.In
	*sql.DB
}

func NewApprovalDocumentRepo(impl ApprovalDocumentRepoImpl) ApprovalDocumentRepo {
	return &impl
}

var (
	ApprovalDocumentTableName = "approval_documents"
	ApprovalDocumentTable     = struct {
		ID             string
		LoanApprovalID string
		DocumentType   string
		FileURL        string
		Description    string
		CreatedAt      string
		UpdatedAt      string
		DeletedAt      string
	}{
		ID:             "id",
		LoanApprovalID: "loan_approval_id",
		DocumentType:   "document_type",
		FileURL:        "file_url",
		Description:    "description",
		CreatedAt:      "created_at",
		UpdatedAt:      "updated_at",
		DeletedAt:      "deleted_at",
	}
)

func (r *ApprovalDocumentRepoImpl) Create(ctx context.Context, approvalDocument *ApprovalDocument) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Build insert query
	builder := sq.
		Insert(ApprovalDocumentTableName).
		Columns(
			ApprovalDocumentTable.LoanApprovalID,
			ApprovalDocumentTable.DocumentType,
			ApprovalDocumentTable.FileURL,
			ApprovalDocumentTable.Description,
			ApprovalDocumentTable.CreatedAt,
			ApprovalDocumentTable.UpdatedAt,
			ApprovalDocumentTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			approvalDocument.LoanApprovalID,
			approvalDocument.DocumentType,
			approvalDocument.FileURL,
			approvalDocument.Description,
			approvalDocument.CreatedAt,
			approvalDocument.UpdatedAt,
			nil,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan id: %v", err)
	}

	return id, nil
}

func (r *ApprovalDocumentRepoImpl) GetByApprovalID(ctx context.Context, approvalID int64) ([]ApprovalDocument, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build the query to get approval documents by approvalID
	builder := sq.
		Select(
			ApprovalDocumentTable.ID,
			ApprovalDocumentTable.LoanApprovalID,
			ApprovalDocumentTable.DocumentType,
			ApprovalDocumentTable.FileURL,
			ApprovalDocumentTable.Description,
			ApprovalDocumentTable.CreatedAt,
			ApprovalDocumentTable.UpdatedAt,
			ApprovalDocumentTable.DeletedAt,
		).
		From(ApprovalDocumentTableName).
		Where(sq.Eq{ApprovalDocumentTable.LoanApprovalID: approvalID}).
		PlaceholderFormat(sq.Dollar)

	// Query the database
	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	// Slice to hold the results
	var approvalDocuments []ApprovalDocument

	// Iterate over the rows and scan each one into an ApprovalDocument
	for rows.Next() {
		var approvalDocument ApprovalDocument
		if err := rows.Scan(
			&approvalDocument.ID,
			&approvalDocument.LoanApprovalID,
			&approvalDocument.DocumentType,
			&approvalDocument.FileURL,
			&approvalDocument.Description,
			&approvalDocument.CreatedAt,
			&approvalDocument.UpdatedAt,
			&approvalDocument.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan approval document: %v", err)
		}
		approvalDocuments = append(approvalDocuments, approvalDocument)
	}

	// Check for errors encountered while iterating
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return approvalDocuments, nil
}
