package sqltest

import (
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
		FirstName: "Aafaq",
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
			TestData:   TestData{whereConditions: ftypes.DMap{{Key: "firstname =", Value: "Zahid"}}},
			Result:     &models.UserTest{},
			Upsert:     true,
			ExpectedResult: UpdateTestResults{
				data: nil,
				err:  nil,
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
				data: nil,
				err:  nil,
			},
			Name: "Testing without upsert",
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
					con.Project(tc.TestData.projections)
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

		})
	}
}
