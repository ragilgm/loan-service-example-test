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
	LoanRequest struct {
		Offset uint64
		Size   uint64
		Status enum.LoanStatus
	}
	Loan struct {
		ID                   int64           `db:"id"`                     // Loan ID
		LoanCode             string          `db:"loan_code"`              // Loan code
		BorrowerID           int64           `db:"borrower_id"`            // Borrower ID
		RequestAmount        float64         `db:"request_amount"`         // Loan request amount
		LoanGrade            string          `db:"loan_grade"`             // Loan grade (A, B, C, D)
		LoanType             enum.LoanType   `db:"loan_type"`              // Type of loan (productive, consumptive, etc.)
		TotalInvestedAmount  float64         `db:"total_invested_amount"`  // Total amount invested
		InvestorCount        int64           `db:"investor_count"`         // Number of investors participating
		FundingDeadline      *time.Time      `db:"funding_deadline"`       // Funding deadline
		LoanStatus           enum.LoanStatus `db:"loan_status"`            // Loan status (proposed, rejected, approved, invested)
		Rate                 float64         `db:"rate"`                   // Interest rate
		Tenures              int64           `db:"tenures"`                // Loan tenure
		TotalInterest        float64         `db:"total_interest"`         // Total repayment amount needed
		TotalRepaymentAmount float64         `db:"total_repayment_amount"` // Total repayment amount needed
		InvestmentPercentage float64         `db:"investment_percentage"`  // Investor profit sharing percentage
		CreatedAt            time.Time       `db:"created_at"`             // Loan creation date
		UpdatedAt            time.Time       `db:"updated_at"`             // Loan status update date
		DeletedAt            *time.Time      `db:"deleted_at"`             // Loan deletion date (if applicable)
	}

	LoanRepo interface {
		Create(context.Context, *Loan) (int64, error)
		Update(ctx context.Context, loan *Loan) error
		GetByID(ctx context.Context, loanID int64) (*Loan, error)
		GetAll(ctx context.Context) ([]Loan, error)
		GetAllPage(ctx context.Context, loanRequest LoanRequest) ([]Loan, int64, error)
	}

	LoanRepoImpl struct {
		dig.In
		*sql.DB
	}
)

var (
	LoanTableName = "loans"
	LoanTable     = struct {
		ID                   string
		LoanCode             string
		BorrowerID           string
		RequestAmount        string
		LoanGrade            string
		LoanType             string
		TotalInvestedAmount  string
		InvestorCount        string
		FundingDeadline      string
		LoanStatus           string
		Rate                 string
		Tenures              string
		TotalInterest        string
		TotalRepaymentAmount string
		InvestmentPercentage string
		CreatedAt            string
		UpdatedAt            string
		DeletedAt            string
	}{
		ID:                   "id",
		LoanCode:             "loan_code",
		BorrowerID:           "borrower_id",
		RequestAmount:        "request_amount",
		LoanGrade:            "loan_grade",
		LoanType:             "loan_type",
		TotalInvestedAmount:  "total_invested_amount",
		InvestorCount:        "investor_count",
		FundingDeadline:      "funding_deadline",
		LoanStatus:           "loan_status",
		Rate:                 "rate",
		Tenures:              "tenures",
		TotalInterest:        "total_interest",
		TotalRepaymentAmount: "total_repayment_amount",
		InvestmentPercentage: "investment_percentage",
		CreatedAt:            "created_at",
		UpdatedAt:            "updated_at",
		DeletedAt:            "deleted_at",
	}
)

func NewLoanRepo(impl LoanRepoImpl) LoanRepo {
	return &impl
}

