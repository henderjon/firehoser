package main

import (
	"net/http"
	"encoding/json"
	"testing"
)

func TestMarshalJson(t *testing.T) {
	s := http.StatusOK

	if j, err := json.Marshal(newResponse(s, 0)); err != nil {
		t.Error("json.Marshal error: ", err)
	}else{
		expected := `{"Status":200,"Message":"OK","Received":0}`
		if string(j) != expected {
			t.Error("Response error: \nexpected\n", expected, "\nactual\n", string(j))
		}
	}
}
