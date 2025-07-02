package models

// BuyerDoc is a struct that represents a buyer document in JSON format
type BuyerDoc struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

// Buyer is a struct that represents a buyer
type Buyer struct {
	// Id is the unique identifier of the buyer
	ID int `json:"id"`

	// BuyerAttributes contains the attributes of a buyer
	BuyerAttributes
}

// BuyerAttributes is a struct that contains the attributes of a buyer
type BuyerAttributes struct {
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
