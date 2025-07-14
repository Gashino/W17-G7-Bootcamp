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

// LocalityBySellerResponse is a struct that represents a Locality in JSON format (Request)
type LocalityBySellerResponse struct {
	ID           int     `json:"id"`
	LocalityName *string `json:"locality_name"`
	SellerCount  *string `json:"seller_count"`
}

// Seller is a struct that represents a Locality
type Locality struct {
	ID int `json:"id"`
	LocalitiesAttributes
}

// ToSellerDoc convierte un Seller en su versión externa (para response JSON).
func ToLocalityDoc(s Locality) LocalitiesDoc {
	return LocalitiesDoc{
		ID:           s.ID,
		LocalityName: s.LocalityName,
		ProvinceName: s.ProvinceName,
		CountryName:  s.CountryName,
	}
}

// FromSellerDoc convierte un LocalitiesDoc (por ejemplo de una base de datos o input externo) a un Locality interno.
func FromSellerDocLocality(d LocalitiesDoc) Locality {
	return Locality{
		ID: d.ID,
		LocalitiesAttributes: LocalitiesAttributes{
			LocalityName: d.LocalityName,
			ProvinceName: d.ProvinceName,
			CountryName:  d.CountryName,
		},
	}
}

// FromCreateRequest convierte una solicitud de creación a un Locality interno.
func FromCreateRequestLocality(r LocalityCreateRequest) Locality {
	return Locality{
		ID: 0, // Se puede reemplazar al insertar
		LocalitiesAttributes: LocalitiesAttributes{
			LocalityName: *r.LocalityName,
			ProvinceName: *r.ProvinceName,
			CountryName:  *r.CountryName,
		},
	}
}

// IsValidCreateRequest valida que los campos obligatorios estén presentes.
func IsValidCreateRequestLocalities(r LocalityCreateRequest) bool {
	return IsNullOrEmpty(r.LocalityName) || IsNullOrEmpty(r.ProvinceName) || IsNullOrEmpty(r.ProvinceName)
}

// IsNullOrEmpty valida que un *string este vacio o nulo
func IsNullOrEmpty(cadena *string) bool {
	if cadena == nil || *cadena == "" {
		return true
	}
	return false
}
