package loader

import (
	"app/pkg/models"
	"encoding/json"
	"os"
)

// NewBuyerJSONFile is a function that returns a new instance of BuyerJSONFile
func NewBuyerJSONFile(path string) *BuyerJSONFile {
	return &BuyerJSONFile{
		path: path,
	}
}

// BuyerJSONFile is a struct that implements the LoaderBuyer interface
type BuyerJSONFile struct {
	// path is the path to the file that contains the buyers in JSON format
	path string
}

// Load is a method that loads the buyers
func (l *BuyerJSONFile) Load() (b map[int]models.Buyer, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var buyersJSON []models.BuyerDoc
	err = json.NewDecoder(file).Decode(&buyersJSON)
	if err != nil {
		return
	}

	// serialize buyers
	b = make(map[int]models.Buyer)
	for _, buyer := range buyersJSON {
		b[buyer.ID] = models.Buyer{
			Id: buyer.ID,
			BuyerAttributes: models.BuyerAttributes{
				CardNumberID: buyer.CardNumberID,
				FirstName:    buyer.FirstName,
				LastName:     buyer.LastName,
			},
		}
	}

	return
}
