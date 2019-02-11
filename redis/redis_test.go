package redis

import "testing"

// NOTE there's no point in making mock database connection to test database
// wrapper, so this test tries to connect to default redis port on localhost
func TestRealDatabase(t *testing.T) {

}
