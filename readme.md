# ZenStore

Dibuat untuk memenuhi technical test dari perusahaan PT. ZEN MULTIMEDIA INDONESIA sebagai GO Developer

## Tech Stack

- **Golang**: Bahasa pemrograman yang digunakan untuk membangun aplikasi ini.
- **GORM**: ORM untuk Golang yang digunakan untuk berinteraksi dengan database.
- **GorillaMux**: Router HTTP yang sangat fleksibel untuk routing.
- **JWT (JSON Web Token)**: Digunakan untuk otentikasi dan otorisasi.
- **Session**: Mengelola sesi pengguna untuk pengalaman pengguna yang lebih baik.
- **Bcrypt**: Algoritma untuk enkripsi password, memastikan keamanan data pengguna.
- **MySQL**: Database relasional yang digunakan untuk menyimpan semua data aplikasi.

## Fitur

- **Manajemen Produk**: Tambahkan, edit, dan hapus produk.
- **Keranjang Belanja**: Tambahkan produk ke keranjang belanja.
- **Otentikasi**: Registrasi dan login pengguna dengan enkripsi password menggunakan bcrypt.
- **Otorisasi**: Menggunakan JWT untuk mengamankan endpoint API.
- **Manajemen Sesi**: Mengelola sesi pengguna untuk menjaga pengalaman pengguna.
- **Swagger Documentation**: Dokumentasi API interaktif yang dapat diakses di [Swagger UI](http://localhost:8000/swagger).

## Instalasi

- go mod tidy
- sesuaikan env dengan setup database
- go run main.go
