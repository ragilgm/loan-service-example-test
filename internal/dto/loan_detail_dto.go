package dto

import (
	"time"
)

type LoanDetailRequestDTO struct {
	BusinessName               string  `json:"business_name" valid:"required"`
	BusinessType               string  `json:"business_type" valid:"required"`
	BusinessAddress            string  `json:"business_address" valid:"required"`
	BusinessPhoneNumber        string  `json:"business_phone_number" valid:"required"`
	BusinessEmail              string  `json:"business_email" valid:"required,email"`
	BusinessRegistrationNumber string  `json:"business_registration_number" valid:"required"`
	BusinessAnnualRevenue      float64 `json:"business_annual_revenue" valid:"required"`
	BusinessExpense            float64 `json:"business_expense" valid:"required"`
	BusinessOwnerName          string  `json:"business_owner_name" valid:"required"`
	BusinessDescription        string  `json:"business_description" valid:"required"`
	LoanPurpose                string  `json:"loan_purpose" valid:"required"`
	BusinessAge                int64   `json:"business_age" valid:"required"`
	BusinessSector             string  `json:"business_sector" valid:"required"`
}

type LoanDetailResponseDTO struct {
	ID                         int64      `json:"id"`                           // Loan detail ID
	LoanID                     int64      `json:"loan_id"`                      // Loan ID, linking to the loans table
	BorrowerID                 int64      `json:"borrower_id"`                  // Borrower ID
	BusinessName               string     `json:"business_name"`                // Business name of the borrower
	BusinessType               string     `json:"business_type"`                // Type of business (e.g., retail, manufacturing)
	BusinessAddress            string     `json:"business_address"`             // Business address
	BusinessPhoneNumber        string     `json:"business_phone_number"`        // Business phone number
	BusinessEmail              string     `json:"business_email"`               // Business email address
	BusinessRegistrationNumber string     `json:"business_registration_number"` // Business registration number
	BusinessAnnualRevenue      float64    `json:"business_annual_revenue"`      // Annual revenue of the business
	BusinessExpense            float64    `json:"business_expense"`             // Annual expenses of the business
	BusinessOwnerName          string     `json:"business_owner_name"`          // Business owner's name
	BusinessDescription        string     `json:"business_description"`         // Business description
	LoanPurpose                string     `json:"loan_purpose"`                 // Purpose of the loan
	BusinessAge                int64      `json:"business_age"`                 // Business age (in years)
	BusinessSector             string     `json:"business_sector"`              // Business sector (e.g., agriculture, technology)
	CreatedAt                  time.Time  `json:"created_at"`                   // Date of loan detail creation
	UpdatedAt                  time.Time  `json:"updated_at"`                   // Loan detail update date
	DeletedAt                  *time.Time `json:"deleted_at,omitempty"`         // Date of loan detail deletion (if applicable)
}
