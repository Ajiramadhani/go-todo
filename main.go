package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID   int    `jaon: "id"`
	Task string `json:"task"`
}

var todos []Todo
var nextID = 1

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	todo.ID = nextID
	nextID++
	todos = append(todos, todo)

	saveTodosToFile()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)

			saveTodosToFile()

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateTodo Todo
	err = json.NewDecoder(r.Body).Decode(&updateTodo)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for i, t := range todos {
		if t.ID == id {
			todos[i].Task = updateTodo.Task

			saveTodosToFile()

			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)
}

func saveTodosToFile() {
	data, _ := json.MarshalIndent(todos, "", "  ")
	ioutil.WriteFile("data.json", data, 0644)
}

func loadTodosFromFile() {
	file, err := os.ReadFile("data.json")
	if err != nil {
		return
	}
	json.Unmarshal(file, &todos)

	// hitung nextID berdasarkan data terakhir
	for _, t := range todos {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
}

func main() {
	loadTodosFromFile()

	r := mux.NewRouter()

	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	r.HandleFunc("/todos/edit/{id}", updateTodo).Methods("PUT")

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
