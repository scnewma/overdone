package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/scnewma/todo/inmem"
	"github.com/scnewma/todo/pkg/tasks"
)

var a App
var tr tasks.Repository
var ts tasks.Service

func TestMain(m *testing.M) {
	tr = inmem.NewRepository()
	ts = tasks.NewService(tr)

	a = App{EnableLogging: false}
	a.Initialize(ts)

	os.Exit(m.Run())
}

func TestNoTasks(t *testing.T) {
	clearTasks()

	req, _ := http.NewRequest("GET", "/tasks", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s\n", body)
	}
}

func TestNonExistentTask(t *testing.T) {
	clearTasks()

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	res := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, res.Code)

	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "Task not found" {
		t.Errorf("Expected 'error' to be 'Task not found'. Got '%s'", m["error"])
	}
}

func TestCreateTask(t *testing.T) {
	clearTasks()

	payload := []byte(`{"content":"new task"}`)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(payload))
	res := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, res.Code)

	var m map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &m)

	if m["content"] != "new task" {
		t.Errorf("Expected 'content' to be 'new task'. Got%s\n", m["content"])
	}

	if m["completed"].(bool) {
		t.Error("Expected 'completed' to be false")
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected 'id' to be '1'. Got %d\n", m["id"])
	}
}

func TestGetTask(t *testing.T) {
	clearTasks()
	addTasks(1)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	res := executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)
}

func TestCompleteTask(t *testing.T) {
	clearTasks()
	addTasks(1)

	// verify current state
	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	res := executeRequest(req)

	var task map[string]interface{}
	json.Unmarshal(res.Body.Bytes(), &task)

	if task["completed"].(bool) {
		t.Error("Expected 'completed' to be false")
	}

	// complete task
	req, _ = http.NewRequest("PUT", "/tasks/1/complete", nil)
	res = executeRequest(req)

	checkResponseCode(t, http.StatusOK, res.Code)

	// verify completed
	json.Unmarshal(res.Body.Bytes(), &task)
	if !task["completed"].(bool) {
		t.Error("Expected 'completed' to be true")
	}
}

func TestCompleteNotExistentTask(t *testing.T) {
	clearTasks()

	req, _ := http.NewRequest("PUT", "/tasks/1/complete", nil)
	res := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, res.Code)

	var m map[string]string
	json.Unmarshal(res.Body.Bytes(), &m)
	if m["error"] != "Task not found" {
		t.Errorf("Expected 'error' to be 'Task not found'. Got '%s'", m["error"])
	}
}

func addTasks(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		ts.Create(fmt.Sprintf("Task %d", count))
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTasks() {
	tr.(*inmem.TaskRepository).Clear()
}
