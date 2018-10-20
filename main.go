package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

//ErrorResponse represents error response json structure
type ErrorResponse struct {
	ErrMessage string `json:"errMessage"`
}

func init() {
	var err error
	db, err = sql.Open("sqlite3", getenv("DBPATH", "ipn.db"))
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/").HandlerFunc(index)
	r.Methods("GET", "HEAD").Path("/static/{file:[a-zA-Z0-9-_./]+}").HandlerFunc(static)

	r.Methods("POST").Path("/note").HandlerFunc(create)
	r.Methods("GET").Path("/note/{id:[0-9]+}").HandlerFunc(read)
	r.Methods("PUT").Path("/note/{id:[0-9]+}").HandlerFunc(update)
	r.Methods("DELETE").Path("/note/{id:[0-9]+}").HandlerFunc(delete)
	r.Methods("GET").Path("/notes").HandlerFunc(list)

	r.HandleFunc("/expire", expire)

	fmt.Println("started http server on port 8080 @", time.Now().Format("15:04:05 02/01/2006"))

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func encodeHTTPResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)

	json.NewEncoder(w).Encode(ErrorResponse{
		ErrMessage: err.Error(),
	})
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
