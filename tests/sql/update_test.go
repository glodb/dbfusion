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

type UpdateTestResults struct {
	data interface{}
	err  error
}

func TestSQLUpdate(t *testing.T) {
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
		TestData       TestData
		Conditions     []int
		Result         interface{}
		ExpectedResult UpdateTestResults
		Upsert         bool
		Name           string
	}{
		{
			Con:        con,
			Data:       users,
			Conditions: []int{WHERE},
			TestData:   TestData{whereConditions: ftypes.DMap{{Key: "firstname =", Value: "Zahid3"}}},
			Result:     &models.UserTest{},
			Upsert:     true,
			ExpectedResult: UpdateTestResults{
				data: models.UserTest{
					FirstName: "Zahid",
					Email:     "aafaqzahid9@gmail.com",
					Username:  "",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing update with upsert",
		},
		{
			Con:        con,
			Data:       users,
			Conditions: []int{WHERE},
			TestData:   TestData{whereConditions: ftypes.DMap{{Key: "firstname =", Value: "Zahid"}}},
			Result:     &models.UserTest{},
			Upsert:     true,
			ExpectedResult: UpdateTestResults{
				data: models.UserTest{
					FirstName: "Zahid",
					Email:     "aafaqzahid9@gmail.com",
					Username:  "",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing without upsert",
		},
		{
			Con:        con,
			Data:       users,
			Conditions: []int{WHERE},
			TestData:   TestData{whereConditions: ftypes.DMap{{Key: "firstname =", Value: "Aafaq"}}},
			Result:     &models.UserTest{},
			Upsert:     true,
			ExpectedResult: UpdateTestResults{
				data: models.UserTest{
					FirstName: "Zahid",
					Email:     "aafaqzahid9@gmail.com",
					Username:  "",
					Password:  "change-me",
					CreatedAt: 0,
					UpdatedAt: 0,
				},
				err: nil,
			},
			Name: "Testing without upsert with existing user",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			groupCounter := 0
			sortCounter := 0
			joinCounter := 0
			for _, condition := range tc.Conditions {
				switch condition {
				case WHERE:
					con.Where(tc.TestData.whereConditions)
				case PROJECT:
					con.Select(tc.TestData.projections)
				case ADDTABLE:
					con.Table(tc.TestData.tableName)
				case GROUPBY:
					con.GroupBy(tc.TestData.groupByFields[groupCounter])
					groupCounter++
				case HAVING:
					con.Having(tc.TestData.havingConditions)
				case SORT:
					con.Sort(tc.TestData.sortValues[sortCounter].key, tc.TestData.sortValues[sortCounter].sortdesc)
					sortCounter++
				case SKIP:
					con.Skip(int64(tc.TestData.skipValues))
				case LIMIT:
					con.Limit(int64(tc.TestData.limitValues))
				case JOIN:
					con.Join(tc.TestData.joinValues[joinCounter])
					joinCounter++
				}
			}
			err := con.UpdateAndFindOne(tc.Data, tc.Result, tc.Upsert)

			// log.Println("Result:", tc.Data)

			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
				return
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
