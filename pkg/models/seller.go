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
	Address     string `json:"address"`
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

// ToSellerDoc converts a Seller to its external version (for JSON response).
func ToSellerDoc(s Seller) SellerDoc {
	return SellerDoc{
		ID:          s.ID,
		CId:         s.CId,
		CompanyName: s.CompanyName,
		Address:     s.Address,
		Telephone:   s.Telephone,
		LocalityID:  s.LocalityID,
	}
}

// FromSellerDoc converts a SellerDoc (e.g., from database or external input) to an internal Seller.
func FromSellerDoc(d SellerDoc) Seller {
	return Seller{
		ID: d.ID,
		SellerAttributes: SellerAttributes{
			CId:         d.CId,
			CompanyName: d.CompanyName,
			Address:     d.Address,
			Telephone:   d.Telephone,
			LocalityID:  d.LocalityID,
		},
	}
}

// FromCreateRequest converts a creation request to an internal Seller.
func FromCreateRequest(r SellerCreateRequest) Seller {
	return Seller{
		ID: 0, // Can be replaced when inserting
		SellerAttributes: SellerAttributes{
			CId:         *r.CId,
			CompanyName: *r.CompanyName,
			Address:     *r.Address,
			Telephone:   *r.Telephone,
			LocalityID:  *r.LocalityID,
		},
	}
}

// IsValidCreateRequest validates that required fields are present.
func IsValidCreateRequest(r SellerCreateRequest) bool {
	return r.CId != nil && r.CompanyName != nil && r.Address != nil && r.Telephone != nil && r.LocalityID != nil
}
