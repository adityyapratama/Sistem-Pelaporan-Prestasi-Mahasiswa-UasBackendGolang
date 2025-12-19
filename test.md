# Panduan Testing API - Sistem Pelaporan Prestasi Mahasiswa

## Prasyarat
1. Jalankan database: `docker compose up -d`
2. Jalankan server: `go run main.go`
3. Import collection Postman: `Sistem_Pelaporan_Prestasi_Mahasiswa.postman_collection.json`

## ⚠️ PENTING: Set Role ID Terlebih Dahulu
Sebelum testing, Anda perlu mendapatkan `role_id` dari database.

**Query SQL untuk mendapatkan Role:**
```sql
SELECT id, name FROM roles;
```

**Set di Postman Variables:**
- `role_id_admin` = UUID untuk role Admin
- `role_id_dosen` = UUID untuk role Dosen  
- `role_id_mahasiswa` = UUID untuk role Mahasiswa

---

## 🔄 Urutan Testing (Step by Step)

### TAHAP 1: SETUP USER & AUTH

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 1 | `/auth/Register` | POST | Register Admin | - |
| 2 | `/auth/Register` | POST | Register Dosen | - |
| 3 | `/auth/Register` | POST | Register Mahasiswa | - |
| 4 | `/auth/Login` | POST | Login sebagai Admin | - |
| 5 | `/auth/profile` | GET | Cek profil Admin | Admin |
| 6 | `/auth/` | GET | Get semua user | Admin |

---

### TAHAP 2: SETUP PROFIL DOSEN & MAHASISWA

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 7 | `/auth/Login` | POST | Login sebagai Dosen | - |
| 8 | `/lectures` | POST | Buat profil Dosen | Dosen |
| 9 | `/lectures/current` | GET | Cek profil dosen sendiri | Dosen |
| 10 | `/auth/Login` | POST | Login sebagai Mahasiswa | - |
| 11 | `/students` | POST | Buat profil Mahasiswa | Mahasiswa |
| 12 | `/students/current` | GET | Cek profil mahasiswa sendiri | Mahasiswa |

---

### TAHAP 3: ADMIN - ASSIGN DOSEN PEMBIMBING

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 13 | `/auth/Login` | POST | Login sebagai Admin | - |
| 14 | `/lectures` | GET | Lihat semua dosen (ambil lecture_id) | Admin |
| 15 | `/students` | GET | Lihat semua mahasiswa (ambil student_id) | Admin |
| 16 | `/students/{id}/advisor` | PUT | Assign dosen pembimbing | Admin |

---

### TAHAP 4: MAHASISWA - LAPOR PRESTASI

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 17 | `/auth/Login` | POST | Login sebagai Mahasiswa | - |
| 18 | `/achievements/upload` | POST | Upload lampiran (opsional) | Mahasiswa |
| 19 | `/achievements` | POST | Buat prestasi (status: draft) | Mahasiswa |
| 20 | `/achievements/achiev` | GET | Lihat prestasi saya | Mahasiswa |
| 21 | `/achievements/{id}` | PUT | Edit prestasi (jika perlu) | Mahasiswa |
| 22 | `/achievements/{id}/submit` | PATCH | Submit prestasi untuk verifikasi | Mahasiswa |

---

### TAHAP 5: DOSEN - VERIFIKASI PRESTASI

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 23 | `/auth/Login` | POST | Login sebagai Dosen | - |
| 24 | `/achievements?status=submitted` | GET | Lihat prestasi yang perlu diverifikasi | Dosen |
| 25 | `/achievements/{id}` | GET | Lihat detail prestasi | Dosen |
| 26 | `/achievements/{id}/verify` | PATCH | Verifikasi (approve/reject) | Dosen |

---

### TAHAP 6: ADMIN - LAPORAN & STATISTIK

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 27 | `/auth/Login` | POST | Login sebagai Admin | - |
| 28 | `/reports/statistics` | GET | Lihat statistik sistem | Admin |
| 29 | `/reports/student/{id}` | GET | Lihat laporan per mahasiswa | Admin |

---

