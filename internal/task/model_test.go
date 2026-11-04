package task

import (
	"errors"
	"testing"
)

func TestTaskStart(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TaskStatus
		wantStatus    TaskStatus
		wantErr       bool
	}{
		{
			name:          "start todo task",
			initialStatus: StatusTodo,
			wantStatus:    StatusInProgress,
		},
		{
			name:          "start completed task",
			initialStatus: StatusDone,
			wantStatus:    StatusDone,
			wantErr:       true,
		}, {
			name:          "start in-progress task",
			initialStatus: StatusInProgress,
			wantStatus:    StatusInProgress,
			wantErr:       true,
		}, {
			name:          "start cancelled task",
			initialStatus: StatusCancelled,
			wantStatus:    StatusCancelled,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				status: tt.initialStatus,
			}
			err := task.Start()
			if tt.wantErr {

				if !errors.Is(err, ErrInvalidStatusTransition) {
					t.Fatalf("expected ErrInvalidStatusTransition, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := task.Status(); got != tt.wantStatus {
				t.Errorf("expected: %v, got:%v", tt.wantStatus, got)
			}
		})
	}

}

func TestTaskComplete(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TaskStatus
		wantStatus    TaskStatus
		wantErr       bool
	}{
		{
			name:          "complete todo task",
			initialStatus: StatusTodo,
			wantStatus:    StatusTodo,
			wantErr:       true,
		},
		{
			name:          "complete in_progress task",
			initialStatus: StatusInProgress,
			wantStatus:    StatusDone,
		}, {
			name:          "complete completed task",
			initialStatus: StatusDone,
			wantStatus:    StatusDone,
			wantErr:       true,
		}, {
			name:          "complete cancelled task",
			initialStatus: StatusCancelled,
			wantStatus:    StatusCancelled,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				status: tt.initialStatus,
			}

			err := task.Complete()
			if tt.wantErr {
				if !errors.Is(err, ErrInvalidStatusTransition) {
					t.Fatalf("expected ErrInvalidStatusTransition, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if got := task.Status(); got != tt.wantStatus {
				t.Errorf("expected status: %v, got: %v", tt.wantStatus, got)
			}
		})
	}
}

func TestTaskCancel(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus TaskStatus
		wantStatus    TaskStatus
		wantErr       bool
	}{
		{
			name:          "cancel todo task",
			initialStatus: StatusTodo,
			wantStatus:    StatusCancelled,
		},
		{
			name:          "cancel in_progress task",
			initialStatus: StatusInProgress,
			wantStatus:    StatusCancelled,
		}, {
			name:          "cancel completed task",
			initialStatus: StatusDone,
			wantStatus:    StatusDone,
			wantErr:       true,
		}, {
			name:          "cancel cancelled task",
			initialStatus: StatusCancelled,
			wantStatus:    StatusCancelled,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				status: tt.initialStatus,
			}

			err := task.Cancel()

			if tt.wantErr {
				if !errors.Is(err, ErrInvalidStatusTransition) {
					t.Fatalf("expected ErrInvalidStatusTransition, got: %v", err)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got := task.Status(); got != tt.wantStatus {
				t.Errorf("expected status: %v, got : %v", tt.wantStatus, got)
			}
		})
	}

}

func TestNewTask(t *testing.T) {
	tests := []struct {
		name         string
		title        string
		priority     TaskPriority
		wantTitle    string
		wantPriority TaskPriority
		wantErr      error
	}{
		{
			name:         "Valid low",
			title:        "Learn Go",
			priority:     PriorityLow,
			wantTitle:    "Learn Go",
			wantPriority: PriorityLow,
		},
		{
			name:         "Valid Medium",
			title:        "Learn Go",
			priority:     PriorityMedium,
			wantTitle:    "Learn Go",
			wantPriority: PriorityMedium,
		}, {
			name:         "Valid High",
			title:        "Learn Go",
			priority:     PriorityHigh,
			wantTitle:    "Learn Go",
			wantPriority: PriorityHigh,
		}, {
			name:         "Default Priority",
			title:        "Learn Go",
			priority:     "",
			wantTitle:    "Learn Go",
			wantPriority: PriorityMedium,
		},
		{
			name:     "Empty Title",
			title:    "",
			priority: PriorityLow,
			wantErr:  ErrEmptyTaskTitle,
		},
		{
			name:     "WhiteSpace Title",
			title:    "      ",
			priority: PriorityLow,
			wantErr:  ErrEmptyTaskTitle,
		},
		{
			name:     "Invalid Priority",
			title:    "Learn Go",
			priority: "banana",
			wantErr:  ErrInvalidTaskPriority,
		},
		{
			name:         "trim title",
			title:        "  Learn Go  ",
			priority:     PriorityLow,
			wantTitle:    "Learn Go",
			wantPriority: PriorityLow,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.title, "", tt.priority)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected err: %v, got: %v", tt.wantErr, err)
				}
				if task != nil {
					t.Fatalf("expected nil task, got:%v", task)
				}
				return

			} else if err != nil {
				t.Fatalf("unexpecter error: %v", err)
			}
			if task == nil {
				t.Fatalf("expected task, got nil")
			}

			if task.title != tt.wantTitle {
				t.Errorf("expected title: %v, got: %v", tt.wantTitle, task.title)
			}
			if task.priority != tt.wantPriority {
				t.Errorf("expected priority: %v, got: %v", tt.wantPriority, task.priority)
			}
			if task.Status() != StatusTodo {
				t.Errorf("expected Status ToDo, got:%v", task.Status())
			}
			if task.createdAt.IsZero() {
				t.Errorf("expected time,got nil")
			}
			if task.updatedAt.IsZero() {
				t.Errorf("expected time, got nil")
			}
			if task.createdAt != task.updatedAt {
				t.Errorf("expected one time, got: %v , %v", task.createdAt, task.updatedAt)
			}
			if task.id != 0 {
				t.Errorf("expected id = 0, got : id = %v", task.id)
			}

		})
	}

}
