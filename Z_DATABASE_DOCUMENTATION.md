# Database Documentation

## Tabel `loans`

Tabel `loans` menyimpan data mengenai pinjaman yang diajukan oleh peminjam (borrower).

| **Kolom**                    | **Tipe Data**          | **Deskripsi**                                                                 |
|------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                           | SERIAL                 | ID pinjaman, auto increment                                                  |
| loan_code                    | VARCHAR(50)            | Kode pinjaman                                                                 |
| borrower_id                  | INT                    | ID peminjam (borrower)                                                        |
| request_amount               | DECIMAL(15, 2)         | Jumlah pinjaman yang diminta oleh peminjam                                    |
| loan_grade                   | VARCHAR(2)             | Kelas pinjaman (A, B, C, D)                                                   |
| loan_type                    | VARCHAR(50)            | Jenis pinjaman (misal: produktif, konsumtif, dll.)                            |
| total_invested_amount        | DECIMAL(15, 2)         | Total dana yang diinvestasikan                                                |
| investor_count               | INT                    | Jumlah investor yang berpartisipasi dalam pinjaman                           |
| funding_deadline             | DATE                   | Tenggat waktu pendanaan                                                      |
| loan_status                  | VARCHAR(50)            | Status pinjaman (proposed, rejected, approved, invested)                      |
| rate                         | DECIMAL(5, 2)          | Suku bunga pinjaman                                                           |
| tenures                      | INT                    | Tenor pinjaman                                                                |
| total_interest               | DECIMAL(15, 2)         | Total bunga yang harus dibayar oleh peminjam                                  |
| total_repayment_amount       | DECIMAL(15, 2)         | Total jumlah yang harus dibayar oleh peminjam (pokok + bunga)                |
| investment_percentage        | DECIMAL(5, 2)          | Persentase bagi hasil untuk investor                                          |
| created_at                   | TIMESTAMP              | Tanggal pembuatan pinjaman                                                   |
| updated_at                   | TIMESTAMP              | Tanggal pembaruan status pinjaman                                             |
| deleted_at                   | TIMESTAMP              | Tanggal penghapusan pinjaman (jika ada)                                       |

## Tabel `loan_details`

Tabel `loan_details` menyimpan informasi lebih rinci tentang pinjaman, termasuk data bisnis peminjam.

| **Kolom**                        | **Tipe Data**          | **Deskripsi**                                                                |
|----------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                               | SERIAL                 | ID detail pinjaman, auto increment                                           |
| loan_id                          | INT                    | ID pinjaman, merujuk ke tabel `loans`                                        |
| borrower_id                      | INT                    | ID peminjam (borrower)                                                       |
| business_name                    | VARCHAR(255)           | Nama bisnis peminjam                                                         |
| business_type                    | VARCHAR(100)           | Jenis bisnis (misal: retail, manufaktur)                                     |
| business_address                 | TEXT                   | Alamat bisnis                                                                |
| business_phone_number            | VARCHAR(20)            | Nomor telepon bisnis                                                          |
| business_email                   | VARCHAR(255)           | Alamat email bisnis                                                           |
| business_registration_number     | VARCHAR(50)            | Nomor registrasi bisnis                                                      |
| business_annual_revenue          | DECIMAL(15, 2)         | Pendapatan tahunan bisnis                                                    |
| business_expense                 | DECIMAL(15, 2)         | Pengeluaran tahunan bisnis                                                   |
| business_owner_name              | VARCHAR(255)           | Nama pemilik bisnis                                                           |
| business_description             | TEXT                   | Deskripsi bisnis                                                             |
| loan_purpose                     | VARCHAR(255)           | Tujuan pinjaman                                                              |
| business_age                     | INT                    | Usia bisnis (dalam tahun)                                                    |
| business_sector                  | VARCHAR(100)           | Sektor bisnis (misal: pertanian, teknologi)                                  |
| created_at                       | TIMESTAMP              | Tanggal pembuatan detail pinjaman                                            |
| updated_at                       | TIMESTAMP              | Tanggal pembaruan detail pinjaman                                            |
| deleted_at                       | TIMESTAMP              | Tanggal penghapusan detail pinjaman (jika ada)                               |

## Tabel `loans_approval`

Tabel `loans_approval` menyimpan informasi terkait persetujuan pinjaman.

| **Kolom**                        | **Tipe Data**          | **Deskripsi**                                                                |
|----------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                               | SERIAL                 | ID persetujuan pinjaman, auto increment                                      |
| loan_id                          | INT                    | ID pinjaman, merujuk ke tabel `loans`                                        |
| approval_number                  | VARCHAR(50)            | Nomor persetujuan                                                             |
| staff_id                         | INT                    | ID staff yang memberikan persetujuan                                         |
| approval_date                    | TIMESTAMP              | Tanggal persetujuan                                                           |
| approval_status                  | VARCHAR(50)            | Status persetujuan (pending, approved, rejected)                              |
| created_at                       | TIMESTAMP              | Tanggal pembuatan record persetujuan                                         |
| updated_at                       | TIMESTAMP              | Tanggal pembaruan record persetujuan                                         |
| deleted_at                       | TIMESTAMP              | Tanggal penghapusan record persetujuan (jika ada)                            |

