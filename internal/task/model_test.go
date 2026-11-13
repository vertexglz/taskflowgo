package task

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name, title, description, wantTitle string
		id                                  TaskID
		priority, wantPriority              TaskPriority
		wantErr                             error
	}{
		{name: "creates task with all fields", id: "task-1", title: "Write tests", description: "Cover the model", priority: PriorityHigh, wantTitle: "Write tests", wantPriority: PriorityHigh},
		{name: "trims title", id: "task-1", title: "  Write tests\t", priority: PriorityLow, wantTitle: "Write tests", wantPriority: PriorityLow},
		{name: "defaults empty priority", id: "task-1", title: "Write tests", wantTitle: "Write tests", wantPriority: PriorityMedium},
		{name: "accepts maximum rune title", id: "task-1", title: strings.Repeat("я", maxTitleLength), wantTitle: strings.Repeat("я", maxTitleLength), wantPriority: PriorityMedium},
		{name: "accepts maximum rune description", id: "task-1", title: "Write tests", description: strings.Repeat("я", maxDescriptionLength), wantTitle: "Write tests", wantPriority: PriorityMedium},
		{name: "rejects empty ID", id: "", title: "Write tests", priority: PriorityLow, wantErr: ErrEmptyTaskID},
		{name: "rejects whitespace ID", id: " \t", title: "Write tests", priority: PriorityLow, wantErr: ErrEmptyTaskID},
		{name: "rejects empty title", id: "task-1", title: "", priority: PriorityLow, wantErr: ErrEmptyTaskTitle},
		{name: "rejects whitespace title", id: "task-1", title: " \t", priority: PriorityLow, wantErr: ErrEmptyTaskTitle},
		{name: "rejects title over rune limit", id: "task-1", title: strings.Repeat("я", maxTitleLength+1), priority: PriorityLow, wantErr: ErrTitleTooLong},
		{name: "rejects description over rune limit", id: "task-1", title: "Write tests", description: strings.Repeat("я", maxDescriptionLength+1), priority: PriorityLow, wantErr: ErrDescriptionTooLong},
		{name: "rejects unknown priority", id: "task-1", title: "Write tests", priority: "urgent", wantErr: ErrInvalidTaskPriority},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTask(tt.id, tt.title, tt.description, tt.priority)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewTask() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				if got != nil {
					t.Fatalf("NewTask() = %v, want nil on error", got)
				}
				return
			}
			if got == nil {
				t.Fatal("NewTask() = nil, want task")
			}
			if got.ID() != tt.id || got.Title() != tt.wantTitle || got.Description() != tt.description || got.Priority() != tt.wantPriority || got.Status() != StatusTodo {
				t.Errorf("task fields = id:%q title:%q description:%q priority:%q status:%q", got.ID(), got.Title(), got.Description(), got.Priority(), got.Status())
			}
			if got.CreatedAt().IsZero() || got.UpdatedAt().IsZero() {
				t.Error("timestamps must be set")
			}
			if !got.CreatedAt().Equal(got.UpdatedAt()) {
				t.Errorf("CreatedAt() = %v, UpdatedAt() = %v; want equal", got.CreatedAt(), got.UpdatedAt())
			}
		})
	}
}

func TestTaskStatusTransitions(t *testing.T) {
	tests := []struct {
		name                string
		initial, wantStatus TaskStatus
		action              func(Task) error
		wantErr             error
	}{
		{name: "starts todo task", initial: StatusTodo, action: Task.Start, wantStatus: StatusInProgress},
		{name: "does not restart in progress task", initial: StatusInProgress, action: Task.Start, wantStatus: StatusInProgress, wantErr: ErrInvalidStatusTransition},
		{name: "does not start done task", initial: StatusDone, action: Task.Start, wantStatus: StatusDone, wantErr: ErrInvalidStatusTransition},
		{name: "does not start cancelled task", initial: StatusCancelled, action: Task.Start, wantStatus: StatusCancelled, wantErr: ErrInvalidStatusTransition},
		{name: "completes in progress task", initial: StatusInProgress, action: Task.Complete, wantStatus: StatusDone},
		{name: "does not complete todo task", initial: StatusTodo, action: Task.Complete, wantStatus: StatusTodo, wantErr: ErrInvalidStatusTransition},
		{name: "does not complete done task", initial: StatusDone, action: Task.Complete, wantStatus: StatusDone, wantErr: ErrInvalidStatusTransition},
		{name: "does not complete cancelled task", initial: StatusCancelled, action: Task.Complete, wantStatus: StatusCancelled, wantErr: ErrInvalidStatusTransition},
		{name: "cancels todo task", initial: StatusTodo, action: Task.Cancel, wantStatus: StatusCancelled},
		{name: "cancels in progress task", initial: StatusInProgress, action: Task.Cancel, wantStatus: StatusCancelled},
		{name: "does not cancel done task", initial: StatusDone, action: Task.Cancel, wantStatus: StatusDone, wantErr: ErrInvalidStatusTransition},
		{name: "does not recancel cancelled task", initial: StatusCancelled, action: Task.Cancel, wantStatus: StatusCancelled, wantErr: ErrInvalidStatusTransition},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := mustNewTask(t)
			setStatus(t, task, tt.initial)
			before := fixedPastTime
			setUpdatedAt(t, task, before)
			err := tt.action(task)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("action error = %v, want %v", err, tt.wantErr)
			}
			if got := task.Status(); got != tt.wantStatus {
				t.Errorf("Status() = %q, want %q", got, tt.wantStatus)
			}
			assertUpdateTime(t, task.UpdatedAt(), before, tt.wantErr == nil)
		})
	}
}

