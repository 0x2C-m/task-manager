package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
	"errors"

	"task-manager/internal/models"
	"task-manager/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	tasks := h.service.GetAll()

	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, "error encoding json", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.Task
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "incorrect json format", http.StatusBadRequest)
	}
	defer r.Body.Close()

	input.ID = int(time.Now().UnixNano() % 100000)
	input.CreatedAt = time.Now()

	finalTask, err := h.service.CreateTask(input)
	if err != nil {
		if errors.Is(err, service.ErrEmptyTitle) || errors.Is(err, service.ErrInvalidStatus) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		http.Error(w, "внутрення ошибка", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(finalTask)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "incorrect task id", http.StatusBadRequest)
		return
	}

	deleted := h.service.DeleteTask(id)
	if !deleted {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "incorrect task id", http.StatusBadRequest)
	}

	var input models.Task
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "incorrect json format", http.StatusBadRequest)
	}
	defer r.Body.Close()

	updatedTask, err := h.service.UpdateTask(id, input)
	if err != nil {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}
