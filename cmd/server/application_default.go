package server

import (
	"app/internal/handler"
	"app/internal/loader"
	"app/internal/repository"
	"app/internal/service"
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

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints

	rt.Route("/products", func(r chi.Router) {
		r.Get("/", prodHandler.GetAll())
		r.Get("/{id}", prodHandler.GetById())
		r.Post("/", prodHandler.Create())
		r.Delete("/{id}", prodHandler.Delete())
		r.Patch("/{id}", prodHandler.Patch())
	})

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
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
