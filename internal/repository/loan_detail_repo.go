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

type (
	LoanDetail struct {
		ID                         int64      `db:"id"`                           // Loan detail ID (Auto Increment)
		LoanID                     int64      `db:"loan_id"`                      // Loan ID, linking to the loans table
		BorrowerID                 int64      `db:"borrower_id"`                  // Borrower ID
		BusinessName               string     `db:"business_name"`                // Business name of the borrower
		BusinessType               string     `db:"business_type"`                // Type of business (e.g., retail, manufacturing)
		BusinessAddress            string     `db:"business_address"`             // Business address
		BusinessPhoneNumber        string     `db:"business_phone_number"`        // Business phone number
		BusinessEmail              string     `db:"business_email"`               // Business email address
		BusinessRegistrationNumber string     `db:"business_registration_number"` // Business registration number
		BusinessAnnualRevenue      float64    `db:"business_annual_revenue"`      // Annual revenue of the business
		BusinessExpense            float64    `db:"business_expense"`             // Annual expenses of the business
		BusinessOwnerName          string     `db:"business_owner_name"`          // Business owner's name
		BusinessDescription        string     `db:"business_description"`         // Business description
		LoanPurpose                string     `db:"loan_purpose"`                 // Purpose of the loan
		BusinessAge                int64      `db:"business_age"`                 // Business age (in years)
		BusinessSector             string     `db:"business_sector"`              // Business sector (e.g., agriculture, technology)
		CreatedAt                  time.Time  `db:"created_at"`                   // Date of loan detail creation
		UpdatedAt                  time.Time  `db:"updated_at"`                   // Loan detail update date
		DeletedAt                  *time.Time `db:"deleted_at"`                   // Date of loan detail deletion (if applicable)
	}
)

var (
	LoanDetailTableName = "loan_details"
	LoanDetailTable     = struct {
		ID                         string
		LoanID                     string
		BorrowerID                 string
		BusinessName               string
		BusinessType               string
		BusinessAddress            string
		BusinessPhoneNumber        string
		BusinessEmail              string
		BusinessRegistrationNumber string
		BusinessAnnualRevenue      string
		BusinessExpense            string
		BusinessOwnerName          string
		BusinessDescription        string
		LoanPurpose                string
		BusinessAge                string
		BusinessSector             string
		CreatedAt                  string
		UpdatedAt                  string
		DeletedAt                  string
	}{
		ID:                         "id",
		LoanID:                     "loan_id",
		BorrowerID:                 "borrower_id",
		BusinessName:               "business_name",
		BusinessType:               "business_type",
		BusinessAddress:            "business_address",
		BusinessPhoneNumber:        "business_phone_number",
		BusinessEmail:              "business_email",
		BusinessRegistrationNumber: "business_registration_number",
		BusinessAnnualRevenue:      "business_annual_revenue",
		BusinessExpense:            "business_expense",
		BusinessOwnerName:          "business_owner_name",
		BusinessDescription:        "business_description",
		LoanPurpose:                "loan_purpose",
		BusinessAge:                "business_age",
		BusinessSector:             "business_sector",
		CreatedAt:                  "created_at",
		UpdatedAt:                  "updated_at",
		DeletedAt:                  "deleted_at",
	}
)

type LoanDetailRepo interface {
	Create(context.Context, *LoanDetail) (int64, error)
	GetByLoanID(ctx context.Context, loanID int64) (*LoanDetail, error)
}

type LoanDetailRepoImpl struct {
	dig.In
	*sql.DB
}

func NewLoanDetailRepo(impl LoanDetailRepoImpl) LoanDetailRepo {
	return &impl
}

