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

// Locator yada yada yada
type Locator interface {
	Locations(email string) ([]Location, error) // READ GET all

	CreateLocation(Loc Location) error  //CREATE
	
	DeleteLocation()

	GeoLoc(Address string) (*MapAddress, error) // get address by lat long
}
