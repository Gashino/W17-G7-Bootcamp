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
	// BuyerLoaderFilePath is the path to the file that contains the buyers
	BuyerLoaderFilePath string
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
		if cfg.BuyerLoaderFilePath != "" {
			defaultConfig.BuyerLoaderFilePath = cfg.BuyerLoaderFilePath
		}
	}

	return &ServerChi{
		serverAddress:       defaultConfig.ServerAddress,
		buyerLoaderFilePath: defaultConfig.BuyerLoaderFilePath,
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
	// dependencies for buyers
	buyerLd := loader.NewBuyerJSONFile(a.buyerLoaderFilePath)
	buyerDb, err := buyerLd.Load()
	if err != nil {
		return
	}
	// - repository
	buyerRp := repository.NewBuyerMap(buyerDb)
	// - service
	buyerSv := service.NewBuyerDefault(buyerRp)
	// - handler
	buyerHd := handler.NewBuyerDefault(buyerSv)

	// router
	rt := chi.NewRouter()
	// - middlewares
	rt.Use(middleware.Logger)
	rt.Use(middleware.Recoverer)
	// - endpoints
	rt.Route("/buyers", func(r chi.Router) {
		r.Get("/", buyerHd.GetAll())
		r.Get("/{id}", buyerHd.GetByID())
		r.Post("/", buyerHd.Create())
		r.Patch("/{id}", buyerHd.Update())
		r.Delete("/{id}", buyerHd.Delete())
	})

	// run server
	err = http.ListenAndServe(a.serverAddress, rt)
	return
}
