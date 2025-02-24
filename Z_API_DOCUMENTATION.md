# API Documentation for Golang Project Test

This document provides the API details for the **Golang Project Test**. Below are the endpoints for different resources such as **Loans**, **Loan Approvals**, **Loan Fundings**, and **Loan Disbursements**.

## **1. Loan API**

### 1.1 Create Loan
- **Description**:
  - API ini digunakan untuk menghasilkan loan baru dengan status awal `purposed`. Ketika permohonan pinjaman diajukan, sistem secara otomatis membuat data untuk loan approval dengan status awal yaitu `pending`. Asumsi dasar dari API ini adalah bahwa begitu pinjaman diajukan, tim operasional akan menerima pemberitahuan untuk segera melakukan survey dan verifikasi terhadap permohonan pinjaman yang diajukan.

- **Method**: `POST`
- **Endpoint**: `/loans`
- **Request Body**:
  ```json
  
   {
    "borrower_id": 123,
    "request_amount": 1000000.00,
    "loan_grade": "A",
    "loan_type": "productive",
    "rate": 5.5,
    "tenures": 12,
    "detail": {
      "business_name": "ABC Manufacturing",
      "business_type": "Manufacturing",
      "business_address": "123 Industrial Road, Cityville, ST 12345",
      "business_phone_number": "62878",
      "business_email": "contact@abcmfg.com",
      "business_registration_number": "REG12345678",
      "business_annual_revenue": 5000000.00,
      "business_expense": 2000000.00,
      "business_owner_name": "John Doe",
      "business_description": "ABC Manufacturing specializes in producing high-quality widgets and gadgets.",
      "loan_purpose": "Expand production capacity and purchase new machinery",
      "business_age": 10,
      "business_sector": "Manufacturing"
    }
  }
  
    ```

### 1.2 Get All Loans
- **Description**:
  - API ini digunakan untuk mengambil daftar semua pinjaman yang ada di sistem. Dengan menggunakan parameter query, Anda dapat memfilter pinjaman berdasarkan status atau mengatur jumlah pinjaman yang ditampilkan per halaman
  
- **Method**: `GET`
- **Endpoint**: `/loans?page=1&size=10`
- **Query Parameters**:
    - `page`: The page number (e.g., 1)
    - `size`: The number of items per page (e.g., 10)
    - `loan_status` (Optional): The status of the loan (Proposed, Rejected, Approved, Invested, Disbursed, Completed)



### 1.3 Get Loan by ID
- **Description**:
  - API ini digunakan untuk mengambil detail informasi tentang sebuah pinjaman berdasarkan ID uniknya. Dengan menggunakan ID pinjaman, pengguna dapat memperoleh informasi lengkap terkait pinjaman tersebut, termasuk statusnya, jumlah pinjaman, dan detail lainnya yang terkait dengan permohonan.
- **Method**: `GET`
- **Endpoint**: `/loans/{id}`


## **2. Loan Approval API**

### 2.1 Get All Loan Approvals

- **Description**:
  - API ini memungkinkan tim approval untuk melihat daftar semua pinjaman yang sedang dalam proses approval Dengan menggunakan parameter query, tim approval dapat melihat pinjaman berdasarkan status persetujuan tertentu seperti "pending" atau mengambil data berdasarkan paginasi (halaman dan ukuran data per halaman).
  - Note : data ini akan ada hanya jika ada loan masuk 
- **Method**: `GET`
- **Endpoint**: `/loans/approvals?page=1&size=10`
- **Query Parameters**:
    - `page`: The page number (e.g., 1)
    - `size`: The number of items per page (e.g., 10)
    - `approval_status` The status of the approval (pending,approved,rejected)

### 2.2 Update Loan Approval

- **Description**:
  - API ini memungkinkan tim approval untuk memperbarui status persetujuan pinjaman dan melampirkan dokumen yang diperlukan untuk mendukung keputusan tersebut. Dengan API ini, tim yang bertanggung jawab dapat mengubah status persetujuan pinjaman, misalnya dari pending menjadi approved, dan menambahkan dokumen terkait sebagai bukti atau referensi, seperti home visit atau store document.
    - Pada saat yang sama jika tim approval menetujui pinjaman , maka status pinjaman akan berubah secara paralel menjadi `approved` untuk menandakan bahwa pinjaman sudah bisa di danai oleh  `lender/investor`  dan sebaliknya, jika pengajuan di tolak oleh tim approval maka status pinjaman akan menjadi `rejected`
- **Method**: `PUT`
- **Endpoint**: `/loans/approvals/{id}`
- **Request Body**:

```json
 {
    "staff_id": 12345,
    "approval_status": "approved",
    "approval_documents": [
      {
        "document_type": "home_visited",
        "file_url": "http://example.com/ktp.pdf",
        "description": "Document about visited location"
      },
      {
        "document_type": "store_document",
        "file_url": "http://example.com/ktp.pdf",
        "description": "Document about real store"
      }
    ]
  }

```


## **3. Loan Funding API**

