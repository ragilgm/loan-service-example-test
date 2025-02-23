-- Table for approval documents related to loan approval
CREATE TABLE approval_documents (
                                    id SERIAL PRIMARY KEY,                             -- Auto-incremented ID (SERIAL is used in PostgreSQL)
                                    loan_approval_id INT NOT NULL,                     -- Loan approval ID, linked to the loans_approval table
                                    document_type VARCHAR(100),                        -- Document type (e.g., 'business_photo', 'product_photo', etc.)
                                    file_url VARCHAR(255) NOT NULL,                     -- File URL for the document (photo or file)
                                    description TEXT,                                  -- Brief description of the document
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Date of document record creation
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    -- Date of document record update (requires manual update in PostgreSQL)
                                    deleted_at TIMESTAMP DEFAULT NULL                  -- Date of document record deletion (if applicable)
);