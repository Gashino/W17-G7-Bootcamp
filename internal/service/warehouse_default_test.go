package service

import (
	"testing"
)

/*
FindAll() (v map[int]models.Warehouse, err error)
FindByID(id int) (v models.Warehouse, err error)
Add(v models.WarehouseDoc) (w models.Warehouse, err error)
Update(id int, v models.WarehouseDoc) (w models.Warehouse, err error)
Delete(id int) (err error)
*/
func TestWarehouseService_FindAll(t *testing.T) {
	t.Run("Si la lista posee “n” elementos devolverá un cantidad de los elementos totales", func(t *testing.T) {
		// Arrange

		// Act

		// Assert

	})

	t.Run("Error", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouseService_FindByID(t *testing.T) {
	t.Run("Si el elemento buscado por id existe devolverá la información del elemento solicitado", func(t *testing.T) {
		// Arrange

		// Act

		// Assert

	})

	t.Run("Si el elemento buscado por id no existe retorna nulo", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouseService_Add(t *testing.T) {
	t.Run("Si contiene los campos necesarios se creará", func(t *testing.T) {
		// Arrange

		// Act

		// Assert

	})

	t.Run("Si el section_number ya existe no podrá ser creado", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouseService_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá la section con la información actualizada", func(t *testing.T) {
		// Arrange

		// Act

		// Assert

	})

	t.Run("Si la section que se desea actualizar no existe se retorna null.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouseService_Delete(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá la section con la información actualizada", func(t *testing.T) {
		// Arrange

		// Act

		// Assert

	})

	t.Run("Si la section que se desea actualizar no existe se retorna null.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}
