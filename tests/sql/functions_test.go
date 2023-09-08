package sqltest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/dbconnections"
	"github.com/glodb/dbfusion/tests/models"
)

func TestSQLCreate(t *testing.T) {
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
			DbType: dbconnections.MYSQL,
			Cache:  &cache,
		}
	con, err := dbfusion.GetInstance().GetConnection(options)

	if err != nil {
		t.Errorf("DBConnection failed with %v", err)
	}

	users := models.UserTest{
		FirstName: "Aafaq",
		Email:     "aafaqzahid9@gmail.com",
		Password:  "change-me",
	}
	nonEntityUser := models.NonEntityUserTest{
		FirstName: "Gul",
		Email:     "gulandaman@gmail.com",
		Username:  "gulandaman",
		Password:  "change-me",
	}

	testCases := []struct {
		Con            dbconnections.DBConnections
		Data           interface{}
		ExpectedResult error
		Name           string
	}{
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Create with Entity Name, Cache, Pre and Post Insert hooks",
		},
		{
			Con:            con,
			Data:           nonEntityUser,
			ExpectedResult: nil,
			Name:           "Create without Entity Name",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := con.Insert(tc.Data)

			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
