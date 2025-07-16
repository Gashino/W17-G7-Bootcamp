package repository

import (
	"app/pkg"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// Variables para controlar el registro del driver txdb para employee tests
var (
	employeeTxdbRegistered bool = false
	employeeTxdbMutex      sync.Mutex
)

// Estructura para leer la configuración de la base de datos desde config.yml
type EmployeeConfig struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Configurar la base de datos para pruebas de employee
func setupEmployeeTxDB() string {
	// Usar un mutex para evitar condiciones de carrera al registrar el driver
	employeeTxdbMutex.Lock()
	defer employeeTxdbMutex.Unlock()

	// Evitar registrar el driver más de una vez
	if !employeeTxdbRegistered {
		// Leer la configuración desde config.yml
		config := loadEmployeeConfig()

		// Configurar la conexión MySQL usando los valores del config.yml
		cfg := mysql.Config{
			User:                 config.Database.User,
			Passwd:               config.Database.Password,
			Net:                  "tcp",
			Addr:                 fmt.Sprintf("%s:%s", config.Database.Host, config.Database.Port),
			DBName:               config.Database.Name,
			ParseTime:            true,
			AllowNativePasswords: true,
		}

		// Verificar si el driver ya está registrado
		for _, driver := range sql.Drivers() {
			if driver == "employee_txdb" {
				// El driver ya está registrado, no necesitamos registrarlo de nuevo
				employeeTxdbRegistered = true
				return "employee_txdb"
			}
		}

		// Registrar el driver si no está registrado
		txdb.Register("employee_txdb", "mysql", cfg.FormatDSN())
		employeeTxdbRegistered = true
	}

	return "employee_txdb"
}

// Cargar la configuración desde config.yml para employee tests
func loadEmployeeConfig() EmployeeConfig {
	var config EmployeeConfig

	// Intentar leer el archivo de configuración
	data, err := os.ReadFile("../../config.yml")
	if err != nil {
		// Si hay un error, mostrar un mensaje y terminar el test
		panic(fmt.Sprintf("Error al leer el archivo config.yml: %v\nAsegúrate de que el archivo config.yml existe en la raíz del proyecto", err))
	}

	// Parsear el archivo YAML
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		// Si hay un error al parsear el YAML, mostrar un mensaje y terminar el test
		panic(fmt.Sprintf("Error al parsear el archivo config.yml: %v\nVerifica que el formato del archivo sea correcto", err))
	}

	// Verificar que la configuración de la base de datos esté completa
	if config.Database.User == "" || config.Database.Host == "" || config.Database.Port == "" || config.Database.Name == "" {
		panic("La configuración de la base de datos en config.yml está incompleta")
	}

	return config
}

// Configurar la base de datos de prueba para cada test de employee
func setupEmployeeTestDB(t *testing.T) (*EmployeeRepositoryMap, *sql.DB) {
	// Usar un nombre único para cada conexión de test
	driver := setupEmployeeTxDB()
	db, err := sql.Open(driver, fmt.Sprintf("employee_test_%s", t.Name()))
	require.NoError(t, err)

	// Verificar que la conexión funciona
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	return &EmployeeRepositoryMap{db: db}, db
}

// Crear un empleado completo para pruebas
func createTestEmployee() models.Employee {
	// Usar timestamp para generar un card number único
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano())
	cardNumber := timestamp[len(timestamp)-8:] // Tomar los últimos 8 dígitos

	return models.Employee{
		CardNumberID: cardNumber,
		FirstName:    "John",
		LastName:     "Doe",
		WarehouseID:  1, // Usar un warehouse_id válido según el script de base de datos
	}
}

func TestEmployeeSQL_FindAll(t *testing.T) {
	repo, db := setupEmployeeTestDB(t)
	defer db.Close()

	// Act
	employees, err := repo.FindAll()

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, employees)
	// Verificamos que se devuelvan los empleados de seed
	assert.GreaterOrEqual(t, len(employees), 1, "Deberían existir al menos algunos empleados en la base de datos de prueba")
}

