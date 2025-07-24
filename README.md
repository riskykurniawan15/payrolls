# Aplikasi Payroll System

Aplikasi sistem penggajian yang dibangun dengan Go (Golang) menggunakan arsitektur clean architecture. Aplikasi ini menangani perhitungan gaji karyawan berdasarkan periode kerja, lembur, dan reimbursement.

## 🚀 Fitur Utama

- **Manajemen Periode**: Membuat dan mengelola periode penggajian
- **Perhitungan Gaji**: Otomatis menghitung gaji berdasarkan kehadiran, lembur, dan reimbursement
- **Slip Gaji**: Generate dan print slip gaji dalam format HTML
- **Laporan Ringkasan**: Laporan ringkasan penggajian untuk manajemen
- **Autentikasi JWT**: Sistem login dan autentikasi yang aman
- **API RESTful**: Endpoint API yang lengkap untuk integrasi

## 📋 Prasyarat

Sebelum menjalankan aplikasi, pastikan sistem Anda memiliki:

- **Go 1.23.0** atau versi yang lebih baru
- **PostgreSQL 12** atau versi yang lebih baru
- **Make** (opsional, untuk menggunakan Makefile)
- **Migrate CLI** (untuk database migration)

## 🛠️ Instalasi

### 1. Clone Repository

```bash
git clone <repository-url>
cd payrolls
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Setup Environment

Buat file `.env` di root directory dengan konfigurasi berikut:

```env
# Server Configuration
SERVER=localhost
PORT=9000
USING_SECURE=false

# Database Configuration
DB_USER=root
DB_PASS=your_password
DB_SERVER=localhost
DB_PORT=5432
DB_NAME=payroll_db
DB_SSL_MODE=disable
DB_TIME_ZONE=Asia/Jakarta
DB_MAX_IDLE_CON=10
DB_MAX_OPEN_CON=100
DB_MAX_LIFE_TIME=10
DB_DEBUG=false

# HASH
HASHING_COST=10

# JWT Configuration
JWT_SECRET_KEY=your_secret_key_here
JWT_EXPIRED=24

# Logger Configuration
LOG_OUTPUT_MODE=both
LOG_LEVEL=debug
LOG_DIR=logger

# Company Configuration
COMPANY_NAME=Your Company Name
```

### 4. Setup Database

#### Menggunakan Makefile:
```bash
# Buat database PostgreSQL terlebih dahulu
# Kemudian jalankan migration
make migrate-up
```

#### Manual:
```bash
# Install migrate CLI jika belum ada
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Jalankan migration
migrate -path database/migrations -database "postgres://username:password@localhost:5432/payroll_db?sslmode=disable&timezone=Asia/Jakarta" up
```

### 5. Seed Data (Opsional)

#### Menggunakan Makefile:
```bash
# Seed semua data
make seed

# Seed hanya user
make seed-users
```

#### Manual:
```bash
go run cmd/seeder/main.go -type=all
```

## 🚀 Cara Menjalankan Aplikasi

### Menggunakan Makefile (Direkomendasikan)

#### 1. Menjalankan Aplikasi
```bash
make run
```

#### 2. Build Aplikasi
```bash
make build
```

#### 3. Menjalankan Test
```bash
# Test biasa
make test

# Test dengan coverage
make test-coverage

# Test dengan verbose output
make test-v
```

#### 4. Database Operations
```bash
# Lihat konfigurasi database
make db-config

# Migration up
make migrate-up

# Migration down (1 step)
make migrate-down

# Migration down (semua)
make migrate-down-all

# Lihat versi migration
make migrate-version

# Buat migration baru
make migrate-create name=nama_migration
```

#### 5. Dependency Management
```bash
# Install dependencies
make deps

# Update dependencies
make deps-update

# Generate wire code
make wire
```

#### 6. Cleanup
```bash
# Hapus file build
make clean
```

### Manual (Tanpa Makefile)

#### 1. Menjalankan Aplikasi
```bash
go run main.go
```

#### 2. Build Aplikasi
```bash
go build -o bin/payrolls main.go
```

#### 3. Menjalankan Binary
```bash
./bin/payrolls
```

#### 4. Menjalankan Test
```bash
# Test semua package
go test ./...

# Test dengan coverage
go test -cover ./...

