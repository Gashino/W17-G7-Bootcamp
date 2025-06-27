package models

type Product struct {
	ID int
	ProductAttributes
	Dimensions
}

type ProductAttributes struct {
	ProductCode                    string
	Description                    string
	NetWeight                      float64
	ExpirationRate                 int
	RecommendedFreezingTemperature float64
	FreezingRate                   int
	ProductTypeId                  int
	SellerId                       int
}

type Dimensions struct {
	Width  float64
	Height float64
	Length float64
}

type ProductDoc struct {
	ID                             int     `json:"id"`
	Description                    string  `json:"description"`
	ExpirationRate                 int     `json:"expiration_rate"`
	FreezingRate                   int     `json:"freezing_rate"`
	Height                         float64 `json:"height"`
	Length                         float64 `json:"length"`
	NetWeight                      float64 `json:"net_weight"`
	ProductCode                    string  `json:"product_code"`
	RecommendedFreezingTemperature float64 `json:"recommended_freezing_temperature"`
	Width                          float64 `json:"width"`
	ProductTypeId                  int     `json:"product_type_id"`
	SellerId                       int     `json:"seller_id"`
}

func (data *ProductDoc) ToStruct() Product {
	return Product{
		ID: data.ID,
		ProductAttributes: ProductAttributes{
			ProductCode:                    data.ProductCode,
			Description:                    data.Description,
			NetWeight:                      data.NetWeight,
			ExpirationRate:                 data.ExpirationRate,
			RecommendedFreezingTemperature: data.RecommendedFreezingTemperature,
			FreezingRate:                   data.FreezingRate,
			ProductTypeId:                  data.ProductTypeId,
			SellerId:                       data.SellerId,
		},
		Dimensions: Dimensions{
			Width:  data.Width,
			Height: data.Height,
			Length: data.Length,
		},
	}
}

func (data *Product) ToJSON() ProductDoc {
	return ProductDoc{
		ID:                             data.ID,
		Description:                    data.Description,
		ExpirationRate:                 data.ExpirationRate,
		FreezingRate:                   data.FreezingRate,
		Height:                         data.Height,
		Length:                         data.Length,
		NetWeight:                      data.NetWeight,
		ProductCode:                    data.ProductCode,
		RecommendedFreezingTemperature: data.RecommendedFreezingTemperature,
		Width:                          data.Width,
		ProductTypeId:                  data.ProductTypeId,
		SellerId:                       data.SellerId,
	}
}
