package repository

import "app/pkg/models"

type SectionRepository interface {
	GetAll() (s map[int]models.Section, err error)
	GetByID(id int) (s models.Section, err error)
	Create(section models.Section) (s models.Section, err error)
	Update(id int, section models.Section) (s models.Section, err error)
	Delete(id int) (err error)
}
