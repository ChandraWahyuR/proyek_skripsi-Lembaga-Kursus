package service

import (
	"errors"
	"testing"
	"time"

	"skripsi/constant"
	"skripsi/features/kursus"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKursusData struct {
	mock.Mock
}

type MockJWT struct {
	mock.Mock
}

func (m *MockKursusData) GetAllKursus() ([]kursus.Kursus, error) {
	args := m.Called()
	return args.Get(0).([]kursus.Kursus), args.Error(1)
}

func (m *MockKursusData) GetAllKursusById(id string) (kursus.Kursus, error) {
	args := m.Called(id)
	return args.Get(0).(kursus.Kursus), args.Error(1)
}

func (m *MockKursusData) AddKursus(data kursus.Kursus) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockKursusData) GetKursusPagination(page, limit int) ([]kursus.Kursus, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]kursus.Kursus), args.Int(1), args.Error(2)
}

func (m *MockKursusData) UpdateKursus(data kursus.Kursus) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockKursusData) DeleteKursus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKursusData) DeleteImageKursus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKursusData) DeleteMateriKursus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockKursusData) DeleteKategoriKursus(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKursusData) GetAllKursusByName(name string, page, limit int) ([]kursus.Kursus, int, error) {
	args := m.Called(name, page, limit)
	return args.Get(0).([]kursus.Kursus), args.Int(1), args.Error(2)
}

func TestGetAllKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	expected := []kursus.Kursus{{ID: "1", Nama: "Kursus A"}, {ID: "2", Nama: "Kursus B"}}
	mockData.On("GetAllKursus").Return(expected, nil)

	result, err := service.GetAllKursus()

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	mockData.AssertExpectations(t)
}

func TestGetAllKursusById(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	t.Run("success", func(t *testing.T) {
		id := "1"
		expected := kursus.Kursus{ID: id, Nama: "Kursus A"}
		mockData.On("GetAllKursusById", id).Return(expected, nil).Once()

		result, err := service.GetAllKursusById(id)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - empty id", func(t *testing.T) {
		id := ""
		result, err := service.GetAllKursusById(id)

		assert.NotNil(t, err)
		assert.Equal(t, kursus.Kursus{}, result)
		assert.Equal(t, constant.ErrGetID, err)
	})
}

func TestAddKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	t.Run("success - valid jadwal time", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		data := kursus.Kursus{
			Nama:         "Kursus A",
			Deskripsi:    "Deskripsi Kursus",
			Harga:        150000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
			Image: []kursus.ImageKursus{
				{
					Url:      "url",
					Name:     "name",
					KursusID: "id-kursus",
					Position: 1,
				},
				{
					Url:      "url2",
					Name:     "name2",
					KursusID: "id-kursus",
					Position: 2,
				},
			},
			Kategori: []kursus.KategoriKursus{
				{
					KursusID:   "id-kursus",
					KategoriID: "id-kategori",
				},
			},
			MateriPembelajaran: []kursus.MateriPembelajaran{
				{
					KursusID:  "id-kursus",
					Deskripsi: "deskripsi",
				},
			},
		}

		// Mocking behavior untuk AddKursus
		mockData.On("AddKursus", data).Return(nil).Once()

		// Panggil fungsi AddKursus
		err := service.AddKursus(data)

		// Assert hasilnya
		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - empty name", func(t *testing.T) {
		data := kursus.Kursus{Nama: ""}

		err := service.AddKursus(data)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNameInstuktor, err)
	})
	t.Run(" failure - empty Deskripsi", func(t *testing.T) {
		mockCourse := kursus.Kursus{
			Nama: "Kursus A",
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrDekripsiKursus, err)
	})
	t.Run(" failure - empty Harga", func(t *testing.T) {
		mockCourse := kursus.Kursus{
			Nama:      "Kursus A",
			Deskripsi: "desc",
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrHargaKursus, err)
	})
	t.Run(" failure - empty instruktur id", func(t *testing.T) {
		mockCourse := kursus.Kursus{
			Nama:      "Kursus A",
			Deskripsi: "desc",
			Harga:     10000,
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInstrukturID, err)
	})
	t.Run(" failure - empty jadwal", func(t *testing.T) {
		mockCourse := kursus.Kursus{
			Nama:         "Kursus A",
			Deskripsi:    "desc",
			Harga:        10000,
			InstrukturID: "1234",
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrJadwal, err)
	})
	t.Run(" failure - empty Image", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		mockCourse := kursus.Kursus{
			Nama:         "Kursus A",
			Deskripsi:    "desc",
			Harga:        10000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrGambarKursus, err)
	})

	t.Run(" failure - empty Kategori", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		mockCourse := kursus.Kursus{
			Nama:         "Kursus A",
			Deskripsi:    "desc",
			Harga:        10000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
			Image: []kursus.ImageKursus{
				{
					Url:      "url",
					Name:     "name",
					KursusID: "id-kursus",
					Position: 1,
				},
				{
					Url:      "url2",
					Name:     "name2",
					KursusID: "id-kursus",
					Position: 2,
				},
			},
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrKategoriKursus, err)
	})
	t.Run(" failure - empty materi pembelajaran", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		mockCourse := kursus.Kursus{
			Nama:         "Kursus A",
			Deskripsi:    "desc",
			Harga:        10000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
			Image: []kursus.ImageKursus{
				{
					Url:      "url",
					Name:     "name",
					KursusID: "id-kursus",
					Position: 1,
				},
				{
					Url:      "url2",
					Name:     "name2",
					KursusID: "id-kursus",
					Position: 2,
				},
			},
			Kategori: []kursus.KategoriKursus{
				{
					KursusID:   "id-kursus",
					KategoriID: "id-kategori",
				},
			},
		}

		err := service.AddKursus(mockCourse)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrMateriPembelajaran, err)
	})

}

