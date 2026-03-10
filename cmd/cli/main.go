package main

import (
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.NewTasksStorage(file)
	args := os.Args[1:]

	HandleCLI(args[0], os.Stdout, store)
}

func HandleCLI(command string, w io.Writer, store *storage.TasksStorage) {
	switch command {
	case "list":
		tasks := store.GetTasks()
		fio.PrintTasks(w, tasks)
	}
}
