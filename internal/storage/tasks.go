package storage

import (
	"encoding/json"
	"errors"
	fio "focus-app/internal/io"
	"io"
	"os"
	"time"
)

const (
	typeNewTask  = "Новая"
	typeDoneTask = "Выполнено"

	ErrNotFoundTask = "не удалось найти задачу"
	ErrIncorrectId  = "некорректный id"
	ErrNotGetTasks  = "не удалось получить задачи"
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

func (ts *TasksStorage) GetTasks() ([]fio.Task, error) {
	if _, err := ts.file.Seek(0, io.SeekStart); err != nil {
		return []fio.Task{}, err
	}

	if err := json.NewDecoder(ts.file).Decode(&ts.tasks); err != nil {
		return []fio.Task{}, err
	}

	return ts.tasks, nil
}

func (ts *TasksStorage) SaveTask(name string) error {
	now := time.Now().Round(time.Second)

	tasks, err := ts.GetTasks()
	if err != nil {
		return errors.New(ErrNotGetTasks)
	}

	nextId := getNextId(tasks)

	newTask := fio.Task{
		Id:        nextId,
		Type:      typeNewTask,
		Name:      name,
		CreatedAt: now,
	}

	tasks = append(tasks, newTask)

	if err := ts.file.Truncate(0); err != nil {
		return err
	}

	if _, err := ts.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := json.NewEncoder(ts.file).Encode(tasks); err != nil {
		return err
	}

	return nil
}

func (ts *TasksStorage) TaskDone(id int) error {
	tasks, err := ts.GetTasks()
	if err != nil {
		return errors.New(ErrNotGetTasks)
	}

	task, err := getTaskById(id, tasks)
	if err != nil {
		return err
	}

	task.Type = typeDoneTask

	if err := ts.file.Truncate(0); err != nil {
		return err
	}

	if _, err := ts.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := json.NewEncoder(ts.file).Encode(ts.tasks); err != nil {
		return err
	}

	return nil
}

func (ts *TasksStorage) TaskDelete(id int) error {
	indexTask := -1

	tasks, err := ts.GetTasks()
	if err != nil {
		return errors.New(ErrNotGetTasks)
	}

	if id < 0 || id > len(tasks) {
		return errors.New(ErrIncorrectId)
	}

	for i, task := range tasks {
		if task.Id == id {
			indexTask = i
			break
		}
	}

	if indexTask == -1 {
		return errors.New(ErrNotFoundTask)
	}

	tasks = append(tasks[:indexTask], tasks[indexTask+1:]...)

	if err := ts.file.Truncate(0); err != nil {
		return err
	}

	if _, err := ts.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := json.NewEncoder(ts.file).Encode(tasks); err != nil {
		return err
	}

	return nil
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

func getTaskById(id int, tasks []fio.Task) (*fio.Task, error) {
	for i := range tasks {
		if tasks[i].Id == id {
			return &tasks[i], nil
		}
	}
	return nil, errors.New(ErrNotFoundTask)
}
