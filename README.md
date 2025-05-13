# Pemanfaatan Algoritma BFS dan DFS dalam Pencarian Recipe pada Permainan Little Alchemy 2

Aplikasi web untuk mencari resep elemen pada permainan Little Alchemy 2 menggunakan algoritma BFS dan DFS. Pengguna dapat mencari cara mendapatkan suatu elemen dari kombinasi elemen-elemen dasar.

## Daftar Isi

- [Fitur](#fitur)
- [Teknologi](#teknologi)
- [Instalasi](#instalasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Struktur Folder](#struktur-folder)
- [Kontributor](#kontributor)
- [Lisensi](#lisensi)

---

## Fitur

- Pencarian elemen Little Alchemy 2 menggunakan algoritma BFS, DFS, Bidirectional BFS, serta jumlah resep yang diinginkan
- Daftar elemen Little Alchemy 2

## Teknologi

- **Frontend:** React.js
- **Backend:** Go (Golang)
- **CSS:** Custom CSS
- **Docker:** Mendukung deployment dengan Docker Compose

## Instalasi

1. **Clone repository**
    ```bash
    git clone https://github.com/username/NamaKelompok.git
    cd NamaKelompok
    ```

2. **Jalankan dengan Docker (opsional, direkomendasikan)**
    ```bash
    docker compose up
    ```
    Untuk menghentikan:
    ```bash
    docker compose down
    ```

3. **Tanpa Docker: Jalankan manual**
    - **Frontend**
        ```bash
        cd src/frontend
        npm install
        npm start
        ```
    - **Backend**
        ```bash
        cd src/backend
        go run main.go
        ```

## Menjalankan Aplikasi

- Jika menggunakan Docker, aplikasi akan berjalan otomatis.
- Jika manual:
    - Frontend: biasanya di `http://localhost:3000`
    - Backend: biasanya di `http://localhost:8080` (atau port lain sesuai konfigurasi)

## Struktur Folder

```
Tubes2_NamaKelompok/
│
├── src/
│   ├── backend/                # Source code backend (Go)
│   │   ├── data/               # Hasil pembentukan tree resep
│   │   ├── scraper/            # Source code scraper
│   │   ├── Dockerfile          # File konfigurasi Docker backend
│   │   └── main.go             # Source code backend utama
│   └── frontend/               # Source code frontend (React.js)
│       ├── public/             # Static assets frontend
│       ├── Dockerfile          # File konfigurasi Docker frontend
│       └── src/                # Source code utama frontend
│           ├── App.css         # Styling utama frontend
│           ├── App.js          # Komponen utama React
│           ├── index.js        # Entry point React
│           ├── pages/          # Halaman website
│           ├── media/          # Asset gambar/icon untuk frontend
│           └── ...             # File JS/komponen lainnya
├── data/                       # Data scraping resep Little Alchemy
├── doc/                        # Dokumentasi projek
├── docker-compose.yml          # Konfigurasi Docker Compose
├── README.md                   # Dokumentasi proyek
└── ...
```

## Kontributor

| Nama                          | NIM        |
|-------------------------------|------------|
| Nicholas Andhika Lucas        | 13523014   |
| Samantha Laqueenna Ginting    | 13523138   |
| Naufarrel Zhafif Abhista      | 13523149   |
