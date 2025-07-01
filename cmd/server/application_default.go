package server

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/models"
	"net/http"

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
	// create product handler with dependences
	prodHandler, err := a.BuildProductHandler()
	if err != nil {
		return err
	}

	// dependencies
	// - loader
	// - repository
	employeeHandler, err := a.BuildemployeeHandler()
	if err != nil {
		return err
	}

	// handler warehouse
	hdWarehouse, err := a.BuildWarehouseHandler()
	if err != nil {
		return err
	}

	// - handler
	sectionHd, err := a.BuildSectionHandler()
	if err != nil {
		return err
	}

	// create seller handler with dependences
	hdSeller, err := a.BuildSellerHandler()
	if err != nil {
		return err
	}

	buyerHd, err := a.BuildBuyerHandler()
	if err != nil {
		return err
	}

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	// Grupo de endpoints para sellers
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

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", prodHandler.GetAll())
		r.Get("/{id}", prodHandler.GetById())
		r.Post("/", prodHandler.Create())
		r.Delete("/{id}", prodHandler.Delete())
		r.Patch("/{id}", prodHandler.Patch())
	})

	rt.Route("/employees", func(r chi.Router) {
		r.Get("/", employeeHandler.GetAllEmployees)
		r.Get("/{id}", employeeHandler.GetEmployee)
		r.Post("/", employeeHandler.CreateEmployee)
		r.Patch("/{id}", employeeHandler.UpdateEmployee)
		r.Delete("/{id}", employeeHandler.DeleteEmployee)
	})
	rt.Route("/api/v1", func(rt chi.Router) {
		rt.Route("/warehouses", func(rt chi.Router) {
			// - GET /warehouses
			rt.Get("/", hdWarehouse.GetAll())
			rt.Get("/{id}", hdWarehouse.GetOne())
			rt.Post("/", hdWarehouse.Add())
			rt.Patch("/{id}", hdWarehouse.Update())
			rt.Delete("/{id}", hdWarehouse.Delete())
		})
	})

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func (a *ServerChi) BuildWarehouseHandler() (*handler.WarehouseDefault, error) {

	ldWarehouse := loader.NewWarehouseJSONFile("docs/db/warehouse_500.json")
	dbWarehouse, err := ldWarehouse.Load()
	if err != nil {
		return nil, err
	}

	// repository
	rpWarehouse := repository.NewWarehouseMap(dbWarehouse)

	// - service
	svWarehouse := service.NewWarehouseDefault(rpWarehouse)

	// handler
	hdWarehouse := handler.NewWarehouseDefault(svWarehouse)
	return hdWarehouse, nil
}

func (a *ServerChi) BuildSectionHandler() (*handler.SectionDefault, error) {
	sectionLd := loader.NewLoaderGeneric[models.Section]()
	sectionDb, err := sectionLd.LoadFromJSON("docs/db/sections_5.json")
	if err != nil {
		return nil, err
	}
	// - repository
	sectionRp := repository.NewSectionMapRepository(sectionDb)
	// - service
	sectionSv := service.NewSectionDefault(sectionRp)
	// - handler
	sectionHd := handler.NewSectionDefault(sectionSv)
	return sectionHd, nil
}

func (a *ServerChi) BuildProductHandler() (*handler.ProductDefault, error) {

	// - loader
	productLoader := loader.NewProductJSONFile("docs/db/products.json")
	productDb, err := productLoader.Load()

	if err != nil {
		return nil, err
	}

	// - repository
	productRp := repository.NewProductMap(productDb)

	// - service
	productSv := service.NewProductDefault(productRp)

	// - handler
	prodHandler := handler.NewProductDefault(productSv)
	return prodHandler, nil
}

func (*ServerChi) BuildemployeeHandler() (*handler.EmployeeHandler, error) {
	loaderEmployee := loader.NewLoaderGeneric[models.EmployeeDocument]()
	employees, err := loaderEmployee.LoadFromJSON("./docs/db/employees.json")
	if err != nil {
		return nil, err
	}
	employeeRepository := repository.NewEmployeeMapRepository(employees)
	// - service
	employeeService := service.NewEmployeeServiceDefault(employeeRepository)
	// - handler
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	return employeeHandler, nil
}

func (a *ServerChi) BuildSellerHandler() (*handler.SellerDefault, error) {

	// - loader
	ldSeller := loader.NewLoaderGeneric[models.SellerDoc]()

	dbSeller, err := ldSeller.LoadFromJSON("./docs/db/sellers.json")
	if err != nil {
		return nil, err
	}

	// - repository
	rpSeller := repository.NewSellerMap(dbSeller)

	// - service
	svSeller := service.NewSellerDefault(rpSeller)

	// - handler
	hdSeller := handler.NewSellerDefault(svSeller)
	return hdSeller, nil
}

func (a *ServerChi) BuildBuyerHandler() (*handler.BuyerDefault, error) {

	// - loader
	ldBuyer := loader.NewBuyerJSONFile("./docs/db/buyers.json")

	dbBuyer, err := ldBuyer.Load()
	if err != nil {
		return nil, err
	}

	// - repository
	rpBuyer := repository.NewBuyerMap(dbBuyer)

	// - service
	svBuyer := service.NewBuyerDefault(rpBuyer)

	// - handler
	hdBuyer := handler.NewBuyerDefault(svBuyer)
	return hdBuyer, nil
}
