package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetCustomers(t *testing.T) {
	req, err := http.NewRequest("GET", "/customers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getCustomers)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":1,"name":"John Doe","role":"Admin","email":"john.doe@gmail.com","phone":"1234567890","contacted":false},{"id":2,"name":"Jane Doe","role":"User","email":"jane.doe@gmail.com","phone":"0987654321","contacted":false},{"id":3,"name":"John Smith","role":"User","email":"john.smith@gmail.com","phone":"1234567890","contacted":false}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
