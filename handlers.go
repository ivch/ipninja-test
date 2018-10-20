package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/go-playground/validator.v9"
)

//Note describes note entity
type Note struct {
	ID        int64  `json:"id"`
	Title     string `validate:"required,max=100" json:"title"`
	Body      string `validate:"required,max=200" json:"body"`
	CreatedAt Time   `json:"created_at"`
	ExpiresAt Time   `json:"expires_at"`
	Canceled  int8   `json:"canceled,omitempty"`
}

func index(w http.ResponseWriter, _ *http.Request) {
	fp, _ := os.Open("static/index.html")
	io.Copy(w, fp)
}

func static(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]

	if _, err := os.Stat("static/" + file); err != nil {
		w.WriteHeader(http.StatusNotFound)
		fp, _ := os.Open("static/404.html")
		io.Copy(w, fp)
		return
	}

	fp, _ := os.Open("static/" + file)
	io.Copy(w, fp)
}

func create(w http.ResponseWriter, r *http.Request) {
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		encodeHTTPError(w, err)
		return
	}

	if err := validator.New().Struct(note); err != nil {
		encodeHTTPError(w, err)
		return
	}

	note.CreatedAt.Time = time.Now()

	if note.ExpiresAt.IsZero() {
		note.ExpiresAt.Time = note.CreatedAt.Add(1 * time.Hour)
	}

	res, err := db.Exec("insert into notes (title, body, created_at, expires_at) values ($1, $2, $3, $4)",
		note.Title, note.Body, note.CreatedAt, note.ExpiresAt)

	if err != nil {
		encodeHTTPError(w, err)
		return
	}

	note.ID, _ = res.LastInsertId()

	encodeHTTPResponse(w, note)
}

func read(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) //we can skip error here because of router settings [0-9]

	if id < 1 {
		encodeHTTPError(w, errors.New("id must be greater than 0"))
		return
	}
	row := db.QueryRow("select * from notes where id = $1 limit 1", id)

	var note Note
	if err := row.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.ExpiresAt, &note.Canceled); err != nil {
		encodeHTTPError(w, err)
		return
	}

	encodeHTTPResponse(w, note)
}

func list(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.Query("select id, title, body, created_at, expires_at, canceled from notes order by created_at desc")
	if err != nil {
		encodeHTTPError(w, err)
		return
	}
	defer rows.Close()

	var notes []Note

	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.Title, &n.Body, &n.CreatedAt, &n.ExpiresAt, &n.Canceled); err != nil {
			//in real life project here should be normal logging
			log.Println(err)
		}
		notes = append(notes, n)
	}

	encodeHTTPResponse(w, notes)
}

func update(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) //we can skip error here because of router settings [0-9]

	if id < 1 {
		encodeHTTPError(w, errors.New("id must be greater than 0"))
		return
	}

	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		encodeHTTPError(w, err)
		return
	}

	if err := validator.New().Struct(note); err != nil {
		encodeHTTPError(w, err)
		return
	}

	if _, err := db.Exec("update notes set title = $1, body = $2, expires_at = $3, canceled = $4 where id = $5",
		note.Title, note.Body, note.ExpiresAt, note.Canceled, id); err != nil {
		encodeHTTPError(w, err)
		return
	}

	row := db.QueryRow("select * from notes where id = $1 limit 1", id)
	if err := row.Scan(&note.ID, &note.Title, &note.Body, &note.CreatedAt, &note.ExpiresAt); err != nil {
		encodeHTTPError(w, err)
		return
	}

	encodeHTTPResponse(w, note)
}

func delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) //we can skip error here because of router settings [0-9]

	if id < 1 {
		encodeHTTPError(w, errors.New("id must be greater than 0"))
		return
	}

	if _, err := db.Exec("delete from notes where id = $1", id); err != nil {
		encodeHTTPError(w, err)
		return
	}

	encodeHTTPResponse(w, nil)
}

func expire(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	go func(conn *websocket.Conn) {
		for {
			rows, err := db.Query("select * from notes where expires_at < datetime('now') and canceled = 0")
			if err != nil {
				encodeHTTPError(w, err)
				return
			}
			var expired []Note

			for rows.Next() {
				var e Note
				if err := rows.Scan(&e.ID, &e.Title, &e.Body, &e.CreatedAt, &e.ExpiresAt, &e.Canceled); err != nil {
					log.Println(err)
				}
				expired = append(expired, e)
			}

			if err = conn.WriteJSON(expired); err != nil {
				fmt.Println(err)
			}

			rows.Close()

			time.Sleep(1 * time.Minute)
		}
	}(conn)
}
