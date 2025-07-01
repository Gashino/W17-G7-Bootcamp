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
	// dependencies

	// create seller handler with dependences
	hdSeller, err := a.BuildSellerHandler()
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

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}

func (a *ServerChi) BuildSellerHandler() (*handler.SellerDefault, error) {

	// - loader
	ldSeller := loader.NewLoaderGeneric[models.SellerDoc]()

	dbSeller, err := ldSeller.LoadFromJSON(a.loaderFilePath)
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
