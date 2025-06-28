package loader

import (
	"app/pkg/models"
	"encoding/json"
	"os"
)

type ProductJsonFile struct {
	path string
}

func (p *ProductJsonFile) Load() (v map[int]models.Product, err error) {
	file, err := os.Open(p.path)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	var productsJSON []models.ProductDoc
	err = json.NewDecoder(file).Decode(&productsJSON)
	if err != nil {
		return
	}

	v = make(map[int]models.Product)

	for _, pd := range productsJSON {
		v[pd.ID] = models.Product{
			ID: pd.ID,
			ProductAttributes: models.ProductAttributes{
				ProductCode:                    *pd.ProductCode,
				Description:                    *pd.Description,
				NetWeight:                      *pd.NetWeight,
				ExpirationRate:                 *pd.ExpirationRate,
				RecommendedFreezingTemperature: *pd.RecommendedFreezingTemperature,
				FreezingRate:                   *pd.FreezingRate,
				ProductTypeId:                  *pd.ProductTypeId,
				SellerId:                       *pd.SellerId,
			},
			Dimensions: models.Dimensions{
				Width:  *pd.Width,
				Height: *pd.Height,
				Length: *pd.Length,
			},
		}
	}

	return

}

func NewProductJSONFile(path string) *ProductJsonFile {
	return &ProductJsonFile{
		path: path,
	}
}
