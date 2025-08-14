package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidCreateRequest(t *testing.T) {
	t.Run("valid create request", func(t *testing.T) {
		cid := "12345"
		company := "Test Company"
		address := "Test Address"
		telephone := "123456789"
		localityID := 1

		request := SellerCreateRequest{
			CId:         &cid,
			CompanyName: &company,
			Address:     &address,
			Telephone:   &telephone,
			LocalityID:  &localityID,
		}

		valid := IsValidCreateRequest(request)
		require.True(t, valid)
	})

	t.Run("missing CId", func(t *testing.T) {
		company := "Test Company"
		address := "Test Address"
		telephone := "123456789"
		localityID := 1

		request := SellerCreateRequest{
			CompanyName: &company,
			Address:     &address,
			Telephone:   &telephone,
			LocalityID:  &localityID,
		}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})

	t.Run("missing CompanyName", func(t *testing.T) {
		cid := "12345"
		address := "Test Address"
		telephone := "123456789"
		localityID := 1

		request := SellerCreateRequest{
			CId:        &cid,
			Address:    &address,
			Telephone:  &telephone,
			LocalityID: &localityID,
		}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})

	t.Run("missing Address", func(t *testing.T) {
		cid := "12345"
		company := "Test Company"
		telephone := "123456789"
		localityID := 1

		request := SellerCreateRequest{
			CId:         &cid,
			CompanyName: &company,
			Telephone:   &telephone,
			LocalityID:  &localityID,
		}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})

	t.Run("missing Telephone", func(t *testing.T) {
		cid := "12345"
		company := "Test Company"
		address := "Test Address"
		localityID := 1

		request := SellerCreateRequest{
			CId:         &cid,
			CompanyName: &company,
			Address:     &address,
			LocalityID:  &localityID,
		}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})

	t.Run("missing LocalityID", func(t *testing.T) {
		cid := "12345"
		company := "Test Company"
		address := "Test Address"
		telephone := "123456789"

		request := SellerCreateRequest{
			CId:         &cid,
			CompanyName: &company,
			Address:     &address,
			Telephone:   &telephone,
		}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})

	t.Run("all fields missing", func(t *testing.T) {
		request := SellerCreateRequest{}

		valid := IsValidCreateRequest(request)
		require.False(t, valid)
	})
}
