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
	// - handler
	sectionHd, err := a.BuildSectionHandler()
	if err != nil {
		return err
	}

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints
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
	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func (a *ServerChi) BuildSectionHandler() (*handler.SectionDefault, error) {
	sectionLd := loader.NewLoaderGeneric[models.Section]()
	sectionDb, err := sectionLd.LoadFromJSON(a.loaderFilePath)
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
