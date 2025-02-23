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

// LoanApproval represents the structure of a loan approval
type (
	LoanApprovalRequest struct {
		Offset uint64
		Size   uint64
		Status enum.LoanStatus
	}
	LoanApproval struct {
		ID             int64               `db:"id"`
		LoanID         int64               `db:"loan_id"`
		ApprovalNumber string              `db:"approval_number"`
		StaffID        *int64              `db:"staff_id"`
		ApprovalDate   *time.Time          `db:"approval_date"`
		ApprovalStatus enum.ApprovalStatus `db:"approval_status"`
		CreatedAt      time.Time           `db:"created_at"`
		UpdatedAt      time.Time           `db:"updated_at"`
		DeletedAt      *time.Time          `db:"deleted_at"`
	}
)

type LoanApprovalRepo interface {
	Create(ctx context.Context, loanApproval *LoanApproval) (int64, error)
	Update(ctx context.Context, loanApproval *LoanApproval) error
	GetByID(ctx context.Context, approvalID int64) (*LoanApproval, error)
	GetAll(ctx context.Context) ([]LoanApproval, error)
	GetAllPage(ctx context.Context, loanRequest LoanApprovalRequest) ([]LoanApproval, int64, error)
}

// LoanApprovalRepoImpl is the implementation of LoanApprovalRepo
type LoanApprovalRepoImpl struct {
	dig.In
	*sql.DB
}

func NewLoanApprovalRepo(impl LoanApprovalRepoImpl) LoanApprovalRepo {
	return &impl
}

var (
	LoanApprovalTableName = "loans_approval"
	LoanApprovalTable     = struct {
		ID             string
		LoanID         string
		ApprovalNumber string
		StaffID        string
		ApprovalDate   string
		ApprovalStatus string
		CreatedAt      string
		UpdatedAt      string
		DeletedAt      string
	}{
		ID:             "id",
		LoanID:         "loan_id",
		ApprovalNumber: "approval_number",
		StaffID:        "staff_id",
		ApprovalDate:   "approval_date",
		ApprovalStatus: "approval_status",
		CreatedAt:      "created_at",
		UpdatedAt:      "updated_at",
		DeletedAt:      "deleted_at",
	}
)

func (r *LoanApprovalRepoImpl) Create(ctx context.Context, loanApproval *LoanApproval) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Build insert query
	builder := sq.
		Insert(LoanApprovalTableName).
		Columns(
			LoanApprovalTable.LoanID,
			LoanApprovalTable.ApprovalNumber,
			LoanApprovalTable.StaffID,
			LoanApprovalTable.ApprovalDate,
			LoanApprovalTable.ApprovalStatus,
			LoanApprovalTable.CreatedAt,
			LoanApprovalTable.UpdatedAt,
			LoanApprovalTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			loanApproval.LoanID,
			loanApproval.ApprovalNumber,
			loanApproval.StaffID,
			loanApproval.ApprovalDate,
			loanApproval.ApprovalStatus,
			loanApproval.CreatedAt,
			loanApproval.UpdatedAt,
			nil,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan id: %v", err)
	}

	return id, nil
}

func (r *LoanApprovalRepoImpl) GetByID(ctx context.Context, approvalID int64) (*LoanApproval, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build the query to get loan approval by ID
	builder := sq.
		Select(
			LoanApprovalTable.ID,
			LoanApprovalTable.LoanID,
			LoanApprovalTable.ApprovalNumber,
			LoanApprovalTable.StaffID,
			LoanApprovalTable.ApprovalDate,
			LoanApprovalTable.ApprovalStatus,
			LoanApprovalTable.CreatedAt,
			LoanApprovalTable.UpdatedAt,
			LoanApprovalTable.DeletedAt,
		).
		From(LoanApprovalTableName).
		Where(sq.Eq{LoanApprovalTable.ID: approvalID}).
		PlaceholderFormat(sq.Dollar)

	var loanApproval LoanApproval
	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err := scanner.Scan(
		&loanApproval.ID,
		&loanApproval.LoanID,
		&loanApproval.ApprovalNumber,
		&loanApproval.StaffID,
		&loanApproval.ApprovalDate,
		&loanApproval.ApprovalStatus,
		&loanApproval.CreatedAt,
		&loanApproval.UpdatedAt,
		&loanApproval.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan loan approval: %v", err)
	}

	return &loanApproval, nil
}

