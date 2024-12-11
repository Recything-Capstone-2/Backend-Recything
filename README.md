# Recything API Documentation

Dokumentasi lengkap API untuk platform Recything.

---

## Registrasi Pengguna

### Endpoint
`POST https://www.recythingtech.my.id/api/v1/register`

### Request Body
| Key             | Tipe    | Keterangan               |
|------------------|---------|--------------------------|
| `nama_lengkap`  | string  | Nama lengkap pengguna.   |
| `email`         | string  | Email pengguna.          |
| `tanggal_lahir` | string  | Tanggal lahir pengguna.  |
| `no_telepon`    | string  | Nomor telepon pengguna.  |
| `password`      | string  | Kata sandi pengguna.     |
| `photo`         | file    | Foto profil pengguna.    |

---

## Login

### Login Admin
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/login`

### Request Body
| Key         | Tipe    | Keterangan                 |
|-------------|---------|----------------------------|
| `email`     | string  | Email admin: `"admin@gmail.com"`. |
| `password`  | string  | Kata sandi admin: `"recything_passwordnya"`. |

---

### Login User
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/login`

### Request Body
| Key         | Tipe    | Keterangan         |
|-------------|---------|--------------------|
| `email`     | string  | Email pengguna.    |
| `password`  | string  | Kata sandi pengguna. |

---

## Logout (User/Admin)

### Endpoint
`GET https://www.recythingtech.my.id/api/v1/logout`

### Authorization
- Bearer token dari sesi login.

---

## Update Foto Profil (User/Admin)

### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/users`

### Request Body
| Key     | Tipe  | Keterangan        |
|---------|-------|-------------------|
| `photo` | file  | Foto profil baru. |

### Authorization
- Bearer token dari sesi login.

---

## Update Data Diri (User/Admin)

### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/user/data/:iduser`

### Request Body
| Key             | Tipe    | Keterangan                                  |
|------------------|---------|---------------------------------------------|
| `nama_lengkap`  | string  | Nama lengkap pengguna.                      |
| `email`         | string  | Email pengguna.                             |
| `tanggal_lahir` | string  | Tanggal lahir pengguna.                     |
| `no_telepon`    | string  | Nomor telepon pengguna.                     |
| `old_password`  | string  | Kata sandi lama (opsional jika tidak diubah). |
| `new_password`  | string  | Kata sandi baru (opsional jika tidak diubah). |

### Authorization
- Bearer token dari sesi login.

---

## Mendapatkan Poin

### Mendapatkan Poin Pengguna (User)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/users/points`

### Authorization
- Bearer token dari sesi login.

### Response
- **200 OK**: Data poin pengguna berhasil didapatkan.

---

### Mendapatkan Semua Poin Pengguna (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/users/points`

### Authorization
- Bearer token admin.

---

## Pengurangan Poin Pengguna (Admin)

### Endpoint
`POST https://www.recythingtech.my.id/api/v1/admin/users/points/deduct`

### Request Body
| Key       | Tipe | Keterangan                     |
|-----------|------|--------------------------------|
| `user_id` | int  | ID pengguna yang poinnya dikurangi. |
| `points`  | int  | Jumlah poin yang dikurangi.    |

### Authorization
- Bearer token admin.

---

## Laporan Sampah

### Menambahkan Laporan (User)
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/report-rubbish`

### Request Body
| Key              | Tipe    | Keterangan                          |
|-------------------|---------|-------------------------------------|
| `location`       | string  | Lokasi laporan.                     |
| `description`    | string  | Deskripsi laporan.                  |
| `photo`          | file    | Foto terkait laporan.               |
| `tanggal_laporan`| date    | Tanggal laporan (format `YYYY-MM-DD`). |
| `category`       | string  | Kategori: `report_rubbish` atau `report_littering`. |

### Authorization
- Bearer token dari sesi login.

---

### Mendapatkan Semua Laporan (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/report-rubbish`

#### Parameter Query
| Key   | Tipe    | Keterangan                     |
|-------|---------|--------------------------------|
| `sort`| string  | Urutan laporan: `asc` atau `desc`. |
| `page`| int     | Nomor halaman untuk pagination.|

### Authorization
- Bearer token admin.

---

### Mendapatkan 10 Laporan Terbaru (Admin)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/admin/latest-report`

### Authorization
- Bearer token admin.

---

### Update Status Laporan (Admin)
#### Endpoint
`PUT https://www.recythingtech.my.id/api/v1/report-rubbish/:idreport`

### Request Body
| Key     | Tipe    | Keterangan                             |
|---------|---------|----------------------------------------|
| `status`| string  | Status laporan: `approved`, `rejected`, atau `completed`. |

### Authorization
- Bearer token admin.

---

## Artikel

### Menambahkan Artikel (Admin)
#### Endpoint
`POST https://www.recythingtech.my.id/api/v1/admin/articles`

### Request Body
| Key           | Tipe    | Keterangan                             |
|---------------|---------|----------------------------------------|
| `judul`       | string  | Judul artikel.                        |
| `author`      | string  | Penulis artikel.                      |
| `konten`      | string  | Isi artikel.                          |
| `link_foto`   | string  | URL untuk foto artikel.               |
| `link_video`  | string  | URL untuk video artikel *(opsional)*. |

### Authorization
- Bearer token admin.

---

### Mendapatkan Semua Artikel (User)
#### Endpoint
`GET https://www.recythingtech.my.id/api/v1/articles`

### Authorization
- Bearer token user.

---

