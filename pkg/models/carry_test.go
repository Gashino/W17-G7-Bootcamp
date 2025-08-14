package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCarryDoc_AreFieldsValid(t *testing.T) {
	t.Run("valid carry doc", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "C001",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		valid := carryDoc.AreFieldsValid()
		require.True(t, valid)
	})

	t.Run("empty cid", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("empty company name", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "C001",
			CompanyName: "",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("empty address", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "C001",
			CompanyName: "Test Company",
			Address:     "",
			Telephone:   "123456789",
			LocalityId:  1,
		}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("empty telephone", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "C001",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "",
			LocalityId:  1,
		}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("zero locality id", func(t *testing.T) {
		carryDoc := CarryDoc{
			Cid:         "C001",
			CompanyName: "Test Company",
			Address:     "Test Address",
			Telephone:   "123456789",
			LocalityId:  0,
		}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})

	t.Run("all fields empty/zero", func(t *testing.T) {
		carryDoc := CarryDoc{}

		valid := carryDoc.AreFieldsValid()
		require.False(t, valid)
	})
}
