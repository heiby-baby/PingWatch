package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type PingResult struct {
	IP          string    `json:"ip"`
	PingTime    time.Time `json:"ping_time"`
	LastSuccess time.Time `json:"last_sccess"`
}

var db *sql.DB

func createDB() {
	var err error

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"db", "5432", "postgres", "postgres", "postgres")

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка ping к базе данных: %v", err)
	}

}

func enableCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func addPingResult(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Вызвался метод addPingResult()")
	enableCORS(&w, r)
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var result PingResult
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO ping_results (ip, ping_time, last_success) VALUES ($1, $2, $3)",
		result.IP, result.PingTime, result.LastSuccess)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func getPingResults(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Вызвался метод getPingResults()")
	enableCORS(&w, r)
	rows, err := db.Query("SELECT ip, ping_time, last_success FROM ping_results")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var results []PingResult
	for rows.Next() {
		var result PingResult
		if err := rows.Scan(&result.IP, &result.PingTime, &result.LastSuccess); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		results = append(results, result)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func main() {
	createDB()

	http.HandleFunc("/ping-results", getPingResults)
	http.HandleFunc("/add-ping-result", addPingResult)

	log.Println("Сервер запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
