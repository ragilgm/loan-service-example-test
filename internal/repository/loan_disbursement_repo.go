package repo

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/test/loan-service/internal/enum"
	"github.com/typical-go/typical-rest-server/pkg/dbtxn"
	"go.uber.org/dig"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=$PROJ/internal/generated/mock/mock_$GOPACKAGE/$GOFILE

type (
	LoanDisbursementRequest struct {
		Offset uint64
		Size   uint64
		Status string // Disbursement Status (Pending, Completed, etc.)
	}

	LoanDisbursement struct {
		ID                 int64                       `db:"id"`                   // Disbursement ID
		LoanID             int64                       `db:"loan_id"`              // Loan ID
		DisburseCode       string                      `db:"disburse_code"`        // Disbursement code
		DisburseAmount     float64                     `db:"disburse_amount"`      // Disbursed amount
		DisbursementStatus enum.LoanDisbursementStatus `db:"disbursement_status"`  // Status (Pending, Completed, etc.)
		DisburseDate       *time.Time                  `db:"disburse_date"`        // Disbursement date
		StaffID            *int64                      `db:"staff_id"`             // Staff ID handling the disbursement
		AgreementURL       string                      `db:"agreement_url"`        // URL template agreement url
		SignedAgreementURL *string                     `db:"signed_agreement_url"` // URL to signed agreement
		CreatedAt          time.Time                   `db:"created_at"`           // Date of creation
		UpdatedAt          time.Time                   `db:"updated_at"`           // Date of last update
		DeletedAt          *time.Time                  `db:"deleted_at"`           // Date of deletion if applicable
	}

	LoanDisbursementRepo interface {
		Create(context.Context, *LoanDisbursement) (int64, error)
		Update(ctx context.Context, disbursement *LoanDisbursement) error
		GetByID(ctx context.Context, disbursementID int64) (*LoanDisbursement, error)
		GetAll(ctx context.Context) ([]LoanDisbursement, error)
		GetAllPage(ctx context.Context, request LoanDisbursementRequest) ([]LoanDisbursement, int64, error)
	}

	LoanDisbursementRepoImpl struct {
		dig.In
		*sql.DB
	}
)

var (
	LoanDisbursementTableName = "loans_disbursement"
	LoanDisbursementTable     = struct {
		ID                 string
		LoanID             string
		DisburseCode       string
		DisburseAmount     string
		DisbursementStatus string
		DisburseDate       string
		StaffID            string
		AgreementURL       string
		SignedAgreementURL string
		CreatedAt          string
		UpdatedAt          string
		DeletedAt          string
	}{
		ID:                 "id",
		LoanID:             "loan_id",
		DisburseCode:       "disburse_code",
		DisburseAmount:     "disburse_amount",
		DisbursementStatus: "disbursement_status",
		DisburseDate:       "disburse_date",
		StaffID:            "staff_id",
		AgreementURL:       "agreement_url",
		SignedAgreementURL: "signed_agreement_url",
		CreatedAt:          "created_at",
		UpdatedAt:          "updated_at",
		DeletedAt:          "deleted_at",
	}
)

func NewLoanDisbursementRepo(impl LoanDisbursementRepoImpl) LoanDisbursementRepo {
	return &impl
}

// Create LoanDisbursement and return last inserted id
func (r *LoanDisbursementRepoImpl) Create(ctx context.Context, disbursement *LoanDisbursement) (int64, error) {
	// Use transaction if any
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Construct insert query
	builder := sq.
		Insert(LoanDisbursementTableName).
		Columns(
			LoanDisbursementTable.LoanID,
			LoanDisbursementTable.DisburseCode,
			LoanDisbursementTable.DisburseAmount,
			LoanDisbursementTable.DisbursementStatus,
			LoanDisbursementTable.DisburseDate,
			LoanDisbursementTable.StaffID,
			LoanDisbursementTable.AgreementURL,
			LoanDisbursementTable.SignedAgreementURL,
			LoanDisbursementTable.CreatedAt,
			LoanDisbursementTable.UpdatedAt,
			LoanDisbursementTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			disbursement.LoanID,
			disbursement.DisburseCode,
			disbursement.DisburseAmount,
			disbursement.DisbursementStatus,
			disbursement.DisburseDate,
			disbursement.StaffID,
			disbursement.AgreementURL,
			disbursement.SignedAgreementURL,
			time.Now(),
			time.Now(),
			nil,
		)

	// Execute the query and scan the ID
	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan id: %v", err)
	}

	return id, nil
}

// Update LoanDisbursement
func (r *LoanDisbursementRepoImpl) Update(ctx context.Context, disbursement *LoanDisbursement) error {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return err
	}

	// Construct update query
	builder := sq.Update(LoanDisbursementTableName).
		Set(LoanDisbursementTable.DisburseAmount, disbursement.DisburseAmount).
		Set(LoanDisbursementTable.DisbursementStatus, disbursement.DisbursementStatus).
		Set(LoanDisbursementTable.DisburseDate, disbursement.DisburseDate).
		Set(LoanDisbursementTable.StaffID, disbursement.StaffID).
		Set(LoanDisbursementTable.AgreementURL, disbursement.AgreementURL).
		Set(LoanDisbursementTable.SignedAgreementURL, disbursement.SignedAgreementURL).
		Set(LoanDisbursementTable.UpdatedAt, time.Now()).
		Where(sq.Eq{LoanDisbursementTable.ID: disbursement.ID}).
		PlaceholderFormat(sq.Dollar)

	// Execute the update
	result, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to update disbursement: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no disbursement found with ID: %d", disbursement.ID)
	}

	return nil
}

