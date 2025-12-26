# 🧺 Laundry Marketplace Backend API

Backend REST API untuk aplikasi Laundry Marketplace menggunakan Golang, Gin, dan PostgreSQL.

## 🛠️ Tech Stack

- **Golang** 1.21+
- **Gin** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing

## 📋 Prerequisites

- Go 1.21 atau lebih baru
- PostgreSQL 12+
- Git

## 🚀 Installation

1. Clone repository:
```bash
git clone <repository-url>
cd laundry-go
```

2. Install dependencies:
```bash
go mod download
# atau
go mod tidy
```

**Note:** Jika ada error "missing go.sum entry", jalankan:
```bash
go mod tidy
go mod download
```

3. Setup database:
```bash
# Buat database PostgreSQL
createdb laundryhub

# Atau menggunakan psql
psql -U postgres -c "CREATE DATABASE laundryhub;"
```

4. Copy environment file:
```bash
cp env.example .env
```

Atau buat file `.env` di root directory dengan konfigurasi berikut:

**Untuk Supabase (Recommended):**
```env
# Server
PORT=8080
ENV=development

# Database - Supabase
DATABASE_URL=postgresql://postgres.mdhmtxtrqzbrqvuusank:yeDgrGt23k1qus4T@aws-1-ap-south-1.pooler.supabase.com:6543/postgres?pgbouncer=true

# JWT
JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRY=24h

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

**Atau untuk local PostgreSQL:**
```env
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=laundryhub
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# CORS
ALLOWED_ORIGINS=http://localhost:3000
```

**Note:** Aplikasi sekarang mendukung `DATABASE_URL` untuk Supabase, Heroku, atau cloud database lainnya. Jika `DATABASE_URL` diset, aplikasi akan menggunakannya secara otomatis.

6. (Optional) Seed dummy data untuk testing:
```bash
go run cmd/seed/main.go
```

Ini akan mengimport data dummy ke database:
- 3 users (user@laundryhub.com, admin@laundryhub.com, test@test.com)
- 6 laundries dengan services
- 3 sample orders

**Credentials untuk testing:**
- Email: `user@laundryhub.com`, Password: `password123`
- Email: `admin@laundryhub.com`, Password: `admin123`
- Email: `test@test.com`, Password: `test123`

## 🏃 Running the Server

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080` (atau sesuai PORT di .env).

## 📚 API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register user baru
- `POST /api/v1/auth/login` - Login user
- `GET /api/v1/auth/me` - Get current user (Protected)

### Laundries

- `GET /api/v1/laundries` - List semua laundry dengan pagination
- `GET /api/v1/laundries/:id` - Get detail laundry

### Orders

- `POST /api/v1/orders` - Create order baru (Protected)
- `GET /api/v1/orders` - List orders user (Protected)
- `GET /api/v1/orders/:id` - Get detail order (Protected)
- `PATCH /api/v1/orders/:id/cancel` - Cancel order (Protected)
- `PATCH /api/v1/orders/:id/status` - Update order status (Protected - Laundry Owner only)

## 🔐 Authentication

Semua endpoint yang protected memerlukan header:
```
Authorization: Bearer <token>
```

Token JWT akan expire dalam 24 jam (atau sesuai JWT_EXPIRY di .env).

## 📊 Database Schema

Database akan otomatis di-migrate saat pertama kali menjalankan aplikasi menggunakan GORM AutoMigrate.

Untuk manual migration, jalankan file SQL di folder `migrations/`:
```bash
psql -U postgres -d laundryhub -f migrations/001_init.sql
```

## 🧪 Testing

Untuk testing API, Anda bisa menggunakan:
- Postman
- cURL
- HTTP client lainnya

### Contoh Register:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "081234567890",
    "address": "Jl. Sudirman No. 123, Jakarta"
  }'
```

### Contoh Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

## 📁 Project Structure

```
laundry-go/
├── cmd/
│   └── server/
│       └── main.go          # Entry point aplikasi
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connection
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware (auth, CORS)
│   ├── models/              # Database models
│   ├── repository/          # Data access layer
│   ├── service/             # Business logic
│   └── utils/               # Utility functions
├── migrations/              # SQL migration files
├── .env.example            # Example environment file
├── go.mod                  # Go modules
└── README.md               # This file
```

## 🔧 Configuration

Semua konfigurasi dilakukan melalui environment variables di file `.env`:

- `PORT` - Server port (default: 8080)
- `ENV` - Environment (development/production)
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_USER` - Database user
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name
- `JWT_SECRET` - Secret key untuk JWT
- `JWT_EXPIRY` - JWT expiration time (default: 24h)
- `ALLOWED_ORIGINS` - CORS allowed origins (comma-separated)

## 📝 Notes

- Password minimum 8 karakter
- Email harus valid format
- Semua ID menggunakan UUID
- Response format konsisten dengan `success`, `message`, dan `data` fields
- Error handling mengikuti HTTP status codes standar

## 🐛 Troubleshooting

### Database connection error
- Pastikan PostgreSQL berjalan
- Periksa credentials di `.env`
- Pastikan database sudah dibuat

### Port already in use
- Ubah `PORT` di `.env`
- Atau hentikan proses yang menggunakan port tersebut

## 📄 License

MIT

