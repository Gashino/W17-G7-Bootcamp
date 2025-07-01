package models

type Section struct {
	ID int `json:"id"`
	SectionAttributes
}

type SectionAttributes struct {
	SectionNumber      int            `json:"section_number"`
	CurrentTemperature float64        `json:"current_temperature"`
	MinimumTemperature float64        `json:"minimum_temperature"`
	CurrentCapacity    int            `json:"current_capacity"`
	MinimumCapacity    int            `json:"minimum_capacity"`
	MaximumCapacity    int            `json:"maximum_capacity"`
	WarehouseID        int            `json:"warehouse_id"`
	ProductTypeID      int            `json:"product_type_id"`
	ProductBatches     []ProductBatch `json:"product_batches"`
}

type ProductBatch struct {
	ID         int   `json:"id"`
	ProductIDs []int `json:"product_ids"`
}

type SectionDoc struct {
	ID                 int     `json:"id"`
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

type SectionPostDoc struct {
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}
