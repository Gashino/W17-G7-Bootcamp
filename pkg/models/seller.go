package models

// SellerAttributes is a struct that represents the attributes of a seller
type SellerAttributes struct {
	CId         string
	CompanyName string
	Adress      string
	Telephone   string
}

// SellerDoc is a struct that represents a Seller in JSON format
type SellerDoc struct {
	ID          int    `json:"id"`
	CId         string `json:"cid"`
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	Telephone   string `json:"telephone"`
}

// SellerCreateRequest is a struct that represents a Seller in JSON format (Request)
type SellerCreateRequest struct {
	CId         *string `json:"cid,omitempty"`
	CompanyName *string `json:"company_name,omitempty"`
	Address     *string `json:"address,omitempty"`
	Telephone   *string `json:"telephone,omitempty"`
}

// Seller is a struct that represents a Seller
type Seller struct {
	Id int
	SellerAttributes
}

// ToSellerDoc convierte un Seller en su versión externa (para response JSON).
func ToSellerDoc(s Seller) SellerDoc {
	return SellerDoc{
		ID:          s.Id,
		CId:         s.CId,
		CompanyName: s.CompanyName,
		Address:     s.Adress,
		Telephone:   s.Telephone,
	}
}

// FromSellerDoc convierte un SellerDoc (por ejemplo de una base de datos o input externo) a un Seller interno.
func FromSellerDoc(d SellerDoc) Seller {
	return Seller{
		Id: d.ID,
		SellerAttributes: SellerAttributes{
			CId:         d.CId,
			CompanyName: d.CompanyName,
			Adress:      d.Address,
			Telephone:   d.Telephone,
		},
	}
}

// FromCreateRequest convierte una solicitud de creación a un Seller interno.
// El campo Id se espera que lo asigne el service o repositorio (por ejemplo, autoincremental).
func FromCreateRequest(r SellerCreateRequest) Seller {
	return Seller{
		Id: 0, // Se puede reemplazar al insertar
		SellerAttributes: SellerAttributes{
			CId:         *r.CId,
			CompanyName: *r.CompanyName,
			Adress:      *r.Address,
			Telephone:   *r.Telephone,
		},
	}
}

// IsValidCreateRequest valida que los campos obligatorios estén presentes.
// Nota: Considera validar formato de teléfono o longitud del texto si lo necesitás.
func IsValidCreateRequest(r SellerCreateRequest) bool {
	return r.CId != nil && r.CompanyName != nil && r.Address != nil && r.Telephone != nil
}
