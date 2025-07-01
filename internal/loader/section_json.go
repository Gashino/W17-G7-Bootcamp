package loader

import (
	"encoding/json"
	"fmt"
	"os"
)

type LoaderGeneric[T any] struct{}

func NewLoaderGeneric[T any]() *LoaderGeneric[T] {
	return &LoaderGeneric[T]{}
}
func (l *LoaderGeneric[T]) LoadFromJSON(filePath string) ([]T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON from %s: %w", filePath, err)
	}
	return items, nil
}
