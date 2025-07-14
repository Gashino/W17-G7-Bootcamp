package service

import (
	"app/internal/repository"
	"app/pkg/models"
)

// NewSectionDefault is a function that returns a new instance of SectionDefault
func NewSectionDefault(rp repository.SectionRepository) *SectionDefault {
	return &SectionDefault{rp: rp}
}

// SectionDefault is a struct that represents the default service for sections
type SectionDefault struct {
	// rp is the repository that will be used by the service
	rp repository.SectionRepository
}

// GetAll is a method that returns all the sections
func (sv *SectionDefault) GetAll() (s map[int]models.Section, err error) {
	return sv.rp.GetAll()
}

// GetByID is a method that returns a section by id
func (sv *SectionDefault) GetByID(id int) (s models.Section, err error) {
	return sv.rp.GetByID(id)
}

// Create is a method that creates a new section
func (sv *SectionDefault) Create(section models.Section) (s models.Section, err error) {
	return sv.rp.Create(section)
}

// Update is a method that updates a section
func (sv *SectionDefault) Update(id int, section models.Section) (s models.Section, err error) {
	sec, err := sv.rp.GetByID(id)
	if err != nil {
		return sec, err
	}

	if section.SectionNumber != 0 {
		sec.SectionNumber = section.SectionNumber
	}
	if section.CurrentTemperature != 0 {
		sec.CurrentTemperature = section.CurrentTemperature
	}
	if section.MinimumTemperature != 0 {
		sec.MinimumTemperature = section.MinimumTemperature
	}
	if section.CurrentCapacity != 0 {
		sec.CurrentCapacity = section.CurrentCapacity
	}
	if section.MinimumCapacity != 0 {
		sec.MinimumCapacity = section.MinimumCapacity
	}
	if section.MaximumCapacity != 0 {
		sec.MaximumCapacity = section.MaximumCapacity
	}
	if section.WarehouseID != 0 {
		sec.WarehouseID = section.WarehouseID
	}
	if section.ProductTypeID != 0 {
		sec.ProductTypeID = section.ProductTypeID
	}

	return sv.rp.Update(id, sec)
}

// Delete is a method that deletes a section
func (sv *SectionDefault) Delete(id int) (err error) {
	return sv.rp.Delete(id)
}

func (sv *SectionDefault) ReportProductsBySection(ptr *int) (s []models.SectionReport, err error) {
	if ptr == nil {
		return sv.rp.GetReportProductsAllSections()
	}
	return sv.rp.GetReportProductsBySection(*ptr)
}
