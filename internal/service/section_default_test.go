package service

import (
	"testing"
)

func TestSectionService_GetAll(t *testing.T) {
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

func TestSectionService_GetByID(t *testing.T) {
	t.Run("Si el elemento buscado por id no existe retorna nulo", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el elemento buscado por id existe devolverá la información del elemento solicitado", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestSectionService_Create(t *testing.T) {
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

func TestSectionService_Update(t *testing.T) {
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

func TestSectionService_Delete(t *testing.T) {
	t.Run("Cuando la section no existe se devolverá null.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si la eliminación es exitosa el elemento no aparecerá en la lista.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}
