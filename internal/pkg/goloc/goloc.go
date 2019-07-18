package goloc

import "time"

//Location yada yada yada
type Location struct {
	ID      int
	Date    time.Time
	Address string
	Lat     float32
	Lng     float32
}

// User yada yada yada
type User struct {
	ID    int
	Email string
}

// MapAddress yada yada yada
type MapAddress struct {
	Address string
	Lat     float32
	Lng     float32
}

// CRUD

// LocationService yada yada yada
type LocationService interface {
	Location(ID int) (*Location, error)                   //GET by id
	Locations(ID int, Time time.Time) ([]Location, error) // GET all id and time

	CreateLocation(Loc Location) error

	GeoLoc(Address string) (*MapAddress, error) // get address by lat long
}
