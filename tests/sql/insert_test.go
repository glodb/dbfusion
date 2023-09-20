package sqltest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
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
			Cache:  &cache,
		}
	con, err := dbfusion.GetInstance().GetMySqlConnection(options)
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

	mapUser := map[string]interface{}{
		"firstName": "Gul",
		"email":     "gulandaman@gmail.com",
		"userName":  "gulandaman",
		"password":  "change-me",
	}
	// userWithAddress := models.UseWithAddress{
	// 	FirstName: "Aafaq",
	// 	Email:     "aafaqzahid9@gmail.com",
	// 	Username:  "aafaqzahid",
	// 	Address:   models.Address{City: "Lahore", PostalCode: "54000", Line1: "DHA Phase 6"},
	// 	Password:  "change-me",
	// }
	testCases := []struct {
		Con            connections.SQLConnection
		Data           interface{}
		ExpectedResult error
		Type           int
		TableName      string
		Name           string
	}{
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Type:           0,
			Name:           "Create with Entity Name, Cache, Pre and Post Insert hooks",
		},
		{
			Con:            con,
			Data:           nonEntityUser,
			ExpectedResult: nil,
			Type:           0,
			Name:           "Create without Entity Name",
		},
		{
			Con:            con,
			TableName:      "users",
			Data:           mapUser,
			ExpectedResult: nil,
			Type:           1,
			Name:           "Insert data from map",
		},
		{
			Con:            con,
			Data:           mapUser,
			ExpectedResult: dbfusionErrors.ErrEntityNameRequired,
			Type:           0,
			Name:           "Insert data from map withour table",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if tc.Type == 0 {
				err = con.InsertOne(tc.Data)
			} else {
				err = con.Table(tc.TableName).InsertOne(tc.Data)
			}
			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
