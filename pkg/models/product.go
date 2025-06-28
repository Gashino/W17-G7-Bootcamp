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
	ID                             int      `json:"id,omitempty"`
	Description                    *string  `json:"description,omitempty"`
	ExpirationRate                 *int     `json:"expiration_rate,omitempty"`
	FreezingRate                   *int     `json:"freezing_rate,omitempty"`
	Height                         *float64 `json:"height,omitempty"`
	Length                         *float64 `json:"length,omitempty"`
	NetWeight                      *float64 `json:"net_weight,omitempty"`
	ProductCode                    *string  `json:"product_code,omitempty"`
	RecommendedFreezingTemperature *float64 `json:"recommended_freezing_temperature,omitempty"`
	Width                          *float64 `json:"width,omitempty"`
	ProductTypeId                  *int     `json:"product_type_id,omitempty"`
	SellerId                       *int     `json:"seller_id,omitempty"`
}

func (data *ProductDoc) ToStruct() Product {
	return Product{
		ID: data.ID,
		ProductAttributes: ProductAttributes{
			ProductCode:                    *data.ProductCode,
			Description:                    *data.Description,
			NetWeight:                      *data.NetWeight,
			ExpirationRate:                 *data.ExpirationRate,
			RecommendedFreezingTemperature: *data.RecommendedFreezingTemperature,
			FreezingRate:                   *data.FreezingRate,
			ProductTypeId:                  *data.ProductTypeId,
			SellerId:                       *data.SellerId,
		},
		Dimensions: Dimensions{
			Width:  *data.Width,
			Height: *data.Height,
			Length: *data.Length,
		},
	}
}

func (data *Product) ToJSON() ProductDoc {
	return ProductDoc{
		ID:                             data.ID,
		Description:                    &data.Description,
		ExpirationRate:                 &data.ExpirationRate,
		FreezingRate:                   &data.FreezingRate,
		Height:                         &data.Height,
		Length:                         &data.Length,
		NetWeight:                      &data.NetWeight,
		ProductCode:                    &data.ProductCode,
		RecommendedFreezingTemperature: &data.RecommendedFreezingTemperature,
		Width:                          &data.Width,
		ProductTypeId:                  &data.ProductTypeId,
		SellerId:                       &data.SellerId,
	}
}

func (data *ProductDoc) Validate() bool {
	if data.Description == nil {
		return false
	}
	if data.ExpirationRate == nil {
		return false
	}
	if data.FreezingRate == nil {
		return false
	}
	if data.Height == nil {
		return false
	}
	if data.Length == nil {
		return false
	}
	if data.NetWeight == nil {
		return false
	}
	if data.ProductCode == nil {
		return false
	}
	if data.RecommendedFreezingTemperature == nil {
		return false
	}
	if data.Width == nil {
		return false
	}
	if data.ProductTypeId == nil {
		return false
	}
	return true
}
