package sqltest

import (
	"reflect"
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
		Password:  "change-me",
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
			Data: ftypes.DMap{{Key: "firstname", Value: "Aafaq"}},
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
				err = con.Where(tc.Data).DeleteOne()
			} else if tc.Type == 2 {
				err = con.DeleteOne(tc.Data)
			}
			// conNew, _ := dbfusion.GetInstance().GetConnection(options)
			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}

			if value, ok := interface{}(tc.Data).(*map[string]interface{}); ok {
				expectedMap := tc.ExpectedResult.data.(*map[string]interface{})
				for key, value1 := range *expectedMap {
					value2, exists := (*value)[key]

					valueInterface1 := reflect.ValueOf(value1)
					valueInterface2 := reflect.ValueOf(value2)

					if !exists || !reflect.DeepEqual(valueInterface1.Interface(), valueInterface2.Interface()) {
						// log.Println(key, value1, value2)
						t.Errorf("Expected userObject %+v, but got %+v ", tc.ExpectedResult.data, tc.Data)
					}
				}
			} else if !reflect.DeepEqual(tc.Data, tc.ExpectedResult.data) {
				t.Errorf("Expected userObject %+v, but got %+v", tc.ExpectedResult.data, tc.Data)
			}
		})
	}
}