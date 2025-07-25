package warehouse

import (
	"app/pkg/models"

	"github.com/stretchr/testify/mock"
)

/*
	Create(v models.Carry) (c models.Carry, err error)
	SearchByLocality(locality_id int) (v map[string]models.CarryByLocality, err error)
*/

type MockCarryService struct {
	mock.Mock
}

func (m *MockCarryService) Create(v models.Carry) (w models.Carry, err error) {
	args := m.Called(v)
	return args.Get(0).(models.Carry), args.Error(1)
}

func (m *MockCarryService) SearchByLocality(locality_id int) (v map[string]models.CarryByLocality, err error) {
	args := m.Called()
	return args.Get(0).(map[string]models.CarryByLocality), args.Error(1)
}
