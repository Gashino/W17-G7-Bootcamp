package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmployee(t *testing.T) {
	t.Run("valid employee with ID", func(t *testing.T) {
		id := 1
		cardNumber := "12345678"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			ID:           &id,
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, true)
		require.NoError(t, err)
	})

	t.Run("valid employee without ID validation", func(t *testing.T) {
		cardNumber := "12345678"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.NoError(t, err)
	})

	t.Run("invalid ID when validation required", func(t *testing.T) {
		id := 0
		cardNumber := "12345678"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			ID:           &id,
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, true)
		require.Error(t, err)
		require.Contains(t, err.Error(), "id is required and must be greater than 0")
	})

	t.Run("missing card number ID", func(t *testing.T) {
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			FirstName:   &firstName,
			LastName:    &lastName,
			WarehouseID: &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "card_number_id is required")
	})

	t.Run("invalid card number length", func(t *testing.T) {
		cardNumber := "123456"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "card_number_id must be 8 digits")
	})

	t.Run("invalid card number format", func(t *testing.T) {
		cardNumber := "1234567a"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "card_number_id must contain only numbers")
	})

	t.Run("missing first name", func(t *testing.T) {
		cardNumber := "12345678"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "first_name is required")
	})

	t.Run("empty first name", func(t *testing.T) {
		cardNumber := "12345678"
		firstName := "   "
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "first_name is required")
	})

	t.Run("missing last name", func(t *testing.T) {
		cardNumber := "12345678"
		firstName := "John"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "last_name is required")
	})

	t.Run("invalid warehouse ID", func(t *testing.T) {
		cardNumber := "12345678"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 0

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployee(employee, false)
		require.Error(t, err)
		require.Contains(t, err.Error(), "warehouse_id is required and must be greater than 0")
	})
}

func TestValidateEmployeeUpdate(t *testing.T) {
	t.Run("valid update with all fields", func(t *testing.T) {
		cardNumber := "12345678"
		firstName := "John"
		lastName := "Doe"
		warehouseID := 1

		employee := Employee{
			CardNumberID: &cardNumber,
			FirstName:    &firstName,
			LastName:     &lastName,
			WarehouseID:  &warehouseID,
		}

		err := ValidateEmployeeUpdate(employee)
		require.NoError(t, err)
	})

	t.Run("valid update with only card number", func(t *testing.T) {
		cardNumber := "12345678"

		employee := Employee{
			CardNumberID: &cardNumber,
		}

		err := ValidateEmployeeUpdate(employee)
		require.NoError(t, err)
	})

	t.Run("no fields provided", func(t *testing.T) {
		employee := Employee{}

		err := ValidateEmployeeUpdate(employee)
		require.Error(t, err)
		require.Contains(t, err.Error(), "at least one field is required for update")
	})

	t.Run("invalid card number in update", func(t *testing.T) {
		cardNumber := "123"

		employee := Employee{
			CardNumberID: &cardNumber,
		}

		err := ValidateEmployeeUpdate(employee)
		require.Error(t, err)
		require.Contains(t, err.Error(), "card_number_id must be 8 digits")
	})
}

func TestStringPtr(t *testing.T) {
	t.Run("string to pointer", func(t *testing.T) {
		str := "test"
		ptr := StringPtr(str)
		require.NotNil(t, ptr)
		require.Equal(t, str, *ptr)
	})
}

func TestIntPtr(t *testing.T) {
	t.Run("int to pointer", func(t *testing.T) {
		val := 123
		ptr := IntPtr(val)
		require.NotNil(t, ptr)
		require.Equal(t, val, *ptr)
	})
}

func TestIsValidStringPtr(t *testing.T) {
	t.Run("valid string pointer", func(t *testing.T) {
		str := "test"
		valid := IsValidStringPtr(&str)
		require.True(t, valid)
	})

	t.Run("nil string pointer", func(t *testing.T) {
		valid := IsValidStringPtr(nil)
		require.False(t, valid)
	})

	t.Run("empty string pointer", func(t *testing.T) {
		str := ""
		valid := IsValidStringPtr(&str)
		require.False(t, valid)
	})
}

func TestIsValidIntPtr(t *testing.T) {
	t.Run("valid int pointer", func(t *testing.T) {
		val := 123
		valid := IsValidIntPtr(&val)
		require.True(t, valid)
	})

	t.Run("nil int pointer", func(t *testing.T) {
		valid := IsValidIntPtr(nil)
		require.False(t, valid)
	})

	t.Run("zero int pointer", func(t *testing.T) {
		val := 0
		valid := IsValidIntPtr(&val)
		require.False(t, valid)
	})
}
