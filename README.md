# Proyek Skripsi Lembaga Kursus
Repository ini membahas bagian backend dari proyek Lembaga Kursus Komputer. Proyek Ini menggunakan Bahasa Pemrograman Go dengan penggunaan kerangka kerja Echo dan ORM berupa GORM. Tujuan proyek ini membuat proses pendaftaran online yang diintegrasikan dengan API pihak ketiga seperti midtrans, untuk pembayaran yang mudah dan aman.

This repository discusses the backend of the Computer Course Institute project. This project uses the Go programming language with the use of the Echo framework and ORM in the form of GORM. The purpose of this project is to create an online registration process that is integrated with third-party APIs such as midtrans, for easy and secure payments.
## Install Echo

```bash
  go get github.com/labstack/echo/v4
```

## Install ORM GORM
```bash
  go get -u gorm.io/gorm
```
## How to Use Http Request
Example how to register with post http request
```bash
curl -X POST "https://url-cloud/api/v1/register" \
     -H "Content-Type: application/json" \
     -d '{
           "username": "John Doe",
           "nomor_hp":"0xxxxxxxxx",
           "email": "johndoe@example.com",
           "password": "johndoePassword123!",
           "confirm_password": "johndoePassword123!"
         }'
```

## ERD
![ERD Skripsi drawio](https://github.com/user-attachments/assets/74e27ab3-ba3a-478e-98c1-210aa5f81ccc)


