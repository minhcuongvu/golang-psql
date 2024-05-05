package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func createConnection() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("CONN_STR"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db, err
}

func getCounter(db *sql.DB) (int, time.Duration, error) {
	startTime := time.Now()

	row := db.QueryRow("SELECT counter FROM counter WHERE id = 1")
	var count int
	err := row.Scan(&count)

	duration := time.Since(startTime)

	return count, duration, err
}

// incrementCounter increments the counter in the database.
func incrementCounter(db *sql.DB) error {
	_, err := db.Exec("UPDATE counter SET counter = counter + 1 WHERE id = 1")
	return err
}

// Increments the counter in the database and fetches the latest value.
func IncrementAndFetch(w http.ResponseWriter, r *http.Request) {
	db, err := createConnection()
	if err != nil {
		http.Error(w, "Failed to connect to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := incrementCounter(db); err != nil {
		http.Error(w, "Failed to increment counter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	count, duration, err := getCounter(db)
	if err != nil {
		http.Error(w, "Failed to fetch counter: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Log and send the results
	fmt.Fprintf(w, "Current Counter Value: %d\n", count)
	logDuration(w, duration)
}

// logDuration logs the duration of the database query to the response writer.
func logDuration(w http.ResponseWriter, duration time.Duration) {
	fmt.Fprintf(w, "Query took [%v ns] or [%v ms] or [%v s]\n", duration.Nanoseconds(), duration.Milliseconds(), duration.Seconds())
}

type NotFoundErr struct {
	Code    int
	Message string
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	notFoundErr := NotFoundErr{
		Code:    http.StatusNotFound,
		Message: "Not Found",
	}

	data := map[string]string{
		fmt.Sprintf("%d", notFoundErr.Code): notFoundErr.Message,
	}

	json, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(json)
}
