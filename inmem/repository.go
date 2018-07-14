package inmem

import "github.com/scnewma/todo/pkg/todo"

type TaskRepository struct {
	Tasks map[int]todo.Task
}

func NewRepository() *TaskRepository {
	return &TaskRepository{
		Tasks: make(map[int]todo.Task),
	}
}

func (tr *TaskRepository) All() ([]todo.Task, error) {
	tasks := make([]todo.Task, 0, len(tr.Tasks))

	for _, value := range tr.Tasks {
		tasks = append(tasks, value)
	}

	return tasks, nil
}

func (tr *TaskRepository) Get(id int) (todo.Task, error) {
	return tr.Tasks[id], nil
}

func (tr *TaskRepository) Save(task todo.Task) error {
	tr.Tasks[task.ID] = task

	return nil
}

func (tr *TaskRepository) NextID() int {
	return len(tr.Tasks) + 1
}
