package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/scnewma/todo/inmem"
	httplogging "github.com/scnewma/todo/pkg/http/logging"
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

	httpAddr := flag.String("http-addr", ":"+port, "HTTP Listen Address")
	timeout := flag.Duration("graceful-timeout", time.Second*15, "The duration the server will wait to for existing connections to finish before shutdown")

	flag.Parse()

	tr := inmem.NewRepository()
	ts := tasks.NewService(tr)

	a := App{EnableLogging: true}
	a.Initialize(ts)
	a.Run(*httpAddr, *timeout)
}

// App bridges the gap between the business logic and the web server by
// listening for HTTP requests and calling the correct application service
type App struct {
	Router        *mux.Router
	Service       tasks.Service
	EnableLogging bool
}

// Initialize sets up the routes for the web server
func (a *App) Initialize(ts tasks.Service) {
	a.Service = ts
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.initializeMiddleware()
}

// Run starts the web server on the given address
func (a *App) Run(addr string, timeout time.Duration) {
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	go func() {
		log.Printf("transport=http address=%s message=listening", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// accept graceful shutdowns via INTERRUPT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // block until signal received

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// block up to timout period for connections to close
	srv.Shutdown(ctx)

	log.Println("shutting down")
	os.Exit(0)
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/tasks", a.getTasks).Methods("GET")
	a.Router.HandleFunc("/tasks", a.createTask).Methods("POST")
	a.Router.HandleFunc("/tasks/{id:[0-9]+}", a.getTask).Methods("GET")
	a.Router.HandleFunc("/tasks/{id:[0-9]+}/complete", a.completeTask).Methods("PUT")
}

func (a *App) initializeMiddleware() {
	if a.EnableLogging {
		a.Router.Use(loggingMiddleware)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return httplogging.NewApacheLoggingHandler(next, os.Stdout)
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
