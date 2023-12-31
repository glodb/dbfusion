package sqltest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/tests/models"
)

type DeleteTestResults struct {
	data interface{}
	err  error
}

func TestSQLDelete(t *testing.T) {
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

	users := models.UserTest{
		FirstName: "Zahid",
		Email:     "aafaqzahid9@gmail.com",
	}
	testCases := []struct {
		Con            connections.SQLConnection
		Data           interface{}
		ExpectedResult DeleteTestResults
		Type           int
		Name           string
	}{
		{
			Con:  con,
			Data: &users,
			ExpectedResult: DeleteTestResults{
				data: nil,
				err:  nil,
			},
			Type: 2,
			Name: "Testing with user object",
		},
		{
			Con:  con,
			Data: ftypes.DMap{{Key: "firstname = ", Value: "Aafaq"}},
			ExpectedResult: DeleteTestResults{
				data: nil,
				err:  nil,
			},
			Type: 1,
			Name: "Testing with where",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			var err error
			if tc.Type == 1 {
				err = con.Table("users").Where(tc.Data).DeleteOne()
			} else if tc.Type == 2 {
				err = con.DeleteOne(tc.Data)
			}
			// conNew, _ := dbfusion.GetInstance().GetConnection(options)
			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}

		})
	}
}
