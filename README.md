# goloc

Writing an android app that will record GEOLOC lat and lng on demand. The intent is to turn off google location services. 

I'm going to start with Mapquest's REST apis. https://developer.mapquest.com/documentation/open/geocoding-api/ 

~~Need to research using PostgreSQL with GORM for simplified database operations in go. https://gorm.io/~~

   GORM is slow and confusing 

Also, need to see what can be accomplished with PostGIS. https://postgis.net/

Working on proper package layout. Root interface package. 

Abstract-ish interface is making mock db testing easier.


//  abstract interface, trying to map NoSQL mongo to this is 

type Locator interface {

	Initialize() error // generic startup routine

	OpenDatabase() error

	Close()

	CreateUsersTable() error
	CreateLocationsTable() error

	CreateUser(user string) (interface{}, error)

	Locations(email string) ([]Location, error) // READ GET all
	Location(id interface{}) (*Location, error) // READ GET one record

	CreateLocation(email string, Loc Location) (interface{}, error) //CREATE

	DeleteLocation() // DELETE todo

	GeoLoc(Address string) (*MapAddress, error) // get address by lat long
}
