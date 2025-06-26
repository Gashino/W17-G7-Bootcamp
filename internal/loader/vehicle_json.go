package loader

import (
	"app/pkg/models"
	"encoding/json"
	"os"
)

// NewVehicleJSONFile is a function that returns a new instance of VehicleJSONFile
func NewVehicleJSONFile(path string) *VehicleJSONFile {
	return &VehicleJSONFile{
		path: path,
	}
}

// VehicleJSONFile is a struct that implements the LoaderVehicle interface
type VehicleJSONFile struct {
	// path is the path to the file that contains the vehicles in JSON format
	path string
}

// Load is a method that loads the vehicles
func (l *VehicleJSONFile) Load() (v map[int]models.Vehicle, err error) {
	// open file
	file, err := os.Open(l.path)
	if err != nil {
		return
	}
	defer file.Close()

	// decode file
	var vehiclesJSON []models.VehicleDoc
	err = json.NewDecoder(file).Decode(&vehiclesJSON)
	if err != nil {
		return
	}

	// serialize vehicles
	v = make(map[int]models.Vehicle)
	for _, vh := range vehiclesJSON {
		v[vh.ID] = models.Vehicle{
			Id: vh.ID,
			VehicleAttributes: models.VehicleAttributes{
				Brand:           vh.Brand,
				Model:           vh.Model,
				Registration:    vh.Registration,
				Color:           vh.Color,
				FabricationYear: vh.FabricationYear,
				Capacity:        vh.Capacity,
				MaxSpeed:        vh.MaxSpeed,
				FuelType:        vh.FuelType,
				Transmission:    vh.Transmission,
				Weight:          vh.Weight,
				Dimensions: models.Dimensions{
					Height: vh.Height,
					Length: vh.Length,
					Width:  vh.Width,
				},
			},
		}
	}

	return
}
