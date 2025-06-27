package loader

import "app/pkg/models"

type ProductLoader interface {
	Load() (v map[int]models.Product, err error)
}
