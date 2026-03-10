package storage

import (
	"encoding/json"
	fio "focus-app/internal/io"
	"io"
	"os"
	"time"
)

const (
	typeNewTask = "Новая"
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

	tasks := ts.GetTasks()
	nextId := getNextId(tasks)

	newTask := fio.Task{
		Id:        nextId,
		Type:      typeNewTask,
		Name:      name,
		CreatedAt: now,
	}

	tasks = append(tasks, newTask)

	ts.file.Truncate(0)
	ts.file.Seek(0, io.SeekStart)
	json.NewEncoder(ts.file).Encode(tasks)

}

func getNextId(tasks []fio.Task) (maxId int) {
	for _, task := range tasks {
		if task.Id > maxId {
			maxId = task.Id
		}
	}
	maxId++

	return
}
