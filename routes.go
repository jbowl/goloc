package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/mux"
	//_ "github.com/lib/pq"

	"net/http"
)

//
func (gApi *geoDB) reqGeoLoc(w http.ResponseWriter, r *http.Request) {

	location := r.URL.Query().Get("location")

	fmt.Println(location)

	// respond to this request with a call to Mapquest API to get lat/lng
	latlng, err := gApi.mq.LatLng(location)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(*latlng)
}

// write one record to Loctions table
func (gApi *geoDB) storeLocation(w http.ResponseWriter, r *http.Request) {

	var ll Location

	json.NewDecoder(r.Body).Decode(&ll)

	ll.Date = time.Now()
	gApi.dbs.Create(&ll)

	err := json.NewEncoder(w).Encode(&ll)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}

// return all rows in Locations table as JSON
func (gApi *geoDB) getLocations(w http.ResponseWriter, r *http.Request) {
	var locations []Location
	gApi.dbs.Find(&locations)
	err := json.NewEncoder(w).Encode(&locations)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
	}
}

func newRouter(gApi *geoDB) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/locations", gApi.storeLocation).Methods("POST")
	r.HandleFunc("/locations", gApi.getLocations).Methods("GET")
	r.HandleFunc("/geoloc", gApi.reqGeoLoc).Methods("GET")

	return r
}
