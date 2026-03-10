package main

import (
	fio "focus-app/internal/io"
	"focus-app/internal/storage"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("tasks.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	TasksStorage := storage.NewTasksStorage(file)

	TasksStorage.SaveTask("Первая")
	TasksStorage.SaveTask("Вторая")
	TasksStorage.SaveTask("Третья")

	tasks := TasksStorage.GetTasks()

	fio.PrintTasks(os.Stdout, tasks)
}
