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

func (data *Product) ToJSON() ProductDoc {
	desc := data.Description
	exp := data.ExpirationRate
	frz := data.FreezingRate
	h := data.Height
	l := data.Length
	w := data.Width
	nw := data.NetWeight
	pc := data.ProductCode
	rft := data.RecommendedFreezingTemperature
	pt := data.ProductTypeId
	si := data.SellerId

	return ProductDoc{
		ID:                             data.ID,
		Description:                    &desc,
		ExpirationRate:                 &exp,
		FreezingRate:                   &frz,
		Height:                         &h,
		Length:                         &l,
		NetWeight:                      &nw,
		ProductCode:                    &pc,
		RecommendedFreezingTemperature: &rft,
		Width:                          &w,
		ProductTypeId:                  &pt,
		SellerId:                       &si,
	}
}

type ProductType struct {
	ID   int
	Name string
}
