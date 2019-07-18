package main

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/jbowl/goloc/internal/pkg/goloc"
	"github.com/jbowl/goloc/internal/pkg/postgres"
)

/*

type geoDB struct {
	db *sql.DB
	mq *geoloc.MqAPI
}
*/

var api goloc.LocationService

func GeoLoc(Address string) (*goloc.MapAddress, error) {

	/*
		// respond to this request with a call to Mapquest API to get lat/lng
		latlng, err := gApi.mq.LatLng(Address)
		if err != nil {
			return nil, err
		}

		return &goloc.MapAddress{Address: latlng.Address,
			Lat: latlng.Lat,
			Lng: latlng.Lng}, nil
	*/

	return nil, nil
}

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

	s, ok := r.Context().Value("interface").(*postgres.LocationService)

	if !ok {
		fmt.Println("s is not type string")
	}

	if nil == s {

	}

	var ll goloc.Location
	/*
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}

		fmt.Println(string(body))
	*/
	json.NewDecoder(r.Body).Decode(&ll)

	//ll.Date = time.Now()
	/*
		gApi.CreateLocation(ll)
		//	gApi.dbs.Create(&ll)

		err := json.NewEncoder(w).Encode(&ll)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	*/
}

// return all rows in Locations table as JSON
func getLocations(w http.ResponseWriter, r *http.Request) {

	/*
		var locations []Location
		gApi.dbs.Find(&locations)
		err := json.NewEncoder(w).Encode(&locations)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
		//w.Header().Set(key"Content-Type", value: "text/plain; charset=utf8")
		w.WriteHeader(http.StatusOK)
	*/
}

/*
func (gAPI *postgres.LocationService) MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		// log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})

}
*/
// NewRouter yada yada yada
//func NewRouter(gAPI *postgres.LocationService) *mux.Router {
func NewRouter(i goloc.LocationService) *mux.Router {

	api = i
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/location", storeLocation).Methods("POST")

	//r.HandleFunc("locations", )

	r.HandleFunc("/locations", getLocations).Methods("GET")

	r.HandleFunc("/geoloc", reqGeoLoc).Methods("GET")

	//r.Use(gAPI.MiddleWare)

	return r
}
