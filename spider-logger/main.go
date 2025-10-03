package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Message struct {
	Component string `json:"component"`
	Error     string `json:"error"`
}

func main() {
	dsn := "user:password@tcp(spider-logger-mysql:3306)/spider-logger?parseTime=true"

	var db *sql.DB
	var err error

	// Próbálkozik csatlakozni a DB-hez 3 másodpercenként
	for {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("Successfully connected to the database!")
			break
		}

		fmt.Printf("Cannot connect to DB, retrying in 3s: %v\n", err)
		time.Sleep(3 * time.Second)
	}
	defer db.Close()

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var msg Message
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		query := "INSERT INTO error_logs (component, error) VALUES (?, ?)"
		_, err := db.Exec(query, msg.Component, msg.Error)
		if err != nil {
			http.Error(w, fmt.Sprintf("DB error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"message saved"}`))
	})

	fmt.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
