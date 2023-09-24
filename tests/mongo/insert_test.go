package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"

	"github.com/glodb/dbfusion/tests/models"
)

func TestMongoCreate(t *testing.T) {
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

	userWithAddress := models.UseWithAddress{
		FirstName: "Aafaq",
		Email:     "aafaqzahid9@gmail.com",
		Username:  "aafaqzahid",
		Address:   models.Address{City: "Lahore", PostalCode: "54000", Line1: "DHA Phase 6"},
		Password:  "change-me",
	}

	mapUser := map[string]interface{}{
		"firstName": "Gul",
		"email":     "gulandaman@gmail.com",
		"userName":  "gulandaman",
		"password":  "change-me",
	}
	testCases := []struct {
		Con            connections.MongoConnection
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
			Name:           "Create with Entity Name, Cache, Pre and Post Insert hooks",
		},
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Create with Entity Name",
		},
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Create with Entity Name, Pre and Post Insert hooks",
		},
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Create with Entity Name, Post Insert hooks",
		},
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Adding more tests to match delete",
		},
		{
			Con:            con,
			Data:           nonEntityUser,
			ExpectedResult: nil,
			Name:           "Create without Entity Name",
		},
		{
			Con:            con,
			Data:           mapUser,
			Type:           1,
			TableName:      "nonStructUser",
			ExpectedResult: nil,
			Name:           "Insert data from map",
		},
		{
			Con:            con,
			Data:           userWithAddress,
			Type:           1,
			TableName:      "users",
			ExpectedResult: nil,
			Name:           "User With Address",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Type == 0 {
				err = con.InsertOne(tc.Data)
			} else {
				err = con.Table(tc.TableName).InsertOne(tc.Data)
			}
			// conNew, _ := dbfusion.GetInstance().GetConnection(options)
			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
