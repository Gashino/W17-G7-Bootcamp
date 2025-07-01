package repository

import (
	"app/pkg"
	"app/pkg/models"
	"fmt"
)

func NewSectionMapRepository(db []models.Section) SectionRepository {
	// Calculate initial maxId
	var maxId int
	dbMap := make(map[int]models.Section)
	for _, e := range db {
		if e.ID > maxId {
			maxId = e.ID
		}
		dbMap[e.ID] = models.Section{
			ID: e.ID,
			SectionAttributes: models.SectionAttributes{
				SectionNumber:      e.SectionNumber,
				CurrentTemperature: e.CurrentTemperature,
				MinimumTemperature: e.MinimumTemperature,
				CurrentCapacity:    e.CurrentCapacity,
				MinimumCapacity:    e.MinimumCapacity,
				MaximumCapacity:    e.MaximumCapacity,
				WarehouseID:        e.WarehouseID,
				ProductTypeID:      e.ProductTypeID,
				ProductBatches:     e.ProductBatches,
			},
		}
	}
	return &SectionRepositoryMap{
		db:     dbMap,
		lastId: maxId,
	}
}

// SectionMap is a struct that represents a section repository
type SectionRepositoryMap struct {
	// db is a map of sections
	db     map[int]models.Section
	lastId int
}

func (r *SectionRepositoryMap) GetAll() (s map[int]models.Section, err error) {
	s = make(map[int]models.Section)

	// copy db
	for key, value := range r.db {
		s[key] = value
	}

	return
}

func (r *SectionRepositoryMap) GetByID(id int) (s models.Section, err error) {
	if s, ok := r.db[id]; !ok {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section with id %d not found", id)
		return models.Section{}, svcErr
	} else {
		return s, nil
	}
}

func (r *SectionRepositoryMap) Create(section models.Section) (s models.Section, err error) {
	for _, value := range r.db {
		if value.SectionNumber == section.SectionNumber {
			svcErr := pkg.ServiceErrors[pkg.ErrConflict]
			svcErr.InternalError = fmt.Errorf("section with section number %d already exists", section.SectionNumber)
			return models.Section{}, svcErr
		}
	}

	r.lastId++
	section.ID = r.lastId
	r.db[section.ID] = section
	return section, nil

}

func (r *SectionRepositoryMap) Update(id int, section models.Section) (s models.Section, err error) {

	existing, ok := r.db[id]
	if !ok {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section with id %d not found", id)
		return models.Section{}, svcErr
	}

	if section.SectionNumber != 0 {
		existing.SectionNumber = section.SectionNumber
	}
	if section.CurrentTemperature != 0.0 {
		existing.CurrentTemperature = section.CurrentTemperature
	}
	if section.MinimumTemperature != 0.0 {
		existing.MinimumTemperature = section.MinimumTemperature
	}
	if section.CurrentCapacity != 0 {
		existing.CurrentCapacity = section.CurrentCapacity
	}
	if section.MinimumCapacity != 0 {
		existing.MinimumCapacity = section.MinimumCapacity
	}
	if section.MaximumCapacity != 0 {
		existing.MaximumCapacity = section.MaximumCapacity
	}
	if section.WarehouseID != 0 {
		existing.WarehouseID = section.WarehouseID
	}
	if section.ProductTypeID != 0 {
		existing.ProductTypeID = section.ProductTypeID
	}
	if section.ProductBatches != nil {
		existing.ProductBatches = section.ProductBatches
	}

	r.db[id] = existing

	return existing, nil
}

func (r *SectionRepositoryMap) Delete(id int) (err error) {

	if _, ok := r.db[id]; !ok {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section with id %d not found", id)
		return svcErr
	}
	delete(r.db, id)
	return nil
}
