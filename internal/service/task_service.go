package service

import (
	"context"
	"fmt"
	"taskflow/internal/task"
)

type IDGenerator interface {
	NewID() task.TaskID
}

type TaskRepository interface {
	Save(ctx context.Context, t task.Task) error
}

type TaskService struct {
	idGenerator IDGenerator
	repository  TaskRepository
}

func NewTaskService(
	idGenerator IDGenerator,
	repository TaskRepository,
) *TaskService {
	return &TaskService{
		idGenerator: idGenerator,
		repository:  repository,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, title string, description string, priority task.TaskPriority) (task.Task, error) {
	id := s.idGenerator.NewID()
	newtask, err := task.NewTask(id, title, description, priority)
	if err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	if err := s.repository.Save(ctx, newtask); err != nil {
		return nil, fmt.Errorf("save task : %w", err)
	}
	return newtask, nil

}
