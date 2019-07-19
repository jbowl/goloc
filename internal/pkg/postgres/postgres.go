package postgres

import (
	//	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

// Locator ptrs to db and mapquest api
type Locator struct {
	Db *sql.DB
	Mq *geoloc.MqAPI
}

// InjectAPIMiddleWare yada
func (ls *Locator) InjectAPIMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			ctx := r.Context()
			ctx = context.WithValue(ctx, , gAPI)
			r = r.WithContext(ctx)
		*/
		next.ServeHTTP(w, r)

	})

}

//Locations -
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

//Location -
func (ls *Locator) Location(id int) (*goloc.Location, error) {

	sqlStatement := `SELECT address, lat, lng, date FROM locations WHERE id=$1;`

	row := ls.Db.QueryRow(sqlStatement, id)
	var address string
	var lat, lng float32
	var time time.Time

	err := row.Scan(&address, &lat, &lng, &time)
	if err != nil {
		return nil, err
	}

	return &goloc.Location{ID: id, Address: address, Lat: lat, Lng: lng, Date: time}, nil
}

// CreateLocation Create a Location record for user with email
func (ls *Locator) CreateLocation(email string, loc goloc.Location) (int64, error) {

	sqlStatement := `SELECT id FROM users WHERE email=$1;`
	var userid int

	row := ls.Db.QueryRow(sqlStatement, email)
	err := row.Scan(&userid)
	if err != nil {
		return -1, err
	}

	// LastInsertedId isn't implemented for postgreSQL
	sqlStatement = `
		INSERT INTO locations (userid, address, lat, lng)
		VALUES ($1, $2, $3, $4) RETURNING id`

	var id int64
	err = ls.Db.QueryRow(sqlStatement, userid, loc.Address, loc.Lat, loc.Lng).Scan(&id)

	if err != nil {
		return -1, err
	}

	return id, nil
}

// GeoLoc - call thru to mapquest api
func (ls *Locator) GeoLoc(Address string) (*goloc.MapAddress, error) {

	ll, err := ls.Mq.LatLng(Address)

	if err != nil {
		return nil, err
	}

	return &goloc.MapAddress{Address: ll.Address, Lat: ll.Lat, Lng: ll.Lng}, nil
}

// DeleteLocation todo
func (ls *Locator) DeleteLocation() {

}

// CreateUsersTable creates if doesn't exist
func (ls *Locator) CreateUsersTable() error {

	const sql = `
	CREATE TABLE IF NOT EXISTS users (
		id serial PRIMARY KEY,
		email TEXT UNIQUE NOT NULL
	)`

	// Exec executes a query without returning any rows.
	res, err := ls.Db.Exec(sql)
	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

// CreateLocationsTable creates if doesn't exist
func (ls *Locator) CreateLocationsTable() error {

	const sql = `
	CREATE TABLE IF NOT EXISTS locations (
		id SERIAL PRIMARY KEY, 
		address text, 
		userid int, 
		lat NUMERIC NOT NULL, 
		lng NUMERIC NOT NULL, 
		date DATE NOT NULL DEFAULT CURRENT_DATE
	)`

	// Exec executes a query without returning any rows.
	_, err := ls.Db.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

func dsn() string {

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable",
		"localhost",
		"5432",
		"postgres",
		"postgres",
		"loc_db")

	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("DB_USERNAME"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_NAME"))
	return dsn
}

//OpenDatabase implementation for generic interface call
func (ls *Locator) OpenDatabase() (*sql.DB, error) {
	return sql.Open("postgres", dsn())
}

// Initialize - not currently called for postrgresql implementation
func (ls *Locator) Initialize() error {
	return nil
}
