-- Table for loan approval information
CREATE TABLE loans_approval (
                                id SERIAL PRIMARY KEY,                             -- Auto-incremented ID (SERIAL is used in PostgreSQL)
                                loan_id INT NOT NULL,                      -- Loan ID
                                approval_number VARCHAR(50) NOT NULL,              -- Approval number
                                staff_id INT,                             -- Staff ID who approved the loan
                                approval_date TIMESTAMP ,                  -- Date of approval
                                approval_status VARCHAR(50) NOT NULL,              -- Approval status (pending, approved, rejected)
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Date of approval record creation
                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Date of record update (requires manual update in PostgreSQL)
                                deleted_at TIMESTAMP DEFAULT NULL                  -- Date of approval record deletion (if applicable)
);
