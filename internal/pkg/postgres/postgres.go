package postgres

import (
	//	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

// LocationService yada yada yada
type LocationService struct {
	Db *sql.DB
	Mq *geoloc.MqAPI
}

// InjectAPIMiddleWare yada
func (gAPI *LocationService) InjectAPIMiddleWare(next http.Handler) http.Handler {
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
func (ls *LocationService) Location(ID int) (*goloc.Location, error) {
	return nil, nil
}
func (ls *LocationService) Locations(ID int, Time time.Time) ([]goloc.Location, error) {
	i := make([]goloc.Location, 0)
	return i, nil
}
func (ls *LocationService) CreateLocation(goloc.Location) error {
	return nil
}
func (ls *LocationService) GeoLoc(Address string) (*goloc.MapAddress, error) {

	ll, err := ls.Mq.LatLng(Address)
	//mapAddress, err := ls.Mq.LatLng(location)

	if err != nil {
		return nil, err
	}

	return &goloc.MapAddress{Address: ll.Address, Lat: ll.Lat, Lng: ll.Lng}, nil

}

/*

func (ls *LocationService) writelocation() {

	sqlStatement := `
INSERT INTO users (age, email, first_name, last_name)
VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(sqlStatement, 30, "jon@calhoun.io", "Jonathan", "Calhoun")
	if err != nil {
		panic(err)
	}
}
*/
