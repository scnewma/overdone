package todo

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type tasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func makeLoadAllTasksEndpoint(ts TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		tasks, err := ts.LoadAll()
		if err != nil {
			return nil, err
		}
		return tasksResponse{tasks}, nil
	}
}

type createTaskRequest struct {
	Content string
}

type createTaskResponse struct {
	ID int `json:"id"`
}

func makeCreateTaskEndpoint(ts TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTaskRequest)
		id, err := ts.Create(req.Content)
		if err != nil {
			return nil, err
		}
		return createTaskResponse{ID: id}, nil
	}
}

type completeTaskRequest struct {
	ID int
}

func makeCompleteTaskEndpoint(ts TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(completeTaskRequest)
		err := ts.MarkComplete(req.ID)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

type retrieveTaskRequest struct {
	ID int
}

func makeRetrieveTaskEndpoint(ts TaskService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(retrieveTaskRequest)
		task, err := ts.LoadByID(req.ID)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}