func TestGetKursusPagination(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	// Mocking data dan return value untuk GetKursusPagination
	mockPage := 2
	mockLimit := 2
	expectedKursus := []kursus.Kursus{
		{ID: "1", Nama: "Kursus A"},
		{ID: "2", Nama: "Kursus B"},
	}
	mockTotal := 4
	mockData.On("GetKursusPagination", mockPage, mockLimit).Return(expectedKursus, mockTotal, nil)

	result, total, err := service.GetKursusPagination(mockPage, mockLimit)

	assert.Nil(t, err)
	assert.Equal(t, expectedKursus, result)
	assert.Equal(t, mockTotal, total)

	mockData.AssertExpectations(t)
}
func TestUpdateKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	t.Run("success", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		data := kursus.Kursus{
			ID:           "valid-id",
			Nama:         "Kursus A",
			Deskripsi:    "Deskripsi Kursus",
			Harga:        150000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
			Image: []kursus.ImageKursus{
				{
					Url:      "url",
					Name:     "name",
					KursusID: "valid-id",
					Position: 1,
				},
				{
					Url:      "url2",
					Name:     "name2",
					KursusID: "valid-id",
					Position: 2,
				},
			},
		}

		// Mock behavior untuk UpdateKursus
		mockData.On("UpdateKursus", data).Return(nil).Once()

		// Panggil fungsi UpdateKursus
		err := service.UpdateKursus(data)

		// Assert hasilnya
		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - get id", func(t *testing.T) {
		validJamMulai := time.Date(2024, time.December, 18, 11, 0, 0, 0, time.UTC)
		validJamSelesai := time.Date(2024, time.December, 18, 13, 0, 0, 0, time.UTC)

		data := kursus.Kursus{
			ID:           "", // Kondisi ID kosong
			Nama:         "Kursus A",
			Deskripsi:    "Deskripsi Kursus",
			Harga:        150000,
			InstrukturID: "1234",
			Jadwal: []kursus.JadwalKursus{
				{
					Hari:       "Senin",
					JamMulai:   validJamMulai,
					JamSelesai: validJamSelesai,
				},
			},
			Image: []kursus.ImageKursus{
				{
					Url:      "url",
					Name:     "name",
					KursusID: "id-kursus",
					Position: 1,
				},
				{
					Url:      "url2",
					Name:     "name2",
					KursusID: "id-kursus",
					Position: 2,
				},
			},
		}

		// Panggil fungsi UpdateKursus
		err := service.UpdateKursus(data)

		// Assert error sesuai dengan constant.ErrEmptyId
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})
	t.Run("fail - update field empty", func(t *testing.T) {
		data := kursus.Kursus{
			ID:   "valid-id", // Kondisi ID kosong
			Nama: "",
		}

		// Panggil fungsi UpdateKursus
		err := service.UpdateKursus(data)

		// Assert error sesuai dengan constant.ErrEmptyId
		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrUpdate, err)
	})
}

func TestDeleteKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	t.Run("success", func(t *testing.T) {
		id := "1"
		mockData.On("DeleteKursus", id).Return(nil).Once()

		err := service.DeleteKursus(id)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("fail - empty id", func(t *testing.T) {
		id := ""

		err := service.DeleteKursus(id)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)

		mockData.AssertExpectations(t)
	})
}
func TestDeleteImageKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	id := "1"
	mockData.On("DeleteImageKursus", id).Return(nil).Once()

	err := service.DeleteImageKursus(id)

	assert.Nil(t, err)
	mockData.AssertExpectations(t)

}

func TestDeleteMateriKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	id := "1"
	mockData.On("DeleteMateriKursus", id).Return(nil).Once()

	err := service.DeleteMateriKursus(id)

	assert.Nil(t, err)
	mockData.AssertExpectations(t)

}
func TestDeleteKategoriKursus(t *testing.T) {
	mockData := new(MockKursusData)
	service := New(mockData, nil)

	id := "1"
	mockData.On("DeleteKategoriKursus", id).Return(nil).Once()

	err := service.DeleteKategoriKursus(id)

	assert.Nil(t, err)
	mockData.AssertExpectations(t)

}

func TestGetAllKursusByName(t *testing.T) {
	mockData := new(MockKursusData) // Mock data repository
	service := New(mockData, nil)   // Inisialisasi service dengan mock data

	t.Run("success - data found", func(t *testing.T) {
		// Mock input dan output
		mockName := "Kursus A"
		mockPage := 1
		mockLimit := 2
		mockKursus := []kursus.Kursus{
			{
				ID:           "1",
				Nama:         "Kursus A",
				Deskripsi:    "Deskripsi Kursus A",
				Harga:        200000,
				InstrukturID: "1234",
			},
			{
				ID:           "2",
				Nama:         "Kursus A - Lanjutan",
				Deskripsi:    "Deskripsi Kursus A Lanjutan",
				Harga:        250000,
				InstrukturID: "5678",
			},
		}
		mockTotal := 2

		// Mocking fungsi GetAllKursusByName
		mockData.On("GetAllKursusByName", mockName, mockPage, mockLimit).Return(mockKursus, mockTotal, nil).Once()

		// Panggil fungsi yang akan diuji
		result, total, err := service.GetAllKursusByName(mockName, mockPage, mockLimit)

		// Validasi hasil
		assert.Nil(t, err)
		assert.Equal(t, mockTotal, total)
		assert.Equal(t, mockKursus, result)

		// Verifikasi bahwa mock dipanggil sesuai ekspektasi
		mockData.AssertExpectations(t)
	})

	t.Run("success - empty result", func(t *testing.T) {
		// Mock input dan output
		mockName := "Kursus Tidak Ada"
		mockPage := 1
		mockLimit := 2
		mockKursus := []kursus.Kursus{}
		mockTotal := 0

		// Mocking fungsi GetAllKursusByName
		mockData.On("GetAllKursusByName", mockName, mockPage, mockLimit).Return(mockKursus, mockTotal, nil).Once()

		// Panggil fungsi yang akan diuji
		result, total, err := service.GetAllKursusByName(mockName, mockPage, mockLimit)

		// Validasi hasil
		assert.Nil(t, err)
		assert.Equal(t, mockTotal, total)
		assert.Equal(t, mockKursus, result)

		// Verifikasi bahwa mock dipanggil sesuai ekspektasi
		mockData.AssertExpectations(t)
	})

	t.Run("failure - error from repository", func(t *testing.T) {
		// Mock input dan output
		mockName := "Error Name"
		mockPage := 1
		mockLimit := 2
		mockError := errors.New("internal error")

		mockData.On("GetAllKursusByName", mockName, mockPage, mockLimit).Return([]kursus.Kursus{}, 0, mockError).Once()

		result, total, err := service.GetAllKursusByName(mockName, mockPage, mockLimit)

		assert.NotNil(t, err)
		assert.Equal(t, mockError, err)
		assert.Equal(t, 0, total)
		assert.Nil(t, result)

		// Verifikasi bahwa mock dipanggil sesuai ekspektasi
		mockData.AssertExpectations(t)
	})
}
