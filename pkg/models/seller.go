package models

// SellerAttributes is a struct that represents the attributes of a seller
type SellerAttributes struct {
	CId         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityID  int    `json:"locality_id"`
}

// SellerDoc is a struct that represents a Seller in JSON format
type SellerDoc struct {
	ID          int    `json:"id"`
	CId         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Adress      string `json:"address"`
	Telephone   string `json:"telephone"`
	LocalityID  int    `json:"locality_id"`
}

// SellerCreateRequest is a struct that represents a Seller in JSON format (Request)
type SellerCreateRequest struct {
	CId         *string `json:"cid,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	Address     *string `json:"address,omitempty"`
	Telephone   *string `json:"telephone,omitempty"`
	LocalityID  *int    `json:"locality_id,omitempty"`
}

// Seller is a struct that represents a Seller
type Seller struct {
	ID int `json:"id"`
	SellerAttributes
}

// ToSellerDoc convierte un Seller en su versión externa (para response JSON).
func ToSellerDoc(s Seller) SellerDoc {
	return SellerDoc{
		ID:          s.ID,
		CId:         s.CId,
		CompanyName: s.CompanyName,
		Adress:      s.Address,
		Telephone:   s.Telephone,
		LocalityID:  s.LocalityID,
	}
}

// FromSellerDoc convierte un SellerDoc (por ejemplo de una base de datos o input externo) a un Seller interno.
func FromSellerDoc(d SellerDoc) Seller {
	return Seller{
		ID: d.ID,
		SellerAttributes: SellerAttributes{
			CId:         d.CId,
			CompanyName: d.CompanyName,
			Address:     d.Adress,
			Telephone:   d.Telephone,
			LocalityID:  d.LocalityID,
		},
	}
}

// FromCreateRequest convierte una solicitud de creación a un Seller interno.
func FromCreateRequest(r SellerCreateRequest) Seller {
	return Seller{
		ID: 0, // Se puede reemplazar al insertar
		SellerAttributes: SellerAttributes{
			CId:         *r.CId,
			CompanyName: *r.CompanyName,
			Address:     *r.Address,
			Telephone:   *r.Telephone,
			LocalityID:  *r.LocalityID,
		},
	}
}

// IsValidCreateRequest valida que los campos obligatorios estén presentes.
func IsValidCreateRequest(r SellerCreateRequest) bool {
	return r.CId != nil && r.CompanyName != nil && r.Address != nil && r.Telephone != nil && r.LocalityID != nil
}
