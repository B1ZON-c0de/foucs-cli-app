package io

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"
)

const (
	NoTasks    = "Нет добавленных задач"
	TimeFormat = "2006-01-02 15:04:05"
)

type Task struct {
	Id        int       `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func PrintTasks(w io.Writer, tasks []Task) {
	if tasks == nil || len(tasks) == 0 {
		if _, err := fmt.Fprint(w, NoTasks); err != nil {
			log.Fatal(err)
		}
		return
	}

	if _, err := fmt.Fprint(w, getStringsTasks(tasks)); err != nil {
		log.Fatal(err)
	}
}

func getStringsTasks(tasks []Task) string {
	var sb strings.Builder

	for _, t := range tasks {
		newTask := fmt.Sprintf("%d. [%s] %s (%s)\n", t.Id, strings.ToUpper(t.Type), t.Name, t.CreatedAt.Format(TimeFormat))
		sb.WriteString(newTask)
	}

	return sb.String()
}
