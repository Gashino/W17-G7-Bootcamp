package repository

import (
	"app/pkg"
	"app/pkg/models"
	"fmt"
)

func NewSectionMapRepository(dbSec *map[int]models.Section, dbType *map[int]models.ProductType, dbWarehouse *map[int]models.Warehouse) SectionRepository {
	// Calculate initial maxId

	return &SectionRepositoryMap{
		dbSections:     dbSec,
		dbProductTypes: dbType,
		dbWarehouses:   dbWarehouse,
		lastId:         len(*dbSec),
	}
}

// SectionMap is a struct that represents a section repository
type SectionRepositoryMap struct {
	// db is a map of sections
	dbSections     *map[int]models.Section
	dbProductTypes *map[int]models.ProductType
	dbWarehouses   *map[int]models.Warehouse
	lastId         int
}

func (r *SectionRepositoryMap) GetAll() (s map[int]models.Section, err error) {
	s = make(map[int]models.Section)

	// copy db
	for key, value := range *r.dbSections {
		s[key] = value
	}

	return
}

func (r *SectionRepositoryMap) GetByID(id int) (s models.Section, err error) {
	if s, ok := (*r.dbSections)[id]; !ok {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section with id %d not found", id)
		return models.Section{}, svcErr
	} else {
		return s, nil
	}
}

func (r *SectionRepositoryMap) Create(section models.Section) (s models.Section, err error) {

	_, productTypeExists := (*r.dbProductTypes)[section.ProductTypeID]

	if !productTypeExists {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("product type with id %d not found", section.ProductTypeID)
		return models.Section{}, svcErr
	}

	for _, value := range *r.dbSections {
		if value.SectionNumber == section.SectionNumber {
			svcErr := pkg.ServiceErrors[pkg.ErrConflict]
			svcErr.InternalError = fmt.Errorf("section with section number %d already exists", section.SectionNumber)
			return models.Section{}, svcErr
		}
	}

	_, warehouseExists := (*r.dbWarehouses)[section.WarehouseID]
	if !warehouseExists {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("warehouse with id %d not found", section.WarehouseID)
		return models.Section{}, svcErr
	}

	r.lastId++
	section.ID = r.lastId
	(*r.dbSections)[section.ID] = section
	return section, nil

}

func (r *SectionRepositoryMap) Update(id int, section models.Section) (s models.Section, err error) {

	existing, ok := (*r.dbSections)[id]
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
		_, warehouseExists := (*r.dbWarehouses)[section.WarehouseID]
		if !warehouseExists {
			svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
			svcErr.InternalError = fmt.Errorf("warehouse with id %d not found", section.WarehouseID)
			return models.Section{}, svcErr
		}
		existing.WarehouseID = section.WarehouseID
	}
	if section.ProductTypeID != 0 {
		existing.ProductTypeID = section.ProductTypeID
	}
	if section.ProductBatches != nil {
		existing.ProductBatches = section.ProductBatches
	}

	(*r.dbSections)[id] = existing

	return existing, nil
}

func (r *SectionRepositoryMap) Delete(id int) (err error) {

	if _, ok := (*r.dbSections)[id]; !ok {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section with id %d not found", id)
		return svcErr
	}
	delete(*r.dbSections, id)
	return nil
}
