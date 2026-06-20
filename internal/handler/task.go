package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"

	"task-manager/internal/repository"
	"task-manager/internal/models"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	tasks := h.repo.GetAll()

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
	if input.Status == "" {
		input.Status = "todo"
	}

	h.repo.Add(input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(input)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "incorrect task id", http.StatusBadRequest)
		return
	}

	deleted := h.repo.Delete(id)
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

	updatedTask, updated := h.repo.Update(id, input)
	if !updated {
		http.Error(w, "task not found", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}
