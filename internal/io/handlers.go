package io

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const NoTasks = "Нет добавленных задач"

type Task struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func PrintTasks(w io.Writer, tasks []Task) {
	if tasks == nil || len(tasks) == 0 {
		fmt.Fprint(w, NoTasks)
		return
	}

	fmt.Fprint(w, getStringsTasks(tasks))
}

func getStringsTasks(tasks []Task) string {
	var sb strings.Builder

	for _, t := range tasks {
		newTask := fmt.Sprintf("%d. [%s] %s (%s)\n", t.Id, strings.ToUpper(t.Type), t.Name, t.CreatedAt.Format("2006-01-02 15:04:05"))
		sb.WriteString(newTask)
	}

	return sb.String()
}
