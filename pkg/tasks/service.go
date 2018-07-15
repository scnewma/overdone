package tasks

import (
	"errors"

	"github.com/scnewma/overdone/pkg/utils"
)

// ErrNotFound is returned by any operation that is performed on a task object,
// but that task object was not found
var ErrNotFound = errors.New("tasks: no task found")

type Service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return Service{repository}
}

func (ts Service) LoadAll() ([]Task, error) {
	return ts.repository.All()
}

func (ts Service) LoadByID(id int) (Task, error) {
	task, err := ts.repository.Get(id)
	if err != nil {
		return Task{}, ErrNotFound
	}
	return task, nil
}

func (ts Service) Create(content string) (Task, error) {
	if utils.IsBlank(content) {
		return Task{}, errors.New("content must be provided to create a task")
	}

	task := Task{
		ID:        ts.repository.NextID(),
		Completed: false,
		Content:   content,
	}

	ts.repository.Save(task)

	return task, nil
}

func (ts Service) MarkComplete(id int) (Task, error) {
	task, err := ts.repository.Get(id)
	if err != nil {
		return Task{}, ErrNotFound
	}

	task.complete()

	err = ts.repository.Save(task)
	if err != nil {
		return Task{}, errors.New("failed to save task")
	}

	return task, nil
}
