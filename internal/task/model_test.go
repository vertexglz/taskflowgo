package task

import (
	"errors"
	"strings"
	"testing"
	"time"
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
			task := task{
				status: tt.initialStatus,
			}
			err := task.Start()
			if tt.wantErr {

				if !errors.Is(err, ErrInvalidStatusTransition) {
					t.Fatalf("expected ErrInvallidStatusTransition, got: %v", err)
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
			task := task{
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
			task := task{
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
		name            string
		description     string
		title           string
		priority        TaskPriority
		wantTitle       string
		wantDescription string
		wantPriority    TaskPriority
		wantErr         error
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
		{
			name:         "title too long",
			title:        strings.Repeat("n", 201),
			priority:     PriorityLow,
			wantPriority: PriorityLow,
			wantErr:      ErrTitleTooLong,
		}, {
			name:         "description too long",
			title:        "test",
			description:  strings.Repeat("n", 501),
			priority:     PriorityLow,
			wantTitle:    "test",
			wantPriority: PriorityLow,
			wantErr:      ErrDescriptionTooLong,
		},
		{
			name:         "title good",
			title:        strings.Repeat("n", 200),
			priority:     PriorityLow,
			wantPriority: PriorityLow,
			wantTitle:    strings.Repeat("n", 200),
		}, {
			name:            "description good",
			title:           "test",
			description:     strings.Repeat("n", 500),
			priority:        PriorityLow,
			wantTitle:       "test",
			wantPriority:    PriorityLow,
			wantDescription: strings.Repeat("n", 500),
		}, {
			name:         "russian title good",
			title:        strings.Repeat("я", 200),
			priority:     PriorityLow,
			wantPriority: PriorityLow,
			wantTitle:    strings.Repeat("я", 200),
		}, {
			name:         "russian title too long",
			title:        strings.Repeat("я", 201),
			priority:     PriorityLow,
			wantPriority: PriorityLow,
			wantErr:      ErrTitleTooLong,
		}, {
			name:            "russian description good",
			title:           "test",
			description:     strings.Repeat("я", 500),
			priority:        PriorityLow,
			wantTitle:       "test",
			wantPriority:    PriorityLow,
			wantDescription: strings.Repeat("я", 500),
		}, {
			name:         "russian description too long",
			title:        "test",
			description:  strings.Repeat("я", 501),
			priority:     PriorityLow,
			wantTitle:    "test",
			wantPriority: PriorityLow,
			wantErr:      ErrDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.title, tt.description, tt.priority)

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
			if task.description != tt.wantDescription {
				t.Errorf("expected description: %v, got:%v", tt.wantDescription, task.description)
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
			if !task.createdAt.Equal(task.updatedAt) {
				t.Errorf("expected one time, got: %v , %v", task.createdAt, task.updatedAt)
			}
			if task.id != 0 {
				t.Errorf("expected id = 0, got : id = %v", task.id)
			}

		})
	}

}

func TestTaskDescription(t *testing.T) {
	tests := []struct {
		name            string
		description     string
		wantDescription string
	}{
		{
			name:            "test description",
			description:     "test",
			wantDescription: "test",
		}, {
			name:            "empty description test",
			description:     "",
			wantDescription: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := task{
				description: tt.description,
			}

			if task.Description() != tt.wantDescription {
				t.Errorf("expected description:%v, got:%v", tt.wantDescription, task.Description())
			}
		})
	}
}

func TestTaskPriority(t *testing.T) {
	tests := []struct {
		name         string
		priority     TaskPriority
		wantPriority TaskPriority
	}{
		{
			name:         "test low priority",
			priority:     PriorityLow,
			wantPriority: PriorityLow,
		}, {
			name:         "test medium priority",
			priority:     PriorityMedium,
			wantPriority: PriorityMedium,
		}, {
			name:         "test high priority",
			priority:     PriorityHigh,
			wantPriority: PriorityHigh,
		},
		{
			name:         "test empty priority",
			priority:     "",
			wantPriority: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := task{
				priority: tt.priority,
			}

			if task.Priority() != tt.wantPriority {
				t.Fatalf("expected priority:%v, got: %v", tt.wantPriority, task.Priority())
			}
		})
	}
}
func TestDescriptionUpdate(t *testing.T) {
	tests := []struct {
		name             string
		description      string
		inputDescription string
		wantDescription  string
		wantErr          error
	}{
		{
			name:             "description normal",
			description:      "test",
			inputDescription: "newtest",
			wantDescription:  "newtest",
		}, {
			name:             "description empty",
			description:      "test",
			inputDescription: "",
			wantDescription:  "",
		}, {
			name:             "description to long",
			description:      "test",
			inputDescription: strings.Repeat("n", 501),
			wantDescription:  "test",
			wantErr:          ErrDescriptionTooLong,
		}, {
			name:             "description good",
			description:      "test",
			inputDescription: strings.Repeat("n", 500),
			wantDescription:  strings.Repeat("n", 500),
		}, {
			name:             "russian description good",
			description:      "",
			inputDescription: strings.Repeat("я", 300),
			wantDescription:  strings.Repeat("я", 300),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			task := task{
				description: tt.description,
				updatedAt:   oldTime,
			}

			err := task.UpdateDescription(tt.inputDescription)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
				}
			}
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if task.description != tt.wantDescription {
				t.Errorf("expected description: %v, got: %v", tt.wantDescription, task.Description())
			}
			if tt.wantErr == nil {
				if !task.updatedAt.After(oldTime) {
					t.Errorf("expected new update time, got: %v", task.updatedAt)
				}
			} else {
				if !task.updatedAt.Equal(oldTime) {
					t.Errorf("expected created time, got:%v", task.updatedAt)
				}
			}

		})
	}
}

