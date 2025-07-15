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
	"gopkg.in/yaml.v2"

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

// ConfigDB representa la estructura de la configuración de la base de datos
type ConfigDB struct {
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"database"`
}

// Lee la configuración desde un archivo YAML cuyo path se obtiene de la variable de entorno CONFIG_PATH
func loadConfig() (*ConfigDB, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yml"
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var cfg ConfigDB
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {

	db, err := initMySQL()
	if err != nil {
		return err
	}

	// create product handler with dependences
	prodHandler, err := a.BuildProductHandler(db)
	if err != nil {
		return err
	}

	// Create warehouse handler with dependencies
	warehouseHd, err := a.BuildWarehouseHandler(db)
	if err != nil {
		return err
	}

	// Create employee handler with dependencies
	employeeHandler, err := a.BuildemployeeHandler(db)
	if err != nil {
		return err
	}

	// Create section handler with dependencies
	sectionHd, err := a.BuildSectionHandler(db)
	if err != nil {
		return err
	}

	productBatchHd, err := a.BuidProductBatchHandler(db)
	if err != nil {
		return err
	}

	carryHd, err := a.BuildCarryHandler(db)
	if err != nil {
		return err
	}

	// create seller handler with dependences
	hdSeller, err := a.BuildSellerHandler(db)
	if err != nil {
		return err
	}

	// Create buyer handler with SQL database (migrated to SQL)
	buyerHd, err := a.BuildBuyerHandler(db)
	if err != nil {
		return err
	}

	// Create purchase order handler with SQL database
	purchaseOrderHd, err := a.BuildPurchaseOrderHandler(db)
	if err != nil {
		return err
	}

	// create seller handler with dependences
	hdLocalities, err := a.BuildLocalityHandler(db)
	if err != nil {
		return err
	}
	// Create inbound order handler with dependencies
	inboundOrderHd, err := a.BuildInboundOrderHandler(db)
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
			rt.Get("/", warehouseHd.GetAll())
			rt.Get("/{id}", warehouseHd.GetOne())
			rt.Post("/", warehouseHd.Add())
			rt.Patch("/{id}", warehouseHd.Update())
			rt.Delete("/{id}", warehouseHd.Delete())
		})

		rt.Route("/carries", func(rt chi.Router) {
			rt.Post("/", carryHd.Create())
			rt.Get("/localities/reportCarries", carryHd.SearchByLocality())
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
			rt.Get("/", sectionHd.GetAll())
			rt.Get("/{id}", sectionHd.GetByID())
			rt.Post("/", sectionHd.PostSection())
			rt.Patch("/{id}", sectionHd.Update())
			rt.Delete("/{id}", sectionHd.Delete())
			// - GET /reportProducts
			rt.Get("/reportProducts", sectionHd.ReportProducts())
		})

		rt.Route("/productBatches", func(rt chi.Router) {
			// - POST /productBatches
			rt.Post("/", productBatchHd.CreateBatch())
		})

		rt.Route("/sellers", func(rt chi.Router) {
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
			r.Get("/reportPurchaseOrders", buyerHd.GetPurchaseOrdersReport())
		})

		rt.Route("/localities", func(rt chi.Router) {
			rt.Get("/reportSellers/{id}", hdLocalities.SellersByLocality())
			rt.Post("/", hdLocalities.Create())
		})

		rt.Route("/purchaseOrders", func(r chi.Router) {
			r.Post("/", purchaseOrderHd.Create())
		})

		rt.Route("/inboundOrders", func(r chi.Router) {
			r.Post("/", inboundOrderHd.CreateInboundOrder)
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
	sectionRp := repository.NewSectionSqlRepository(db)
	// - service
	sectionSv := service.NewSectionDefault(sectionRp)
	// - handler
	sectionHd := handler.NewSectionDefault(sectionSv)
	return sectionHd, nil
}

func (a *ServerChi) BuidProductBatchHandler(db *sql.DB) (*handler.ProductBatchDefault, error) {
	// - repository
	productBatchRp := repository.NewProductBatchSqlRepository(db)
	// - service
	productBatchSv := service.NewProductBatchDefault(productBatchRp)
	// - handler
	productBatchHd := handler.NewProductBatchDefault(productBatchSv)
	return productBatchHd, nil
}

func (a *ServerChi) BuildWarehouseHandler(db *sql.DB) (*handler.WarehouseDefault, error) {

	// - repository
	warehouseRp := repository.NewWarehouseSql(db)
	// - service
	warehouseSv := service.NewWarehouseDefault(warehouseRp)
	// - handler
	warehouseHd := handler.NewWarehouseDefault(warehouseSv)
	return warehouseHd, nil
}

func (a *ServerChi) BuildCarryHandler(db *sql.DB) (*handler.CarryDefault, error) {
	// - repository
	carryRp := repository.NewCarrySql(db)
	// - service
	carrySv := service.NewCarryDefault(carryRp)
	// - handler
	carryHd := handler.NewCarryDefault(carrySv)
	return carryHd, nil
}

func (a *ServerChi) BuildProductHandler(db *sql.DB) (prodHandler *handler.ProductDefault, err error) {

	// - repository
	productRp := repository.NewProductSqlRepository(db)

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

func (a *ServerChi) BuildBuyerHandler(db *sql.DB) (*handler.BuyerHandler, error) {
	// - repository (now using SQL instead of map)
	rpBuyer := repository.NewBuyerSQL(db)

	// - service
	svBuyer := service.NewBuyerDefault(rpBuyer)

	// - handler
	hdBuyer := handler.NewBuyerHandler(svBuyer)
	return hdBuyer, nil
}

func (a *ServerChi) BuildPurchaseOrderHandler(db *sql.DB) (*handler.PurchaseOrderHandler, error) {
	// - repository
	purchaseOrderRp := repository.NewPurchaseOrderSQL(db)

	// - service
	purchaseOrderSv := service.NewPurchaseOrderDefault(purchaseOrderRp)

	// - handler
	purchaseOrderHd := handler.NewPurchaseOrderHandler(purchaseOrderSv)
	return purchaseOrderHd, nil
}

// Mejorar la función initMySQL existente
func initMySQL() (*sql.DB, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	mysqlCfg := mysql.Config{
		User:                 cfg.Database.User,
		Passwd:               cfg.Database.Password,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", cfg.Database.Host, cfg.Database.Port),
		DBName:               cfg.Database.Name,
		ParseTime:            true,
		AllowNativePasswords: true,
	}
	db, err := sql.Open("mysql", mysqlCfg.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error pinging database: %w", err)
	}
	log.Println("Connected to MySQL database successfully")
	return db, nil
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

func (a *ServerChi) BuildInboundOrderHandler(db *sql.DB) (*handler.InboundOrderHandler, error) {
	// - repository
	inboundOrderRp := repository.NewInboundOrderSQL(db)

	// - service
	inboundOrderSv := service.NewInboundOrderServiceDefault(inboundOrderRp)

	// - handler
	inboundOrderHd := handler.NewInboundOrderHandler(inboundOrderSv)
	return inboundOrderHd, nil
}
