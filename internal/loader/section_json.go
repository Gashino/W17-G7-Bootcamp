package loader

import (
	"encoding/json"
	"fmt"
	"os"
)

type LoaderGeneric2[T any] struct{}

func NewLoaderGeneric2[T any]() *LoaderGeneric2[T] {
	return &LoaderGeneric2[T]{}
}
func (l *LoaderGeneric2[T]) LoadFromJSON(filePath string) ([]T, error) {
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
