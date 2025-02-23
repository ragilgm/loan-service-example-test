# Struktur Folder Proyek

Berikut adalah struktur direktori proyek Golang yang relevan dengan pengembangan aplikasi:

## Root Project
- **.env**: File konfigurasi environment untuk variabel-variabel yang digunakan oleh aplikasi.
- **go.mod**: File untuk mengelola dependensi modul Go.
- **go.sum**: File yang menyimpan checksum untuk setiap dependensi yang digunakan dalam proyek.
- **Makefile**: File yang berisi instruksi untuk menjalankan perintah build atau task otomatis lainnya.

## Direktori `cmd/`
- **cmd**: Berisi file-file entry point untuk menjalankan aplikasi.
    - **main.go**: File utama yang menjadi titik masuk aplikasi.

## Direktori `database/`
- **database**: Menyimpan skrip atau file yang berkaitan dengan pengelolaan dan setup database.
    - **/pg/migrations/**: Folder yang berisi file skrip migrasi database.
        - **1_loan.up.sql**


## Direktori `deploy/`
- **deploy**: Berisi file docker yang perlu di jalankan agar aplikasi dapat berjalan
    - **kafka.yml**:
    - **pg.yml**



## Struktur Folder `internal/`

## Folder `consts/`
- **`consts/`**: Folder ini berisi file yang mendefinisikan konstanta-konstanta yang digunakan di seluruh aplikasi. Konstanta ini sering digunakan untuk nilai-nilai tetap yang tidak berubah selama runtime, seperti status pinjaman, kode kesalahan, atau konfigurasi lainnya.

## Folder `dto/`
- **`dto/`**: Folder ini berisi **Data Transfer Objects (DTO)** yang digunakan untuk memodelkan data yang dikirimkan melalui API. DTO berfungsi untuk mendefinisikan struktur data yang akan digunakan untuk komunikasi antar lapisan aplikasi atau antara aplikasi dengan pengguna.

## Folder `enum/`
- **`enum/`**: Folder ini berisi file yang mendefinisikan enumerasi (enum) untuk berbagai status atau kategori yang digunakan dalam aplikasi. Enums memungkinkan Anda untuk mendefinisikan nilai tetap yang digunakan dalam berbagai konteks, seperti status pinjaman, jenis pinjaman, dll.

## Folder `handler/`
- **`handler/`**: Folder ini berisi file untuk **HTTP handlers**, yang menangani permintaan dan respons dari client. Di sini, Anda akan menemukan logika yang menangani API routes dan proses permintaan untuk fungsi tertentu, seperti pembuatan pinjaman, penanganan persetujuan, atau pengelolaan pinjaman.

## Folder `infra/`
- **`infra/`**: Folder ini berisi kode yang berkaitan dengan **infrastruktur** aplikasi, seperti koneksi database dan konfigurasi lainnya. Semua yang berhubungan dengan pengelolaan infrastruktur dan integrasi dengan sistem lain ditempatkan di sini.

## Folder `repository/`
- **`repository/`**: Folder ini berisi file yang bertanggung jawab untuk **akses data** dan interaksi dengan database. Repository bertindak sebagai lapisan penghubung antara aplikasi dan penyimpanan data, menyediakan API untuk mengambil, menambah, memperbarui, atau menghapus data.

## Folder `service/`
- **`service/`**: Folder ini berisi file yang mendefinisikan **logika bisnis** aplikasi. Service sering kali memanggil repository untuk mengakses data dan kemudian memprosesnya berdasarkan kebutuhan aplikasi, seperti perhitungan bunga pinjaman atau logika investasi.

## Folder `utils/`
- **`utils/`**: Folder ini berisi file utilitas yang menyediakan berbagai fungsi **bantuan** yang sering digunakan di seluruh aplikasi, seperti pengolahan string, perhitungan waktu, atau pengaturan validasi umum.

## File `shutdown.go`
- **`shutdown.go`**: File ini berisi logika untuk menangani proses **shutdown** aplikasi dengan benar. File ini memastikan bahwa aplikasi dapat dihentikan dengan aman, membersihkan sumber daya yang digunakan, seperti koneksi database atau service lain yang sedang berjalan.

## File `start.go`
- **`start.go`**: File ini berfungsi sebagai titik awal aplikasi. Biasanya, file ini berisi kode untuk menginisialisasi dan memulai server HTTP (seperti `echo`), mengkonfigurasi middleware, dan memastikan semua service yang diperlukan berjalan dengan baik saat aplikasi dimulai.
