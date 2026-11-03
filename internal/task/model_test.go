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
				ID:     1,
				Title:  "Test",
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
				ID:     1,
				Title:  "test",
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
				ID:     1,
				Title:  "test",
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
