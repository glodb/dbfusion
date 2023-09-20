package mongotest

import (
	"log"
	"reflect"
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/tests/models"
)

type FindTestResults struct {
	data interface{}
	err  error
}

func TestMongoFind(t *testing.T) {
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
	con, err := dbfusion.GetInstance().GeMongoConnection(options)

	if err != nil {
		t.Errorf("DBConnection failed with %v", err)
	}

	// textQuery := query.QMap{"$text": query.QMap{"$search": "some search"}}

	testCases := []struct {
		Con            connections.MongoConnection
		Query          interface{}
		Data           interface{}
		Options        queryoptions.FindOptions
		Type           int
		TableName      string
		ExpectedResult FindTestResults
		Name           string
	}{
		{
			Con:     con,
			Query:   ftypes.QMap{"email": "aafaqzahid9@gmail.com"},
			Data:    &models.UserTest{},
			Options: queryoptions.FindOptions{ForceDB: true, CacheResult: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
					Email:     "aafaqzahid9@gmail.com",
					Password:  "0f0bf2567ec111697671d2fd76af0d6c",
					UpdatedAt: 0,
					CreatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing force db query wth user hook",
		},
		{
			Con:     con,
			Query:   ftypes.QMap{"email": "aafaqzahid9@gmail.com"},
			Data:    &models.UserTest{},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
					Email:     "aafaqzahid9@gmail.com",
					Password:  "0f0bf2567ec111697671d2fd76af0d6c",
					UpdatedAt: 0,
					CreatedAt: 0,
				},
				err: nil,
			},
			Name: "Reading Cache",
		},
		{
			Con:     con,
			Query:   ftypes.QMap{"email": "gulandaman@gmail.com", "password": "change-me"},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					UpdatedAt: 0,
					CreatedAt: 0,
				},
				err: nil,
			},
			Name: "Reading Default Cache",
		},
		{
			Con:     con,
			Query:   ftypes.QMap{"firstname": "Gul"},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing force db query with non hook",
		},
		{
			Con:       con,
			Type:      1,
			TableName: "users",
			Query:     ftypes.QMap{"email": "aafaqzahid9@gmail.com"},
			Data:      &map[string]interface{}{},
			Options:   queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &map[string]interface{}{
					"createdAt": int64(0),
					"email":     "aafaqzahid9@gmail.com",
					"firstname": "Aafaq",
					"password":  "0f0bf2567ec111697671d2fd76af0d6c",
					"updatedAt": int64(0),
					"username":  "",
				},
				err: nil,
			},
			Name: "Testing force db query with map",
		},
		// {
		// 	Con:     con.Table("users"),
		// 	Query:   textQuery,
		// 	Data:    &map[string]interface{}{},
		// 	Options: queryoptions.FindOptions{ForceDB: true},
		// 	ExpectedResult: FindTestResults{
		// 		data: &map[string]interface{}{},
		// 		err:  err.(mongo.CommandError),
		// 	},
		// 	Name: "Testing force db text query with map",
		// },
		{
			Con:       con,
			Type:      1,
			TableName: "users",
			Query:     ftypes.QMap{"email": "aafaqzahid9@gmail.com"},
			Data:      &map[string]interface{}{},
			Options:   queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &map[string]interface{}{
					"createdAt": int64(0),
					"email":     "aafaqzahid9@gmail.com",
					"firstname": "Aafaq",
					"password":  "0f0bf2567ec111697671d2fd76af0d6c",
					"updatedAt": int64(0),
					"username":  "",
				},
			},
			Name: "Testing force db query with qmap",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{{Key: "email", Value: "gulandaman@gmail.com"}, {Key: "password", Value: "change-me"}},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing force db query with non hook and forceDB to false",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{{Key: "email", Value: "gulandaman@gmail.com"}, {Key: "firstname", Value: "Gul"}},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false, CacheResult: true},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing force db query with non hook and forceDB and save query result to cache",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{{Key: "email", Value: "gulandaman@gmail.com"}, {Key: "firstname", Value: "Gul"}},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false, CacheResult: true},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Reading query result from cache as saved previously",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{{Key: "username", Value: "gulandaman"}},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false, CacheResult: false},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "DB query when nothing found in cache",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{},
			Data:    &models.NonEntityUserTest{},
			Options: queryoptions.FindOptions{ForceDB: false, CacheResult: false},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Empty conditions test",
		},
		{
			Con:     con,
			Query:   ftypes.DMap{},
			Data:    &models.NonEntityUserTest{},
			Type:    2,
			Options: queryoptions.FindOptions{ForceDB: false, CacheResult: false},
			ExpectedResult: FindTestResults{
				data: &models.NonEntityUserTest{
					FirstName: "Gul",
					Email:     "gulandaman@gmail.com",
					Username:  "gulandaman",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Empty conditions test",
		},
	}

	// { $text: { $search: "search query" } } TODO: Test this query also

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var err error
			if tc.Type == 0 {
				err = con.Where(tc.Query).FindOne(tc.Data, tc.Options)
			} else if tc.Type == 1 {
				err = con.Table(tc.TableName).Where(tc.Query).FindOne(tc.Data, tc.Options)
			} else if tc.Type == 2 {
				err = con.Where(tc.Query).Sort("password").FindOne(tc.Data, tc.Options)
			}
			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
			}

			// // Use type assertions to check the type of the result
			// switch v := tc.Data.(type) {
			// case map[string]interface{}:
			// 	// Check if it's a map[string]interface{}
			// 	// Add validation logic for maps if needed
			// default:

			if value, ok := interface{}(tc.Data).(*map[string]interface{}); ok {
				expectedMap := tc.ExpectedResult.data.(*map[string]interface{})
				for key, value1 := range *expectedMap {
					value2, exists := (*value)[key]

					valueInterface1 := reflect.ValueOf(value1)
					valueInterface2 := reflect.ValueOf(value2)

					if !exists || !reflect.DeepEqual(valueInterface1.Interface(), valueInterface2.Interface()) {
						log.Println(key, value1, value2)
						t.Errorf("Expected userObject %+v, but got %+v ", tc.ExpectedResult.data, tc.Data)
					}
				}
			} else if !reflect.DeepEqual(tc.Data, tc.ExpectedResult.data) {
				t.Errorf("Expected userObject %+v, but got %+v", tc.ExpectedResult.data, tc.Data)
			}
			// }
		})
	}
}
