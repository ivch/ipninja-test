package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {

	fixture := Note{
		Title:    "TestTitle",
		Body:     "TestBody",
		Canceled: 0,
	}

	fixture.CreatedAt.Time = time.Now()
	fixture.CreatedAt.Time = time.Now().Add(1 * time.Hour)
	data, _ := json.Marshal(fixture)

	req, err := http.NewRequest("POST", "/create", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(create)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var result Note
	if err := json.Unmarshal(rr.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}

	if result.Title != fixture.Title {
		t.Errorf("handler returned wrongtitle: got %v want %v",
			result.Title, fixture.Title)
	}

	//etc
}
