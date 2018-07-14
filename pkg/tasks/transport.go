package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(ts Service) http.Handler {
	loadTasksHandler := kithttp.NewServer(
		makeLoadAllTasksEndpoint(ts),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
	)

	createTaskHandler := kithttp.NewServer(
		makeCreateTaskEndpoint(ts),
		decodeCreateTaskRequest,
		kithttp.EncodeJSONResponse,
	)

	completeTaskHandler := kithttp.NewServer(
		makeCompleteTaskEndpoint(ts),
		decodeCompleteTaskRequest,
		kithttp.EncodeJSONResponse,
	)

	retrieveTaskHandler := kithttp.NewServer(
		makeRetrieveTaskEndpoint(ts),
		decodeRetrieveTaskRequest,
		kithttp.EncodeJSONResponse,
	)

	r := mux.NewRouter()

	r.Handle("/tasks/", loadTasksHandler).Methods("GET")
	r.Handle("/tasks/", createTaskHandler).Methods("POST")
	r.Handle("/tasks/{id}", retrieveTaskHandler).Methods("GET")
	r.Handle("/tasks/{id}/complete", completeTaskHandler).Methods("PUT")

	return r
}

var errBadRoute = errors.New("bad route")

func decodeCreateTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return createTaskRequest{Content: body.Content}, nil
}

func decodeCompleteTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := decodeTaskID(r)
	if err != nil {
		return nil, err
	}

	return completeTaskRequest{ID: id}, nil
}

func decodeRetrieveTaskRequest(_ context.Context, r *http.Request) (interface{}, error) {
	id, err := decodeTaskID(r)
	if err != nil {
		return nil, err
	}

	return retrieveTaskRequest{ID: id}, nil
}

func decodeTaskID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idstr, ok := vars["id"]
	if !ok {
		return -1, errBadRoute
	}

	id, err := strconv.Atoi(idstr)
	if err != nil {
		return -1, errBadRoute
	}

	return id, nil
}