### TAHAP 7: ADMIN - PERMISSION MANAGEMENT (Opsional)

| No | Endpoint | Method | Deskripsi | Login Sebagai |
|----|----------|--------|-----------|---------------|
| 30 | `/permissions` | POST | Buat permission baru | Admin |
| 31 | `/permissions` | GET | Lihat semua permission | Admin |
| 32 | `/permissions/assign` | POST | Assign permission ke role | Admin |

---

## 📋 Data Test untuk Setiap Endpoint

### 1. Register Admin
```json
{
    "username": "admin",
    "email": "admin@university.ac.id",
    "password": "Admin123!",
    "full_name": "Administrator Sistem",
    "role_id": "{{role_id_admin}}"
}
```

### 2. Register Dosen
```json
{
    "username": "dosen_budi",
    "email": "budi.santoso@university.ac.id",
    "password": "Dosen123!",
    "full_name": "Dr. Budi Santoso, M.Kom",
    "role_id": "{{role_id_dosen}}"
}
```

### 3. Register Mahasiswa
```json
{
    "username": "mhs_andi",
    "email": "andi.pratama@student.university.ac.id",
    "password": "Mahasiswa123!",
    "full_name": "Andi Pratama",
    "role_id": "{{role_id_mahasiswa}}"
}
```

### 4. Login Admin
```json
{
    "username": "admin",
    "password": "Admin123!"
}
```

### 7. Login Dosen
```json
{
    "username": "dosen_budi",
    "password": "Dosen123!"
}
```

### 8. Buat Profil Dosen
```json
{
    "lecturer_id": "DSN001",
    "department": "Fakultas Ilmu Komputer"
}
```

### 10. Login Mahasiswa
```json
{
    "username": "mhs_andi",
    "password": "Mahasiswa123!"
}
```

### 11. Buat Profil Mahasiswa
```json
{
    "student_id": "2024001001",
    "name": "Andi Pratama",
    "program_study": "Teknik Informatika",
    "academic_year": "2024/2025"
}
```

### 16. Assign Dosen Pembimbing
```json
{
    "advisor_id": "{{lecture_id}}"
}
```

### 19. Buat Prestasi - Kompetisi
```json
{
    "type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024",
    "description": "Memenangkan kompetisi hackathon tingkat nasional dengan tema Smart City.",
    "details": {
        "event_name": "Hackathon Indonesia 2024",
        "organizer": "Kementerian Kominfo RI",
        "location": "Jakarta Convention Center",
        "date": "2024-11-15",
        "rank": "Juara 1",
        "prize": "Rp 50.000.000",
        "team_members": ["Andi Pratama", "Budi Setiawan", "Citra Dewi"],
        "project_name": "SmartCity Dashboard"
    },
    "tags": ["hackathon", "nasional", "juara 1", "smart city", "2024"],
    "attachments": [
        {
            "fileName": "sertifikat_juara.pdf",
            "fileType": "application/pdf",
            "fileUrl": "/uploads/sertifikat_juara.pdf",
            "uploadedAt": "2024-11-20T10:00:00Z"
        }
    ]
}
```

### 19b. Buat Prestasi - Akademik (Alternatif)
```json
{
    "type": "academic",
    "title": "Paper Publikasi di Jurnal Internasional",
    "description": "Publikasi paper penelitian di IEEE Access.",
    "details": {
        "journal_name": "IEEE Access",
        "paper_title": "Machine Learning for Smart Agriculture",
        "doi": "10.1109/ACCESS.2024.123456",
        "publication_date": "2024-10-01",
        "impact_factor": 3.9,
        "authors": ["Andi Pratama", "Dr. Budi Santoso"],
        "indexed": ["Scopus", "Web of Science"]
    },
    "tags": ["publikasi", "jurnal internasional", "IEEE", "machine learning"]
}
```

### 19c. Buat Prestasi - Organisasi (Alternatif)
```json
{
    "type": "organization",
    "title": "Ketua BEM Fakultas 2024",
    "description": "Terpilih sebagai Ketua BEM Fakultas Ilmu Komputer.",
    "details": {
        "organization": "BEM Fakultas Ilmu Komputer",
        "position": "Ketua",
        "period": "2024-2025",
        "election_date": "2024-01-15",
        "programs": ["Pelatihan Coding untuk SMA", "IT Career Fair"]
    },
    "tags": ["organisasi", "BEM", "leadership"]
}
```

