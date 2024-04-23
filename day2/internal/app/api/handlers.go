package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/narekps/go2/day2/internal/app/models"
	"net/http"
	"strconv"
)

type SuccessResponse struct {
	Tasks []*models.Task `json:"tasks"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (api *API) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - CreateTaskHandler called.", r.Method, r.URL.Path)
	var task *models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		api.sendErrorResponse(w, 400, err.Error())
		return
	}

	t, err := api.storage.Task().Create(task)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}

	api.sendSuccessResponse(w, 201, t)
}

func (api *API) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - GetTaskHandler called.", r.Method, r.URL.Path)
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.sendErrorResponse(w, 400, err.Error())
		return
	}

	task, ok, err := api.storage.Task().FindById(id)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}

	if !ok {
		api.sendErrorResponse(w, 404, fmt.Sprintf("Task with id %d not found", id))
		return
	}

	api.sendSuccessResponse(w, 200, task)
}

func (api *API) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - DeleteTaskHandler called.", r.Method, r.URL.Path)
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.sendErrorResponse(w, 400, err.Error())
		return
	}

	_, ok, err := api.storage.Task().FindById(id)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}
	if !ok {
		api.sendErrorResponse(w, 404, fmt.Sprintf("Task with id %d not found", id))
		return
	}

	task, err := api.storage.Task().DeleteTask(id)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}

	api.sendSuccessResponse(w, 200, task)
}

func (api *API) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - UpdateTaskHandler called.", r.Method, r.URL.Path)
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		api.sendErrorResponse(w, 400, err.Error())
		return
	}

	_, ok, err := api.storage.Task().FindById(id)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}
	if !ok {
		api.sendErrorResponse(w, 404, fmt.Sprintf("Task with id %d not found", id))
		return
	}

	var task *models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		api.sendErrorResponse(w, 400, err.Error())
		return
	}

	task.ID = int64(id)

	task, err = api.storage.Task().UpdateTask(task)
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}

	api.sendSuccessResponse(w, 200, task)
}

func (api *API) GetAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - GetAllTasks called.", r.Method, r.URL.Path)

	tasks, err := api.storage.Task().SelectAllTasks()
	if err != nil {
		api.sendErrorResponse(w, 500, fmt.Sprintf("GetAllTasks error: %v", err))
		return
	}

	api.sendSuccessResponse(w, 200, SuccessResponse{Tasks: tasks})
}

func (api *API) DeleteAllTasksHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - DeleteAllTasksHandler called.", r.Method, r.URL.Path)

	err := api.storage.Task().DeleteAllTasks()
	if err != nil {
		api.sendErrorResponse(w, 500, err.Error())
		return
	}

	tasks := make([]*models.Task, 0, 0)
	api.sendSuccessResponse(w, 200, SuccessResponse{Tasks: tasks})
}

func (api *API) GetAllTasksByTagHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - GetAllTasksByTagHandler called.", r.Method, r.URL.Path)
	tag := mux.Vars(r)["tag"]
	tasks, err := api.storage.Task().FindTasksByTag(tag)
	if err != nil {
		api.sendErrorResponse(w, 500, fmt.Sprintf("GetAllTasks error: %v", err))
		return
	}

	api.sendSuccessResponse(w, 200, SuccessResponse{Tasks: tasks})
}

func (api *API) GetAllTasksByDateHandler(w http.ResponseWriter, r *http.Request) {
	api.logger.Infof("%s %s - GetAllTasksByDateHandler called.", r.Method, r.URL.Path)
	year := mux.Vars(r)["year"]
	month := mux.Vars(r)["month"]
	day := mux.Vars(r)["day"]

	date := fmt.Sprintf("%v-%v-%v", year, month, day)

	tasks, err := api.storage.Task().SelectTasksByDate(date)
	if err != nil {
		api.sendErrorResponse(w, 500, fmt.Sprintf("GetAllTasksByDateHandler error: %v", err))
		return
	}

	api.sendSuccessResponse(w, 200, SuccessResponse{Tasks: tasks})
}

func (api *API) sendErrorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	api.logger.Error(msg)
	response := ErrorResponse{
		Code:    code,
		Message: msg,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func (api *API) sendSuccessResponse(w http.ResponseWriter, code int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
