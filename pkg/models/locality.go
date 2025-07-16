package models

// LocalitiesAttributes is a struct that represents the attributes of a Locality
type LocalitiesAttributes struct {
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}

// LocalitiesDoc is a struct that represents a Locality in JSON format
type LocalitiesDoc struct {
	ID           int    `json:"id"`
	LocalityName string `json:"locality_name"`
	ProvinceName string `json:"province_name"`
	CountryName  string `json:"country_name"`
}

// LocalityCreateRequest is a struct that represents a Locality in JSON format (Request)
type LocalityCreateRequest struct {
	LocalityName *string `json:"locality_name,omitempty"`
	ProvinceName *string `json:"province_name,omitempty"`
	CountryName  *string `json:"country_name,omitempty"`
}

// LocalityBySellerResponse is a struct that represents a Locality response with seller count
type LocalityBySellerResponse struct {
	ID           int     `json:"id"`
	LocalityName *string `json:"locality_name"`
	SellerCount  *string `json:"seller_count"`
}

// Locality is a struct that represents a Locality
type Locality struct {
	ID int `json:"id"`
	LocalitiesAttributes
}

// ToLocalityDoc converts a Locality to its external version (for JSON response).
func ToLocalityDoc(s Locality) LocalitiesDoc {
	return LocalitiesDoc{
		ID:           s.ID,
		LocalityName: s.LocalityName,
		ProvinceName: s.ProvinceName,
		CountryName:  s.CountryName,
	}
}

// FromLocalityDoc converts a LocalitiesDoc (e.g., from database or external input) to an internal Locality.
func FromLocalityDoc(d LocalitiesDoc) Locality {
	return Locality{
		ID: d.ID,
		LocalitiesAttributes: LocalitiesAttributes{
			LocalityName: d.LocalityName,
			ProvinceName: d.ProvinceName,
			CountryName:  d.CountryName,
		},
	}
}

// FromCreateRequestLocality converts a creation request to an internal Locality.
func FromCreateRequestLocality(r LocalityCreateRequest) Locality {
	return Locality{
		ID: 0, // Can be replaced when inserting
		LocalitiesAttributes: LocalitiesAttributes{
			LocalityName: *r.LocalityName,
			ProvinceName: *r.ProvinceName,
			CountryName:  *r.CountryName,
		},
	}
}

// IsValidCreateRequestLocalities validates that required fields are present.
func IsValidCreateRequestLocalities(r LocalityCreateRequest) bool {
	return !IsNullOrEmpty(r.LocalityName) && !IsNullOrEmpty(r.ProvinceName) && !IsNullOrEmpty(r.CountryName)
}

// IsNullOrEmpty validates that a *string is empty or null
func IsNullOrEmpty(cadena *string) bool {
	if cadena == nil || *cadena == "" {
		return true
	}
	return false
}
