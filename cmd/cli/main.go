package main

import (
	"fmt"
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewTasksStorage(file)
	args := os.Args[1:]

	HandleCLI(args, os.Stdout, store)
}

func HandleCLI(args []string, w io.Writer, store *storage.TasksStorage) {
	if len(args) == 0 {
		fmt.Fprint(w, "Команда не указана")
	}

	command := args[0]

	switch command {
	case "list":
		tasks := store.GetTasks()
		fio.PrintTasks(w, tasks)
	case "add":
		if len(args) < 2 {
			fmt.Fprint(w, "Укажите название задачи")
		}
		store.SaveTask(args[1])
	case "done":
		if len(args) < 2 {
			fmt.Fprint(w, "Укажите id задачи")
		}
		taskId, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprint(w, "Не удлось преобразовать строку в число")
		}
		store.TaskDone(taskId)
	default:
		fmt.Fprint(w, "Неверное использование без аргументов")
	}
}
