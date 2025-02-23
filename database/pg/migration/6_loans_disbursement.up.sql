CREATE TABLE loans_disbursement (
                                    id SERIAL PRIMARY KEY,                              -- Disbursement ID, SERIAL menggantikan AUTO_INCREMENT
                                    loan_id INT NOT NULL,                        -- Loan ID
                                    disburse_code VARCHAR(50) NOT NULL,                  -- Disbursement code
                                    disburse_amount DECIMAL(15, 2),                      -- Disbursed amount
                                    disbursement_status VARCHAR(50),                     -- Disbursement status (pending, completed, etc.)
                                    disburse_date TIMESTAMP DEFAULT NULL,                -- Disbursement date
                                    staff_id INT NULL,                              -- Staff ID handling the disbursement
                                    agreement_url VARCHAR(255) NOT NULL,                -- URL template agreement url
                                    signed_agreement_url VARCHAR(255)  NULL,          -- URL to the signed loan agreement letter
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- Date of disbursement record creation
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,      -- Date of disbursement record update
                                    deleted_at TIMESTAMP DEFAULT NULL                    -- Date of disbursement record deletion (if applicable)
);
