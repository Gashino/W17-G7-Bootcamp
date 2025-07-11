package repository

import (
	"app/pkg/models"
	"database/sql"
)

func NewSectionSqlRepository(db *sql.DB) SectionRepository {

	return &SectionRepositorySql{
		db: db,
	}
}

// SectionMap is a struct that represents a section repository
type SectionRepositorySql struct {
	// db is a map of sections
	db *sql.DB
}

func (r *SectionRepositorySql) GetAll() (s map[int]models.Section, err error) {
	return
}

func (r *SectionRepositorySql) GetByID(id int) (s models.Section, err error) {
	return
}

func (r *SectionRepositorySql) Create(section models.Section) (s models.Section, err error) {
	return
}

func (r *SectionRepositorySql) Update(id int, section models.Section) (s models.Section, err error) {
	return
}

func (r *SectionRepositorySql) Delete(id int) (err error) {
	return
}
