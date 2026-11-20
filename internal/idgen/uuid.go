package idgen

import (
	"taskflow/internal/task"
	"uuid"
)

type UUIDGenerator struct{}

func (UUIDGenerator) NewID() task.TaskID {
	id := uuid.NewV4()
	return task.TaskID(id.String())
}
