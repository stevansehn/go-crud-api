package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    "time"
	"encoding/json"
	"net/http"

    _ "github.com/go-sql-driver/mysql"
    "github.com/joho/godotenv"
	"github.com/gorilla/mux"
)

type User struct {
    ID        int
    Username  string
    Password  string
    CreatedAt time.Time
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

var db *sql.DB

func connectToDb() {
	var err error
    err = godotenv.Load()
    checkErr(err)

    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbName := os.Getenv("DB_NAME")
    parseTime := os.Getenv("DB_PARSE_TIME")

    dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=%s", dbUser, dbPassword, dbHost, dbPort, dbName, parseTime)

    db, err = sql.Open("mysql", dsn)
    checkErr(err)

    checkErr(db.Ping())

    query := `
        CREATE TABLE IF NOT EXISTS users (
            id INT AUTO_INCREMENT,
            username TEXT NOT NULL,
            password TEXT NOT NULL,
            created_at DATETIME,
            PRIMARY KEY (id)
        );`
    _, err = db.Exec(query)
    checkErr(err)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var newUser User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	checkErr(err)

	newUser.CreatedAt = time.Now()

	result, err := db.Exec(`INSERT INTO users (username, password, created_at) VALUES (?, ?, ?)`, newUser.Username, newUser.Password, newUser.CreatedAt)
	checkErr(err)

	_, err = result.LastInsertId()
	checkErr(err)

	w.WriteHeader(http.StatusCreated) 
	json.NewEncoder(w).Encode(newUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	err := db.QueryRow(`SELECT id, username, password, created_at FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			checkErr(err)
		}
		return
	}

	json.NewEncoder(w).Encode(u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var u User
	err := db.QueryRow(`SELECT id, username, password, created_at FROM users WHERE id = ?`, id).Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
		} else {
			checkErr(err)
		}
		return
	}

	var newUser User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	checkErr(err)

	newUser.CreatedAt = time.Now()

	result, err := db.Exec(`UPDATE users SET username = ?, password = ?, created_at = ? WHERE id = ?`, newUser.Username, newUser.Password, newUser.CreatedAt, id)
	checkErr(err)

	rowsAffected, err := result.RowsAffected()
	checkErr(err)

	if rowsAffected > 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"] 

	result, err := db.Exec(`DELETE FROM users WHERE id = ?`, id)
	checkErr(err)

	rowsAffected, err := result.RowsAffected()
	checkErr(err)

	if rowsAffected > 0 {
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, username, password, created_at FROM users")
    checkErr(err)
    defer rows.Close()

    var users []User 

    for rows.Next() {
        var u User
        checkErr(rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt))
        users = append(users, u) 
    }
    checkErr(rows.Err())

	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
}

func main() {
	connectToDb()
	defer db.Close()

	r := mux.NewRouter()

    usersRouter := r.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("/", createUser).Methods("POST")
	usersRouter.HandleFunc("/{id}", getUser).Methods("GET")
	usersRouter.HandleFunc("/{id}", updateUser).Methods("PUT")
	usersRouter.HandleFunc("/{id}", deleteUser).Methods("DELETE")
	usersRouter.HandleFunc("/", getAllUsers).Methods("GET")

	log.Println("Listening on port 8000...")
    log.Fatal(http.ListenAndServe(":8000", r))
}
