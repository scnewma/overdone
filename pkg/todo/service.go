package todo

import (
	"errors"

	"github.com/scnewma/todo/pkg/utils"
)

type TaskService struct {
	repository TaskRepository
}

func NewService(repository TaskRepository) *TaskService {
	return &TaskService{repository}
}

func (ts *TaskService) LoadAll() ([]*Task, error) {
	return ts.repository.All()
}

func (ts *TaskService) LoadByID(id int) (*Task, error) {
	task, err := ts.repository.Get(id)
	if err != nil {
		return nil, errors.New("failed to retrieve task")
	}
	return task, nil
}

func (ts *TaskService) Create(content string) (int, error) {
	if utils.IsBlank(content) {
		return -1, errors.New("content must be provided to create a task")
	}

	task := &Task{
		ID:        ts.repository.NextID(),
		Completed: false,
		Content:   content,
	}

	ts.repository.Save(task)

	return task.ID, nil
}

func (ts *TaskService) MarkComplete(id int) error {
	task, err := ts.repository.Get(id)
	if err != nil {
		return errors.New("could not retrieve task")
	}

	task.complete()

	err = ts.repository.Save(task)
	if err != nil {
		return errors.New("failed to save task")
	}

	return nil
}
