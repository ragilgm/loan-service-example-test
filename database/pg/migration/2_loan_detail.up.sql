CREATE TABLE loan_details (
                              id SERIAL PRIMARY KEY,                          -- Loan detail ID (Auto Increment)
                              loan_id INT,                                    -- Loan ID, linking to the loans table
                              borrower_id INT,                                -- Borrower ID
                              business_name VARCHAR(255),                      -- Business name of the borrower
                              business_type VARCHAR(100),                      -- Type of business (e.g., retail, manufacturing)
                              business_address TEXT,                           -- Business address
                              business_phone_number VARCHAR(20),              -- Business phone number
                              business_email VARCHAR(255),                     -- Business email address
                              business_registration_number VARCHAR(50),        -- Business registration number
                              business_annual_revenue DECIMAL(15, 2),          -- Annual revenue of the business
                              business_expense DECIMAL(15, 2),                 -- Annual expenses of the business
                              business_owner_name VARCHAR(255),                -- Business owner's name
                              business_description TEXT,                       -- Business description
                              loan_purpose VARCHAR(255),                       -- Purpose of the loan
                              business_age INT,                               -- Business age (in years)
                              business_sector VARCHAR(100),                    -- Business sector (e.g., agriculture, technology)
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Date of loan detail creation
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Loan detail update date (will be updated via trigger)
                              deleted_at TIMESTAMP DEFAULT NULL                -- Date of loan detail deletion (if applicable)
);