func TestTaskUpdates(t *testing.T) {
	tests := []struct {
		name    string
		update  func(Task) error
		verify  func(*testing.T, Task)
		wantErr error
		changed bool
	}{
		{name: "updates and trims title", update: func(task Task) error { return task.UpdateTitle("  New title  ") }, verify: func(t *testing.T, task Task) {
			if task.Title() != "New title" {
				t.Errorf("Title() = %q, want New title", task.Title())
			}
		}, changed: true},
		{name: "rejects empty title", update: func(task Task) error { return task.UpdateTitle(" \t") }, verify: assertUnchanged, wantErr: ErrEmptyTaskTitle},
		{name: "rejects too long title", update: func(task Task) error { return task.UpdateTitle(strings.Repeat("я", maxTitleLength+1)) }, verify: assertUnchanged, wantErr: ErrTitleTooLong},
		{name: "keeps time for equal title", update: func(task Task) error { return task.UpdateTitle("Initial title") }, verify: assertUnchanged},
		{name: "updates description", update: func(task Task) error { return task.UpdateDescription("New description") }, verify: func(t *testing.T, task Task) {
			if task.Description() != "New description" {
				t.Errorf("Description() = %q, want New description", task.Description())
			}
		}, changed: true},
		{name: "updates description to empty", update: func(task Task) error { return task.UpdateDescription("") }, verify: func(t *testing.T, task Task) {
			if task.Description() != "" {
				t.Errorf("Description() = %q, want empty", task.Description())
			}
		}, changed: true},
		{name: "rejects too long description", update: func(task Task) error { return task.UpdateDescription(strings.Repeat("я", maxDescriptionLength+1)) }, verify: assertUnchanged, wantErr: ErrDescriptionTooLong},
		{name: "keeps time for equal description", update: func(task Task) error { return task.UpdateDescription("Initial description") }, verify: assertUnchanged},
		{name: "updates priority", update: func(task Task) error { return task.UpdatePriority(PriorityHigh) }, verify: func(t *testing.T, task Task) {
			if task.Priority() != PriorityHigh {
				t.Errorf("Priority() = %q, want %q", task.Priority(), PriorityHigh)
			}
		}, changed: true},
		{name: "rejects empty priority", update: func(task Task) error { return task.UpdatePriority("") }, verify: assertUnchanged, wantErr: ErrInvalidTaskPriority},
		{name: "rejects unknown priority", update: func(task Task) error { return task.UpdatePriority("urgent") }, verify: assertUnchanged, wantErr: ErrInvalidTaskPriority},
		{name: "keeps time for equal priority", update: func(task Task) error { return task.UpdatePriority(PriorityMedium) }, verify: assertUnchanged},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := mustNewTask(t)
			before := fixedPastTime
			setUpdatedAt(t, task, before)
			err := tt.update(task)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("update error = %v, want %v", err, tt.wantErr)
			}
			tt.verify(t, task)
			assertUpdateTime(t, task.UpdatedAt(), before, tt.changed)
		})
	}
}

func TestIsValidPriority(t *testing.T) {
	for _, priority := range []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh} {
		if !isValidPriority(priority) {
			t.Errorf("isValidPriority(%q) = false, want true", priority)
		}
	}
	for _, priority := range []TaskPriority{"", "urgent"} {
		if isValidPriority(priority) {
			t.Errorf("isValidPriority(%q) = true, want false", priority)
		}
	}
}

func mustNewTask(t *testing.T) Task {
	t.Helper()
	task, err := NewTask("task-1", "Initial title", "Initial description", PriorityMedium)
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	return task
}

var fixedPastTime = time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)

func setUpdatedAt(t *testing.T, value Task, updatedAt time.Time) {
	t.Helper()
	concrete, ok := value.(*task)
	if !ok {
		t.Fatalf("Task implementation is %T, want *task", value)
	}
	concrete.updatedAt = updatedAt
}

func setStatus(t *testing.T, task Task, status TaskStatus) {
	t.Helper()
	switch status {
	case StatusTodo:
	case StatusInProgress:
		if err := task.Start(); err != nil {
			t.Fatal(err)
		}
	case StatusDone:
		if err := task.Start(); err != nil {
			t.Fatal(err)
		}
		if err := task.Complete(); err != nil {
			t.Fatal(err)
		}
	case StatusCancelled:
		if err := task.Cancel(); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unsupported test status %q", status)
	}
}

func assertUpdateTime(t *testing.T, got, before time.Time, wantChanged bool) {
	t.Helper()
	if wantChanged && !got.After(before) {
		t.Errorf("UpdatedAt() = %v, want after %v", got, before)
	}
	if !wantChanged && !got.Equal(before) {
		t.Errorf("UpdatedAt() changed from %v to %v", before, got)
	}
}

func assertUnchanged(t *testing.T, task Task) {
	t.Helper()
	if task.Title() != "Initial title" || task.Description() != "Initial description" || task.Priority() != PriorityMedium {
		t.Errorf("task was changed: title=%q description=%q priority=%q", task.Title(), task.Description(), task.Priority())
	}
}
