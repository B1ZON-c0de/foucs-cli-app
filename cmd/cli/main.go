package main

import (
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"log"
	"os"
)

func main() {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	TasksStorage := storage.NewTasksStorage(file)

	tasks := TasksStorage.GetTasks()

	fio.PrintTasks(os.Stdout, tasks)
}
