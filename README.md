# macha-auth-service

Auth service untuk autentikasi (login), otorisasi, dan manajemen user via HTTP + gRPC.

## Tujuan Service
Service ini berfungsi sebagai pusat autentikasi dan otorisasi untuk aplikasi, termasuk:
- Verifikasi kredensial pengguna dan penerbitan token.
- Validasi token untuk akses endpoint.
- Manajemen user (tambah, ubah, nonaktifkan, dan list).

## Konfigurasi
Template: [cmd/config/config.example.yaml](cmd/config/config.example.yaml)

Copy dulu:

```bash
cp cmd/config/config.example.yaml cmd/config/config.yaml
```

Contoh konfigurasi PostgreSQL:
- database.host: localhost
- database.port: 5432
- database.user: admin
- database.password: password123
- database.dbname: appdb
- appport: :8011
- grpcport: :50053

## Menjalankan Service

```bash
go run ./cmd/main.go -config cmd/config/config.yaml -log.file logs
```

## Endpoint HTTP

Base URL: http://localhost:8011

### GET /ping
Response: `Pong!`

### POST /login
Body:
```json
{"username":"loginuser","password":"login123"}
```
Response sukses:
```json
{"code":0,"message":"Success","data":{"token":"<token>"}}
```

### GET /user
Header: `Token: <token>`
Response sukses:
```json
{"code":0,"message":"Success","data":[{"id":1,"username":"admin","fullName":"Admin User","role":1}]}
```

### POST /user
Header: `Token: <admin-token>`
Body:
```json
{"username":"newuser","password":"newpass","fullName":"New User","role":2}
```

### PUT /user
Header: `Token: <admin-token>`
Body:
```json
{"user_id":2,"username":"staff_edit","password":"staff123","fullName":"Staff Edit","role":2}
```

### DELETE /user/{userID}
Header: `Token: <admin-token>`

### DELETE /logout
Header: `Token: <token>`

## gRPC
Service gRPC berjalan pada `grpcport` di config.

## Environment Variables
Jika tidak memakai file config, bisa pakai env dengan prefix `AUTH_`:
- `AUTH_DATABASE_HOST`
- `AUTH_DATABASE_PORT`
- `AUTH_DATABASE_USER`
- `AUTH_DATABASE_PASSWORD`
- `AUTH_DATABASE_DBNAME`
- `AUTH_DATABASE_DEBUG`
- `AUTH_APPPORT`
- `AUTH_GRPCPORT`

## Kode Error
- 0: Success
- 101: User not found
- 102: Invalid Password
- 201: Invalid Format Request
- 202: Invalid Token
- 203: Invalid Request
- 211: User Not Active (atau not allowed)
- 212: User Not Active
- 213: Duplicate User
- 301: Error query database
- 302: Error Update database

## Hasil Test (ringkas)
Integration test dan curl test memvalidasi:
- login sukses/gagal (user tidak ada, password salah)
- add user sukses/duplicate/invalid token/non-admin/inactive/expired
- edit user sukses/invalid request/invalid token/non-admin/inactive
- deactivate user sukses/invalid id
- get all user sukses/invalid token/non-admin/expired
- logout sukses/invalid token