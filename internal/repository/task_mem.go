package repository

import (
	"sync"
	"task-manager/internal/models"
)

type TaskRepository struct {
	mu sync.Mutex
	tasks []models.Task
}


func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: []models.Task{},
	}
}


func (r *TaskRepository) GetAll() []models.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.tasks
}


func (r *TaskRepository) Add(task models.Task) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, task)
}


func (r *TaskRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return true
		}
	}
	return false
}

func (r *TaskRepository) Update(id int, updateTask models.Task) (models.Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for idx, task := range r.tasks {
		if task.ID == id {
			if updateTask.Title != "" {
				r.tasks[idx].Title = updateTask.Title
			}
			if updateTask.Description != "" {
				r.tasks[idx].Description = updateTask.Description
			}
			if updateTask.Status != "" {
				r.tasks[idx].Status = updateTask.Status
			}
			return r.tasks[idx], true
		}
	}
	return models.Task{}, false
}
