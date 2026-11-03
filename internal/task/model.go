package task

import (
	"errors"
	"fmt"
	"time"
)

type TaskStatus string
type TaskPriority string

var ErrInvalidStatusTransition = errors.New("invalid status transition")

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
	StatusCancelled  TaskStatus = "cancelled"
)
const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	status      TaskStatus
	Priority    TaskPriority
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateTaskRequest struct {
	Title       string       `json:"title"`
	Description string       `json:"description,omitempty"`
	Priority    TaskPriority `json:"priority"`
}

func (t Task) Status() TaskStatus {
	return t.status
}

func (t *Task) Start() error {
	if t.Status() == StatusTodo {
		t.status = StatusInProgress
		return nil
	}
	return fmt.Errorf("cannot start task from status %q : %w", t.status, ErrInvalidStatusTransition)

}

func (t *Task) Complete() error {
	if t.status == StatusInProgress {
		t.status = StatusDone
		return nil
	}
	return fmt.Errorf("cannot complete task from status %q : %w", t.status, ErrInvalidStatusTransition)
}

func (t *Task) Cancel() error {
	if t.status == StatusTodo || t.status == StatusInProgress {
		t.status = StatusCancelled
		return nil
	}
	return fmt.Errorf("cannot cancel task from status %q : %w", t.status, ErrInvalidStatusTransition)
}