func TestEmployeeSQL_FindById(t *testing.T) {
	repo, db := setupEmployeeTestDB(t)
	defer db.Close()

	t.Run("existing employee", func(t *testing.T) {
		// Act
		employee, err := repo.FindById(1)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, employee)
		assert.Equal(t, 1, employee.ID)
	})

	t.Run("non-existing employee", func(t *testing.T) {
		// Act
		employee, err := repo.FindById(999)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Employee{}, employee)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})
}

func TestEmployeeSQL_Save(t *testing.T) {
	t.Run("create new employee successfully", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange
		newEmployee := createTestEmployee()

		// Act
		createdEmployee, err := repo.Save(newEmployee)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, createdEmployee)
		assert.NotZero(t, createdEmployee.ID)
		assert.Equal(t, newEmployee.CardNumberID, createdEmployee.CardNumberID)
		assert.Equal(t, newEmployee.FirstName, createdEmployee.FirstName)
		assert.Equal(t, newEmployee.LastName, createdEmployee.LastName)
		assert.Equal(t, newEmployee.WarehouseID, createdEmployee.WarehouseID)

		// Verificar que el empleado fue creado en la base de datos
		savedEmployee, err := repo.FindById(createdEmployee.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdEmployee.CardNumberID, savedEmployee.CardNumberID)
	})

	t.Run("create employee with duplicate card number", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un primer empleado
		employee1 := createTestEmployee()
		createdEmployee, err := repo.Save(employee1)
		require.NoError(t, err)
		require.NotEmpty(t, createdEmployee)

		// Intentar crear un segundo empleado con el mismo card number
		employee2 := createTestEmployee()
		employee2.CardNumberID = createdEmployee.CardNumberID // Usar el mismo card number

		// Act
		duplicateEmployee, err := repo.Save(employee2)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Employee{}, duplicateEmployee)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, err.(pkg.ServiceError).Code)
	})

	t.Run("create employee with invalid warehouse id", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un empleado con un warehouse_id inválido
		invalidEmployee := createTestEmployee()
		invalidEmployee.WarehouseID = 999 // ID que no existe

		// Act
		createdEmployee, err := repo.Save(invalidEmployee)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, models.Employee{}, createdEmployee)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})
}

func TestEmployeeSQL_Update(t *testing.T) {
	t.Run("update existing employee", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un empleado primero
		newEmployee := createTestEmployee()
		createdEmployee, err := repo.Save(newEmployee)
		require.NoError(t, err)
		require.NotEmpty(t, createdEmployee)

		// Actualizar el empleado
		updatedEmployee := models.Employee{
			CardNumberID: "87654321",
			FirstName:    "Jane",
			LastName:     "Smith",
			WarehouseID:  1,
		}

		// Act
		resultEmployee, err := repo.Update(updatedEmployee, createdEmployee.ID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, updatedEmployee.CardNumberID, resultEmployee.CardNumberID)
		assert.Equal(t, updatedEmployee.FirstName, resultEmployee.FirstName)
		assert.Equal(t, updatedEmployee.LastName, resultEmployee.LastName)

		// Verificar que el empleado fue actualizado
		savedEmployee, err := repo.FindById(createdEmployee.ID)
		assert.NoError(t, err)
		assert.Equal(t, updatedEmployee.CardNumberID, savedEmployee.CardNumberID)
		assert.Equal(t, updatedEmployee.FirstName, savedEmployee.FirstName)
		assert.Equal(t, updatedEmployee.LastName, savedEmployee.LastName)
	})

	t.Run("update with duplicate card number", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear dos empleados
		employee1 := createTestEmployee()
		employee2 := createTestEmployee()

		createdEmployee1, err := repo.Save(employee1)
		require.NoError(t, err)
		createdEmployee2, err := repo.Save(employee2)
		require.NoError(t, err)

		// Intentar actualizar el empleado 2 con el card number del empleado 1
		updatedEmployee := models.Employee{
			CardNumberID: createdEmployee1.CardNumberID, // Card number duplicado (mismo que createdEmployee1)
			FirstName:    "Updated",
			LastName:     "Name",
			WarehouseID:  1,
		}

		// Act
		_, err = repo.Update(updatedEmployee, createdEmployee2.ID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrConflict].Code, err.(pkg.ServiceError).Code)

		// Verificar que el empleado no fue actualizado
		savedEmployee, err := repo.FindById(createdEmployee2.ID)
		assert.NoError(t, err)
		assert.Equal(t, createdEmployee2.CardNumberID, savedEmployee.CardNumberID)

		// Verificar que el empleado 1 sigue existiendo
		_, err = repo.FindById(createdEmployee1.ID)
		assert.NoError(t, err)
	})

	t.Run("update with invalid warehouse id", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un empleado primero
		newEmployee := createTestEmployee()
		createdEmployee, err := repo.Save(newEmployee)
		require.NoError(t, err)

		// Intentar actualizar con warehouse_id inválido
		updatedEmployee := models.Employee{
			CardNumberID: "87654321",
			FirstName:    "Jane",
			LastName:     "Smith",
			WarehouseID:  999, // ID que no existe
		}

		// Act
		_, err = repo.Update(updatedEmployee, createdEmployee.ID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})
}

