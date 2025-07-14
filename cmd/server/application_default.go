package server

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress string
}

// NewServerChi is a function that returns a new instance of ServerChi
func NewServerChi(cfg *ConfigServerChi) *ServerChi {
	// default values
	defaultConfig := &ConfigServerChi{
		ServerAddress: ":8080",
	}
	if cfg != nil {
		if cfg.ServerAddress != "" {
			defaultConfig.ServerAddress = cfg.ServerAddress
		}
	}

	return &ServerChi{
		serverAddress: defaultConfig.ServerAddress,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress string
	// buyerLoaderFilePath is the path to the file that contains the buyers
	buyerLoaderFilePath string
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {

	db, err := initMySQL()
	if err != nil {
		return err
	}

	/*	// create product handler with dependences
		prodHandler, err := a.BuildProductHandler(&productDb, &productTypeDb, &sectionDb, &sellerDb)
		if err != nil {
			return err
		}

		// dependencies
		// - loader
		// - repository
		employeeHandler, err := a.BuildemployeeHandler(&employeeDb, &warehouseDb)
		if err != nil {
			return err
		}

		// handler warehouse

		hdWarehouse, err := a.BuildWarehouseHandler(&sectionDb, &employeeDb, &warehouseDb)
		if err != nil {
			return err
		}*/

	// - handler
	sectionHd, err := a.BuildSectionHandler(db)
	if err != nil {
		return err
	}

	// create seller handler with dependences
	hdSeller, err := a.BuildSellerHandler(db)
	if err != nil {
		return err
	}
	/*
		buyerHd, err := a.BuildBuyerHandler(&buyerDb)
		if err != nil {
			return err
		}*/

	// create seller handler with dependences
	hdLocalities, err := a.BuildLocalityHandler(db)
	if err != nil {
		return err
	}
	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	rt.Route("/api/v1", func(rt chi.Router) {
		rt.Route("/warehouses", func(rt chi.Router) {
			// - GET /warehouses
			rt.Get("/", hdWarehouse.GetAll())
			rt.Get("/{id}", hdWarehouse.GetOne())
			rt.Post("/", hdWarehouse.Add())
			rt.Patch("/{id}", hdWarehouse.Update())
			rt.Delete("/{id}", hdWarehouse.Delete())
		})

		rt.Route("/employees", func(r chi.Router) {
			r.Get("/", employeeHandler.GetAllEmployees)
			r.Get("/{id}", employeeHandler.GetEmployee)
			r.Post("/", employeeHandler.CreateEmployee)
			r.Patch("/{id}", employeeHandler.UpdateEmployee)
			r.Delete("/{id}", employeeHandler.DeleteEmployee)
		})

		rt.Route("/products", func(r chi.Router) {
			r.Get("/", prodHandler.GetAll())
			r.Get("/{id}", prodHandler.GetById())
			r.Post("/", prodHandler.Create())
			r.Delete("/{id}", prodHandler.Delete())
			r.Patch("/{id}", prodHandler.Patch())
		})

		rt.Route("/sections", func(rt chi.Router) {
			// - GET /vehicles
			rt.Get("/", sectionHd.GetAll())
			// - GET /vehicles/{id}
			rt.Get("/{id}", sectionHd.GetByID())
			// - POST /vehicles
			rt.Post("/", sectionHd.PostSection())
			// - PUT /vehicles/{id}
			rt.Patch("/{id}", sectionHd.Update())
			// - DELETE /vehicles/{id}
			rt.Delete("/{id}", sectionHd.Delete())
		})

		rt.Route("/sellers", func(rt chi.Router) {
			// - GET /sellers
			rt.Get("/", hdSeller.GetAll())
			rt.Get("/{id}", hdSeller.GetById())
			rt.Post("/", hdSeller.Create())
			rt.Patch("/{id}", hdSeller.Update())
			rt.Delete("/{id}", hdSeller.Delete())
		})

		rt.Route("/buyers", func(r chi.Router) {
			r.Get("/", buyerHd.GetAll())
			r.Get("/{id}", buyerHd.GetByID())
			r.Post("/", buyerHd.Create())
			r.Patch("/{id}", buyerHd.Update())
			r.Delete("/{id}", buyerHd.Delete())
		})

		rt.Route("/localities", func(rt chi.Router) {
			rt.Post("/", hdLocalities.Create())
		})

	})
	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func (a *ServerChi) createMaps() (map[int]models.Section, map[int]models.Product, map[int]models.Employee, map[int]models.ProductType, map[int]models.Warehouse, map[int]models.Seller, map[int]models.Buyer, error) {
	sectionLd := loader.NewLoaderGeneric[models.Section]()
	sectionDb, err := sectionLd.LoadFromJSON("docs/db/sections.json")
	if err != nil {
		fmt.Println("err1")
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	productLd := loader.NewLoaderGeneric[models.Product]()
	productDb, err := productLd.LoadFromJSON("docs/db/products.json")
	if err != nil {
		fmt.Println("err2")
		return nil, nil, nil, nil, nil, nil, nil, err
	}
	employeeLd := loader.NewLoaderGeneric[models.Employee]()
	employeeDb, err := employeeLd.LoadFromJSON("docs/db/employees.json")
	if err != nil {
		fmt.Println("err3")
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	productTypeLd := loader.NewLoaderGeneric[models.ProductType]()
	productTypeDb, err := productTypeLd.LoadFromJSON("docs/db/productTypes.json")
	if err != nil {
		fmt.Println("err4")
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	warehouseLd := loader.NewLoaderGeneric[models.Warehouse]()
	warehouseDb, err := warehouseLd.LoadFromJSON("docs/db/warehouse_500.json")
	if err != nil {
		fmt.Println("err5")
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	sellerLd := loader.NewLoaderGeneric[models.Seller]()
	sellerDb, err := sellerLd.LoadFromJSON("docs/db/sellers.json")
	if err != nil {
		fmt.Println("err6")
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	buyerLd := loader.NewLoaderGeneric[models.Buyer]()
	buyerDb, err := buyerLd.LoadFromJSON("docs/db/buyers.json")
	if err != nil {
		fmt.Println("err7")
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	return sectionDb, productDb, employeeDb, productTypeDb, warehouseDb, sellerDb, buyerDb, nil
}

func (a *ServerChi) BuildSectionHandler(db *sql.DB) (*handler.SectionDefault, error) {

	// - repository
	sectionRp := repository.NewSectionMapRepository(db)
	// - service
	sectionSv := service.NewSectionDefault(sectionRp)
	// - handler
	sectionHd := handler.NewSectionDefault(sectionSv)
	return sectionHd, nil
}

func (a *ServerChi) BuildWarehouseHandler(sectionDb *map[int]models.Section, employeeDb *map[int]models.Employee, warehouseDb *map[int]models.Warehouse) (*handler.WarehouseDefault, error) {

	// - repository
	warehouseRp := repository.NewWarehouseMap(sectionDb, employeeDb, warehouseDb)
	// - service
	warehouseSv := service.NewWarehouseDefault(warehouseRp)
	// - handler
	warehouseHd := handler.NewWarehouseDefault(warehouseSv)
	return warehouseHd, nil
}

func (a *ServerChi) BuildProductHandler(productDb *map[int]models.Product, productTypeDb *map[int]models.ProductType, dbSection *map[int]models.Section, dbSellers *map[int]models.Seller) (prodHandler *handler.ProductDefault, err error) {

	// - repository
	productRp := repository.NewProductMap(productDb, productTypeDb, dbSection, dbSellers)

	// - service
	productSv := service.NewProductDefault(productRp)

	// - handler
	prodHandler = handler.NewProductDefault(productSv)
	return prodHandler, nil
}

func (*ServerChi) BuildemployeeHandler(db *sql.DB) (*handler.EmployeeHandler, error) {

	employeeRepository := repository.NewEmployeeRepository(db)
	// - service
	employeeService := service.NewEmployeeServiceDefault(employeeRepository)
	// - handler
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	return employeeHandler, nil
}

func (a *ServerChi) BuildSellerHandler(sellerDb *sql.DB) (*handler.SellerDefault, error) {
	// - repository
	rpSeller := repository.NewSellerSql(sellerDb)

	// - service
	svSeller := service.NewSellerDefault(rpSeller)

	// - handler
	hdSeller := handler.NewSellerDefault(svSeller)
	return hdSeller, nil
}

func (a *ServerChi) BuildBuyerHandler(buyerDb *map[int]models.Buyer) (*handler.BuyerDefault, error) {
	// - repository
	rpBuyer := repository.NewBuyerMap(buyerDb)

	// - service
	svBuyer := service.NewBuyerDefault(rpBuyer)

	// - handler
	hdBuyer := handler.NewBuyerDefault(svBuyer)
	return hdBuyer, nil
}

// Mejorar la función initMySQL existente
func initMySQL() (*sql.DB, error) {

	// Obtener las credenciales desde variables de entorno
	user := getEnvOrDefault("DB_USER", "root")
	password := getEnvOrDefault("DB_PASSWORD", "asda1125")
	host := getEnvOrDefault("DB_HOST", "localhost")
	port := getEnvOrDefault("DB_PORT", "3306")
	dbname := getEnvOrDefault("DB_NAME", "db_test")

	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", host, port),
		DBName:               dbname,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	// Verificar que la conexión funciona
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Println("Connected to MySQL database successfully")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (a *ServerChi) BuildLocalityHandler(localityDb *sql.DB) (*handler.LocalityDefault, error) {
	// - repository
	rpLocality := repository.NewLocalitySql(localityDb)

	// - service
	svLocality := service.NewLocalityDefault(rpLocality)

	// - handler
	hdLocality := handler.NewLocalityDefault(svLocality)
	return hdLocality, nil
}
