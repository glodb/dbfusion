package redistest

import (
	"testing"

	"github.com/glodb/dbfusion/caches"
)

func TestRedisConnections(t *testing.T) {
	// Create a test server with the handler

	testCases := []struct {
		ConnectionURI string
		Name          string
		Password      string
		ErrorExpected bool
	}{
		{
			ConnectionURI: "localhost:6379",
			Name:          "Redis Valid Connection",
			Password:      "",
			ErrorExpected: false,
		},
		{
			ConnectionURI: "localhost/63791",
			Name:          "Failed Redis Valid Connection",
			Password:      "",
			ErrorExpected: true,
		},
		{
			ConnectionURI: "localhost:6379",
			Name:          "Redis Valid Connection with Password",
			Password:      "redis",
			ErrorExpected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			var cache caches.Cache
			cache = &caches.RedisCache{}
			if tc.Password != "" {
				err = cache.ConnectCache(tc.ConnectionURI, tc.Password)
			} else {
				err = cache.ConnectCache(tc.ConnectionURI)
			}

			if err != nil {
				if !tc.ErrorExpected {
					t.Errorf("Error occurred %v", err)
				} else {
					t.Logf("Expected Error occurred %v", err)
				}
			} else {
				cache.DisconnectCache()
			}
		})
	}
}