// Create Loan and return last inserted id
func (r *LoanRepoImpl) Create(ctx context.Context, loan *Loan) (int64, error) {
	// // use transaction if any
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Menyusun query untuk insert
	builder := sq.
		Insert(LoanTableName).
		Columns(
			LoanTable.LoanCode,
			LoanTable.BorrowerID,
			LoanTable.RequestAmount,
			LoanTable.LoanGrade,
			LoanTable.LoanType,
			LoanTable.TotalInvestedAmount,
			LoanTable.InvestorCount,
			LoanTable.FundingDeadline,
			LoanTable.LoanStatus,
			LoanTable.Rate,
			LoanTable.Tenures,
			LoanTable.TotalInterest,
			LoanTable.TotalRepaymentAmount,
			LoanTable.InvestmentPercentage,
			LoanTable.CreatedAt,
			LoanTable.UpdatedAt,
			LoanTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			loan.LoanCode,
			loan.BorrowerID,
			loan.RequestAmount,
			loan.LoanGrade,
			loan.LoanType,
			loan.TotalInvestedAmount,
			loan.InvestorCount,
			nil,
			loan.LoanStatus,
			loan.Rate,
			loan.Tenures,
			loan.TotalInterest,
			loan.TotalRepaymentAmount,
			loan.InvestmentPercentage,
			time.Now(),
			time.Now(),
			nil,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		// dbtxn akan secara otomatis melakukan rollback jika ada error
		return -1, fmt.Errorf("failed to scan id: %v", err)
	}

	// Commit transaksi setelah berhasil
	// dbtxn akan secara otomatis melakukan commit atau rollback jika terjadi error
	return id, nil
}
func (r *LoanRepoImpl) Update(ctx context.Context, loan *Loan) error {
	// Use transaction if needed
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return err
	}

	// Construct the query to update the loan by ID
	builder := sq.Update(LoanTableName).
		Set(LoanTable.TotalInvestedAmount, loan.TotalInvestedAmount).
		Set(LoanTable.InvestorCount, loan.InvestorCount).
		Set(LoanTable.FundingDeadline, loan.FundingDeadline).
		Set(LoanTable.LoanStatus, loan.LoanStatus).
		Set(LoanTable.TotalInterest, loan.TotalInterest).
		Set(LoanTable.TotalRepaymentAmount, loan.TotalRepaymentAmount).
		Set(LoanTable.InvestmentPercentage, loan.InvestmentPercentage).
		Set(LoanTable.UpdatedAt, time.Now()). // Update the `UpdatedAt` field
		Where(sq.Eq{LoanTable.ID: loan.ID}).
		PlaceholderFormat(sq.Dollar)

	// Execute the update query
	result, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to update loan: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no loan found with ID: %d", loan.ID)
	}

	return nil
}

func (r *LoanRepoImpl) GetAllPage(ctx context.Context, loanRequest LoanRequest) ([]Loan, int64, error) {
	// Use transaction if needed
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, 0, err
	}

	// Query untuk mengambil data pinjaman dengan limit dan offset
	builder := sq.
		Select(
			LoanTable.ID,
			LoanTable.LoanCode,
			LoanTable.BorrowerID,
			LoanTable.RequestAmount,
			LoanTable.LoanGrade,
			LoanTable.LoanType,
			LoanTable.TotalInvestedAmount,
			LoanTable.InvestorCount,
			LoanTable.FundingDeadline,
			LoanTable.LoanStatus,
			LoanTable.Rate,
			LoanTable.Tenures,
			LoanTable.TotalInterest,
			LoanTable.TotalRepaymentAmount,
			LoanTable.InvestmentPercentage,
			LoanTable.CreatedAt,
			LoanTable.UpdatedAt,
			LoanTable.DeletedAt,
		).
		From(LoanTableName).
		Limit(loanRequest.Size).
		Offset(loanRequest.Offset).
		PlaceholderFormat(sq.Dollar)

	if loanRequest.Status != "" {
		builder = builder.Where(sq.Eq{LoanTable.LoanStatus: loanRequest.Status})

	}

	// Execute the query and scan the result
	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var loans []Loan
	for rows.Next() {
		var loan Loan
		if err := rows.Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.BorrowerID,
			&loan.RequestAmount,
			&loan.LoanGrade,
			&loan.LoanType,
			&loan.TotalInvestedAmount,
			&loan.InvestorCount,
			&loan.FundingDeadline,
			&loan.LoanStatus,
			&loan.Rate,
			&loan.Tenures,
			&loan.TotalInterest,
			&loan.TotalRepaymentAmount,
			&loan.InvestmentPercentage,
			&loan.CreatedAt,
			&loan.UpdatedAt,
			&loan.DeletedAt,
		); err != nil {
			return nil, 0, err
		}
		loans = append(loans, loan)
	}

	// Get the total record count (for pagination metadata)
	countQuery := sq.Select("COUNT(*)").
		From(LoanTableName).
		PlaceholderFormat(sq.Dollar)

	var totalRecords int64
	countScanner := countQuery.RunWith(txn).QueryRowContext(ctx)
	if err := countScanner.Scan(&totalRecords); err != nil {
		return nil, 0, err
	}

	return loans, totalRecords, nil
}

