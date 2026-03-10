package main

import (
	"bytes"
	"fmt"
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHandleCLI(t *testing.T) {
	now, timeStr := getTime()
	displayedTime := now.Format("2006-01-02 15:04:05")

	t.Run("Аргумент list", func(t *testing.T) {
		var buf bytes.Buffer

		mockData := fmt.Sprintf(`[
	{"id":1, "type":"Новая", "name":"Первая задача", "created_at":"%s"},
	{"id":2, "type":"Новая", "name":"Вторая задача", "created_at":"%s"}
	]`, timeStr, timeStr)

		tempFile, clearFile := createTempFile(t, mockData)
		defer clearFile()

		store := storage.NewTasksStorage(tempFile)

		args := []string{"list"}
		HandleCLI(args, &buf, store)

		want := fmt.Sprintf("1. [НОВАЯ] Первая задача (%s)\n2. [НОВАЯ] Вторая задача (%s)\n",
			displayedTime, displayedTime)

		assertBuf(t, &buf, want)
	})

	t.Run("Аргумент add", func(t *testing.T) {
		var buf bytes.Buffer

		tempFile, clearFile := createTempFile(t, "[]")
		defer clearFile()

		store := storage.NewTasksStorage(tempFile)
		args := []string{"add", "Новая задача"}
		HandleCLI(args, &buf, store)

		got, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		want := []fio.Task{
			{1, "Новая", "Новая задача", now},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Ожидали %v Получили %v", want, got)
		}
	})

	t.Run("Аргумент done", func(t *testing.T) {
		var buf bytes.Buffer
		tempFile, clearFile := createTempFile(t, "[]")
		defer clearFile()

		store := storage.NewTasksStorage(tempFile)

		if err := store.SaveTask("New Task"); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveTask("New Task"); err != nil {
			t.Fatal(err)
		}

		args := []string{"done", "1"}
		HandleCLI(args, &buf, store)

		tasks, err := store.GetTasks()
		if err != nil {
			t.Fatal(err)
		}

		got := getTaskById(1, tasks).Type
		want := "Выполнено"

		if got != want {
			t.Errorf("Ожидали %v Получили %v", want, got)
		}
	})
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

func getTime() (now time.Time, timeStr string) {
	now = time.Now().Round(time.Second)
	timeStr = now.Format(time.RFC3339)

	return
}

func assertBuf(t *testing.T, buf *bytes.Buffer, want string) {
	t.Helper()

	if want != buf.String() {
		t.Errorf("Ожидали %q получили %q", want, buf.String())
	}
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