func TestTitleUpdate(t *testing.T) {
	tests := []struct {
		name             string
		title            string
		inputTitle       string
		wantTitle        string
		wantErr          error
		wantUpdateChange bool
	}{
		{
			name:             "title normal",
			title:            "test",
			inputTitle:       "newtest",
			wantTitle:        "newtest",
			wantUpdateChange: true,
		}, {
			name:             "title empty",
			title:            "test",
			inputTitle:       "",
			wantTitle:        "test",
			wantErr:          ErrEmptyTaskTitle,
			wantUpdateChange: false,
		}, {
			name:             "title to long",
			title:            "test",
			inputTitle:       strings.Repeat("n", 201),
			wantTitle:        "test",
			wantErr:          ErrTitleTooLong,
			wantUpdateChange: false,
		}, {
			name:             "title good",
			title:            "test",
			inputTitle:       strings.Repeat("n", 200),
			wantTitle:        strings.Repeat("n", 200),
			wantUpdateChange: true,
		}, {
			name:             "russian title good",
			title:            "тест",
			inputTitle:       strings.Repeat("я", 200),
			wantTitle:        strings.Repeat("я", 200),
			wantUpdateChange: true,
		},
		{
			name:             "check the update time",
			title:            "test",
			inputTitle:       "test",
			wantTitle:        "test",
			wantUpdateChange: false,
		},
		{
			name:             "check the trimspace",
			title:            "test",
			inputTitle:       "   123test   ",
			wantTitle:        "123test",
			wantUpdateChange: true,
		},
		{
			name:             "check the trimspace empty",
			title:            "test",
			inputTitle:       "        ",
			wantTitle:        "test",
			wantErr:          ErrEmptyTaskTitle,
			wantUpdateChange: false,
		},
		{
			name:             "russian title too long",
			title:            "тест",
			inputTitle:       strings.Repeat("я", 201),
			wantTitle:        "тест",
			wantErr:          ErrTitleTooLong,
			wantUpdateChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			task := task{
				title:     tt.title,
				updatedAt: oldTime,
			}

			err := task.UpdateTitle(tt.inputTitle)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
				}
			}
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if task.title != tt.wantTitle {
				t.Errorf("expected title: %v, got: %v", tt.wantTitle, task.title)
			}
			if tt.wantUpdateChange {
				if !task.updatedAt.After(oldTime) {
					t.Errorf("expected new update time, got:%v", task.updatedAt)
				}
			} else {
				if !task.updatedAt.Equal(oldTime) {
					t.Errorf("expected old time, got:%v", task.updatedAt)
				}
			}

		})
	}
}

func TestPriorityUpdate(t *testing.T) {
	tests := []struct {
		name             string
		initialPriority  TaskPriority
		inputPriority    TaskPriority
		wantPriority     TaskPriority
		wantErr          error
		wantUpdateChange bool
	}{
		{
			name:             "low to high",
			initialPriority:  PriorityLow,
			inputPriority:    PriorityHigh,
			wantPriority:     PriorityHigh,
			wantUpdateChange: true,
		}, {
			name:             "high to medium",
			initialPriority:  PriorityHigh,
			inputPriority:    PriorityMedium,
			wantPriority:     PriorityMedium,
			wantUpdateChange: true,
		}, {
			name:             "medium to low",
			initialPriority:  PriorityMedium,
			inputPriority:    PriorityLow,
			wantPriority:     PriorityLow,
			wantUpdateChange: true,
		}, {
			name:             "invalid priority",
			initialPriority:  PriorityLow,
			inputPriority:    "banana",
			wantPriority:     PriorityLow,
			wantErr:          ErrInvalidTaskPriority,
			wantUpdateChange: false,
		}, {
			name:             "medium to medium",
			initialPriority:  PriorityMedium,
			inputPriority:    PriorityMedium,
			wantPriority:     PriorityMedium,
			wantUpdateChange: false,
		}, {
			name:             "empty priority",
			initialPriority:  PriorityLow,
			inputPriority:    "",
			wantPriority:     PriorityLow,
			wantErr:          ErrInvalidTaskPriority,
			wantUpdateChange: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
			task := task{
				priority:  tt.initialPriority,
				updatedAt: oldTime,
			}

			err := task.UpdatePriority(tt.inputPriority)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error: %v, got: %v", tt.wantErr, err)
				}
			}
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error:%v", err)
			}

			if task.priority != tt.wantPriority {
				t.Errorf("expected priority: %v, got: %v", tt.wantPriority, task.priority)
			}
			if tt.wantUpdateChange {
				if !task.updatedAt.After(oldTime) {
					t.Errorf("expected new update time, got:%v", task.updatedAt)
				}
			} else {
				if !task.updatedAt.Equal(oldTime) {
					t.Errorf("expected old time, got:%v", task.updatedAt)
				}
			}

		})
	}
}

func TestPriorityUpdateZeroValue(t *testing.T) {
	task := task{}

	err := task.UpdatePriority("")

	if !errors.Is(err, ErrInvalidTaskPriority) {
		t.Fatalf("expected error:%v, got:%v", ErrInvalidTaskPriority, err)
	}
	if task.priority != "" {
		t.Errorf("expected empty priority, got:%v", task.priority)
	}
	if !task.updatedAt.IsZero() {
		t.Errorf("expected zero updateAt, got:%v", task.updatedAt)
	}
}
