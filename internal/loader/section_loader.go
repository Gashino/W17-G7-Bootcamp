package loader

import "app/pkg/models"

type SectionLoader interface {
	Load() (s map[int]models.Section, err error)
}
