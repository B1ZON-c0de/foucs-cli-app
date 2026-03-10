package storage_test

import (
	"fmt"
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"os"
	"reflect"
	"testing"
	"time"
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
		got := store.GetTasks()

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

		store.SaveTask("Новая задача")
		store.SaveTask("Новая задача")
		store.SaveTask("Новая задача")

		got := store.GetTasks()
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

		store.TaskDone(1)
		task := getTaskById(1, store.GetTasks())
		got := task.Type
		want := "Выполнено"

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Ожидали %q Получил  %q", want, got)
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

	tempFile.Write([]byte(initialData))
	clearFile := func() {
		os.Remove(tempFile.Name())
		tempFile.Close()
	}

	return tempFile, clearFile
}

func getTaskById(id int, tasks []fio.Task) fio.Task {
	var foundedTask fio.Task
	for _, task := range tasks {
		if task.Id == id {
			foundedTask = task
		}
	}

	return foundedTask
}
