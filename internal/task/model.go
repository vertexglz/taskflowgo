package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type TaskStatus string
type TaskPriority string

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrEmptyTaskTitle          = errors.New("title cannot be empty")
	ErrInvalidTaskPriority     = errors.New("invalid task priority")
)

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
	id          int64
	title       string
	description string
	status      TaskStatus
	priority    TaskPriority
	createdAt   time.Time
	updatedAt   time.Time
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

func NewTask(title, description string, priority TaskPriority) (*Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTaskTitle
	}
	if priority == "" {
		priority = PriorityMedium
	} else if priority != PriorityLow && priority != PriorityMedium && priority != PriorityHigh {
		return nil, ErrInvalidTaskPriority
	}
	now := time.Now()
	task := &Task{
		title:       title,
		description: description,
		status:      StatusTodo,
		priority:    priority,
		createdAt:   now,
		updatedAt:   now}
	return task, nil
}
