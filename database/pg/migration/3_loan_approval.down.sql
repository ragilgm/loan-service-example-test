-- Menghapus index yang sudah dibuat
DROP INDEX IF EXISTS idx_loans_approval_loan_id;
DROP INDEX IF EXISTS idx_loans_approval_staff_id;
DROP INDEX IF EXISTS idx_loans_approval_approval_status;

-- Menghapus tabel loans_approval
DROP TABLE IF EXISTS loans_approval;

