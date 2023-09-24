package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"

	"github.com/glodb/dbfusion/tests/models"
)

func TestMongoCreateIndexes(t *testing.T) {
	validDBName := "testDBFusion"
	validUri := "mongodb://localhost:27017"
	cache := caches.RedisCache{}
	err := cache.ConnectCache("localhost:6379")
	cache.FlushAll()
	if err != nil {
		t.Errorf("Error in redis connection, occurred %v", err)
	}
	options :=
		dbfusion.Options{
			DbName: &validDBName,
			Uri:    &validUri,
			Cache:  &cache,
		}
	con, err := dbfusion.GetInstance().GeMongoConnection(options)

	if err != nil {
		t.Errorf("DBConnection failed with %v", err)
	}
	testCases := []struct {
		Con            connections.MongoConnection
		Data           interface{}
		ExpectedResult error
		IfExists       bool
		TableName      string
		Name           string
	}{
		{
			Con:            con,
			Data:           models.UserCreateTable{},
			ExpectedResult: nil,
			IfExists:       true,
			Name:           "Create checking if exists",
		},
		{
			Con:            con,
			Data:           models.UserCreateTable{},
			ExpectedResult: nil,
			IfExists:       false,
			Name:           "Create without checking if exists",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			con.CreateIndexes(tc.Data)
			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
