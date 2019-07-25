package mocklocator

import (
	"github.com/jbowl/goloc/internal/pkg/goloc"
)

// MockLocator -
type MockLocator struct {
}

// Initialize -
func (api *MockLocator) Initialize() error {
	return nil
}

// OpenDatabase -
func (api *MockLocator) OpenDatabase() error {
	return nil
}

// Close -
func (api *MockLocator) Close() error {
	return nil
}

// CreateUser -
func (api *MockLocator) CreateUser(email string) (interface{}, error) {
	return -1, nil
}

// CreateUsersTable -
func (api *MockLocator) CreateUsersTable() error {
	return nil
}

// CreateLocationsTable -
func (api *MockLocator) CreateLocationsTable() error {
	return nil
}

//Locations -
func (api *MockLocator) Locations(string) ([]goloc.Location, error) { // READ GET all
	return make([]goloc.Location, 0), nil
}

//Location -
func (api *MockLocator) Location(interface{}) (*goloc.Location, error) {

    


	return nil, nil
}

//CreateLocation -
func (api *MockLocator) CreateLocation(email string, loc goloc.Location) (interface{}, error) {
	return -1, nil
}

//DeleteLocation -
func (api *MockLocator) DeleteLocation() {}

//GeoLoc -
func (api *MockLocator) GeoLoc(Address string) (*goloc.MapAddress, error) {
	return nil, nil
}
