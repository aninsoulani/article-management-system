# Project Artikel (Go & Next.js)

Ini adalah repository untuk aplikasi artikel sederhana yang dibuat menggunakan Go untuk backend API dan Next.js untuk dashboard frontend.

## 📂 Struktur Folder

Repository ini adalah monorepo yang berisi dua proyek utama:

- `/article-services-go`: Folder untuk layanan backend yang dibuat dengan bahasa Go.
- `/article-dashboard-nextjs`: Folder untuk aplikasi frontend dashboard yang dibuat dengan Next.js.

---

## 🛠️ Prasyarat (Prerequisites)

Pastikan perangkat Anda sudah terinstall:

- [Go](https://go.dev/doc/install) (versi 1.18 atau lebih baru)
- [Node.js](https://nodejs.org/en/) (versi 16 atau lebih baru)
- NPM / Yarn

---

## 🚀 Cara Menjalankan Proyek

Proyek ini terdiri dari dua bagian (backend dan frontend) yang harus dijalankan secara terpisah di terminal yang berbeda.

### 1. Menjalankan Backend (Go Service)

a. Buka terminal baru, lalu masuk ke direktori backend:

```bash
cd article-services-go
```

b. Download semua dependency yang dibutuhkan:

```bash
go mod tidy
```

c. Jalankan server backend:

```bash
go run main.go
```

✅ Server backend akan berjalan di `http://localhost:8080` (atau port lain sesuai konfigurasi Anda).

### 2. Menjalankan Frontend (Next.js Dashboard)

a. Buka terminal **kedua** (biarkan terminal backend tetap berjalan), lalu masuk ke direktori frontend:

```bash
cd article-dashboard-nextjs
```

b. Install semua dependency yang dibutuhkan:

```bash
npm install
```

c. Jalankan server development frontend:

```bash
npm run dev
```

✅ Aplikasi frontend akan berjalan dan bisa diakses di `http://localhost:3000`.
