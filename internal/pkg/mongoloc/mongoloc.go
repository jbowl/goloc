package mongoloc

import (
	"context"
	"fmt"
	"time"

	"github.com/jbowl/goloc/internal/pkg/geoloc"
	"github.com/jbowl/goloc/internal/pkg/goloc"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Locator ptrs to db and mapquest api
type Locator struct {
	Client *mongo.Client
	DB     *mongo.Database
	Mq     *geoloc.MqAPI
}

// User -
type User struct {
	ID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string             `json:"email,omitempty" bson:"email,omitempty"`
}

type Location struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID  primitive.ObjectID `json:"_userid,omitempty" bson:"_userid,omitempty"`
	Date    time.Time
	Address string
	Lat     float32
	Lng     float32
}

func dsn() string {

	return "mongodb://localhost:27017"
}

// Close -
func (ls *Locator) Close() error { return nil }

//OpenDatabase implementation for generic interface call
func (ls *Locator) OpenDatabase() error {

	clientOptions := options.Client().ApplyURI(dsn())
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	ls.Client = client

	ls.DB = ls.Client.Database("loc_db")

	return err
}

// Initialize - not currently called for postrgresql implementation
func (ls *Locator) Initialize() error {
	return nil
}

//Locations -
func (ls *Locator) Locations(email string) ([]goloc.Location, error) {

	users := ls.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	filter := bson.M{"email": email}

	var user User

	err := users.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, err
	}
	cancel()

	filter = bson.M{"_userid": user.ID}

	locations := ls.DB.Collection("locations")
	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	cur, err := locations.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	i := make([]goloc.Location, 0)

	for cur.Next(ctx) {
		var mloc Location
		err = cur.Decode(&mloc)
		if err != nil {
			return i, nil
		}
		loc := goloc.Location{ID: mloc.ID,
			UserID:  mloc.UserID,
			Date:    mloc.Date,
			Address: mloc.Address,
			Lat:     mloc.Lat,
			Lng:     mloc.Lng}

		i = append(i, loc)
	}

	return i, nil
}

//Location -
func (ls *Locator) Location(interface{}) (*goloc.Location, error) {
	return nil, nil
}

// DeleteLocation -
func (ls *Locator) DeleteLocation(id interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	//	oid := fmt.Sprintf(`ObjectId("%s")`, id)

	oid, err := primitive.ObjectIDFromHex(id.(string))
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	locations := ls.DB.Collection("locations")

	result, err := locations.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	fmt.Println(result.DeletedCount)

	return nil
}

// CreateLocation Create a Location record for user with email
func (ls *Locator) CreateLocation(email string, loc goloc.Location) (interface{}, error) {

	users := ls.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	filter := bson.M{"email": email}

	var user User

	err := users.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		return nil, err
	}
	cancel()

	// map goloc.Location to mongoloc.Location
	mgloc := Location{UserID: user.ID,
		Date:    loc.Date,
		Address: loc.Address,
		Lat:     loc.Lat,
		Lng:     loc.Lng}

	ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()

	locations := ls.DB.Collection("locations")

	result, err := locations.InsertOne(ctx, mgloc)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil

}

// GeoLoc - call thru to mapquest api
func (ls *Locator) GeoLoc(Address string) (*goloc.MapAddress, error) {

	ll, err := ls.Mq.LatLng(Address)

	if err != nil {
		return nil, err
	}

	return &goloc.MapAddress{Address: ll.Address, Lat: ll.Lat, Lng: ll.Lng}, nil
}

// CreateUser -
func (ls *Locator) CreateUser(email string) (interface{}, error) {

	user := User{Email: email}

	//	collection := ls.Client.Database("loc_db").Collection("users")
	collection := ls.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

// CreateUsersTable -
func (ls *Locator) CreateUsersTable() error {

	//ls.Client.Database("loc_db").Collection("users")

	ls.DB.Collection("users")

	return nil
}

// CreateLocationsTable creates if doesn't exist
func (ls *Locator) CreateLocationsTable() error {
	//ls.Client.Database("loc_db").Collection("locations")

	ls.DB.Collection("locations")

	return nil
}
