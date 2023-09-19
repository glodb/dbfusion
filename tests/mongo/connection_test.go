package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/implementations"
)

const (
	REDIS      = 1
	MEM_CACHED = 2
)

func TestMongoConnections(t *testing.T) {
	// Create a test server with the handler

	validDBName := "testDBFusion"
	validUri := "mongodb://localhost:27017"

	testCases := []struct {
		Option         dbfusion.Options
		ConnectCache   bool
		CacheType      int
		ExpectedResult struct {
			Connection connections.Connection
			Error      error
		}
		Name string
	}{
		{
			Option: dbfusion.Options{
				DbName: &validDBName,
			},
			Name: "InvalidURI",
			ExpectedResult: struct {
				Connection connections.Connection
				Error      error
			}{
				Connection: nil,
				Error:      dbfusionErrors.ErrUriRequiredForConnection, // You can set the expected error value here.
			},
		},
		{
			Option: dbfusion.Options{
				DbName: &validDBName,
				Uri:    &validUri,
			},
			Name: "Valid Connection",
			ExpectedResult: struct {
				Connection connections.Connection
				Error      error
			}{
				Connection: &implementations.MongoConnection{},
				Error:      nil, // You can set the expected error value here.
			},
		},
		{
			Option: dbfusion.Options{
				DbName: &validDBName,
				Uri:    &validUri,
			},
			ConnectCache: true,
			CacheType:    REDIS,
			Name:         "Valid Connection with Redis Cache",
			ExpectedResult: struct {
				Connection connections.Connection
				Error      error
			}{
				Connection: &implementations.MongoConnection{},
				Error:      nil, // You can set the expected error value here.
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.ConnectCache {
				if tc.CacheType == REDIS {
					cache := caches.RedisCache{}
					err := cache.ConnectCache("localhost:6379")
					if err != nil {
						t.Errorf("Error in redis connection, occurred %v", err)
					} else {
						tc.Option.Cache = &cache
					}
				}
			}
			con, err := dbfusion.GetInstance().GeMongoConnection(tc.Option)

			if err != tc.ExpectedResult.Error {
				t.Errorf("Expected status code %v, but got %v", tc.ExpectedResult.Error, err)
			}
			if tc.ExpectedResult.Connection == nil && con != nil {
				t.Error("Expected no connection, but got connection")
			}
		})
	}
}
