package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

var api goloc.Locator

//  done
func reqGeoLoc(w http.ResponseWriter, r *http.Request) {
	/* not a good idea
	ls, ok := r.Context().Value("interface").(*postgres.LocationService)
	if !ok {

	}
	*/
	location := r.URL.Query().Get("location")

	fmt.Println(location)

	//  api.GeoLoc(location)

	//mapAddress, err := ls.Mq.LatLng(location)

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

func CreateLocation(Loc goloc.Location) error {

	return nil
}

// write one record to Loctions table
func storeLocation(w http.ResponseWriter, r *http.Request) {
	var loc goloc.Location

	err := json.NewDecoder(r.Body).Decode(&loc)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = api.CreateLocation(loc)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Location", "TODO")
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("TODO return body created"))
}

// return all rows in Locations table as JSON
func getLocations(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")

	fmt.Println(email)

	//  api.GeoLoc(location)

	//mapAddress, err := ls.Mq.LatLng(location)

	locs, err := api.Locations(email)

	// respond to this request with a call to Mapquest API to get lat/lng
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&locs)
}

func middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})

}

// NewRouter set handler funcs and middleware
func NewRouter(service goloc.Locator) *mux.Router {

	api = service
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/location", storeLocation).Methods("POST")

	r.HandleFunc("/locations", getLocations).Methods("GET")

	r.HandleFunc("/geoloc", reqGeoLoc).Methods("GET")

	r.Use(middleWare)

	return r
}
