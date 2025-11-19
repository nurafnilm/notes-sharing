# Notes Sharing App

Aplikasi berbagi notes dengan autentikasi JWT, CRUD notes, upload gambar (bonus), dan logging request/response ke DB. Stack: Next.js 16 (Frontend), Golang Fiber v2 (Backend), PostgreSQL 16 (DB), Docker Compose.

## Fitur
- Register & Login dengan JWT (simpan di localStorage).
- CRUD Notes: Create/Read/List/Delete (hanya owner bisa delete).
- Upload gambar pada notes (multipart/form-data, simpan ke /uploads).
- Logging: Setiap request/response dicatat ke table `logs` di DB (datetime, method/endpoint, masked headers, payload, response, status).
- Deployment: Docker Compose dengan persistent DB + uploads.

## Tech Stack & Arsitektur
- Frontend: Next.js App Router, Axios dengan JWT interceptor.
- Backend: Fiber routes, GORM ORM, Bcrypt hash, JWT HS256.
- DB: Postgres tables (users, notes, logs).
- Docker: Volumes untuk persistence (postgres-data, uploads).

## Setup

### Local (Tanpa Docker)
#### Backend
1. `cd backend && go mod tidy`
2. Jalankan Postgres local (misal Docker: `docker run -p 5432:5432 -e POSTGRES_PASSWORD=rahasia123 -e POSTGRES_DB=notesdb postgres:16`).
3. `cp .env.example .env` → isi vars.
4. `go run main.go` (port 8000).

#### Frontend
1. `cd notes-frontend && npm install`
2. `cp .env.example .env.local` → isi `NEXT_PUBLIC_API_URL=http://localhost:8000`.
3. `npm run dev` (port 3000).

### Docker (Recommended)
1. `cp .env.example .env` → isi vars.
2. `mkdir uploads` (folder gambar).
3. `docker-compose up --build`.
4. Akses: http://localhost:3000.

## Contoh Log (Query DB)
      timestamp      | method | endpoint | status_code
---------------------+--------+----------+-------------
 2025-11-19 13:54:30 | GET    | /notes   |         200
 2025-11-19 13:54:29 | GETION | /notes/2 |         200
 2025-11-19 13:53:58 | GET    | /notes   |         200

## Screenshots
![Register Page](screenshots/tampilan-registrasi.png)  
![Login Page](screenshots/tampilan-login.png) 
![Page Awal](screenshots/tampilan-note-awal.png)  
![Upload Note Tulisan](screenshots/tampilan-new-note.png)  
![Upload Note Gambar](screenshots/tampilan-new-note-gambar.png)  
![Tampilan Setelah Ditambah Notes](screenshots/tampilan-notes.png)  
![Isi Note Tulisan](screenshots/isi-note-tulisan.png)  
![Isi Note Gambar](screenshots/isi-note-gambar.png) 
![Hapus Note](screenshots/hapus-note.png)  