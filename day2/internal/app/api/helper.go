package api

import (
	"github.com/narekps/go2/day2/storage"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	prefix string = "/api/v1"
)

// Configure logger level
func (api *API) configureLogger() error {
	logLevel, err := logrus.ParseLevel(api.config.LogLevel)
	if err != nil {
		return err
	}

	api.logger.SetLevel(logLevel)

	return nil

}

// Configure router
func (api *API) configureRouter() {
	api.router.StrictSlash(false)

	api.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello! This is rest api."))
	})

	api.router.HandleFunc(prefix+"/task", api.CreateTaskHandler).Methods("POST")
	api.router.HandleFunc(prefix+"/task/{id:[0-9]+}", api.GetTaskHandler).Methods("GET")
	api.router.HandleFunc(prefix+"/task/{id:[0-9]+}", api.DeleteTaskHandler).Methods("DELETE")
	api.router.HandleFunc(prefix+"/task/{id:[0-9]+}", api.UpdateTaskHandler).Methods("PUT")
	api.router.HandleFunc(prefix+"/task", api.GetAllTasksHandler).Methods("GET")
	api.router.HandleFunc(prefix+"/task", api.DeleteAllTasksHandler).Methods("DELETE")
	api.router.HandleFunc(prefix+"/task/{tag}", api.GetAllTasksByTagHandler).Methods("GET")
	api.router.HandleFunc(prefix+"/task/{year:[0-9]+}/{month:[0-9]+}/{day:[0-9]+}", api.GetAllTasksByDateHandler).Methods("GET")
}

func (api *API) configureStorage() error {
	s := storage.New(api.config.Storage)

	if err := s.Open(); err != nil {
		return err
	}

	api.storage = s

	return nil
}
