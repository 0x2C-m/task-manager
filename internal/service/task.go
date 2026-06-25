package service

import (
	"errors"
	"strings"

	"task-manager/internal/models"
	"task-manager/internal/repository"
)

var (
	ErrEmptyTitle = errors.New("task name can't be empty")
	ErrInvalidStatus = errors.New("недопустимый статус задачи")
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task models.Task) (models.Task, error) {
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		return models.Task{}, ErrEmptyTitle
	}

	if task.Status == "" {
		task.Status = "todo"
	}

	return s.repo.Add(task), nil
}

func (s *TaskService) GetAll() []models.Task {
	return s.repo.GetAll()
}

func (s *TaskService) DeleteTask(id int) bool {
	return s.repo.Delete(id)
}

func (s *TaskService) UpdateTask(id int, input models.Task) (models.Task, error) {
	if input.Status != "" {
		if input.Status != "todo" && input.Status != "in_progress" {
			return models.Task{}, ErrInvalidStatus
		}
	}

	task, updated := s.repo.Update(id, input)
	if !updated {
		return models.Task{}, errors.New("task not found")
	}

	return task, nil
}
