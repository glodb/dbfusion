package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/tests/models"
)

const (
	WHERE = iota
	ADDTABLE
	SETLIMIT
)

type TestData struct {
	whereConditions interface{}
	tableName       string
	pageValues      int
}

func TestMongoUpdate(t *testing.T) {
	validDBName := "testDBFusion"
	validUri := "mongodb://localhost:27017"
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
	con, err := dbfusion.GetInstance().GetMongoConnection(options)

	if err != nil {
		t.Errorf("DBConnection failed with %v", err)
	}

	users := models.UserTest{
		FirstName: "Aafaq",
		Email:     "aafaqzahid9@gmail.com",
		Password:  "change-me",
	}
	result := models.UserTest{}
	// result := primitive.M{}

	testCases := []struct {
		Con            connections.MongoConnection
		Data           interface{}
		ExpectedResult error
		Result         interface{}
		TestData       TestData
		Conditions     []int
		Name           string
	}{
		{
			Con:            con,
			Data:           users,
			ExpectedResult: nil,
			Name:           "Update User with object",
			Conditions:     []int{WHERE},
			TestData:       TestData{whereConditions: ftypes.DMap{{"email", "aafaqzahid9@gmail.com"}}},
			Result:         &result,
		},
		{
			Con:            con,
			Data:           ftypes.QMap{"email": "aafaqzahid9+1@gmail.com"},
			ExpectedResult: nil,
			Name:           "Update User map",
			Conditions:     []int{ADDTABLE, WHERE},
			TestData:       TestData{tableName: "users", whereConditions: ftypes.DMap{{"email", "aafaqzahid9@gmail.com"}}},
			Result:         &ftypes.QMap{},
		},
		{
			Con:            con,
			Data:           ftypes.DMap{{"email", "aafaqzahid9+2@gmail.com"}},
			ExpectedResult: nil,
			Name:           "Update User dmap",
			Conditions:     []int{ADDTABLE, WHERE},
			TestData:       TestData{tableName: "users", whereConditions: ftypes.DMap{{"email", "aafaqzahid9@gmail.com"}}},
			Result:         &models.UserTest{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			for _, condition := range tc.Conditions {
				switch condition {
				case WHERE:
					con.Where(tc.TestData.whereConditions)
				case ADDTABLE:
					con.Table(tc.TestData.tableName)
				}
			}
			con.UpdateAndFindOne(tc.Data, tc.Result, false)
			// conNew, _ := dbfusion.GetInstance().GetConnection(options)
			if err != tc.ExpectedResult {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}
		})
	}
}
