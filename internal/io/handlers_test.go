package io_test

import (
	"bytes"
	"fmt"
	fio "focus-app/internal/io"
	"testing"
	"time"
)

func TestTasksHandlers(t *testing.T) {
	t.Run("Если длина среза 0", func(t *testing.T) {
		var buf bytes.Buffer
		var data []fio.Task

		fio.PrintTasks(&buf, data)

		want := fio.NoTasks

		assertBuf(t, &buf, want)
	})

	t.Run("Вывыод среза задач в консоль", func(t *testing.T) {
		var buf bytes.Buffer
		tasks := []fio.Task{
			{1, "Новая", "Первая задача", time.Now()},
			{2, "Новая", "Вторая задача", time.Now()},
		}

		now := time.Now().Format("2006-01-02 15:04:05")

		fio.PrintTasks(&buf, tasks)

		want := fmt.Sprintf(`1. [НОВАЯ] Первая задача (%s)
2. [НОВАЯ] Вторая задача (%s)
`, now, now)

		assertBuf(t, &buf, want)
	})

}

func assertBuf(t *testing.T, buf *bytes.Buffer, want string) {
	t.Helper()

	if want != buf.String() {
		t.Errorf("Ожидали %q получили %q", want, buf.String())
	}
}
