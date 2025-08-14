package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWarehouseDoc_AreFieldsValid(t *testing.T) {
	t.Run("valid warehouse doc", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "Test Address",
			Telephone:      "123456789",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.True(t, valid)
	})

	t.Run("empty warehouse code", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "",
			Address:        "Test Address",
			Telephone:      "123456789",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("empty address", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "",
			Telephone:      "123456789",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("empty telephone", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "Test Address",
			Telephone:      "",
			MinCapacity:    100,
			MinTemperature: -10,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("zero min capacity", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "Test Address",
			Telephone:      "123456789",
			MinCapacity:    0,
			MinTemperature: -10,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("zero min temperature", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{
			WarehouseCode:  "WH001",
			Address:        "Test Address",
			Telephone:      "123456789",
			MinCapacity:    100,
			MinTemperature: 0,
		}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("all fields empty/zero", func(t *testing.T) {
		warehouseDoc := WarehouseDoc{}

		valid := warehouseDoc.AreFieldsValid()
		require.False(t, valid)
	})
}
