-- Menghapus index yang sudah dibuat
DROP INDEX IF EXISTS edx_loans_disbursement_loan_id;
DROP INDEX IF EXISTS edx_loans_disbursement_disbursement_status;
DROP INDEX IF EXISTS edx_loans_disbursement_disburse_code;
DROP INDEX IF EXISTS edx_loans_disbursement_staff_id;

-- Menghapus tabel loans_disbursement
DROP TABLE IF EXISTS loans_disbursement;
