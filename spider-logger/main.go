package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ErrorLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Component string    `json:"component"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	dsn := "user:password@tcp(spider-logger-mysql:3306)/spider-logger?parseTime=true"

	var db *gorm.DB
	var err error

	// Próbálkozik csatlakozni 3 másodpercenként, amíg nem sikerül
	for {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			sqlDB, _ := db.DB()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				fmt.Println("Successfully connected to the database!")
				break
			} else {
				err = pingErr
			}
		}

		fmt.Printf("Cannot connect to DB, retrying in 3s: %v\n", err)
		time.Sleep(3 * time.Second)
	}

	// Automatikus migráció — létrehozza az `error_logs` táblát, ha nem létezik
	if err := db.AutoMigrate(&ErrorLog{}); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var msg ErrorLog
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
			return
		}

		if msg.Component == "" || msg.Error == "" {
			http.Error(w, "Missing fields: component and error are required", http.StatusBadRequest)
			return
		}

		// Mentés az adatbázisba GORM-mal
		if err := db.Create(&msg).Error; err != nil {
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
