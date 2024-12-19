package service

import (
	"skripsi/constant"
	"skripsi/features/kategori"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKategoriData struct {
	mock.Mock
}

func (m *MockKategoriData) GetAllKategori() ([]kategori.Kategori, error) {
	args := m.Called()
	return args.Get(0).([]kategori.Kategori), args.Error(1)
}

func (m *MockKategoriData) GetKategoriById(id string) (kategori.Kategori, error) {
	args := m.Called(id)
	return args.Get(0).(kategori.Kategori), args.Error(1)
}

func (m *MockKategoriData) CreateKategori(data kategori.Kategori) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockKategoriData) UpdateKategori(data kategori.Kategori) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockKategoriData) DeleteKategori(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKategoriData) GetKategoriWithPagination(page, limit int) ([]kategori.Kategori, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]kategori.Kategori), args.Int(1), args.Error(2)
}

func TestKategoriService(t *testing.T) {
	mockData := new(MockKategoriData)
	service := New(mockData, nil)

	t.Run("GetAllKategori success", func(t *testing.T) {
		mockResponse := []kategori.Kategori{{ID: "1", Nama: "Kategori A"}}
		mockData.On("GetAllKategori").Return(mockResponse, nil).Once()

		result, err := service.GetAllKategori()

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	t.Run("GetKategoriById success", func(t *testing.T) {
		mockResponse := kategori.Kategori{ID: "1", Nama: "Kategori A"}
		mockData.On("GetKategoriById", "1").Return(mockResponse, nil).Once()

		result, err := service.GetKategoriById("1")

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	// t.Run("GetKategoriById fail - empty ID", func(t *testing.T) {
	// 	_, err := service.GetKategoriById("")

	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, err.Error(), err)
	// })

	t.Run("CreateKategori success", func(t *testing.T) {
		mockRequest := kategori.Kategori{Nama: "Kategori A", Deskripsi: "Deskripsi A", ImageUrl: "image.jpg"}
		mockData.On("CreateKategori", mockRequest).Return(nil).Once()

		err := service.CreateKategori(mockRequest)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("CreateKategori fail - empty name", func(t *testing.T) {
		mockRequest := kategori.Kategori{Nama: "", Deskripsi: "Deskripsi A", ImageUrl: "image.jpg"}

		err := service.CreateKategori(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNamaKategori, err)
	})

	t.Run("UpdateKategori success", func(t *testing.T) {
		mockRequest := kategori.Kategori{ID: "1", Nama: "Kategori Updated"}
		mockData.On("UpdateKategori", mockRequest).Return(nil).Once()

		err := service.UpdateKategori(mockRequest)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("UpdateKategori fail - empty ID", func(t *testing.T) {
		mockRequest := kategori.Kategori{Nama: "Kategori Updated"}

		err := service.UpdateKategori(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})

	t.Run("DeleteKategori success", func(t *testing.T) {
		mockData.On("DeleteKategori", "1").Return(nil).Once()

		err := service.DeleteKategori("1")

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("DeleteKategori fail - empty ID", func(t *testing.T) {
		err := service.DeleteKategori("")

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})

	t.Run("GetKategoriWithPagination success", func(t *testing.T) {
		mockResponse := []kategori.Kategori{{ID: "1", Nama: "Kategori A"}}
		mockTotal := 10
		mockData.On("GetKategoriWithPagination", 1, 5).Return(mockResponse, mockTotal, nil).Once()

		result, total, err := service.GetKategoriWithPagination(1, 5)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		assert.Equal(t, mockTotal, total)
		mockData.AssertExpectations(t)
	})
}
