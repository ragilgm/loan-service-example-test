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

// LoanFunding represents the structure of the loan funding records
type (
	LoanFunding struct {
		ID                 int64                  `db:"id"`
		LoanOrderNumber    string                 `db:"loan_order_number"`
		OrderNumber        string                 `db:"order_number"`
		LoanID             int64                  `db:"loan_id"`
		LenderID           int64                  `db:"lender_id"`
		LenderEmail        string                 `db:"lender_email"`
		InvestmentAmount   float64                `db:"investment_amount"`
		Rate               float64                `db:"rate"`
		Interest           float64                `db:"interest"`
		ROI                float64                `db:"roi"`
		InterestPaid       float64                `db:"interest_paid"`
		CapitalAmountPaid  float64                `db:"capital_amount_paid"`
		TotalAmountPaid    float64                `db:"total_amount_paid"`
		InvestmentDate     time.Time              `db:"investment_date"`
		Status             enum.LoanFundingStatus `db:"status"`
		LenderAgreementURL string                 `db:"lender_agreement_url"`
		CreatedAt          time.Time              `db:"created_at"`
		UpdatedAt          time.Time              `db:"updated_at"`
		DeletedAt          *time.Time             `db:"deleted_at"`
	}
)

type LoanFundingRepo interface {
	Create(ctx context.Context, loanFunding *LoanFunding) (int64, error)
	Update(ctx context.Context, loanFunding *LoanFunding) error
	GetByID(ctx context.Context, id int64) (*LoanFunding, error)
	GetByLoanID(ctx context.Context, loanID int64) ([]LoanFunding, error)
	GetByLoanOrderNumber(ctx context.Context, loanOrderNumber string) (*LoanFunding, error)
	GetByLenderID(ctx context.Context, lenderID int64) ([]LoanFunding, error)
}

// LoanFundingRepoImpl is the implementation of LoanFundingRepo
type LoanFundingRepoImpl struct {
	dig.In
	*sql.DB
}

func NewLoanFundingRepo(impl LoanFundingRepoImpl) LoanFundingRepo {
	return &impl
}

var (
	LoanFundingTableName = "loan_funding"
	LoanFundingTable     = struct {
		ID                 string
		LoanOrderNumber    string
		OrderNumber        string
		LoanID             string
		LenderID           string
		LenderEmail        string
		InvestmentAmount   string
		Rate               string
		Interest           string
		ROI                string
		InterestPaid       string
		CapitalAmountPaid  string
		TotalAmountPaid    string
		InvestmentDate     string
		Status             string
		LenderAgreementURL string
		CreatedAt          string
		UpdatedAt          string
		DeletedAt          string
	}{
		ID:                 "id",
		LoanOrderNumber:    "loan_order_number",
		OrderNumber:        "order_number",
		LoanID:             "loan_id",
		LenderID:           "lender_id",
		LenderEmail:        "lender_email",
		InvestmentAmount:   "investment_amount",
		Rate:               "rate",
		Interest:           "interest",
		ROI:                "roi",
		InterestPaid:       "interest_paid",
		CapitalAmountPaid:  "capital_amount_paid",
		TotalAmountPaid:    "total_amount_paid",
		InvestmentDate:     "investment_date",
		Status:             "status",
		LenderAgreementURL: "lender_agreement_url",
		CreatedAt:          "created_at",
		UpdatedAt:          "updated_at",
		DeletedAt:          "deleted_at",
	}
)

func (r *LoanFundingRepoImpl) Create(ctx context.Context, loanFunding *LoanFunding) (int64, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Build insert query
	builder := sq.
		Insert(LoanFundingTableName).
		Columns(
			LoanFundingTable.LoanOrderNumber,
			LoanFundingTable.OrderNumber,
			LoanFundingTable.LoanID,
			LoanFundingTable.LenderID,
			LoanFundingTable.LenderEmail,
			LoanFundingTable.InvestmentAmount,
			LoanFundingTable.Rate,
			LoanFundingTable.Interest,
			LoanFundingTable.ROI,
			LoanFundingTable.InterestPaid,
			LoanFundingTable.CapitalAmountPaid,
			LoanFundingTable.TotalAmountPaid,
			LoanFundingTable.InvestmentDate,
			LoanFundingTable.Status,
			LoanFundingTable.LenderAgreementURL,
			LoanFundingTable.CreatedAt,
			LoanFundingTable.UpdatedAt,
			LoanFundingTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			loanFunding.LoanOrderNumber,
			loanFunding.OrderNumber,
			loanFunding.LoanID,
			loanFunding.LenderID,
			loanFunding.LenderEmail,
			loanFunding.InvestmentAmount,
			loanFunding.Rate,
			loanFunding.Interest,
			loanFunding.ROI,
			loanFunding.InterestPaid,
			loanFunding.CapitalAmountPaid,
			loanFunding.TotalAmountPaid,
			loanFunding.InvestmentDate,
			loanFunding.Status,
			loanFunding.LenderAgreementURL,
			loanFunding.CreatedAt,
			loanFunding.UpdatedAt,
			nil, // Deleting record, set to nil by default
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	var id int64
	if err := scanner.Scan(&id); err != nil {
		return -1, fmt.Errorf("failed to scan id: %v", err)
	}

	return id, nil
}

func (r *LoanFundingRepoImpl) Update(ctx context.Context, loanFunding *LoanFunding) error {
	// Use transaction if needed
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return err
	}

	builder := sq.Update(LoanFundingTableName).
		Set(LoanFundingTable.InvestmentAmount, loanFunding.InvestmentAmount).
		Set(LoanFundingTable.Rate, loanFunding.Rate).
		Set(LoanFundingTable.Interest, loanFunding.Interest).
		Set(LoanFundingTable.ROI, loanFunding.ROI).
		Set(LoanFundingTable.InterestPaid, loanFunding.InterestPaid).
		Set(LoanFundingTable.CapitalAmountPaid, loanFunding.CapitalAmountPaid).
		Set(LoanFundingTable.TotalAmountPaid, loanFunding.TotalAmountPaid).
		Set(LoanFundingTable.InvestmentDate, loanFunding.InvestmentDate).
		Set(LoanFundingTable.Status, loanFunding.Status).
		Set(LoanFundingTable.UpdatedAt, time.Now()). // Update the `UpdatedAt` field
		Where(sq.Eq{LoanFundingTable.ID: loanFunding.ID}).
		PlaceholderFormat(sq.Dollar)

	// Execute the update query
	result, err := builder.RunWith(txn).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to update loan funding: %v", err)
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no loan funding found with ID: %d", loanFunding.ID)
	}

	return nil
}

