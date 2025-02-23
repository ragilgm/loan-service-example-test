CREATE TABLE loans (
                       id SERIAL PRIMARY KEY,  -- Loan ID, auto increment using SERIAL in PostgreSQL
                       loan_code VARCHAR(50) NOT NULL,     -- Loan code
                       borrower_id INT,                    -- Borrower ID
                       request_amount DECIMAL(15, 2) DEFAULT 0,  -- Loan request amount
                       loan_grade VARCHAR(2),              -- Loan grade (A, B, C, D)
                       loan_type VARCHAR(50),              -- Type of loan (productive, consumptive, etc.)
                       total_invested_amount DECIMAL(15, 2) DEFAULT 0, -- Total amount invested
                       investor_count INT DEFAULT 0,       -- Number of investors participating
                       funding_deadline DATE,              -- Funding deadline
                       loan_status VARCHAR(50),            -- Loan status (proposed, rejected, approved, invested)
                       rate DECIMAL(5, 2),                 -- Interest rate
                       tenures INT DEFAULT 0,  -- tenure loan
                       total_interest DECIMAL(15, 2) DEFAULT 0, -- total interest which need borrower pay
                       total_repayment_amount DECIMAL(15, 2) DEFAULT 0, -- total repayment which need borrower pay
                       investment_percentage DECIMAL(5, 2), -- Investor profit sharing percentage
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Loan creation date
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Loan status update date
                       deleted_at TIMESTAMP DEFAULT NULL   -- Date of loan deletion (if applicable)
);
