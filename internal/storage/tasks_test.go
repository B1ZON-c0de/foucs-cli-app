package storage_test

import (
	"errors"
	"fmt"
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	ErrNotFound = "задача не найдена"
)

func TestStorage(t *testing.T) {
	t.Run("Получение задач из json файла", func(t *testing.T) {
		now, timeStr := getTime()

		mockTasks := fmt.Sprintf(`[
{"id":1, "type":"Новая", "name":"Первая задача", "created_at":"%s"},
{"id":2, "type":"Новая", "name":"Вторая задача", "created_at":"%s"}
]`, timeStr, timeStr)

		tempFile, clearFile := createTempFile(t, mockTasks)
		defer clearFile()
		store := storage.NewTasksStorage(tempFile)
		got, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		want := []fio.Task{
			{1, "Новая", "Первая задача", now},
			{2, "Новая", "Вторая задача", now},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Ожидали %v получили %v", want, got)
		}
	})

	t.Run("Добавление задачи", func(t *testing.T) {
		now, _ := getTime()

		tempFile, clearFile := createTempFile(t, "[]")
		defer clearFile()

		store := storage.NewTasksStorage(tempFile)

		if err := store.SaveTask("Новая задача"); err != nil {
			t.Fatal(err)
		}

		if err := store.SaveTask("Новая задача"); err != nil {
			t.Fatal(err)
		}

		if err := store.SaveTask("Новая задача"); err != nil {
			t.Fatal(err)
		}

		got, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		want := []fio.Task{
			{1, "Новая", "Новая задача", now},
			{2, "Новая", "Новая задача", now},
			{3, "Новая", "Новая задача", now},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Ожидали %v получили %v", want, got)
		}

	})

	t.Run("Выполнение задачи", func(t *testing.T) {
		_, timeStr := getTime()

		mockTasks := fmt.Sprintf(`[
{"id":1, "type":"Новая", "name":"Первая задача", "created_at":"%s"},
{"id":2, "type":"Новая", "name":"Вторая задача", "created_at":"%s"}
]`, timeStr, timeStr)

		tempFile, clearFile := createTempFile(t, mockTasks)
		defer clearFile()
		store := storage.NewTasksStorage(tempFile)

		if err := store.TaskDone(1); err != nil {
			t.Fatal(err)
		}

		tasks, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		task, err := getTaskById(1, tasks)
		if err != nil {
			t.Errorf("Ошибки быть не должно но она появилась")
		}

		got := task.Type
		want := "Выполнено"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Ожидали %q Получил  %q", want, got)
		}

	})

	t.Run("Удаление задачи", func(t *testing.T) {
		now, timeStr := getTime()

		mockTasks := fmt.Sprintf(`[
{"id":1, "type":"Новая", "name":"Первая задача", "created_at":"%s"},
{"id":2, "type":"Новая", "name":"Вторая задача", "created_at":"%s"}
]`, timeStr, timeStr)

		tempFile, clearFile := createTempFile(t, mockTasks)
		defer clearFile()
		store := storage.NewTasksStorage(tempFile)

		if err := store.TaskDelete(1); err != nil {
			t.Fatal(err)
		}

		got, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		want := []fio.Task{
			{2, "Новая", "Вторая задача", now},
		}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("Ожидали %v Получил  %v", want, got)
		}

	})
}

func getTime() (now time.Time, timeStr string) {
	now = time.Now().Round(time.Second)
	timeStr = now.Format(time.RFC3339)

	return
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()

	tempFile, err := os.CreateTemp("", "temp-tasks.json")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tempFile.Write([]byte(initialData)); err != nil {
		t.Fatal(err)
	}

	clearFile := func() {
		if err := tempFile.Close(); err != nil {
			t.Fatal(err)
		}

		if err := os.Remove(tempFile.Name()); err != nil {
			t.Fatal(err)
		}
	}

	return tempFile, clearFile
}

func getTaskById(id int, tasks []fio.Task) (fio.Task, error) {
	for _, task := range tasks {
		if task.Id == id {
			return task, nil
		}
	}

	return fio.Task{}, errors.New(ErrNotFound)
}
