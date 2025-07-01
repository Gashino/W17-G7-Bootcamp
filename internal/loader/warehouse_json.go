package loader

import (
	"app/pkg/models"
	"encoding/json"
	"os"
)

// NewVehicleJSONFile is a function that returns a new instance of VehicleJSONFile
func NewWarehouseJSONFile(path string) *WarehouseJSONFile {
	return &WarehouseJSONFile{
		path: path,
	}
}

// VehicleJSONFile is a struct that implements the LoaderVehicle interface
type WarehouseJSONFile struct {
	// path is the path to the file that contains the vehicles in JSON format
	path string
}

// Load is a method that loads the warehouses
func (l *WarehouseJSONFile) Load() (v map[int]models.Warehouse, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var warehouseJSON []models.WarehouseDoc
	err = json.NewDecoder(file).Decode(&warehouseJSON)
	if err != nil {
		return
	}

	// serialize warehouses
	v = make(map[int]models.Warehouse)
	for _, wh := range warehouseJSON {
		v[wh.ID] = models.Warehouse{
			ID:             wh.ID,
			WarehouseCode:  wh.WarehouseCode,
			Address:        wh.Address,
			Telephone:      wh.Telephone,
			MinCapacity:    wh.MinCapacity,
			MinTemperature: wh.MinTemperature,
		}
	}

	return
}