## Tabel `approval_documents`

Tabel `approval_documents` menyimpan dokumen-dokumen terkait persetujuan pinjaman.

| **Kolom**                        | **Tipe Data**          | **Deskripsi**                                                                |
|----------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                               | SERIAL                 | ID dokumen persetujuan pinjaman, auto increment                              |
| loan_approval_id                 | INT                    | ID persetujuan pinjaman, merujuk ke tabel `loans_approval`                   |
| document_type                    | VARCHAR(100)           | Jenis dokumen (misal: 'business_photo', 'product_photo', dll.)               |
| file_url                         | VARCHAR(255)           | URL file dokumen (foto atau file)                                            |
| description                      | TEXT                   | Deskripsi singkat dokumen                                                    |
| created_at                       | TIMESTAMP              | Tanggal pembuatan dokumen                                                   |
| updated_at                       | TIMESTAMP              | Tanggal pembaruan dokumen                                                   |
| deleted_at                       | TIMESTAMP              | Tanggal penghapusan dokumen (jika ada)                                      |

## Tabel `loan_funding`

Tabel `loan_funding` menyimpan informasi mengenai pendanaan pinjaman yang dilakukan oleh lender.

| **Kolom**                        | **Tipe Data**          | **Deskripsi**                                                                |
|----------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                               | SERIAL                 | ID pendanaan, auto increment                                                |
| loan_order_number                | VARCHAR(50)            | Nomor urut pinjaman                                                           |
| order_number                     | VARCHAR(50)            | Nomor urut saat lender melakukan pendanaan                                   |
| loan_id                          | INT                    | ID pinjaman, merujuk ke tabel `loans`                                        |
| lender_id                        | INT                    | ID lender                                                                   |
| lender_email                     | VARCHAR(255)           | Email lender                                                                 |
| investment_amount                | DECIMAL(15, 2)         | Jumlah investasi oleh lender                                                 |
| rate                             | DECIMAL(5, 2)          | Suku bunga yang diterapkan                                                    |
| interest                         | DECIMAL(15, 2)         | Jumlah bunga yang didapatkan oleh lender                                      |
| roi                              | DECIMAL(15, 2)         | Return on Investment (ROI) yang didapatkan oleh lender                       |
| interest_paid                    | DECIMAL(15, 2)         | Bunga yang sudah dibayar oleh borrower                                        |
| capital_amount_paid              | DECIMAL(15, 2)         | Jumlah pokok yang sudah dibayar                                              |
| total_amount_paid                | DECIMAL(15, 2)         | Total jumlah yang sudah dibayar (pokok + bunga)                              |
| investment_date                  | TIMESTAMP              | Tanggal pendanaan                                                           |
| status                           | VARCHAR(50)            | Status pendanaan (misal: invested, ongoing, completed)                       |
| lender_agreement_url             | VARCHAR(255)           | URL perjanjian lender, diunggah ke cloud                                     |
| created_at                       | TIMESTAMP              | Tanggal pembuatan record pendanaan                                          |
| updated_at                       | TIMESTAMP              | Tanggal pembaruan record pendanaan                                          |
| deleted_at                       | TIMESTAMP              | Tanggal penghapusan record pendanaan (jika ada)                             |

## Tabel `loans_disbursement`

Tabel `loans_disbursement` menyimpan informasi mengenai pencairan pinjaman yang dilakukan.

| **Kolom**                        | **Tipe Data**          | **Deskripsi**                                                                |
|----------------------------------|------------------------|-------------------------------------------------------------------------------|
| id                               | SERIAL                 | ID pencairan, auto increment                                                 |
| loan_id                          | INT                    | ID pinjaman, merujuk ke tabel `loans`                                        |
| disburse_code                    | VARCHAR(50)            | Kode pencairan                                                               |
| disburse_amount                  | DECIMAL(15, 2)         | Jumlah yang dicairkan                                                         |
| disbursement_status              | VARCHAR(50)            | Status pencairan (misal: pending, completed)                                 |
| disburse_date                    | TIMESTAMP              | Tanggal pencairan                                                            |
| staff_id                         | INT                    | ID staff yang menangani pencairan                                            |
| agreement_url                    | VARCHAR(255)           | URL template perjanjian pinjaman                                            |
| signed_agreement_url             | VARCHAR(255)           | URL untuk perjanjian yang sudah ditandatangani                               |
| created_at                       | TIMESTAMP              | Tanggal pembuatan pencairan                                                  |
| updated_at                       | TIMESTAMP              | Tanggal pembaruan pencairan                                                  |
| deleted_at                       | TIMESTAMP              | Tanggal penghapusan pencairan (jika ada)                                     |

---

Dokumentasi ini memberikan gambaran tentang struktur tabel yang digunakan untuk menangani pinjaman, detail pinjaman, persetujuan pinjaman, pendanaan pinjaman, serta pencairan pinjaman. Pastikan Anda menyesuaikan relasi dan field tambahan sesuai dengan kebutuhan aplikasi Anda.
