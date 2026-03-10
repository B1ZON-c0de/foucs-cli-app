package storage

import (
	"encoding/json"
	fio "focus-app/internal/io"
	"io"
	"os"
	"time"
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

func (ts *TasksStorage) SaveTask(name string) {
	now := time.Now().Round(time.Second)
	newTask := fio.Task{
		Id:        1,
		Type:      "Новая",
		Name:      name,
		CreatedAt: now,
	}

	tasks := ts.GetTasks()
	tasks = append(tasks, newTask)

	ts.file.Seek(0, io.SeekStart)
	json.NewEncoder(ts.file).Encode(tasks)

}
