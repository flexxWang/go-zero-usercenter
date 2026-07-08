package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Todo struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type CreateTodoRequest struct {
	Text string `json:"text"`
}

type UpdateTodoRequest struct {
	Text *string `json:"text"`
	Done *bool   `json:"done"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TodoStore struct {
	mu     sync.Mutex
	nextID int
	items  []Todo
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		nextID: 3,
		items: []Todo{
			{ID: 1, Text: "Learn basic Go syntax", Done: true},
			{ID: 2, Text: "Build a tiny HTTP API", Done: false},
		},
	}
}

func (s *TodoStore) List() []Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	items := make([]Todo, len(s.items))
	copy(items, s.items)
	return items
}

func (s *TodoStore) Create(text string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	text, err := validateTodoText(text)
	if err != nil {
		return Todo{}, err
	}

	todo := Todo{
		ID:   s.nextID,
		Text: text,
		Done: false,
	}

	s.items = append(s.items, todo)
	s.nextID++

	return todo, nil
}

func (s *TodoStore) GetByID(id int) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, item := range s.items {
		if item.ID == id {
			return item, nil
		}
	}

	return Todo{}, errors.New("todo not found")
}

func (s *TodoStore) Update(id int, req UpdateTodoRequest) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.Text == nil && req.Done == nil {
		return Todo{}, errors.New("at least one field must be provided")
	}

	for i, item := range s.items {
		if item.ID != id {
			continue
		}

		if req.Text != nil {
			text, err := validateTodoText(*req.Text)
			if err != nil {
				return Todo{}, err
			}
			item.Text = text
		}

		if req.Done != nil {
			item.Done = *req.Done
		}

		s.items[i] = item
		return item, nil
	}

	return Todo{}, errors.New("todo not found")
}

func (s *TodoStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.items {
		if item.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}

	return errors.New("todo not found")
}

type App struct {
	store *TodoStore
}

func NewApp() *App {
	return &App{
		store: NewTodoStore(),
	}
}

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", a.handleHealth)
	mux.HandleFunc("/todos", a.handleTodos)
	mux.HandleFunc("/todos/", a.handleTodoByID)
	return mux
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos := a.store.List()
		writeJSON(w, http.StatusOK, todos)
	case http.MethodPost:
		var req CreateTodoRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid json body")
			return
		}

		todo, err := a.store.Create(req.Text)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusCreated, todo)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (a *App) handleTodoByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseTodoID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		todo, err := a.store.GetByID(id)
		if err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, todo)
	case http.MethodPatch:
		var req UpdateTodoRequest

		if err := decodeJSONBody(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		todo, err := a.store.Update(id, req)
		if err != nil {
			status := http.StatusBadRequest
			if err.Error() == "todo not found" {
				status = http.StatusNotFound
			}
			writeError(w, status, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, todo)
	case http.MethodDelete:
		if err := a.store.Delete(id); err != nil {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func parseTodoID(path string) (int, error) {
	idPart := strings.TrimPrefix(path, "/todos/")
	if idPart == "" || strings.Contains(idPart, "/") {
		return 0, errors.New("invalid todo id")
	}

	id, err := strconv.Atoi(idPart)
	if err != nil {
		return 0, errors.New("todo id must be a number")
	}

	return id, nil
}

func decodeJSONBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return errors.New("invalid json body")
	}

	return nil
}

func validateTodoText(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", errors.New("text is required")
	}

	if len(text) > 120 {
		return "", errors.New("text must be at most 120 characters")
	}

	return text, nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write json response failed: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func main() {
	app := NewApp()

	server := &http.Server{
		Addr:    ":8080",
		Handler: app.routes(),
	}

	log.Println("server started at http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
