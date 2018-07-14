package tasks

import (
	"errors"

	"github.com/scnewma/todo/pkg/utils"
)

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
		return Task{}, errors.New("failed to retrieve task")
	}
	return task, nil
}

func (ts Service) Create(content string) (int, error) {
	if utils.IsBlank(content) {
		return -1, errors.New("content must be provided to create a task")
	}

	task := Task{
		ID:        ts.repository.NextID(),
		Completed: false,
		Content:   content,
	}

	ts.repository.Save(task)

	return task.ID, nil
}

func (ts Service) MarkComplete(id int) error {
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
