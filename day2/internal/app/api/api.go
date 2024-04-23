package api

import (
	"github.com/gorilla/mux"
	"github.com/narekps/go2/day2/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

type API struct {
	config  *Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *Config) *API {
	return &API{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (api *API) Start() error {
	if err := api.configureLogger(); err != nil {
		return err
	}

	api.configureRouter()

	api.logger.Info("Starting API server at address ", api.config.BindAddr)

	if err := api.configureStorage(); err != nil {
		return err
	}

	return http.ListenAndServe(api.config.BindAddr, api.router)
}
