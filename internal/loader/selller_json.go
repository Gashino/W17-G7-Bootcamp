package loader

import (
	"app/pkg/models"
	"encoding/json"
	"os"
)

// NewSellerJSONFile is a function that returns a new instance of SellerJSONFile
func NewSellerJSONFile(path string) *SellerJSONFile {
	return &SellerJSONFile{
		path: path,
	}
}

// SellerJSONFile is a struct that implements the LoaderSeller interface
type SellerJSONFile struct {
	// path is the path to the file that contains the Sellers in JSON format
	path string
}

// Load is a method that loads the Sellers
func (l *SellerJSONFile) Load() (s map[int]models.Seller, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var SellersJSON []models.SellerDoc
	err = json.NewDecoder(file).Decode(&SellersJSON)
	if err != nil {
		return
	}

	// serialize Sellers
	s = make(map[int]models.Seller)
	for _, seller := range SellersJSON {
		s[seller.ID] = models.Seller{
			Id: seller.ID,
			SellerAttributes: models.SellerAttributes{
				CId:         seller.CId,
				CompanyName: seller.CompanyName,
				Adress:      seller.Address,
				Telephone:   seller.Telephone,
			},
		}
	}

	return
}
