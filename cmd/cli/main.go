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
		if _, err := fmt.Fprint(w, "Команда не указана"); err != nil {
			log.Fatal(err)
		}
	}

	command := args[0]

	switch command {
	case "list":
		tasks, err := store.GetTasks()
		if err != nil {
			log.Fatal(err)
		}
		fio.PrintTasks(w, tasks)
	case "add":
		if len(args) < 2 {
			if _, err := fmt.Fprint(w, "Укажите название задачи"); err != nil {
				log.Fatal(err)
			}
		}

		if err := store.SaveTask(args[1]); err != nil {
			log.Fatal(err)
		}
	case "done":
		if len(args) < 2 {
			if _, err := fmt.Fprint(w, "Укажите id задачи"); err != nil {
				log.Fatal(err)
			}
		}
		taskId, err := strconv.Atoi(args[1])
		if err != nil {
			if _, err := fmt.Fprint(w, "Не удлось преобразовать строку в число"); err != nil {
				log.Fatal(err)
			}
		}

		if err := store.TaskDone(taskId); err != nil {
			log.Fatal(err)
		}
	default:
		if _, err := fmt.Fprint(w, "Неверное использование без аргументов"); err != nil {
			log.Fatal(err)
		}
	}
}