func (r *LoanFundingRepoImpl) GetByLoanOrderNumber(ctx context.Context, loanOrderNumber string) (*LoanFunding, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build the query to get loan funding by ID
	builder := sq.
		Select(
			LoanFundingTable.ID,
			LoanFundingTable.LoanOrderNumber,
			LoanFundingTable.OrderNumber,
			LoanFundingTable.LoanID,
			LoanFundingTable.LenderID,
			LoanFundingTable.LenderEmail,
			LoanFundingTable.InvestmentAmount,
			LoanFundingTable.Rate,
			LoanFundingTable.Interest,
			LoanFundingTable.ROI,
			LoanFundingTable.InterestPaid,
			LoanFundingTable.CapitalAmountPaid,
			LoanFundingTable.TotalAmountPaid,
			LoanFundingTable.InvestmentDate,
			LoanFundingTable.Status,
			LoanFundingTable.LenderAgreementURL,
			LoanFundingTable.CreatedAt,
			LoanFundingTable.UpdatedAt,
			LoanFundingTable.DeletedAt,
		).
		From(LoanFundingTableName).
		Where(sq.Eq{LoanFundingTable.LoanOrderNumber: loanOrderNumber}).
		PlaceholderFormat(sq.Dollar)

	var loanFunding LoanFunding
	err = builder.RunWith(txn).QueryRowContext(ctx).Scan(
		&loanFunding.ID,
		&loanFunding.LoanOrderNumber,
		&loanFunding.OrderNumber,
		&loanFunding.LoanID,
		&loanFunding.LenderID,
		&loanFunding.LenderEmail,
		&loanFunding.InvestmentAmount,
		&loanFunding.Rate,
		&loanFunding.Interest,
		&loanFunding.ROI,
		&loanFunding.InterestPaid,
		&loanFunding.CapitalAmountPaid,
		&loanFunding.TotalAmountPaid,
		&loanFunding.InvestmentDate,
		&loanFunding.Status,
		&loanFunding.LenderAgreementURL,
		&loanFunding.CreatedAt,
		&loanFunding.UpdatedAt,
		&loanFunding.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &loanFunding, nil
}

func (r *LoanFundingRepoImpl) GetByID(ctx context.Context, id int64) (*LoanFunding, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build the query to get loan funding by ID
	builder := sq.
		Select(
			LoanFundingTable.ID,
			LoanFundingTable.LoanOrderNumber,
			LoanFundingTable.OrderNumber,
			LoanFundingTable.LoanID,
			LoanFundingTable.LenderID,
			LoanFundingTable.LenderEmail,
			LoanFundingTable.InvestmentAmount,
			LoanFundingTable.Rate,
			LoanFundingTable.Interest,
			LoanFundingTable.ROI,
			LoanFundingTable.InterestPaid,
			LoanFundingTable.CapitalAmountPaid,
			LoanFundingTable.TotalAmountPaid,
			LoanFundingTable.InvestmentDate,
			LoanFundingTable.Status,
			LoanFundingTable.LenderAgreementURL,
			LoanFundingTable.CreatedAt,
			LoanFundingTable.UpdatedAt,
			LoanFundingTable.DeletedAt,
		).
		From(LoanFundingTableName).
		Where(sq.Eq{LoanFundingTable.ID: id}).
		PlaceholderFormat(sq.Dollar)

	var loanFunding LoanFunding
	err = builder.RunWith(txn).QueryRowContext(ctx).Scan(
		&loanFunding.ID,
		&loanFunding.LoanOrderNumber,
		&loanFunding.OrderNumber,
		&loanFunding.LoanID,
		&loanFunding.LenderID,
		&loanFunding.LenderEmail,
		&loanFunding.InvestmentAmount,
		&loanFunding.Rate,
		&loanFunding.Interest,
		&loanFunding.ROI,
		&loanFunding.InterestPaid,
		&loanFunding.CapitalAmountPaid,
		&loanFunding.TotalAmountPaid,
		&loanFunding.InvestmentDate,
		&loanFunding.Status,
		&loanFunding.LenderAgreementURL,
		&loanFunding.CreatedAt,
		&loanFunding.UpdatedAt,
		&loanFunding.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	return &loanFunding, nil
}

func (r *LoanFundingRepoImpl) GetByLoanID(ctx context.Context, loanID int64) ([]LoanFunding, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build query to get loan funding by lenderID
	builder := sq.
		Select(
			LoanFundingTable.ID,
			LoanFundingTable.LoanOrderNumber,
			LoanFundingTable.OrderNumber,
			LoanFundingTable.LoanID,
			LoanFundingTable.LenderID,
			LoanFundingTable.LenderEmail,
			LoanFundingTable.InvestmentAmount,
			LoanFundingTable.Rate,
			LoanFundingTable.Interest,
			LoanFundingTable.ROI,
			LoanFundingTable.InterestPaid,
			LoanFundingTable.CapitalAmountPaid,
			LoanFundingTable.TotalAmountPaid,
			LoanFundingTable.InvestmentDate,
			LoanFundingTable.Status,
			LoanFundingTable.LenderAgreementURL,
			LoanFundingTable.CreatedAt,
			LoanFundingTable.UpdatedAt,
			LoanFundingTable.DeletedAt,
		).
		From(LoanFundingTableName).
		Where(sq.Eq{LoanFundingTable.LoanID: loanID}).
		PlaceholderFormat(sq.Dollar)

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var loanFundings []LoanFunding
	for rows.Next() {
		var loanFunding LoanFunding
		if err := rows.Scan(
			&loanFunding.ID,
			&loanFunding.LoanOrderNumber,
			&loanFunding.OrderNumber,
			&loanFunding.LoanID,
			&loanFunding.LenderID,
			&loanFunding.LenderEmail,
			&loanFunding.InvestmentAmount,
			&loanFunding.Rate,
			&loanFunding.Interest,
			&loanFunding.ROI,
			&loanFunding.InterestPaid,
			&loanFunding.CapitalAmountPaid,
			&loanFunding.TotalAmountPaid,
			&loanFunding.InvestmentDate,
			&loanFunding.Status,
			&loanFunding.LenderAgreementURL,
			&loanFunding.CreatedAt,
			&loanFunding.UpdatedAt,
			&loanFunding.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		loanFundings = append(loanFundings, loanFunding)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return loanFundings, nil
}

func (r *LoanFundingRepoImpl) GetByLenderID(ctx context.Context, lenderID int64) ([]LoanFunding, error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return nil, err
	}

	// Build query to get loan funding by lenderID
	builder := sq.
		Select(
			LoanFundingTable.ID,
			LoanFundingTable.LoanOrderNumber,
			LoanFundingTable.OrderNumber,
			LoanFundingTable.LoanID,
			LoanFundingTable.LenderID,
			LoanFundingTable.LenderEmail,
			LoanFundingTable.InvestmentAmount,
			LoanFundingTable.Rate,
			LoanFundingTable.Interest,
			LoanFundingTable.ROI,
			LoanFundingTable.InterestPaid,
			LoanFundingTable.CapitalAmountPaid,
			LoanFundingTable.TotalAmountPaid,
			LoanFundingTable.InvestmentDate,
			LoanFundingTable.Status,
			LoanFundingTable.LenderAgreementURL,
			LoanFundingTable.CreatedAt,
			LoanFundingTable.UpdatedAt,
			LoanFundingTable.DeletedAt,
		).
		From(LoanFundingTableName).
		Where(sq.Eq{LoanFundingTable.LenderID: lenderID}).
		PlaceholderFormat(sq.Dollar)

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var loanFundings []LoanFunding
	for rows.Next() {
		var loanFunding LoanFunding
		if err := rows.Scan(
			&loanFunding.ID,
			&loanFunding.LoanOrderNumber,
			&loanFunding.OrderNumber,
			&loanFunding.LoanID,
			&loanFunding.LenderID,
			&loanFunding.LenderEmail,
			&loanFunding.InvestmentAmount,
			&loanFunding.Rate,
			&loanFunding.Interest,
			&loanFunding.ROI,
			&loanFunding.InterestPaid,
			&loanFunding.CapitalAmountPaid,
			&loanFunding.TotalAmountPaid,
			&loanFunding.InvestmentDate,
			&loanFunding.Status,
			&loanFunding.LenderAgreementURL,
			&loanFunding.CreatedAt,
			&loanFunding.UpdatedAt,
			&loanFunding.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		loanFundings = append(loanFundings, loanFunding)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %v", err)
	}

	return loanFundings, nil
}
