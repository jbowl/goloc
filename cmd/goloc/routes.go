package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

var api goloc.Locator // global to this file, TODO handle a better way

func reqGeoLoc(w http.ResponseWriter, r *http.Request) {
	/* not a good idea
	ls, ok := r.Context().Value("interface").(*postgres.LocationService)
	if !ok {

	}
	*/
	location := r.URL.Query().Get("location")

	mapAddress, err := api.GeoLoc(location)

	// respond to this request with a call to Mapquest API to get lat/lng
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*mapAddress)

}

func storeUser(w http.ResponseWriter, r *http.Request) {

	var user goloc.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := api.CreateUser(user.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var loc string

	switch id.(type) {
	case int:
		loc = fmt.Sprintf("%s/user/%d", r.Host, id.(int))
	case string:
		loc = fmt.Sprintf("%s/user/%s", r.Host, id.(string))
	default:
		loc = "unknown id type"
	}

	w.Header().Set("Location", loc)
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("TODO return body for create request"))
}

// write one record to Loctions table
func storeLocation(w http.ResponseWriter, r *http.Request) {

	user := r.URL.Query().Get("user")

	if len(user) < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var loc goloc.Location
	err := json.NewDecoder(r.Body).Decode(&loc)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := api.CreateLocation(user, loc)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	var location string

	switch id.(type) {
	case int:
		location = fmt.Sprintf("%s/location/%d", r.Host, id.(int))
	case string:
		location = fmt.Sprintf("%s/locattion/%s", r.Host, id.(string))
	default:
		location = "unknown id type"
	}

	w.Header().Set("Location", location)
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("TODO return body for create requst"))
}

// return all rows in Locations table as JSON
func getLocations(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")

	fmt.Println(email)

	locs, err := api.Locations(email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&locs)
}

// return all rows in Locations table as JSON
func getLocation(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id argument"))
		return
	}

	// respond to this request with a call to Mapquest API to get lat/lng
	locs, err := api.Location(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&locs)
}

func deleteLocation(w http.ResponseWriter, r *http.Request) {

	// can't pass mongo id like this string
	vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["id"])
	id := vars["id"]
	fmt.Println(id)

	err := api.DeleteLocation(id)

	// respond to this request with a call to Mapquest API to get lat/lng
	//	locs, err := api.Location(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	//	json.NewEncoder(w).Encode(&locs)
}

// TODO add logging messages
func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// SetAPI -
func SetAPI(service goloc.Locator) {
	api = service
}

// NewRouter set handler funcs and middleware
func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/user", storeUser).Methods("POST") //create
	r.HandleFunc("/location", storeLocation).Methods("POST")

	r.HandleFunc("/location/{id}", getLocation).Methods("GET")
	r.HandleFunc("/locations", getLocations).Methods("GET")

	r.HandleFunc("/location/{id}", deleteLocation).Methods("DELETE")

	r.HandleFunc("/geoloc", reqGeoLoc).Methods("GET")

	r.Use(middleWare)

	return r
}
