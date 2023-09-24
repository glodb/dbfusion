package sqltest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/tests/models"
)

func TestSQLCreateTable(t *testing.T) {
	validDBName := "dbfusion"
	validUri := "root:change-me@tcp(localhost:3306)/dbfusion"
	cache := caches.RedisCache{}
	err := cache.ConnectCache("localhost:6379")
	if err != nil {
		t.Errorf("Error in redis connection, occurred %v", err)
	}
	options :=
		dbfusion.Options{
			DbName: &validDBName,
			Uri:    &validUri,
			Cache:  &cache,
		}
	con, err := dbfusion.GetInstance().GetMySqlConnection(options)
	if err != nil {
		t.Errorf("DBConnection failed with %v", err)
	}

	testCases := []struct {
		Con            connections.SQLConnection
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
			con.CreateTable(tc.Data, tc.IfExists)
			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
