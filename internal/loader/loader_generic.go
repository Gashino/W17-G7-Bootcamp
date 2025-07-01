package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

type LoaderGeneric[T any] struct{}

func NewLoaderGeneric[T any]() *LoaderGeneric[T] {
	return &LoaderGeneric[T]{}
}

func (l *LoaderGeneric[T]) LoadFromJSON(path string) (map[int]T, error) {
	// Leer el archivo
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Deserializar el JSON a un slice de T
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	// Crear el mapa de resultados
	result := make(map[int]T)

	// Usar reflexión para obtener el campo ID de cada item
	for _, item := range items {
		// Usar reflexión para obtener el valor del campo ID
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		// Buscar el campo ID
		idField := val.FieldByName("ID")
		if !idField.IsValid() || !idField.CanInterface() {
			return nil, fmt.Errorf("item does not have a valid ID field")
		}

		// Convertir el ID a int
		var id int
		switch idVal := idField.Interface().(type) {
		case int:
			id = idVal
		case int64:
			id = int(idVal)
		case float64:
			id = int(idVal)
		default:
			return nil, fmt.Errorf("ID field is not a valid integer type")
		}

		// Agregar al mapa de resultados
		result[id] = item
	}

	return result, nil
}