// Get LoanDisbursement by ID
func (r *LoanDisbursementRepoImpl) GetByID(ctx context.Context, disbursementID int64) (*LoanDisbursement, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Construct select query to get LoanDisbursement by ID
	builder := sq.Select(
		LoanDisbursementTable.ID,
		LoanDisbursementTable.LoanID,
		LoanDisbursementTable.DisburseCode,
		LoanDisbursementTable.DisburseAmount,
		LoanDisbursementTable.DisbursementStatus,
		LoanDisbursementTable.DisburseDate,
		LoanDisbursementTable.StaffID,
		LoanDisbursementTable.AgreementURL,
		LoanDisbursementTable.SignedAgreementURL,
		LoanDisbursementTable.CreatedAt,
		LoanDisbursementTable.UpdatedAt,
		LoanDisbursementTable.DeletedAt,
	).
		From(LoanDisbursementTableName).
		Where(sq.Eq{LoanDisbursementTable.ID: disbursementID}).
		PlaceholderFormat(sq.Dollar)

	// Execute the query and scan the result
	var disbursement LoanDisbursement
	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err := scanner.Scan(
		&disbursement.ID,
		&disbursement.LoanID,
		&disbursement.DisburseCode,
		&disbursement.DisburseAmount,
		&disbursement.DisbursementStatus,
		&disbursement.DisburseDate,
		&disbursement.StaffID,
		&disbursement.AgreementURL,
		&disbursement.SignedAgreementURL,
		&disbursement.CreatedAt,
		&disbursement.UpdatedAt,
		&disbursement.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan disbursement: %v", err)
	}

	return &disbursement, nil
}

// Get all LoanDisbursements
func (r *LoanDisbursementRepoImpl) GetAll(ctx context.Context) ([]LoanDisbursement, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Construct select query to get all LoanDisbursements
	builder := sq.Select(
		LoanDisbursementTable.ID,
		LoanDisbursementTable.LoanID,
		LoanDisbursementTable.DisburseCode,
		LoanDisbursementTable.DisburseAmount,
		LoanDisbursementTable.DisbursementStatus,
		LoanDisbursementTable.DisburseDate,
		LoanDisbursementTable.StaffID,
		LoanDisbursementTable.AgreementURL,
		LoanDisbursementTable.SignedAgreementURL,
		LoanDisbursementTable.CreatedAt,
		LoanDisbursementTable.UpdatedAt,
		LoanDisbursementTable.DeletedAt,
	).
		From(LoanDisbursementTableName).
		PlaceholderFormat(sq.Dollar)

	// Execute query and scan results
	var disbursements []LoanDisbursement
	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var disbursement LoanDisbursement
		if err := rows.Scan(
			&disbursement.ID,
			&disbursement.LoanID,
			&disbursement.DisburseCode,
			&disbursement.DisburseAmount,
			&disbursement.DisbursementStatus,
			&disbursement.DisburseDate,
			&disbursement.StaffID,
			&disbursement.AgreementURL,
			&disbursement.SignedAgreementURL,
			&disbursement.CreatedAt,
			&disbursement.UpdatedAt,
			&disbursement.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		disbursements = append(disbursements, disbursement)
	}

	return disbursements, nil
}

// Get paginated LoanDisbursements
func (r *LoanDisbursementRepoImpl) GetAllPage(ctx context.Context, request LoanDisbursementRequest) ([]LoanDisbursement, int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, 0, err
	}

	// Query to count the total number of records
	countBuilder := sq.Select("COUNT(*)").
		From(LoanDisbursementTableName).
		Where(sq.Eq{LoanDisbursementTable.DeletedAt: nil}).
		PlaceholderFormat(sq.Dollar)

	// Execute the count query
	var totalRecords int64
	scanner := countBuilder.RunWith(txn).QueryRowContext(ctx)
	if err := scanner.Scan(&totalRecords); err != nil {
		return nil, 0, fmt.Errorf("failed to count records: %v", err)
	}

	// Query to get paginated data
	pageBuilder := sq.Select(
		LoanDisbursementTable.ID,
		LoanDisbursementTable.LoanID,
		LoanDisbursementTable.DisburseCode,
		LoanDisbursementTable.DisburseAmount,
		LoanDisbursementTable.DisbursementStatus,
		LoanDisbursementTable.DisburseDate,
		LoanDisbursementTable.StaffID,
		LoanDisbursementTable.AgreementURL,
		LoanDisbursementTable.SignedAgreementURL,
		LoanDisbursementTable.CreatedAt,
		LoanDisbursementTable.UpdatedAt,
		LoanDisbursementTable.DeletedAt,
	).
		From(LoanDisbursementTableName).
		Limit(request.Size).
		Offset(request.Offset).
		PlaceholderFormat(sq.Dollar)

	if request.Status != "" {
		pageBuilder = pageBuilder.Where(sq.Eq{LoanDisbursementTable.DisbursementStatus: request.Status})

	}

	// Execute the query and fetch the paginated data
	var disbursements []LoanDisbursement
	rows, err := pageBuilder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var disbursement LoanDisbursement
		if err := rows.Scan(
			&disbursement.ID,
			&disbursement.LoanID,
			&disbursement.DisburseCode,
			&disbursement.DisburseAmount,
			&disbursement.DisbursementStatus,
			&disbursement.DisburseDate,
			&disbursement.StaffID,
			&disbursement.AgreementURL,
			&disbursement.SignedAgreementURL,
			&disbursement.CreatedAt,
			&disbursement.UpdatedAt,
			&disbursement.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %v", err)
		}
		disbursements = append(disbursements, disbursement)
	}

	return disbursements, totalRecords, nil
}
