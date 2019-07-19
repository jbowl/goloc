package postgres

import (
	//	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

// Locator yada yada yada
type Locator struct {
	Db *sql.DB
	Mq *geoloc.MqAPI
}

// InjectAPIMiddleWare yada
func (gAPI *Locator) InjectAPIMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			ctx := r.Context()
			ctx = context.WithValue(ctx, , gAPI)
			r = r.WithContext(ctx)
		*/
		next.ServeHTTP(w, r)

	})

}

/*

type LocationService interface {
	Location(ID int) (*Location, error)                   //GET by id
	Locations(ID int, Time time.Time) ([]Location, error) // GET all id and time

	CreateLocation(Loc Location) error

	GeoLoc(Address string) (*MapAddress, error) // get address by lat long
}

*/

// Location yada yada yada
func (ls *Locator) Location(ID int) (*goloc.Location, error) {
	return nil, nil
}
func (ls *Locator) Locations(email string) ([]goloc.Location, error) {

	i := make([]goloc.Location, 0)

	rows, err := ls.Db.Query("select locations.id, locations.date, locations.address,"+
		" locations.lat, locations.lng from locations, users where users.email = $1",
		email)

	if err != nil {
		return i, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var time time.Time
		var address string
		var lat, lng float32
		err = rows.Scan(&id, &time, &address, &lat, &lng)
		if err != nil {
			// handle this error
			return i, err
		}
		i = append(i, goloc.Location{ID: id, Date: time, Address: address, Lat: lat, Lng: lng})
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return i, err
	}

	return i, nil
}
func (ls *Locator) CreateLocation(loc goloc.Location) error {

	sqlStatement := `
		INSERT INTO locations (id, address, lat, lng)
		VALUES ($1, $2, $3, $4)`

	_, err := ls.Db.Exec(sqlStatement, 1, loc.Address, loc.Lat, loc.Lng)

	return err
}
func (ls *Locator) GeoLoc(Address string) (*goloc.MapAddress, error) {

	ll, err := ls.Mq.LatLng(Address)
	//mapAddress, err := ls.Mq.LatLng(location)

	if err != nil {
		return nil, err
	}

	return &goloc.MapAddress{Address: ll.Address, Lat: ll.Lat, Lng: ll.Lng}, nil

}
