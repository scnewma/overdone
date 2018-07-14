package tasks_test

import (
	"testing"

	"github.com/scnewma/todo/inmem"
	"github.com/scnewma/todo/pkg/tasks"
)

func fixture() tasks.Service {
	inMemRepository := inmem.NewRepository()

	return tasks.NewService(inMemRepository)
}

func TestCreateTask(t *testing.T) {
	s := fixture()

	content := "new task"
	id, _ := s.Create(content)

	task, _ := s.LoadByID(id)
	if task.Completed {
		t.Error("expected new task to not be completed")
	}

	if task.Content != content {
		t.Errorf("expected Content to be %s but was %s", content, task.Content)
	}
}

func TestLoadAll(t *testing.T) {
	s := fixture()

	s.Create("task1")
	s.Create("task2")

	tasks, _ := s.LoadAll()

	if len(tasks) != 2 {
		t.Errorf("expected %d tasks to be loaded but only found %d", 2, len(tasks))
	}
}

func TestLoadByID(t *testing.T) {
	s := fixture()

	content := "new task"
	id, _ := s.Create(content)

	task, _ := s.LoadByID(id)

	if task.Content != content {
		t.Errorf("expected Content to be %s but was %s", content, task.Content)
	}
}

func TestCompleteTask(t *testing.T) {
	s := fixture()

	content := "new task"
	id, _ := s.Create(content)

	err := s.MarkComplete(id)
	if err != nil {
		t.Errorf("did not expect an error %v", err)
	}

	task, _ := s.LoadByID(id)

	if !task.Completed {
		t.Error("expected task to be marked completed")
	}
}
