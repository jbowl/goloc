package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jbowl/goloc/internal/pkg/mocklocator"
	"github.com/stretchr/testify/assert"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

//https://lanre.wtf/blog/2017/04/08/testing-http-handlers-go/
func TestGetLocation(t *testing.T) {

	api := &mocklocator.MockLocator{}

	SetAPI(api)

	fmt.Println(api)
	req, err := http.NewRequest("GET", "/location/4", nil)

	checkError(err, t)

	rr := httptest.NewRecorder()

	r := mux.NewRouter()

	r.HandleFunc("/location/{id}", getLocation).Methods("GET")

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d.\n Got %d", http.StatusOK, status)
	}

	expected := string(`{"id":2,"title":"Go is cool","content":"Yeah i have been told that multiple times"}`)

	assert.JSONEq(t, expected, rr.Body.String(), "Response body differs")

}
