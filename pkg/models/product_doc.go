package models

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

func (data *ProductDoc) MapDataToStruct(product *Product) {
	if data.Width != nil {
		product.Width = *data.Width
	}
	if data.Height != nil {
		product.Height = *data.Height
	}
	if data.Length != nil {
		product.Length = *data.Length
	}
	if data.NetWeight != nil {
		product.NetWeight = *data.NetWeight
	}
	if data.ProductCode != nil {
		product.ProductCode = *data.ProductCode
	}
	if data.ProductTypeId != nil {
		product.ProductTypeId = *data.ProductTypeId
	}
	if data.RecommendedFreezingTemperature != nil {
		product.RecommendedFreezingTemperature = *data.RecommendedFreezingTemperature
	}
	if data.Description != nil {
		product.Description = *data.Description
	}
	if data.ExpirationRate != nil {
		product.ExpirationRate = *data.ExpirationRate
	}
	if data.FreezingRate != nil {
		product.FreezingRate = *data.FreezingRate
	}
	if data.SellerId != nil {
		product.FreezingRate = *data.FreezingRate
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
