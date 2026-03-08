package storage

import (
	"encoding/json"
	fio "focus-app/internal/io"
	"io"
	"os"
)

type TasksStorage struct {
	tasks []fio.Task
	file  *os.File
}

func NewTasksStorage(file *os.File) *TasksStorage {
	return &TasksStorage{
		tasks: []fio.Task{},
		file:  file,
	}
}

func (ts *TasksStorage) GetTasks() []fio.Task {
	ts.file.Seek(0, io.SeekStart)
	json.NewDecoder(ts.file).Decode(&ts.tasks)

	return ts.tasks
}
