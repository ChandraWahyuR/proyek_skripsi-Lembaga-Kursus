package service

import (
	"skripsi/constant"
	"skripsi/features/instruktur"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInstrukturData struct {
	mock.Mock
}

func (m *MockInstrukturData) GetInstrukturWithPagination(page, limit int) ([]instruktur.Instruktur, int, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]instruktur.Instruktur), args.Int(1), args.Error(2)
}

func (m *MockInstrukturData) GetAllInstruktur() ([]instruktur.Instruktur, error) {
	args := m.Called()
	return args.Get(0).([]instruktur.Instruktur), args.Error(1)
}

func (m *MockInstrukturData) GetAllInstrukturByID(id string) (instruktur.Instruktur, error) {
	args := m.Called(id)
	return args.Get(0).(instruktur.Instruktur), args.Error(1)
}

func (m *MockInstrukturData) PostInstruktur(data instruktur.Instruktur) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockInstrukturData) UpdateInstruktur(data instruktur.UpdateInstruktur) error {
	args := m.Called(data)
	return args.Error(0)
}

func (m *MockInstrukturData) DeleteInstruktur(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockInstrukturData) GetInstruktorByName(name string, page, limit int) ([]instruktur.Instruktur, int, error) {
	args := m.Called(name, page, limit)
	return args.Get(0).([]instruktur.Instruktur), args.Int(1), args.Error(2)
}

func TestInstrukturService(t *testing.T) {
	mockData := new(MockInstrukturData)
	service := New(mockData, nil)

	t.Run("GetInstrukturWithPagination success", func(t *testing.T) {
		mockResponse := []instruktur.Instruktur{{ID: "1", Name: "John"}}
		mockTotal := 10
		mockData.On("GetInstrukturWithPagination", 1, 5).Return(mockResponse, mockTotal, nil).Once()

		result, total, err := service.GetInstrukturWithPagination(1, 5)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		assert.Equal(t, mockTotal, total)
		mockData.AssertExpectations(t)
	})

	t.Run("GetAllInstruktur success", func(t *testing.T) {
		mockResponse := []instruktur.Instruktur{{ID: "1", Name: "John"}}
		mockData.On("GetAllInstruktur").Return(mockResponse, nil).Once()

		result, err := service.GetAllInstruktur()

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	t.Run("GetAllInstrukturByID success", func(t *testing.T) {
		mockResponse := instruktur.Instruktur{ID: "1", Name: "John"}
		mockData.On("GetAllInstrukturByID", "1").Return(mockResponse, nil).Once()

		result, err := service.GetAllInstrukturByID("1")

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		mockData.AssertExpectations(t)
	})

	t.Run("PostInstruktur success", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "John", Email: "john@example.com", Alamat: "Street 123", NoHp: "08123456789", Gender: "M", NIK: "123456789", NomorIndukPendidikan: "987654321", UrlImage: "image.jpg"}
		mockData.On("PostInstruktur", mockRequest).Return(nil).Once()

		err := service.PostInstruktur(mockRequest)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("PostInstruktur fail - empty name", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "", Email: "john@example.com"}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNameInstuktor, err)
	})

	t.Run("PostInstruktur fail - empty email", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyEmailInstuktor, err)
	})
	t.Run("PostInstruktur fail - empty email", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyAlamatInstuktor, err)
	})
	t.Run("PostInstruktur fail - empty no hp", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: "abc", NoHp: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyNumbertelponInstuktor, err)
	})
	t.Run("PostInstruktur fail - empty gender", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: "abc", NoHp: "01234567812", Gender: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrGenderInstruktorRmpty, err)
	})

	t.Run("PostInstruktur fail - empty NIK", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: "abc", NoHp: "01234567812", Gender: "Perempuan", NIK: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInstrukturNIKEmpty, err)
	})

	t.Run("PostInstruktur fail - empty NIP", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: "abc", NoHp: "01234567812", Gender: "Perempuan", NIK: "12345", NomorIndukPendidikan: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInstrukturNIPEmpty, err)
	})
	t.Run("PostInstruktur fail - empty image", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: "abc", NoHp: "01234567812", Gender: "Perempuan", NIK: "12345", NomorIndukPendidikan: "avq1231", UrlImage: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrInstrukturImageEmpty, err)
	})
	t.Run("PostInstruktur fail - empty email", func(t *testing.T) {
		mockRequest := instruktur.Instruktur{Name: "abc", Email: "abc@gmail.com", Alamat: ""}

		err := service.PostInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyAlamatInstuktor, err)
	})

	t.Run("UpdateInstruktur success", func(t *testing.T) {
		mockRequest := instruktur.UpdateInstruktur{ID: "1", Name: "John Updated", NoHp: "01234567812"}
		mockData.On("UpdateInstruktur", mockRequest).Return(nil).Once()

		err := service.UpdateInstruktur(mockRequest)

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("UpdateInstruktur fail - empty ID", func(t *testing.T) {
		mockRequest := instruktur.UpdateInstruktur{Name: "John Updated"}

		err := service.UpdateInstruktur(mockRequest)

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})

	t.Run("DeleteInstruktur success", func(t *testing.T) {
		mockData.On("DeleteInstruktur", "1").Return(nil).Once()

		err := service.DeleteInstruktur("1")

		assert.Nil(t, err)
		mockData.AssertExpectations(t)
	})

	t.Run("DeleteInstruktur fail - empty ID", func(t *testing.T) {
		err := service.DeleteInstruktur("")

		assert.NotNil(t, err)
		assert.Equal(t, constant.ErrEmptyId, err)
	})

	t.Run("GetInstruktorByName success", func(t *testing.T) {
		mockResponse := []instruktur.Instruktur{{ID: "1", Name: "John"}}
		mockTotal := 1
		mockData.On("GetInstruktorByName", "John", 1, 5).Return(mockResponse, mockTotal, nil).Once()

		result, total, err := service.GetInstruktorByName("John", 1, 5)

		assert.Nil(t, err)
		assert.Equal(t, mockResponse, result)
		assert.Equal(t, mockTotal, total)
		mockData.AssertExpectations(t)
	})
}
