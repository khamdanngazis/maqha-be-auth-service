# Maqha Auth Service - Dokumentasi

## Daftar Isi
1. [Overview](#overview)
2. [Arsitektur](#arsitektur)
3. [Setup & Installation](#setup--installation)
4. [Database](#database)
5. [API Endpoints](#api-endpoints)
6. [Testing](#testing)
7. [Troubleshooting](#troubleshooting)

---

## Overview

**Maqha Auth Service** adalah microservice untuk autentikasi dan otorisasi pengguna dalam sistem Maqha. Service ini menyediakan REST API dan gRPC interface untuk login, validasi token, dan manajemen user.

### Tujuan Service
- Memverifikasi kredensial pengguna dan menerbitkan token.
- Memvalidasi token untuk akses endpoint.
- Menyediakan manajemen user (tambah, ubah, nonaktifkan, list).

### Fitur Utama
- Autentikasi berbasis token
- Manajemen user (Create, Update, Deactivate, List)
- Logging terstruktur
- Unit tests dan integration tests
- gRPC support untuk komunikasi inter-service

### Tech Stack
- **Language**: Go 1.22
- **Database**: PostgreSQL 16
- **API Framework**: Gorilla Mux (HTTP), gRPC
- **ORM**: GORM
- **Logging**: Logrus
- **Testing**: testify, standard Go testing

---

## Arsitektur

### Struktur Direktori
```
maqha-be-auth-service/
├── cmd/
│   ├── main.go                 # Entry point aplikasi
│   └── config/                 # Configuration files
│       ├── config.yaml         # Development config
│       ├── config-test.yaml    # Testing config
│       └── config-prod.yaml    # Production config
├── internal/
│   ├── app/
│   │   ├── entity/             # Domain models
│   │   │   ├── client.go
│   │   │   └── user.go
│   │   ├── model/              # Request/response models
│   │   │   ├── login_model.go
│   │   │   ├── response.go
│   │   │   └── user_model.go
│   │   ├── repository/         # Data access layer
│   │   │   └── user_repo.go
│   │   ├── service/            # Business logic
│   │   │   ├── auth_service.go
│   │   │   └── errors_utils.go
│   ├── config/                 # Configuration management
│   ├── database/               # Database connection
│   └── interface/              # External interfaces
│       ├── http/               # HTTP handlers & routing
│       │   ├── handler/
│       │   └── router/
│       └── grpc/               # gRPC handlers
├── doc/
│   └── migrations/             # Database migrations (SQL)
├── tests/                      # Test files
│   ├── integration/            # Integration tests
│   └── unit/                   # Unit tests (jika ada)
├── go.mod & go.sum             # Dependencies
└── README.md

```

### Alur Data
```
HTTP/gRPC Request
    ↓
Handler (HTTP/gRPC)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Database (PostgreSQL)
```

---

## Setup & Installation

### Prerequisites
- Go 1.22+
- PostgreSQL 16+
- Docker (opsional, untuk database)
- Git

### Clone Repository
```bash
git clone https://github.com/khamdanngazis/maqha-be-auth-service.git
cd maqha-be-auth-service
```

### Install Dependencies
```bash
go mod download
go mod tidy
```

### Setup Database

#### Menggunakan Docker
```bash
docker run -d \
  --name postgresql \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=password123 \
  -e POSTGRES_DB=appdb \
  -p 5432:5432 \
  -v postgresql_data:/var/lib/postgresql/data \
  postgres:16
```

#### Create Database & Run Migrations (SQL)
```bash
# Create test database
 docker exec postgresql createdb -U admin appdb_test

# Run migrations untuk development database
 docker exec -i postgresql psql -U admin -d appdb < doc/migrations/001_init_postgres.sql

# Run migrations untuk test database
 docker exec -i postgresql psql -U admin -d appdb_test < doc/migrations/001_init_postgres.sql
```

### Configuration
Copy template lalu edit file config:

```bash
cp cmd/config/config.example.yaml cmd/config/config.yaml
```

Edit [cmd/config/config.yaml](cmd/config/config.yaml):
```yaml
database:
  host: localhost
  port: 5432
  user: admin
  password: password123
  dbname: appdb
  debug: true
externalconnection:
  authservice:
    host: localhost:50051
appport: :8011
grpcport: :50053
imagepath: "../../public/images"
```

### Run Server
```bash
cd cmd
go run main.go -config config/config.yaml -log.file ../logs
```

Server akan listening di:
- HTTP: `http://localhost:8011`
- gRPC: `localhost:50053`

### Environment Variables
Jika tidak memakai file config, bisa pakai env dengan prefix `AUTH_`:
- `AUTH_DATABASE_HOST`
- `AUTH_DATABASE_PORT`
- `AUTH_DATABASE_USER`
- `AUTH_DATABASE_PASSWORD`
- `AUTH_DATABASE_DBNAME`
- `AUTH_DATABASE_DEBUG`
- `AUTH_APPPORT`
- `AUTH_GRPCPORT`

---

## Database

### Schema

#### Client Table
```sql
CREATE TABLE client (
    id BIGSERIAL PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(50) NOT NULL,
    address TEXT NOT NULL,
    owner_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    token VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

#### User Table
```sql
CREATE TABLE "user" (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role BIGINT NOT NULL,
    token TEXT NOT NULL,
    token_expired TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_client
        FOREIGN KEY (client_id)
        REFERENCES client (id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);
```

### Insert Sample Data
```sql
INSERT INTO client (company_name, email, phone_number, address, owner_name, is_active, token, created_at)
VALUES ('Test Coffee', 'info@testcoffee.com', '+1234567890', '123 Main St', 'John Doe', true, 'clienttoken', NOW());

INSERT INTO "user" (client_id, username, password, full_name, role, token, token_expired, is_active, created_at)
VALUES (1, 'admin', '<bcrypt_hash>', 'Admin User', 1, 'admintoken', NOW() + interval '1 day', TRUE, NOW());
```

---

## API Endpoints

### Authentication
Semua endpoint (kecuali `/ping`) memerlukan header:
```
Token: <token>
```

### 1. Ping
```
GET /ping
Response: Pong!
```

### 2. Login
```
POST /login
Body:
{
  "username": "loginuser",
  "password": "login123"
}

Response (200):
{
  "code": 0,
  "message": "Success",
  "data": { "token": "<token>" }
}
```

### 3. Get All Users
```
GET /user
Headers: Token: <token>

Response (200):
{
  "code": 0,
  "message": "Success",
  "data": [
    {"id":1,"username":"admin","fullName":"Admin User","role":1}
  ]
}
```

### 4. Add User
```
POST /user
Headers:
  - Token: <admin-token>
  - Content-Type: application/json

Body:
{
  "username": "newuser",
  "password": "newpass",
  "fullName": "New User",
  "role": 2
}
```

### 5. Edit User
```
PUT /user
Headers: Token: <admin-token>

Body:
{
  "user_id": 2,
  "username": "staff_edit",
  "password": "staff123",
  "fullName": "Staff Edit",
  "role": 2
}
```

### 6. Deactivate User
```
DELETE /user/{userID}
Headers: Token: <admin-token>
```

### 7. Logout
```
DELETE /logout
Headers: Token: <token>
```

---

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Test
```bash
# Integration tests
 go test ./tests/integration -v
```

### Test Results (Current Status)
```
ok      maqhaa/auth_service/tests/integration   4.292s
```

### Integration Tests
Test coverage mencakup:
- ✅ Login sukses/gagal (user tidak ada, password salah)
- ✅ Add user sukses/duplicate/invalid token/non-admin/inactive/expired
- ✅ Edit user sukses/invalid request/invalid token/non-admin/inactive
- ✅ Deactivate user sukses/invalid id
- ✅ Get all user sukses/invalid token/non-admin/expired
- ✅ Logout sukses/invalid token

---

## Troubleshooting

### Error: Database Connection Failed
**Solusi:**
1. Pastikan PostgreSQL running: `docker ps | grep postgresql`
2. Cek konfigurasi di `cmd/config/config.yaml`
3. Cek credentials: user=admin, password=password123

### Error: Log File Not Found
**Solusi:**
```bash
mkdir -p logs
```

### Error: Invalid Token
**Solusi:**
1. Pastikan token ada di database
2. Pastikan user dengan token tersebut ada di table `user`
3. Test dengan token yang valid

### First Deploy (DB kosong)
Gunakan SQL migration: [doc/migrations/001_init_postgres.sql](doc/migrations/001_init_postgres.sql)

### Error: Port Already in Use
**Solusi:**
```bash
# Kill process di port 8011
lsof -ti:8011 | xargs kill -9

# Atau ganti port di config
appport: :8012
```

---

## References

- [Go Documentation](https://golang.org/doc/)
- [GORM Documentation](https://gorm.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [gRPC Documentation](https://grpc.io/docs/)
- [Gorilla Mux Documentation](https://github.com/gorilla/mux)

---

## License
Proprietary - Maqha

## Support
Untuk pertanyaan atau bug report, hubungi tim development.

---

**Last Updated:** February 15, 2026
**Version:** 1.0.0