// Insert LoanDetail and return last inserted id
func (r *LoanDetailRepoImpl) Create(ctx context.Context, loanDetail *LoanDetail) (int64, error) {
	// Memulai transaksi
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return -1, err
	}

	// Menyusun query untuk insert
	builder := sq.
		Insert(LoanDetailTableName).
		Columns(
			LoanDetailTable.LoanID,
			LoanDetailTable.BorrowerID,
			LoanDetailTable.BusinessName,
			LoanDetailTable.BusinessType,
			LoanDetailTable.BusinessAddress,
			LoanDetailTable.BusinessPhoneNumber,
			LoanDetailTable.BusinessEmail,
			LoanDetailTable.BusinessRegistrationNumber,
			LoanDetailTable.BusinessAnnualRevenue,
			LoanDetailTable.BusinessExpense,
			LoanDetailTable.BusinessOwnerName,
			LoanDetailTable.BusinessDescription,
			LoanDetailTable.LoanPurpose,
			LoanDetailTable.BusinessAge,
			LoanDetailTable.BusinessSector,
			LoanDetailTable.CreatedAt,
			LoanDetailTable.UpdatedAt,
			LoanDetailTable.DeletedAt,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		Values(
			loanDetail.LoanID,
			loanDetail.BorrowerID,
			loanDetail.BusinessName,
			loanDetail.BusinessType,
			loanDetail.BusinessAddress,
			loanDetail.BusinessPhoneNumber,
			loanDetail.BusinessEmail,
			loanDetail.BusinessRegistrationNumber,
			loanDetail.BusinessAnnualRevenue,
			loanDetail.BusinessExpense,
			loanDetail.BusinessOwnerName,
			loanDetail.BusinessDescription,
			loanDetail.LoanPurpose,
			loanDetail.BusinessAge,
			loanDetail.BusinessSector,
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

func (r *LoanDetailRepoImpl) GetByLoanID(ctx context.Context, loanID int64) (*LoanDetail, error) {
	// Menyusun query untuk mengambil data berdasarkan LoanID
	builder := sq.
		Select(
			LoanDetailTable.ID,
			LoanDetailTable.LoanID,
			LoanDetailTable.BorrowerID,
			LoanDetailTable.BusinessName,
			LoanDetailTable.BusinessType,
			LoanDetailTable.BusinessAddress,
			LoanDetailTable.BusinessPhoneNumber,
			LoanDetailTable.BusinessEmail,
			LoanDetailTable.BusinessRegistrationNumber,
			LoanDetailTable.BusinessAnnualRevenue,
			LoanDetailTable.BusinessExpense,
			LoanDetailTable.BusinessOwnerName,
			LoanDetailTable.BusinessDescription,
			LoanDetailTable.LoanPurpose,
			LoanDetailTable.BusinessAge,
			LoanDetailTable.BusinessSector,
			LoanDetailTable.CreatedAt,
			LoanDetailTable.UpdatedAt,
			LoanDetailTable.DeletedAt,
		).
		From(LoanDetailTableName).
		Where(sq.Eq{LoanDetailTable.LoanID: loanID}).
		PlaceholderFormat(sq.Dollar)

	// Jalankan query dan ambil hasilnya
	var loanDetail LoanDetail
	err := builder.RunWith(r.DB).QueryRowContext(ctx).Scan(
		&loanDetail.ID,
		&loanDetail.LoanID,
		&loanDetail.BorrowerID,
		&loanDetail.BusinessName,
		&loanDetail.BusinessType,
		&loanDetail.BusinessAddress,
		&loanDetail.BusinessPhoneNumber,
		&loanDetail.BusinessEmail,
		&loanDetail.BusinessRegistrationNumber,
		&loanDetail.BusinessAnnualRevenue,
		&loanDetail.BusinessExpense,
		&loanDetail.BusinessOwnerName,
		&loanDetail.BusinessDescription,
		&loanDetail.LoanPurpose,
		&loanDetail.BusinessAge,
		&loanDetail.BusinessSector,
		&loanDetail.CreatedAt,
		&loanDetail.UpdatedAt,
		&loanDetail.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Jika tidak ada data ditemukan
			return nil, fmt.Errorf("loan detail not found for LoanID %d", loanID)
		}
		// Error lain saat melakukan query
		return nil, fmt.Errorf("failed to get loan details for LoanID %d: %v", loanID, err)
	}

	return &loanDetail, nil
}
