package mongoloc

import (
	"context"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Locator ptrs to db and mapquest api
type Locator struct {
	Client *mongo.Client
	Mq     *geoloc.MqAPI
}

// User -
type User struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email,omitempty" bson:"email,omitempty"`
}

func dsn() string {

	return "mongodb://localhost:27017"
}

// Close -
func (ls *Locator) Close() {}

//OpenDatabase implementation for generic interface call
func (ls *Locator) OpenDatabase() error {

	clientOptions := options.Client().ApplyURI(dsn())
	client, err := mongo.Connect(context.TODO(), clientOptions)

	ls.Client = client
	return err
}

// Initialize - not currently called for postrgresql implementation
func (ls *Locator) Initialize() error {
	return nil
}

//Locations -
func (ls *Locator) Locations(email string) ([]goloc.Location, error) {

	i := make([]goloc.Location, 0)
	return i, nil
}

//Location -
func (ls *Locator) Location(interface{}) (*goloc.Location, error) {
	return nil, nil
}

// CreateLocation Create a Location record for user with email
func (ls *Locator) CreateLocation(email string, loc goloc.Location) (interface{}, error) {
	return -1, nil
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

// CreateUser -
func (ls *Locator) CreateUser(email string) (interface{}, error) {

	user := User{Email: email}

	collection := ls.Client.Database("loc_db").Collection("users")

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

// CreateUsersTable
func (ls *Locator) CreateUsersTable() error {

	ls.Client.Database("loc_db").Collection("users")

	return nil
}

// CreateLocationsTable creates if doesn't exist
func (ls *Locator) CreateLocationsTable() error {
	ls.Client.Database("loc_db").Collection("locations")

	return nil
}
