package handler

import "testing"

func TestWarehouse_GetAll(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá un listado de todas los warehouses existentes", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_GetOne(t *testing.T) {
	t.Run("Cuando la petición sea exitosa el backend devolverá la información del warehouse solicitado", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
	t.Run("Cuando el warehouse no exista se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_Add(t *testing.T) {
	t.Run("Cuando el ingreso de datos sea exitoso se devolverá un código 201 junto con el objeto ingresado.", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el objeto JSON no contiene los campos necesarios se devolverá un código 422", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el warehouse_code ya existe devuelve un error 409 Conflict", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_Update(t *testing.T) {
	t.Run("Cuando la actualización de datos sea exitosa se devolverá el warehouse con la información actualizada junto con un código 200", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Si el warehouse que se desea actualizar no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}

func TestWarehouse_Delete(t *testing.T) {
	t.Run("Cuando el warehouse no existe se devolverá un código 404", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})

	t.Run("Cuando la eliminación sea exitosa se devolverá un código 204", func(t *testing.T) {
		// Arrange

		// Act

		// Assert
	})
}
