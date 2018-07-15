package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/scnewma/todo/inmem"
	"github.com/scnewma/todo/pkg/tasks"
	"github.com/scnewma/todo/pkg/utils"
)

const (
	defaultPort = "8080"
)

func main() {
	port := os.Getenv("PORT")
	if utils.IsBlank(port) {
		port = defaultPort
	}

	httpAddr := flag.String("http.addr", ":"+port, "HTTP Listen Address")

	flag.Parse()

	tr := inmem.NewRepository()
	ts := tasks.NewService(tr)

	a := App{}
	a.Initialize(ts)
	a.Run(*httpAddr)
}

// App bridges the gap between the business logic and the web server by
// listening for HTTP requests and calling the correct application service
type App struct {
	Router  *mux.Router
	Service tasks.Service
}

// Initialize sets up the routes for the web server
func (a *App) Initialize(ts tasks.Service) {
	a.Service = ts
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run starts the web server on the given address
func (a *App) Run(addr string) {
	errs := make(chan error, 2)

	go func() {
		log.Printf("transport=http address=%s message=listening", addr)
		errs <- http.ListenAndServe(addr, a.Router)
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	log.Printf("terminated %v", <-errs)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/tasks", a.getTasks).Methods("GET")
	a.Router.HandleFunc("/tasks", a.createTask).Methods("POST")
	a.Router.HandleFunc("/tasks/{id:[0-9]+}", a.getTask).Methods("GET")
	a.Router.HandleFunc("/tasks/{id:[0-9]+}/complete", a.completeTask).Methods("PUT")
}

func (a *App) getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := a.Service.LoadAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (a *App) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	t, err := a.Service.LoadByID(id)
	if err != nil {
		switch err {
		case tasks.ErrNotFound:
			respondWithError(w, http.StatusNotFound, "Task not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func (a *App) createTask(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	task, err := a.Service.Create(body.Content)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (a *App) completeTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	t, err := a.Service.MarkComplete(id)
	if err != nil {
		switch err {
		case tasks.ErrNotFound:
			respondWithError(w, http.StatusNotFound, "Task not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, t)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
