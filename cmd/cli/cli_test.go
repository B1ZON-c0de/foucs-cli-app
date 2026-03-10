package main

import (
	"bytes"
	"fmt"
	"focus-app/internal/storage"
	"os"
	"testing"
	"time"
)

func TestHandleCLI(t *testing.T) {
	t.Run("Аргумент list", func(t *testing.T) {
		var buf bytes.Buffer

		now, timeStr := getTime()
		displayedTime := now.Format("2006-01-02 15:04:05")

		mockData := fmt.Sprintf(`[
	{"id":1, "type":"Новая", "name":"Первая задача", "created_at":"%s"},
	{"id":2, "type":"Новая", "name":"Вторая задача", "created_at":"%s"}
	]`, timeStr, timeStr)

		tempFile, clearFile := createTempFile(t, mockData)
		defer clearFile()

		store := storage.NewTasksStorage(tempFile)

		HandleCLI("list", &buf, store)

		want := fmt.Sprintf("1. [НОВАЯ] Первая задача (%s)\n2. [НОВАЯ] Вторая задача (%s)\n",
			displayedTime, displayedTime)

		assertBuf(t, &buf, want)
	})
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
