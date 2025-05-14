package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal("Error al abrir la base de datos:", err)
	}
	defer db.Close()

	err = initDatabase()
	if err != nil {
		log.Fatal("Error al inicializar la base de datos:", err)
	}

	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetTasks(w, r)
		case http.MethodPost:
			handleCreateTask(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDatabase() error {
	log.Println("initDatabase ejecutándose")

	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL CHECK (completed IN (0,1))
	);`
	_, err := db.Exec(query)
	return err
}

func getAllTasks() ([]Task, error) {
	log.Println("Ejecutando getAllTasks")

	rows, err := db.Query("SELECT id, title, completed FROM tasks")
	if err != nil {
		log.Println("Error al ejecutar la consulta SQL:", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Completed)
		if err != nil {
			log.Println("Error al escanear fila:", err)
			continue
		}
		tasks = append(tasks, t)
	}
	log.Println("Se obtuvieron tareas:", tasks)
	return tasks, nil
}

func createTask(title string) error {
	log.Println("createTask ejecutándose con título:", title)

	_, err := db.Exec("INSERT INTO tasks (title, completed) VALUES (?, ?)", title, false)
	return err
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	log.Println("handleGetTasks ejecutándose")

	tasks, err := getAllTasks()
	if err != nil {
		http.Error(w, "Error al obtener tareas", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func handleCreateTask(w http.ResponseWriter, r *http.Request) {
	log.Println("handleCreateTask ejecutándose")

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	err = createTask(task.Title)
	if err != nil {
		http.Error(w, "No se pudo crear la tarea", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
