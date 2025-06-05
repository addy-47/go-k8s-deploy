package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Task represents a task with ID and description
type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

// TaskStore holds tasks with a mutex for thread-safe access
type TaskStore struct {
	sync.Mutex
	tasks  map[int]Task
	nextID int
}

// Global logger for file and terminal
var logger *log.Logger

// Initialize TaskStore and logger
var store = &TaskStore{
	tasks:  make(map[int]Task),
	nextID: 1,
}

func main() {
	if err := os.MkdirAll("./logs", 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	file, err := os.OpenFile("logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()
	// Assuming you set up a logger, e.g.:
	multiWriter := io.MultiWriter(file, os.Stdout)
	logger = log.New(multiWriter, "", log.LstdFlags)

	// Register HTTP handlers
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", taskByIDHandler)
	http.HandleFunc("/health", healthHandler)

	logger.Println("Server starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}

// tasksHandler handles POST and GET requests for /tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	switch r.Method {
	case http.MethodPost:
		var task Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			logger.Printf("Error: Invalid request body: %v", err)
			return
		}
		store.Lock()
		task.ID = store.nextID
		store.tasks[task.ID] = task
		store.nextID++
		store.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)
		logger.Printf("Created task: %+v", task)
	case http.MethodGet:
		store.Lock()
		tasks := make([]Task, 0, len(store.tasks))
		for _, task := range store.tasks {
			tasks = append(tasks, task)
		}
		store.Unlock()
		json.NewEncoder(w).Encode(tasks)
		logger.Printf("Returned %d tasks", len(tasks))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Printf("Error: Method %s not allowed", r.Method)
	}
}

// taskByIDHandler handles GET requests for /tasks/{id}
func taskByIDHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		logger.Printf("Error: Method %s not allowed", r.Method)
		return
	}
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		logger.Printf("Error: Invalid task ID: %v", err)
		return
	}
	store.Lock()
	task, exists := store.tasks[id]
	store.Unlock()
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		logger.Printf("Error: Task %d not found", id)
		return
	}
	json.NewEncoder(w).Encode(task)
	logger.Printf("Returned task: %+v", task)
}

// healthHandler for Kubernetes liveness/readiness probes
func healthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("Received %s request for %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Healthy")
}
