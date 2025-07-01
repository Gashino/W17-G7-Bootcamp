package server

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
	"app/pkg/models"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// ConfigServerChi is a struct that represents the configuration for ServerChi
type ConfigServerChi struct {
	// ServerAddress is the address where the server will be listening
	ServerAddress string
	// LoaderFilePath is the path to the file that contains the vehicles
	LoaderFilePath string
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
		if cfg.LoaderFilePath != "" {
			defaultConfig.LoaderFilePath = cfg.LoaderFilePath
		}
	}

	return &ServerChi{
		serverAddress:  defaultConfig.ServerAddress,
		loaderFilePath: defaultConfig.LoaderFilePath,
	}
}

// ServerChi is a struct that implements the Application interface
type ServerChi struct {
	// serverAddress is the address where the server will be listening
	serverAddress string
	// loaderFilePath is the path to the file that contains the vehicles
	loaderFilePath string
}

// Run is a method that runs the server
func (a *ServerChi) Run() (err error) {

	sectionDb, productDb, employeeDb, productTypeDb, warehouseDb, err := a.createMaps()
	if err != nil {
		return err
	}

	// create product handler with dependences
	prodHandler, err := a.BuildProductHandler(&productDb, &productTypeDb)
	if err != nil {
		return err
	}

	// dependencies
	// - loader

	// - repository
	//employeeHandler, err := a.BuildemployeeHandler()
	//if err != nil {
	//	return err
	//}

	// handler warehouse

	hdWarehouse, err := a.BuildWarehouseHandler(&sectionDb, &employeeDb, &warehouseDb)
	if err != nil {
		return err
	}

	// - handler
	sectionHd, err := a.BuildSectionHandler(&sectionDb, &productTypeDb)
	if err != nil {
		return err
	}

	/*
		// create seller handler with dependences
		hdSeller, err := a.BuildSellerHandler()
		if err != nil {
			return err
		}
	*/

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)

	// - endpoints
	// Grupo de endpoints para sellers

	/*
		rt.Route("/sellers", func(rt chi.Router) {
			// - GET /sellers
			rt.Get("/", hdSeller.GetAll())
			rt.Get("/{id}", hdSeller.GetById())
			rt.Post("/", hdSeller.Create())
			rt.Patch("/{id}", hdSeller.Update())
			rt.Delete("/{id}", hdSeller.Delete())
		})
	*/

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

	/*rt.Route("/employees", func(r chi.Router) {
		r.Get("/", employeeHandler.GetAllEmployees)
		r.Get("/{id}", employeeHandler.GetEmployee)
		r.Post("/", employeeHandler.CreateEmployee)
		r.Patch("/{id}", employeeHandler.UpdateEmployee)
		r.Delete("/{id}", employeeHandler.DeleteEmployee)
	})
	})*/

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

func (a *ServerChi) createMaps() (map[int]models.Section, map[int]models.Product, map[int]models.Employee, map[int]models.ProductType, map[int]models.Warehouse, error) {
	sectionLd := loader.NewLoaderGeneric[models.Section]()
	sectionDb, err := sectionLd.LoadFromJSON("docs/db/sections.json")
	if err != nil {
		fmt.Println("err1")
		return nil, nil, nil, nil, nil, err
	}
	productLd := loader.NewLoaderGeneric[models.Product]()
	productDb, err := productLd.LoadFromJSON("docs/db/products.json")
	if err != nil {
		fmt.Println("err2")
		return nil, nil, nil, nil, nil, err
	}
	employeeLd := loader.NewLoaderGeneric[models.Employee]()
	employeeDb, err := employeeLd.LoadFromJSON("docs/db/employees.json")
	if err != nil {
		fmt.Println("err3")
		return nil, nil, nil, nil, nil, err
	}

	productTypeLd := loader.NewLoaderGeneric[models.ProductType]()
	productTypeDb, err := productTypeLd.LoadFromJSON("docs/db/productTypes.json")
	if err != nil {
		fmt.Println("err4")
		return nil, nil, nil, nil, nil, err
	}

	warehouseLd := loader.NewLoaderGeneric[models.Warehouse]()
	warehouseDb, err := warehouseLd.LoadFromJSON("docs/db/warehouse_500.json")
	if err != nil {
		fmt.Println("err5")
		return nil, nil, nil, nil, nil, err
	}

	return sectionDb, productDb, employeeDb, productTypeDb, warehouseDb, nil
}

func (a *ServerChi) BuildSectionHandler(sectionDb *map[int]models.Section, productTypeDb *map[int]models.ProductType) (*handler.SectionDefault, error) {

	// - repository
	sectionRp := repository.NewSectionMapRepository(sectionDb, productTypeDb)
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

func (a *ServerChi) BuildProductHandler(productDb *map[int]models.Product, productTypeDb *map[int]models.ProductType) (prodHandler *handler.ProductDefault, err error) {

	// - repository
	productRp := repository.NewProductMap(productDb, productTypeDb)

	// - service
	productSv := service.NewProductDefault(productRp)

	// - handler
	prodHandler = handler.NewProductDefault(productSv)
	return prodHandler, nil
}

/*
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
*/