# Test dengan verbose
go test -v ./...

# Test package tertentu
go test ./repositories/period_detail -v
```

#### 5. Database Migration Manual
```bash
# Migration up
migrate -path database/migrations -database "postgres://username:password@localhost:5432/payroll_db?sslmode=disable&timezone=Asia/Jakarta" up

# Migration down
migrate -path database/migrations -database "postgres://username:password@localhost:5432/payroll_db?sslmode=disable&timezone=Asia/Jakarta" down 1

# Buat migration baru
migrate create -ext sql -dir database/migrations nama_migration
```

## 📚 Struktur Aplikasi

```
payrolls/
├── cmd/                    # Command line tools
│   └── seeder/            # Database seeder
├── config/                # Konfigurasi aplikasi
├── constant/              # Konstanta aplikasi
├── database/              # Database migrations
│   └── migrations/
├── driver/                # Database driver
├── infrastructure/        # Infrastructure layer
│   └── http/             # HTTP server
│       ├── engine/       # Echo server setup
│       ├── entities/     # Response entities
│       ├── handler/      # HTTP handlers
│       ├── middleware/   # HTTP middlewares
│       └── router/       # Route definitions
├── mocks/                # Mock files untuk testing
├── models/               # Data models
│   ├── audit_trail/     # Audit trail models
│   ├── attendance/      # Attendance models
│   ├── health/          # Health check models
│   ├── overtime/        # Overtime models
│   ├── payslip/         # Payslip models
│   ├── period/          # Period models
│   ├── period_detail/   # Period detail models
│   ├── reimbursement/   # Reimbursement models
│   └── user/            # User models
├── repositories/         # Data access layer
│   ├── audit_trail/     # Audit trail repository
│   ├── attendance/      # Attendance repository
│   ├── health/          # Health check repository
│   ├── instance/        # Database instance
│   ├── overtime/        # Overtime repository
│   ├── period/          # Period repository
│   ├── period_detail/   # Period detail repository
│   ├── reimbursement/   # Reimbursement repository
│   └── user/            # User repository
├── services/             # Business logic layer
│   ├── audit_trail/     # Audit trail service
│   ├── attendance/      # Attendance service
│   ├── health/          # Health check service
│   ├── overtime/        # Overtime service
│   ├── payslip/         # Payslip service
│   ├── period/          # Period service
│   ├── period_detail/   # Period detail service
│   ├── reimbursement/   # Reimbursement service
│   └── user/            # User service
├── utils/                # Utility functions
│   ├── bcrypt/          # Password hashing
│   ├── code_generator/  # Code generation utilities
│   ├── data_tipes/      # Custom data types
│   ├── env/             # Environment utilities
│   ├── jwt/             # JWT utilities
│   ├── logger/          # Logging utilities
│   └── validator/       # Validation utilities
├── main.go              # Entry point
├── go.mod               # Go modules
├── go.sum               # Go modules checksum
├── Makefile             # Build automation
└── README.md            # Dokumentasi ini
```

## 🔧 Konfigurasi

### Environment Variables

| Variable | Deskripsi | Default |
|----------|-----------|---------|
| `SERVER` | Host server | `localhost` |
| `PORT` | Port server | `9000` |
| `USING_SECURE` | Gunakan HTTPS | `false` |
| `DB_USER` | Username database | `root` |
| `DB_PASS` | Password database | `` |
| `DB_SERVER` | Host database | `localhost` |
| `DB_PORT` | Port database | `5432` |
| `DB_NAME` | Nama database | `public` |
| `DB_SSL_MODE` | Mode SSL database | `disable` |
| `DB_TIME_ZONE` | Timezone database | `Asia/Jakarta` |
| `JWT_SECRET_KEY` | Secret key JWT | `` |
| `JWT_EXPIRED` | Expired JWT (jam) | `24` |
| `LOG_OUTPUT_MODE` | Mode output log | `both` |
| `LOG_LEVEL` | Level log | `debug` |
| `LOG_DIR` | Directory log | `logger` |
| `COMPANY_NAME` | Nama perusahaan | `Blank Company` |

## 📡 API Endpoints

### Authentication
- `POST /auth/login` - Login user

### Health Check
- `GET /health` - Health check endpoint

### User Management
- `GET /user/profile` - Get user profile (protected)

### Period Management (Admin only)
- `POST /periods` - Create new period
- `GET /periods` - List all periods
- `GET /periods/:id` - Get period by ID
- `PUT /periods/:id` - Update period
- `DELETE /periods/:id` - Delete period
- `POST /periods/:id/run-payroll` - Run payroll for period

### Attendance (Employee only)
- `GET /attendances` - Get user attendances
- `GET /attendances/:id` - Get attendance by ID
- `POST /attendances/check-in` - Check in
- `POST /attendances/check-out` - Check out (only latest record)
- `POST /attendances/check-out/:id` - Check out by ID

### Overtime (Employee only)
- `GET /overtimes` - Get user overtimes
- `GET /overtimes/:id` - Get overtime by ID
- `POST /overtimes` - Create overtime
- `PUT /overtimes/:id` - Update overtime
- `DELETE /overtimes/:id` - Delete overtime

### Reimbursement (Employee only)
- `GET /reimbursements` - Get user reimbursements
- `GET /reimbursements/:id` - Get reimbursement by ID
- `POST /reimbursements` - Create reimbursement
- `PUT /reimbursements/:id` - Update reimbursement
- `DELETE /reimbursements/:id` - Delete reimbursement

### Payslip (Employee only)
- `GET /payslip` - List user payslips
- `POST /payslip/generate/:id` - Generate payslip
- `GET /payslip/print?token=xxx` - Print payslip

### Payslip (Admin only)
- `GET /payslip/summary/generate/:id` - Generate summary payslip
- `GET /payslip/summary/print?token=xxx` - Generate payslip summary

## 📝 Catatan Penggunaan API

### Authentication
- Semua endpoint (kecuali `/auth/login` dan `/health`) memerlukan JWT token
- Token harus dikirim dalam header: `Authorization: Bearer YOUR_JWT_TOKEN`
- Token expired dalam 24 jam (dapat dikonfigurasi)

### Error Response Format
```json
{
  "status": 400,
  "error": "Error message description"
}
```

### Pagination
Endpoint yang mendukung pagination akan mengembalikan response dengan format:
```json
{
  "status": 200,
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Query Parameters
Beberapa endpoint mendukung query parameters:
- `page`: Halaman (default: 1)
- `limit`: Jumlah item per halaman (default: 10, max: 100)
- `sort_by`: Field untuk sorting
- `sort_desc`: Boolean untuk descending sort (default: false)

## 🧪 Testing

### Menjalankan Test

```bash
# Test semua package
make test

# Test dengan coverage
make test-coverage

# Test dengan verbose
make test-v

# Test package tertentu
go test ./repositories/period_detail -v
```

### Coverage Report

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out
```

## 📝 Logging

Aplikasi menggunakan Zap logger dengan konfigurasi:

- **Output Mode**: `terminal`, `file`, atau `both`
- **Log Level**: `debug`, `info`, `warn`, `error`
- **Log Directory**: `logger/` (default)

Log akan disimpan dalam format:
```
logger/
├── app.log
├── app.log.2024-01-01
└── app.log.2024-01-02
```

## 🔒 Security

### JWT Authentication
- Secret key harus diatur di environment variable `JWT_SECRET_KEY`
- Token expired dalam 24 jam (dapat dikonfigurasi)
- Middleware JWT diterapkan pada semua endpoint protected

## 🚨 Troubleshooting

### Masalah Umum

1. **Database Connection Error**
   - Pastikan PostgreSQL berjalan
   - Periksa konfigurasi database di `.env`
   - Pastikan database sudah dibuat

2. **Migration Error**
   - Pastikan migrate CLI terinstall
   - Periksa connection string database
   - Jalankan `make migrate-version` untuk cek status

3. **Port Already in Use**
   - Ganti port di `.env` file
   - Atau hentikan aplikasi yang menggunakan port tersebut

4. **JWT Error**
   - Pastikan `JWT_SECRET_KEY` sudah diatur
   - Periksa format token di header Authorization

## 📮 Import ke Postman
Klik tombol di bawah untuk langsung import collection ke Postman:

[![Import ke Postman](https://run.pstmn.io/button.svg)](https://raw.githubusercontent.com/riskykurniawan15/payrolls/refs/heads/main/payrolls.postman_collection.json)
