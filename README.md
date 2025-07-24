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
DB_TIMEZONE=Asia/Jakarta
DB_MAX_IDLE_CON=10
DB_MAX_OPEN_CON=100
DB_MAX_LIFE_TIME=10
DB_DEBUG=false

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
├── repositories/         # Data access layer
├── services/             # Business logic layer
├── utils/                # Utility functions
│   ├── env/             # Environment utilities
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
| `DB_TIMEZONE` | Timezone database | `Asia/Jakarta` |
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
- `POST /attendances/check-out` - Check out
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
- `GET /payslips` - List user payslips
- `POST /payslips/:id/generate` - Generate payslip
- `GET /payslips/print?token=xxx` - Print payslip
- `POST /payslips/:id/summary` - Generate payslip summary
- `GET /payslips/summary/print?token=xxx` - Print payslip summary

## 🔌 Contoh Penggunaan API (cURL)

### Authentication

#### Login User
```bash
curl -X POST http://localhost:9000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123"
  }'
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "roles": "admin"
    }
  }
}
```

### Health Check

#### Health Check
```bash
curl -X GET http://localhost:9000/health
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "message": "Service is healthy",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

### User Management

#### Get User Profile (Protected)
```bash
curl -X GET http://localhost:9000/user/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "roles": "admin",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

### Period Management (Admin only)

#### Create New Period
```bash
curl -X POST http://localhost:9000/periods \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Januari 2024",
    "start_date": "2024-01-01",
    "end_date": "2024-01-31"
  }'
```

**Response:**
```json
{
  "status": 201,
  "data": {
    "id": 1,
    "name": "Januari 2024",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T00:00:00Z",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

#### List All Periods
```bash
curl -X GET http://localhost:9000/periods \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Period by ID
```bash
curl -X GET http://localhost:9000/periods/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Update Period
```bash
curl -X PUT http://localhost:9000/periods/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Januari 2024 (Updated)",
    "start_date": "2024-01-01",
    "end_date": "2024-01-31"
  }'
```

#### Delete Period
```bash
curl -X DELETE http://localhost:9000/periods/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Run Payroll for Period
```bash
curl -X POST http://localhost:9000/periods/1/run-payroll \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "status": "Payroll processing started",
    "job_id": "payroll_1_abc123"
  }
}
```

### Attendance (Employee only)

#### Get User Attendances
```bash
curl -X GET http://localhost:9000/attendances \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Attendance by ID
```bash
curl -X GET http://localhost:9000/attendances/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Check In
```bash
curl -X POST http://localhost:9000/attendances/check-in \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-01",
    "check_in_time": "08:00:00"
  }'
```

**Response:**
```json
{
  "status": 201,
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-01T00:00:00Z",
    "check_in_time": "08:00:00",
    "check_out_time": null,
    "created_at": "2024-01-01T08:00:00Z"
  }
}
```

#### Check Out
```bash
curl -X POST http://localhost:9000/attendances/check-out \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2024-01-01",
    "check_out_time": "17:00:00"
  }'
```

#### Check Out by ID
```bash
curl -X POST http://localhost:9000/attendances/1/check-out \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "check_out_time": "17:00:00"
  }'
```

### Overtime (Employee only)

#### Get User Overtimes
```bash
curl -X GET http://localhost:9000/overtimes \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Overtime by ID
```bash
curl -X GET http://localhost:9000/overtimes/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create Overtime
```bash
curl -X POST http://localhost:9000/overtimes \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "overtimes_date": "2024-01-01",
    "start_time": "18:00:00",
    "end_time": "21:00:00",
    "description": "Lembur project deadline"
  }'
```

**Response:**
```json
{
  "status": 201,
  "data": {
    "id": 1,
    "user_id": 1,
    "overtimes_date": "2024-01-01T00:00:00Z",
    "start_time": "18:00:00",
    "end_time": "21:00:00",
    "total_hours_time": 3.0,
    "description": "Lembur project deadline",
    "created_at": "2024-01-01T18:00:00Z"
  }
}
```

#### Update Overtime
```bash
curl -X PUT http://localhost:9000/overtimes/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "start_time": "19:00:00",
    "end_time": "22:00:00",
    "description": "Lembur project deadline (updated)"
  }'
```

#### Delete Overtime
```bash
curl -X DELETE http://localhost:9000/overtimes/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Reimbursement (Employee only)

#### Get User Reimbursements
```bash
curl -X GET http://localhost:9000/reimbursements \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Get Reimbursement by ID
```bash
curl -X GET http://localhost:9000/reimbursements/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Create Reimbursement
```bash
curl -X POST http://localhost:9000/reimbursements \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Biaya transportasi meeting",
    "date": "2024-01-01",
    "amount": 50000,
    "description": "Biaya taksi untuk meeting dengan client"
  }'
```

**Response:**
```json
{
  "status": 201,
  "data": {
    "id": 1,
    "user_id": 1,
    "title": "Biaya transportasi meeting",
    "date": "2024-01-01T00:00:00Z",
    "amount": 50000,
    "description": "Biaya taksi untuk meeting dengan client",
    "status": 1,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

#### Update Reimbursement
```bash
curl -X PUT http://localhost:9000/reimbursements/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Biaya transportasi meeting (updated)",
    "amount": 75000,
    "description": "Biaya taksi untuk meeting dengan client (updated)"
  }'
```

#### Delete Reimbursement
```bash
curl -X DELETE http://localhost:9000/reimbursements/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Payslip (Employee only)

#### List User Payslips
```bash
curl -X GET http://localhost:9000/payslips \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "status": 200,
  "data": [
    {
      "id": 1,
      "periods_id": 1,
      "period_name": "Januari 2024",
      "start_date": "2024-01-01T00:00:00Z",
      "end_date": "2024-01-31T00:00:00Z",
      "take_home_pay": 2900000,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

#### Generate Payslip
```bash
curl -X POST http://localhost:9000/payslips/1/generate \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "print_url": "http://localhost:9000/payslips/print?token=abc123",
    "token": "abc123",
    "expires_at": "2024-01-02T12:00:00Z"
  }
}
```

#### Print Payslip
```bash
curl -X GET "http://localhost:9000/payslips/print?token=abc123"
```

**Response:** HTML content untuk print

#### Generate Payslip Summary
```bash
curl -X POST http://localhost:9000/payslips/1/summary \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "status": 200,
  "data": {
    "print_url": "http://localhost:9000/payslips/summary/print?token=def456",
    "token": "def456",
    "expires_at": "2024-01-02T12:00:00Z"
  }
}
```

#### Print Payslip Summary
```bash
curl -X GET "http://localhost:9000/payslips/summary/print?token=def456"
```

**Response:** HTML content untuk print summary

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
