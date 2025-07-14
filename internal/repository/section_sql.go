package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
)

func NewSectionSqlRepository(db *sql.DB) SectionRepository {

	return &SectionRepositorySql{
		db: db,
	}
}

// SectionSql is a struct that represents a section repository
type SectionRepositorySql struct {
	// db is a map of sections
	db *sql.DB
}

func (r *SectionRepositorySql) GetAll() (s map[int]models.Section, err error) {
	s = make(map[int]models.Section)
	rows, err := r.db.Query(`
		SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id
		FROM sections
	`)
	if errors.Is(err, sql.ErrNoRows) {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("no se encontraron sections")
		return map[int]models.Section{}, svcErr
	}
	if err != nil {
		return map[int]models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	defer rows.Close()

	for rows.Next() {
		var section models.Section
		err := rows.Scan(
			&section.ID,
			&section.SectionNumber,
			&section.CurrentTemperature,
			&section.CurrentCapacity,
			&section.MinimumTemperature,
			&section.MinimumCapacity,
			&section.ProductTypeID,
			&section.WarehouseID,
		)
		if err != nil {
			return nil, err
		}
		s[section.ID] = section
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *SectionRepositorySql) GetByID(id int) (s models.Section, err error) {
	err = r.db.QueryRow(`
		SELECT id, section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id
		FROM sections
		WHERE id = ?
	`, id).Scan(
		&s.ID,
		&s.SectionNumber,
		&s.CurrentTemperature,
		&s.CurrentCapacity,
		&s.MinimumTemperature,
		&s.MinimumCapacity,
		&s.ProductTypeID,
		&s.WarehouseID,
	)
	if errors.Is(err, sql.ErrNoRows) {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section con id %d no encontrada", id)
		return models.Section{}, svcErr
	}
	if err != nil {
		return models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	return s, nil
}

func (r *SectionRepositorySql) Create(section models.Section) (s models.Section, err error) {
	result, err := r.db.Exec(`
		INSERT INTO sections (
			section_number, current_temperature, current_capacity, minimum_temperature, minimum_capacity, product_type_id, warehouse_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		section.SectionNumber,
		section.CurrentTemperature,
		section.CurrentCapacity,
		section.MinimumTemperature,
		section.MinimumCapacity,
		section.ProductTypeID,
		section.WarehouseID,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			svcErr := pkg.ServiceErrors[pkg.ErrConflict]
			svcErr.InternalError = fmt.Errorf("valor duplicado: %v", mysqlErr.Message)
			return models.Section{}, svcErr
		}
		return models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	section.ID = int(id)
	return section, nil
}

func (r *SectionRepositorySql) Update(id int, section models.Section) (s models.Section, err error) {
	result, err := r.db.Exec(`
  UPDATE sections SET
   section_number = ?,
   current_temperature = ?,
   current_capacity = ?,
   minimum_temperature = ?,
   minimum_capacity = ?,
   product_type_id = ?,
   warehouse_id = ?
  WHERE id = ?
 `,
		section.SectionNumber,
		section.CurrentTemperature,
		section.CurrentCapacity,
		section.MinimumTemperature,
		section.MinimumCapacity,
		section.ProductTypeID,
		section.WarehouseID,
		id,
	)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			svcErr := pkg.ServiceErrors[pkg.ErrConflict]
			svcErr.InternalError = fmt.Errorf("valor duplicado: %v", mysqlErr.Message)
			return models.Section{}, svcErr
		}
		return models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Section{}, pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	if rowsAffected == 0 {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section con id %d no encontrada", id)
		return models.Section{}, svcErr
	}

	section.ID = id
	return section, nil
}

func (r *SectionRepositorySql) Delete(id int) (err error) {
	result, err := r.db.Exec(`
  DELETE FROM sections WHERE id = ?
 `, id)
	if err != nil {
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return pkg.ServiceErrors[pkg.ErrInternalServer]
	}
	if rowsAffected == 0 {
		svcErr := pkg.ServiceErrors[pkg.ErrNotFound]
		svcErr.InternalError = fmt.Errorf("section con id %d no encontrada", id)
		return svcErr
	}
	return nil
}
