package main

import (
	"encoding/json"
	"testing"
)

type withTime struct {
	T  string `json:"t"`
	Tm Time   `json:"tm"`
}

func TestTime_UnmarshalJSON(t *testing.T) {
	js := `{"t": "sadasd", "tm": "2017-09-00 15:12:01"}`
	var wt withTime
	if err := json.Unmarshal([]byte(js), &wt); err == nil {
		t.Error("No error on wrong time")
		return
	}

	js = `{"t": "sadasd", "tm": "2017-09-01 15:12:01"}`
	if err := json.Unmarshal([]byte(js), &wt); err != nil {
		t.Error(err)
		return
	}

	js = `{"t": "sadasd", "tm": "0000-00-00 00:00:00"}`
	if err := json.Unmarshal([]byte(js), &wt); err != nil {
		t.Error(err)
		return
	}
}