func (r *LoanApprovalRepoImpl) GetAll(ctx context.Context) ([]LoanApproval, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build the query to get all loan approvals
	builder := sq.
		Select(
			LoanApprovalTable.ID,
			LoanApprovalTable.LoanID,
			LoanApprovalTable.ApprovalNumber,
			LoanApprovalTable.StaffID,
			LoanApprovalTable.ApprovalDate,
			LoanApprovalTable.ApprovalStatus,
			LoanApprovalTable.CreatedAt,
			LoanApprovalTable.UpdatedAt,
			LoanApprovalTable.DeletedAt,
		).
		From(LoanApprovalTableName).
		PlaceholderFormat(sq.Dollar)

	// Execute query and map result to loanApproval
	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query loan approvals: %v", err)
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var loanApprovals []LoanApproval
	for rows.Next() {
		var loanApproval LoanApproval
		if err := rows.Scan(
			&loanApproval.ID,
			&loanApproval.LoanID,
			&loanApproval.ApprovalNumber,
			&loanApproval.StaffID,
			&loanApproval.ApprovalDate,
			&loanApproval.ApprovalStatus,
			&loanApproval.CreatedAt,
			&loanApproval.UpdatedAt,
			&loanApproval.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan loan approval row: %v", err)
		}
		loanApprovals = append(loanApprovals, loanApproval)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %v", err)
	}

	return loanApprovals, nil
}

func (r *LoanApprovalRepoImpl) GetAllPage(ctx context.Context, approvalRequest LoanApprovalRequest) ([]LoanApproval, int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, 0, err
	}

	// Build the query with pagination (limit and offset)
	builder := sq.
		Select(
			LoanApprovalTable.ID,
			LoanApprovalTable.LoanID,
			LoanApprovalTable.ApprovalNumber,
			LoanApprovalTable.StaffID,
			LoanApprovalTable.ApprovalDate,
			LoanApprovalTable.ApprovalStatus,
			LoanApprovalTable.CreatedAt,
			LoanApprovalTable.UpdatedAt,
			LoanApprovalTable.DeletedAt,
		).
		From(LoanApprovalTableName)

	if approvalRequest.Status != "" {
		builder = builder.Where(sq.Eq{LoanApprovalTable.ApprovalStatus: approvalRequest.Status})
	}

	builder = builder.Limit(approvalRequest.Size).
		Offset(approvalRequest.Offset).
		PlaceholderFormat(sq.Dollar)

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var loanApprovals []LoanApproval
	for rows.Next() {
		var loanApproval LoanApproval
		if err := rows.Scan(
			&loanApproval.ID,
			&loanApproval.LoanID,
			&loanApproval.ApprovalNumber,
			&loanApproval.StaffID,
			&loanApproval.ApprovalDate,
			&loanApproval.ApprovalStatus,
			&loanApproval.CreatedAt,
			&loanApproval.UpdatedAt,
			&loanApproval.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("failed to scan loan approval: %v", err)
		}
		loanApprovals = append(loanApprovals, loanApproval)
	}

	// Get the total record count (for pagination)
	countQuery := sq.Select("COUNT(*)").
		From(LoanApprovalTableName).
		PlaceholderFormat(sq.Dollar)

	var totalRecords int64
	countScanner := countQuery.RunWith(txn).QueryRowContext(ctx)
	if err := countScanner.Scan(&totalRecords); err != nil {
		return nil, 0, fmt.Errorf("failed to get total records: %v", err)
	}

	return loanApprovals, totalRecords, nil
}

func (r *LoanApprovalRepoImpl) Update(ctx context.Context, loanApproval *LoanApproval) error {
	// Memulai transaksi
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return err
	}

	// Membangun query untuk melakukan update
	builder := sq.
		Update(LoanApprovalTableName).
		Set(LoanApprovalTable.ApprovalStatus, loanApproval.ApprovalStatus).
		Set(LoanApprovalTable.ApprovalDate, loanApproval.ApprovalDate).
		Set(LoanApprovalTable.StaffID, loanApproval.StaffID).
		Where(sq.Eq{LoanApprovalTable.ID: loanApproval.ID}).
		PlaceholderFormat(sq.Dollar)

	// Menjalankan query update
	_, err = builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to update loan approval: %v", err)
	}

	return nil
}
