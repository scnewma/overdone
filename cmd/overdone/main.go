package main

import (
	"fmt"
	"log"

	"github.com/scnewma/todo/inmem"
	"github.com/scnewma/todo/pkg/todo"
)

func main() {
	inmem := inmem.NewRepository()
	taskService := todo.NewService(inmem)

	id, err := taskService.Create("do something")
	if err != nil {
		log.Fatal(err)
	}

	task, err := taskService.LoadByID(id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("task retrieved %s - %t\n", task.Content, task.Completed)

	taskService.MarkComplete(task.ID)
	task, err = taskService.LoadByID(task.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("task retrieved %s - %t\n", task.Content, task.Completed)
}
