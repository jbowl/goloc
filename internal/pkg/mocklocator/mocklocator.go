package mocklocator

import (
	"database/sql"

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
func (api *MockLocator) OpenDatabase() (*sql.DB, error) {
	return nil, nil
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
func (api *MockLocator) Location(int) (*goloc.Location, error) {
	return nil, nil
}

//CreateLocation -
func (api *MockLocator) CreateLocation(email string, loc goloc.Location) (int64, error) {
	return -1, nil
}

//DeleteLocation -
func (api *MockLocator) DeleteLocation() {}

//GeoLoc -
func (api *MockLocator) GeoLoc(Address string) (*goloc.MapAddress, error) {
	return nil, nil
}
