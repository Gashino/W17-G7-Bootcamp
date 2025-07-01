package models

type Warehouse struct {
	Id             int
	WarehouseCode  string
	Address        string
	Telephone      string
	MinCapacity    int
	MinTemperature int
}

type WarehouseDoc struct {
	Id             int    `json:"id"`
	WarehouseCode  string `json:"warehouse_code"`
	Address        string `json:"address"`
	Telephone      string `json:"telephone"`
	MinCapacity    int    `json:"minimun_capacity"`
	MinTemperature int    `json:"minimun_temperature"`
}

func (v *WarehouseDoc) AreFieldsValid() bool {
	if v.WarehouseCode == "" || v.Address == "" || v.Telephone == "" ||
		v.MinCapacity == 0 || v.MinTemperature == 0 {
		return false
	}
	return true
}
