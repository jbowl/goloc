package goloc

import (
	"time"
)

// Define Root Domain pkg

// User just email for now
type User struct {
	ID    interface{}
	Email string
}

// Location Date - time recorded,
//          Address - describe Lat Lng could be street address just city,state etc...
//          Lat - lattitude
//          Lng - longitude
type Location struct {
	ID      interface{}
	UserID  interface{}
	Date    time.Time
	Address string
	Lat     float32
	Lng     float32
}

// MapAddress -
type MapAddress struct {
	Address string
	Lat     float32
	Lng     float32
}

// Locator - abstract interface, eventually fully CRUD
type Locator interface {
	Initialize() error // generic startup routine

	OpenDatabase() error

	Close() error

	CreateUsersTable() error
	CreateLocationsTable() error

	CreateUser(user string) (interface{}, error)

	Locations(email string) ([]Location, error) // READ GET all
	Location(id interface{}) (*Location, error) // READ GET one record

	CreateLocation(email string, Loc Location) (interface{}, error) //CREATE

	DeleteLocation(id interface{}) error // DELETE todo

	GeoLoc(Address string) (*MapAddress, error) // get address by lat long
}
