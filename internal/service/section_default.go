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
	return sv.rp.Update(id, section)
}

// Delete is a method that deletes a section
func (sv *SectionDefault) Delete(id int) (err error) {
	return sv.rp.Delete(id)
}