func TestEmployeeSQL_Delete(t *testing.T) {
	t.Run("delete existing employee", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un empleado primero
		newEmployee := createTestEmployee()
		createdEmployee, err := repo.Save(newEmployee)
		require.NoError(t, err)
		require.NotEmpty(t, createdEmployee)

		// Act
		err = repo.Delete(createdEmployee.ID)

		// Assert
		assert.NoError(t, err)

		// Verificar que el empleado fue eliminado
		_, err = repo.FindById(createdEmployee.ID)
		assert.Error(t, err)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound], err)
	})

	t.Run("delete non-existing employee", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Act
		err := repo.Delete(9999) // ID que no existe

		// Assert
		// Note: La implementación actual no verifica si el empleado existe antes de eliminar
		// por lo que no devuelve error incluso si el ID no existe
		assert.NoError(t, err)
	})
}

func TestEmployeeSQL_ReportInboundOrdersCountByEmployee(t *testing.T) {
	t.Run("get all employees report", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(nil)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reports)
		// Los empleados de seed deberían estar en el reporte
		assert.GreaterOrEqual(t, len(reports), 1, "Deberían existir al menos algunos empleados en el reporte")
	})

	t.Run("get report for specific employee", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange - Crear un empleado primero
		newEmployee := createTestEmployee()
		createdEmployee, err := repo.Save(newEmployee)
		require.NoError(t, err)
		require.NotEmpty(t, createdEmployee)

		// Act
		employeeID := createdEmployee.ID
		reports, err := repo.ReportInboundOrdersCountByEmployee(&employeeID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, reports)
		assert.Len(t, reports, 1)
		assert.Equal(t, createdEmployee.ID, reports[0].ID)
		assert.Equal(t, createdEmployee.CardNumberID, reports[0].CardNumberID)
		assert.Equal(t, createdEmployee.FirstName, reports[0].FirstName)
		assert.Equal(t, createdEmployee.LastName, reports[0].LastName)
		// El empleado recién creado no tiene inbound orders asociadas
		assert.Equal(t, 0, reports[0].InboundOrdersCount)
	})

	t.Run("get report for non-existing employee", func(t *testing.T) {
		repo, db := setupEmployeeTestDB(t)
		defer db.Close()

		// Arrange
		nonExistingID := 9999

		// Act
		reports, err := repo.ReportInboundOrdersCountByEmployee(&nonExistingID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, reports)
		assert.Equal(t, pkg.ServiceErrors[pkg.ErrNotFound].Code, err.(pkg.ServiceError).Code)
	})
}
