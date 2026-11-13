package task

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

type TaskStatus string
type TaskPriority string
type TaskID string

var (
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrEmptyTaskTitle          = errors.New("title cannot be empty")
	ErrInvalidTaskPriority     = errors.New("invalid task priority")
	ErrDescriptionTooLong      = errors.New("description too long: maximum 500 runes")
	ErrTitleTooLong            = errors.New("title too long: maximum 200 runes")
	ErrEmptyTaskID             = errors.New("task ID cannot be empty")
)
var _ Task = (*task)(nil)

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

const (
	maxTitleLength       = 200
	maxDescriptionLength = 500
)

type task struct {
	id          TaskID
	title       string
	description string
	status      TaskStatus
	priority    TaskPriority
	createdAt   time.Time
	updatedAt   time.Time
}

type Task interface {
	ID() TaskID
	Title() string
	Description() string
	Status() TaskStatus
	Priority() TaskPriority

	CreatedAt() time.Time
	UpdatedAt() time.Time

	UpdateTitle(string) error
	UpdateDescription(string) error
	UpdatePriority(TaskPriority) error

	Start() error
	Complete() error
	Cancel() error
	isTask()
}

func (t task) Status() TaskStatus {
	return t.status
}

func (t task) ID() TaskID {
	return t.id
}

func (t task) CreatedAt() time.Time {
	return t.createdAt
}

func (t task) UpdatedAt() time.Time {
	return t.updatedAt

}

func (t *task) Start() error {
	if t.Status() == StatusTodo {
		t.status = StatusInProgress
		t.updatedAt = time.Now()
		return nil
	}
	return fmt.Errorf("cannot start task from status %q : %w", t.status, ErrInvalidStatusTransition)

}

func (*task) isTask() {

}

func isValidPriority(priority TaskPriority) bool {
	switch priority {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return true
	default:
		return false
	}

}

func (t *task) Complete() error {
	if t.status == StatusInProgress {
		t.status = StatusDone
		t.updatedAt = time.Now()
		return nil
	}
	return fmt.Errorf("cannot complete task from status %q : %w", t.status, ErrInvalidStatusTransition)
}

func (t *task) Cancel() error {
	if t.status == StatusTodo || t.status == StatusInProgress {
		t.status = StatusCancelled
		t.updatedAt = time.Now()
		return nil

	}
	return fmt.Errorf("cannot cancel task from status %q : %w", t.status, ErrInvalidStatusTransition)
}
func (t task) Title() string {
	return t.title
}

func (t task) Description() string {
	return t.description
}
func (t task) Priority() TaskPriority {
	return t.priority
}
func NewTask(id TaskID, title string, description string, priority TaskPriority) (Task, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrEmptyTaskID
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrEmptyTaskTitle
	}
	if utf8.RuneCountInString(title) > maxTitleLength {
		return nil, ErrTitleTooLong
	}
	if priority == "" {
		priority = PriorityMedium
	} else if !isValidPriority(priority) {
		return nil, ErrInvalidTaskPriority
	}
	if utf8.RuneCountInString(description) > maxDescriptionLength {
		return nil, ErrDescriptionTooLong
	}
	now := time.Now()
	task := &task{
		id:          id,
		title:       title,
		description: description,
		status:      StatusTodo,
		priority:    priority,
		createdAt:   now,
		updatedAt:   now}
	return task, nil
}

func (t *task) UpdateDescription(description string) error {
	if utf8.RuneCountInString(description) > maxDescriptionLength {
		return ErrDescriptionTooLong
	}
	if t.description == description {
		return nil
	}
	t.description = description
	t.updatedAt = time.Now()
	return nil
}

func (t *task) UpdateTitle(title string) error {
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrEmptyTaskTitle
	}
	if utf8.RuneCountInString(title) > maxTitleLength {
		return ErrTitleTooLong
	}
	if t.title == title {
		return nil
	}
	t.title = title
	t.updatedAt = time.Now()
	return nil

}

func (t *task) UpdatePriority(priority TaskPriority) error {
	if !isValidPriority(priority) {
		return ErrInvalidTaskPriority
	}
	if t.priority == priority {
		return nil
	}
	t.priority = priority
	t.updatedAt = time.Now()
	return nil
}
