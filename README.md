# Recything API Documentation

Recything adalah platform yang bertujuan untuk mendorong pelaporan dan pengelolaan sampah, serta meningkatkan kesadaran lingkungan. Berikut adalah dokumentasi API untuk aplikasi Recything.

## Base URL
```
https://www.recythingtech.my.id/api/v1
```

---

## Endpoints Overview

| Endpoint                                      | Method | Description                              |
|----------------------------------------------|--------|------------------------------------------|
| `/register`                                  | POST   | Register a new user                      |
| `/login`                                     | POST   | Login for admin or user                  |
| `/logout`                                    | GET    | Logout admin or user                     |
| `/users`                                     | PUT    | Update profile photo                    |
| `/user/data/:iduser`                         | PUT    | Update personal data                     |
| `/users/points`                              | GET    | Get user points                          |
| `/admin/users/points`                        | GET    | Get all users' points                    |
| `/report-rubbish`                            | POST   | Add a rubbish report                     |
| `/admin/report-rubbish`                      | GET    | Get all rubbish reports (with filters)   |
| `/admin/report-rubbish/:id`                  | GET    | Get a specific rubbish report by ID      |
| `/admin/report-rubbish/:id`                  | DELETE | Delete a rubbish report                  |
| `/admin/reports/statistics`                  | GET    | Get report statistics                    |
| `/admin/articles`                            | POST   | Add a new article                        |
| `/admin/articles/:id`                        | PUT    | Update an existing article               |
| `/admin/articles/:id`                        | DELETE | Delete an article                        |
| `/articles`                                  | GET    | Get all articles                         |
| `/articles/:id`                              | GET    | Get a specific article by ID             |
| `/report-rubbish/history`                    | GET    | Get user's rubbish report history        |

---

## Endpoint Details

### 1. Register
**Method:** `POST`
- **URL:** `/register`
- **Request Body:**
  ```json
  {
    "nama_lengkap": "string",
    "email": "string",
    "tanggal_lahir": "string",
    "no_telepon": "string",
    "password": "string",
    "photo": "string"
  }
  ```
- **Response:**
  Informasi status registrasi pengguna.

---

### 2. Login
**Method:** `POST`
- **URL:** `/login`
- **Request Body:**
  ```json
  {
    "email": "string",
    "password": "string"
  }
  ```
- **Response:**
  Token untuk autentikasi.

---

### 3. Logout
**Method:** `GET`
- **URL:** `/logout`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Response: {
    "meta": {
        "message": "Login successful",
        "code": 200,
        "status": "success"
    },
    "data": {
        "id_user": 1,
        "nama_lengkap": "Admin Recything",
        "tanggal_lahir": "2003-03-20",
        "no_telepon": "085357549320",
        "email": "admin@gmail.com",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiQWRtaW4gUmVjeXRoaW5nIiwidXNlcklEIjoxLCJyb2xlIjoiYWRtaW4iLCJleHAiOjE3MzQyODAxNTB9.B-78LOPIbRlBMZXXH1b5UqXKeMe2emdVKRoh9FuNvgk",
        "role": "admin",
        "photo": ""
    }
}**
  Status logout.

---

### 4. Update Foto Profil
**Method:** `PUT`
- **URL:** `/users`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Request Body:** (Form-Data)
  - `photo`: file
- **Response:**
  Status pembaruan foto profil.

---

### 5. Update Data Diri
**Method:** `PUT`
- **URL:** `/user/data/:iduser`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Request Body:** (Form-Data)
  ```json
  {
    "nama_lengkap": "string",
    "email": "string",
    "tanggal_lahir": "string",
    "no_telepon": "string",
    "old_password": "string",
    "new_password": "string"
  }
  ```
- **Response:**
  Status pembaruan data diri.

---

### 6. Get Point User
**Method:** `GET`
- **URL:** `/users/points`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Response:**
  Poin pengguna.

---

### 7. Get All Point User
**Method:** `GET`
- **URL:** `/admin/users/points`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Response:**
  Daftar poin seluruh pengguna.

---

### 8. Add Report Rubbish
**Method:** `POST`
- **URL:** `/report-rubbish`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Request Body:** (Form-Data)
  ```json
  {
    "location": "string",
    "description": "string",
    "photo": "file",
    "tanggal_laporan": "string",
    "category": "report_rubbish/report_littering"
  }
  ```
- **Response:**
  Status laporan sampah.

---

### 9. Get All Report Rubbish
**Method:** `GET`
- **URL:** `/admin/report-rubbish`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Query Parameters:**
  - `page`: int (optional)
  - `sort`: `asc` or `desc` (optional)
  - `status`: `process`, `rejected`, `completed` (optional)
- **Response:**
  Daftar laporan sampah dengan paginasi dan filter.

---

### 10. Delete Report
**Method:** `DELETE`
- **URL:** `/admin/report-rubbish/:id`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Response:**
  Status penghapusan laporan.

---

### 11. Add Article
**Method:** `POST`
- **URL:** `/admin/articles`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Request Body:** (JSON)
  ```json
  {
    "judul": "string",
    "author": "string",
    "konten": "string",
    "link_foto": "string",
    "link_video": "string" // optional
  }
  ```
- **Response:**
  Status penambahan artikel.

---

### 12. Statistik
**Method:** `GET`
- **URL:** `/admin/reports/statistics`
- **Headers:**
  ```
  Authorization: Bearer <token>
  ```
- **Response:**
  Statistik laporan sampah.

---

## Authentication
Semua endpoint (kecuali login dan register) memerlukan header berikut:
```plaintext
Authorization: Bearer <token>
```

## Catatan Penting
- Gunakan endpoint sesuai dengan peran Anda (admin/user).
- Pastikan data yang dikirim sesuai dengan tipe data yang diminta untuk menghindari error.

## Kontak
Jika Anda memiliki pertanyaan lebih lanjut, hubungi kami di [support@recythingtech.my.id](mailto:support@recythingtech.my.id).
