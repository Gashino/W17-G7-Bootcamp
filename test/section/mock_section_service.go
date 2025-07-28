package section

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

type MockSectionService struct {
	mock.Mock
}

func (m *MockSectionService) GetAll() (map[int]models.Section, error) {
	args := m.Called()
	return args.Get(0).(map[int]models.Section), args.Error(1)
}

func (m *MockSectionService) GetByID(id int) (models.Section, error) {
	args := m.Called(id)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionService) Create(section models.Section) (models.Section, error) {
	args := m.Called(section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionService) Update(id int, section models.Section) (models.Section, error) {
	args := m.Called(id, section)
	return args.Get(0).(models.Section), args.Error(1)
}

func (m *MockSectionService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSectionService) ReportProductsAllSections() (s []models.SectionReport, err error) {
	args := m.Called()
	return args.Get(0).([]models.SectionReport), args.Error(1)
}

func (m *MockSectionService) ReportProductsBySection(ptr *int) (s []models.SectionReport, err error) {
	args := m.Called(ptr)
	return args.Get(0).([]models.SectionReport), args.Error(1)
}

func (m *MockSectionService) GetReportProductsBySection(id int) ([]models.SectionReport, error) {
	args := m.Called(id)
	return args.Get(0).([]models.SectionReport), args.Error(1)
}

func (m *MockSectionService) GetReportProductsAllSections() ([]models.SectionReport, error) {
	args := m.Called()
	return args.Get(0).([]models.SectionReport), args.Error(1)
}