### 21. Update Prestasi
```json
{
    "type": "competition",
    "title": "Juara 1 Hackathon Nasional 2024 - Updated",
    "description": "Memenangkan kompetisi hackathon tingkat nasional dengan tema Smart City. Tim kami berhasil mengembangkan solusi inovatif.",
    "details": {
        "event_name": "Hackathon Indonesia 2024",
        "organizer": "Kementerian Kominfo RI",
        "location": "Jakarta Convention Center",
        "date": "2024-11-15",
        "rank": "Juara 1",
        "prize": "Rp 50.000.000",
        "team_members": ["Andi Pratama", "Budi Setiawan", "Citra Dewi"],
        "project_name": "SmartCity Waste Management Dashboard",
        "github_url": "https://github.com/team/smartcity-dashboard"
    },
    "tags": ["hackathon", "nasional", "juara 1", "smart city", "2024", "waste management"]
}
```

### 26a. Verifikasi Prestasi - APPROVE
```json
{
    "status": "verified",
    "notes": "Prestasi telah diverifikasi. Sertifikat dan bukti pendukung sudah lengkap dan valid."
}
```

### 26b. Verifikasi Prestasi - REJECT
```json
{
    "status": "rejected",
    "notes": "Prestasi ditolak karena bukti sertifikat tidak jelas. Mohon upload ulang dengan resolusi yang lebih baik."
}
```

### 30. Buat Permission
```json
{
    "name": "manage_achievements",
    "resource": "achievements",
    "action": "create,read,update,delete",
    "description": "Izin untuk mengelola data prestasi mahasiswa"
}
```

### 32. Assign Permission ke Role
```json
{
    "role_id": "{{role_id_dosen}}",
    "permission_id": "{{permission_id}}"
}
```

---

## 🔑 Tips Testing

1. **Simpan Token**: Setelah login, copy `access_token` ke variable Postman
2. **Simpan ID**: Setelah create data, copy ID untuk digunakan di endpoint selanjutnya
3. **Refresh Token**: Jika token expired, gunakan endpoint `/auth/refresh`
4. **Cek Status**: Prestasi memiliki flow: `draft` → `submitted` → `verified`/`rejected`

---

## 🎯 Expected Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    ALUR SISTEM                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. ADMIN register semua user (Admin, Dosen, Mahasiswa)     │
│                          ↓                                  │
│  2. DOSEN login & buat profil dosen                         │
│                          ↓                                  │
│  3. MAHASISWA login & buat profil mahasiswa                 │
│                          ↓                                  │
│  4. ADMIN assign dosen pembimbing ke mahasiswa              │
│                          ↓                                  │
│  5. MAHASISWA buat prestasi (status: draft)                 │
│                          ↓                                  │
│  6. MAHASISWA submit prestasi (status: submitted)           │
│                          ↓                                  │
│  7. DOSEN verifikasi prestasi (status: verified/rejected)   │
│                          ↓                                  │
│  8. ADMIN lihat laporan & statistik                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 Status Prestasi

| Status | Deskripsi | Siapa yang set |
|--------|-----------|----------------|
| `draft` | Baru dibuat, belum disubmit | Mahasiswa (otomatis) |
| `submitted` | Sudah disubmit, menunggu verifikasi | Mahasiswa |
| `verified` | Sudah diverifikasi/disetujui | Dosen |
| `rejected` | Ditolak, perlu diperbaiki | Dosen |

---

## ❌ Troubleshooting

| Error | Solusi |
|-------|--------|
| 401 Unauthorized | Login ulang, pastikan token di header |
| 404 Not Found | Cek ID yang digunakan sudah benar |
| 400 Bad Request | Cek format JSON body |
| 409 Conflict | Data sudah ada (duplicate) |
| 500 Internal Error | Cek log server untuk detail error |
