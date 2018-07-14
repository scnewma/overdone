package inmem

import (
	"github.com/scnewma/todo/pkg/tasks"
)

type TaskRepository struct {
	Tasks map[int]tasks.Task
}

func NewRepository() *TaskRepository {
	return &TaskRepository{
		Tasks: make(map[int]tasks.Task),
	}
}

func (tr *TaskRepository) All() ([]tasks.Task, error) {
	tasks := make([]tasks.Task, 0, len(tr.Tasks))

	for _, value := range tr.Tasks {
		tasks = append(tasks, value)
	}

	return tasks, nil
}

func (tr *TaskRepository) Get(id int) (tasks.Task, error) {
	return tr.Tasks[id], nil
}

func (tr *TaskRepository) Save(task tasks.Task) error {
	tr.Tasks[task.ID] = task

	return nil
}

func (tr *TaskRepository) NextID() int {
	return len(tr.Tasks) + 1
}