func (r *LoanRepoImpl) GetByID(ctx context.Context, loanID int64) (*Loan, error) {
	// use transaction if any
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Construct the query to get the loan by ID
	builder := sq.
		Select(
			LoanTable.ID,
			LoanTable.LoanCode,
			LoanTable.BorrowerID,
			LoanTable.RequestAmount,
			LoanTable.LoanGrade,
			LoanTable.LoanType,
			LoanTable.TotalInvestedAmount,
			LoanTable.InvestorCount,
			LoanTable.FundingDeadline,
			LoanTable.LoanStatus,
			LoanTable.Rate,
			LoanTable.Tenures,
			LoanTable.TotalInterest,
			LoanTable.TotalRepaymentAmount,
			LoanTable.InvestmentPercentage,
			LoanTable.CreatedAt,
			LoanTable.UpdatedAt,
			LoanTable.DeletedAt,
		).
		From(LoanTableName).
		Where(sq.Eq{LoanTable.ID: loanID}).
		PlaceholderFormat(sq.Dollar)

	// Execute the query and scan the result into the Loan struct
	var loan Loan
	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	// Scan the values into the loan struct
	if err := scanner.Scan(
		&loan.ID,
		&loan.LoanCode,
		&loan.BorrowerID,
		&loan.RequestAmount,
		&loan.LoanGrade,
		&loan.LoanType,
		&loan.TotalInvestedAmount,
		&loan.InvestorCount,
		&loan.FundingDeadline,
		&loan.LoanStatus,
		&loan.Rate,
		&loan.Tenures,
		&loan.TotalInterest,
		&loan.TotalRepaymentAmount,
		&loan.InvestmentPercentage,
		&loan.CreatedAt,
		&loan.UpdatedAt,
		&loan.DeletedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan loan: %v", err)
	}

	// Return the loan details
	return &loan, nil
}

func (r *LoanRepoImpl) GetAll(ctx context.Context) ([]Loan, error) {
	// use transaction if any
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Construct the query to get all loans
	builder := sq.
		Select(
			LoanTable.ID,
			LoanTable.LoanCode,
			LoanTable.BorrowerID,
			LoanTable.RequestAmount,
			LoanTable.LoanGrade,
			LoanTable.LoanType,
			LoanTable.TotalInvestedAmount,
			LoanTable.InvestorCount,
			LoanTable.FundingDeadline,
			LoanTable.LoanStatus,
			LoanTable.Rate,
			LoanTable.Tenures,
			LoanTable.TotalInterest,
			LoanTable.TotalRepaymentAmount,
			LoanTable.InvestmentPercentage,
			LoanTable.CreatedAt,
			LoanTable.UpdatedAt,
			LoanTable.DeletedAt,
		).
		From(LoanTableName).
		PlaceholderFormat(sq.Dollar)

	// Execute the query and scan the result into a slice of Loan structs
	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query loans: %v", err)
	}

	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	// Iterate over the rows and scan the result into the Loan struct
	var loans []Loan
	for rows.Next() {
		var loan Loan
		if err = rows.Scan(
			&loan.ID,
			&loan.LoanCode,
			&loan.BorrowerID,
			&loan.RequestAmount,
			&loan.LoanGrade,
			&loan.LoanType,
			&loan.TotalInvestedAmount,
			&loan.InvestorCount,
			&loan.FundingDeadline,
			&loan.LoanStatus,
			&loan.Rate,
			&loan.Tenures,
			&loan.TotalInterest,
			&loan.TotalRepaymentAmount,
			&loan.InvestmentPercentage,
			&loan.CreatedAt,
			&loan.UpdatedAt,
			&loan.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan loan row: %v", err)
		}
		loans = append(loans, loan)
	}

	// Check for errors that occurred during iteration
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during row iteration: %v", err)
	}

	// Return the slice of loans
	return loans, nil
}