### 3.1 Create Loan Funding
- **Description**:
  - API ini digunakan oleh lender (pemberi pinjaman) untuk mendanai pinjaman yang tersedia di platform. Melalui API ini, lender dapat mengajukan jumlah dana yang ingin mereka investasikan dalam pinjaman tertentu. API ini juga memungkinkan lender untuk mengupload URL perjanjian sebagai bukti kesepakatan pendanaan.
  - disclaimer kenapa saya pilih URL bukan file , karna menurut saya ukuran file itu cukup besar sehingga jika di upload ke server performa nya akan berkurang , jadi dalam case ini saya asumsikan bahwa file agreement sudah di setujui dan sudah di tanda tangan oleh lender
  - Di lain proses , pada saat lender mendanai , system akan terus mengkalkulasi total dana yang berhasil di investasikan oleh lender kepada peminjam, prosess nya menggunakan Kafka/Asyc pertimbangan nya adalah karna disini sangat rawan sekali terjadi inkonsistensi data, maka dari itu proses di API ini async jadi lender belum dapat memastikan apakah investasi nya sudah berhasil di masukan atau gagal, untuk keputusan nya itu akan di infokan melalui email, jika gagal maka asumsi saya dana akan di kembalikan kepada lender
  - Jika pendanaan berhasil maka porsi lender langsung di hitung pada saat itu
  - #### Rumus:
  - 1. tenureYears = Tenor / 12 bulan
  - 2. interest = investAmount * (rate/100) * tenureYears
  - Contoh:
    Jika:
    - `principal = 100,000`
    - `annualRate = 10%`
    - `tenureMonths = 24`
    - `Maka bunga adalah = 20,000`
    - `Maka ROI yang di terima lender adalah : investAmount + 20,000 = 120,000`
  - system juga akan mengehcek pada setiap kali pendanaan masuk , apakah total pinjaman sudah sama dengan total yang di investasikan , jika sudah sama maka status pinjaman loan akan berubah menjadi `disbursed`
  - jika pinjaman status nya sudah menjadi `invested` maka sistem akan menggenerate initial `loan_disburse` dengan status `pending`
- **Method**: `POST`
- **Endpoint**: `/loans`
- **Request Body**:

```json
{
  "order_number": "LN1234567890",
  "loan_id": 1,
  "lender_id": 67894,
  "lender_email": "lender@example.com",
  "investment_amount": 100000.00,
  "lender_agreement_url": "http://example.com/agreement"
}

```


### 3.2 Get Loan Funding by Lender ID
- **Description**:
  - API ini digunakan untuk mendapatkan informasi tentang dana yang telah diinvestasikan oleh lender berdasarkan **lender_id** yang diberikan. Melalui API ini, lender dapat melihat daftar semua pinjaman yang berhasil mereka danai atau tidak, dan lender juga dapat mendapatkan informasi informasi mengenai ROI nya.


- **Method**: `GET`
- **Endpoint**: `/loan-fundings/lender/{lender_id}`
- **Request Header**:
  - `Content-Type: application/json`


## **4. Loan Disbursement API**

### 4.1 Get All Loan Disbursements
- **Description**:
  - API ini untuk membantu tim approval untuk mendapatkan daftar disbursement pinjaman baik itu yang belum di prosess `pending` sudah di prosess `completed` atau yang di batalkan `canceled`
  - Note : data ini akan ada hanya jika data loan sudah berhasil di invest oleh lender, untuk mencapai hal ini , loan perlu di invest oleh lender sebanyak x ( yang di butuhkan oleh borrower )

  
- **Method**: `GET`
- **Endpoint**: `/loan-disbursements?page=1&size=10&approval_status=pending`
- **Query Parameters**:
    - `page`: The page number (e.g., 1)
    - `size`: The number of items per page (e.g., 10)
    - `approval_status`: The approval status of the disbursement ( pending,completed,cancelled)


### 4.2 Update Loan Disbursement State
- **Description**:
  - API ini digunakan untuk memperbarui status pencairan dana (disbursement) untuk pinjaman yang telah disetujui. Setelah dana disalurkan kepada peminjam, status disbursement perlu diperbarui menjadi "completed" (selesai).
  - Selain itu, API ini juga memerlukan ID pegawai (staff_id) yang bertanggung jawab atas pembaruan status disbursement dan URL perjanjian yang telah ditandatangani (signed_agreement_url). Perjanjian ini adalah bukti resmi bahwa dana telah dicairkan sesuai dengan kesepakatan.
  - Di prosess lain ketika disbursement sudah berhasil di lakukan , system akan mengupdate status loan menjadi `disbursed` dan akan mengkalkulasi mengenai bunga , dan total yang harus di bayar borrower terhadap pinjaman nya

- **Method**: `PUT`
- **Endpoint**: `/loan-disbursements/{id}`
- **Request Body**:
 
```json

 {
  "loan_id": 1,
  "disbursement_status": "completed",
  "staff_id": 123,
  "signed_agreement_url": "http://google.com"
  }

```

## **Base URL**
All endpoints should be tested on the following base URL:
- `localhost:9090` 
