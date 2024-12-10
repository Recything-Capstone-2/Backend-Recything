# Recything API Documentation

Dokumentasi lengkap API untuk platform Recything.

---

## Registrasi Pengguna

### Endpoint
`POST https://www.recythingtech.my.id/api/v1/register`

### Request Body (JSON)
- `nama_lengkap` (string)
- `email` (string)
- `tanggal_lahir` (string)
- `no_telepon` (string)
- `password` (string)
- `photo` (file)

### Response
- **201 Created**: Registrasi berhasil.
- **400 Bad Request**: Data tidak valid.

---

## Login

### Login Admin
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/login`

#### Request Body (JSON)
- `email` (string): `"admin@gmail.com"`
- `password` (string): `"recything_passwordnya"`

### Login User
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/login`

#### Request Body (JSON)
- `email` (string)
- `password` (string)

### Response
- **200 OK**: Login berhasil, mengembalikan token.
- **401 Unauthorized**: Kredensial salah.

---

## Logout (User/Admin)

### Endpoint
`GET https://www.recythingtech.my.id/api/v1/logout`

### Authorization
- Bearer token dari sesi login.

### Response
- **200 OK**: Logout berhasil.

---

## Update Foto Profil (User/Admin)

### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/users`

### Request Body (Form-Data)
- `photo` (file)

### Authorization
- Bearer token dari sesi login.

### Response
- **200 OK**: Foto profil berhasil diperbarui.

---

## Update Data Diri (User/Admin)

### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/user/data/:iduser`

### Request Body (Form-Data)
- `nama_lengkap` (string)
- `email` (string)
- `tanggal_lahir` (string)
- `no_telepon` (string)
- `old_password` (string)
- `new_password` (string)

### Catatan
- Update tetap berhasil meskipun `password` atau `email` tidak diubah.

### Authorization
- Bearer token dari sesi login.

### Response
- **200 OK**: Data berhasil diperbarui.

---

## Mendapatkan Poin

### Mendapatkan Poin Pengguna (User)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/users/points`

### Mendapatkan Semua Poin Pengguna (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/users/points`

### Authorization
- Bearer token dari sesi login.

### Response
- **200 OK**: Data poin berhasil didapatkan.

---

## Pengurangan Poin Pengguna (Admin)

### Endpoint
`POST https://www.recythingtech.my.id/api/v1/admin/users/points/deduct`

### Request Body (JSON)
- `user_id` (int)
- `points` (int): Jumlah poin yang dikurangi.

### Authorization
- Bearer token admin.

### Response
- **200 OK**: Poin berhasil dikurangi.

---

## Laporan Sampah

### Menambahkan Laporan (User)
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/report-rubbish`

#### Request Body (Form-Data)
- `location` (string)
- `description` (string)
- `photo` (file)
- `tanggal_laporan` (date): Format `YYYY-MM-DD`
- `category` (string): `report_rubbish` atau `report_littering`.

#### Authorization
- Bearer token dari sesi login.

#### Response
- **201 Created**: Laporan berhasil ditambahkan.

### Mendapatkan Laporan (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/report-rubbish?page=1`

#### Filtering dan Pagination
- Sort Descending: `?sort=desc`
- Sort Ascending: `?sort=asc`
- Pagination: `?page=1`

#### Authorization
- Bearer token admin.

#### Response
- **200 OK**: Data laporan berhasil didapatkan.

### Mendapatkan 10 Laporan Terbaru (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/latest-report`

#### Authorization
- Bearer token admin.

#### Response
- **200 OK**: Data laporan terbaru berhasil didapatkan.

### Update Status Laporan (Admin)
#### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/report-rubbish/:idreport`

#### Request Body (JSON)
- `status` (string): `approved`, `rejected`, atau `completed`.

#### Authorization
- Bearer token admin.

#### Response
- **200 OK**: Status laporan berhasil diperbarui.

---

## Artikel

### Menambahkan Artikel (Admin)
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/admin/articles`

#### Request Body (JSON)
- `judul` (string)
- `author` (string)
- `konten` (string)
- `link_foto` (string)
- `link_video` (string) *(opsional)*

#### Authorization
- Bearer token admin.

#### Response
- **201 Created**: Artikel berhasil ditambahkan.

### Mendapatkan Semua Artikel (User)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/articles`

#### Authorization
- Bearer token user.

#### Response
- **200 OK**: Data artikel berhasil didapatkan.

---

Jika ada bagian yang kurang jelas atau Anda ingin memperbarui informasi lebih lanjut, beri tahu saya!
