package constant

import "errors"

const Unauthorized = "Unauthorized"
const InternalServerError = "Internal Server Error"
const BadInput = "Format data not valid"

var ErrBadRequest = errors.New("bad request")
var ErrUnauthorized = errors.New("unauthorized")

var ErrEmptyOtp = errors.New("otp cannot be empty")
var ErrDataNotfound = errors.New("data kosong")
var ErrGetData = errors.New("gagal saat mengambil data")
var ErrEmptyId = errors.New("id cannot be empty")
var ErrFailedDelete = errors.New("tidak ada data yang terhapus")

// JWT
var ErrGenerateJWT = errors.New("failed to generate jwt token")
var ErrValidateJWT = errors.New("failed to validate jwt token")

// Update
var ErrUpdate = errors.New("update satu data minimal")

// Validator
var ErrHashPassword = errors.New("failed to hash password")

// Empty Register
var ErrEmptyEmailRegister = errors.New("email cannot be empty")
var ErrEmptyNameRegister = errors.New("username cannot be empty")
var ErrEmptyPasswordRegister = errors.New("password cannot be empty")
var ErrPasswordNotMatch = errors.New("password not match")

// Register Format Not Valid
var ErrInvalidEmail = errors.New("email is not valid")
var ErrInvalidUsername = errors.New("username formating not valid")
var ErrInvalidPhone = errors.New("phone formating not valid")

// Login
var ErrEmptyLogin = errors.New("email or Password cannot be empty")
var ErrUserNotFound = errors.New("user not found")
var ErrLenPassword = errors.New("password must be at least 8 characters")
var ErrInvalidPassword = errors.New("password must contain at least 1 number, 1 uppercase letter, one punctuation symbol and 1 lowercase letter")
var ErrEmptyPasswordLogin = errors.New("password cannot be empty")

// Admin
var ErrAdminNotFound = errors.New("data admin tidak ada")
var ErrAdminUserNameEmpty = errors.New("username tidak boleh kosong")
var ErrAdminPasswordEmpty = errors.New("password tidak boleh kosong")
var ErrEmptyGender = errors.New("gender tidak boleh kosong")
var ErrGenderChoice = errors.New("pilih gender antara laki-laki atau perempuan")

// Instruktur
var ErrInstrukturNotFound = errors.New("data instruktur tidak ada")
var ErrGetInstruktur = errors.New("error saat mengambil data instruktur")
var ErrInstrukturID = errors.New("error id instruktor tidak ada")
var ErrInstrukturNIKEmpty = errors.New("error NIK tidak boleh kosong")
var ErrInstrukturNIPEmpty = errors.New("error NIP tidak boleh kosong")
var ErrInstrukturImageEmpty = errors.New("error image instruktor tidak boleh kosong")

// post

var ErrEmptyNameInstuktor = errors.New("name cannot be empty")
var ErrEmptyEmailInstuktor = errors.New("email cannot be empty")
var ErrEmptyAlamatInstuktor = errors.New("alamat cannot be empty")
var ErrEmptyNumbertelponInstuktor = errors.New("number telpon cannot be empty")
var ErrEmptyDescriptionInstuktor = errors.New("description cannot be empty")

// Kategori
var ErrKategoriNotFound = errors.New("kategori tidak ditemukan")

// post
var ErrEmptyNamaKategori = errors.New("nama kategori tidak boleh kosong")
var ErrEmptyImageUrlKategori = errors.New("gambar kategori tidak boleh kosong")
var ErrEmptyDeskripsiKategori = errors.New("deskripsi kategori tidak boleh kosong")

// Kursus
var ErrKursusNotFound = errors.New("data kursus tidak ada")
var ErrImageKursusNotFound = errors.New("gambar kursus tidak ada")
var ErrKategoriKursusNotFound = errors.New("kategori kursus tidak ada")
var ErrMateriKursusNotFound = errors.New("materi kursus tidak ada")
var ErrJadwal = errors.New("jadwal harus diisi")
var ErrJadwalFormat = errors.New("format jadwal salah")
var ErrGambarKursus = errors.New("gambar tidak boleh kosong")
var ErrKategoriKursus = errors.New("kategori tidak boleh kosong")
var ErrMateriPembelajaran = errors.New("materi tidak boleh kosong")
var ErrKursusNotfound = errors.New("kursus tidak ditemukan")

// GCS
var ErrOpeningFile = errors.New("gagal membuka file")
var ErrUploadGCS = errors.New("gagal upload gambar")

// ID
var ErrGetID = errors.New("id tidak ditemukan")

// Voucher
var ErrVoucherNotFound = errors.New("data voucher tidak ada")
var ErrVoucherIDNotFound = errors.New("id voucher tidak ditemukan")
var ErrVoucherUsed = errors.New("voucher telah digunakan")

// Voucher data
var ErrVoucherFailedCreate = errors.New("data voucher gagal dibuat")

// Transaksi
var ErrTransaksiNotFound = errors.New("data transaksi belum ada")
var ErrValidateDokumenUser = errors.New("user belum mengisi seluruh dokumen penting")
