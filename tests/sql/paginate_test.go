package sqltest

import (
	"log"
	"reflect"
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/tests/models"
)

type PaginateTestResults struct {
	data interface{}
	err  error
}

func TestSQLPaginate(t *testing.T) {
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
		Conditions     []int
		Options        queryoptions.FindOptions
		ExpectedResult PaginateTestResults
		TestData       TestData
		PageNumber     int
		Type           int
		Name           string
	}{
		{
			Con:        con,
			Data:       &[]models.UserTest{},
			Conditions: []int{ADDTABLE, SETLIMIT},
			TestData:   TestData{tableName: "users", pageValues: 2},
			PageNumber: 2,
			Options:    queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: PaginateTestResults{
				data: &[]models.UserTest{
					{FirstName: "Gul", Email: "gulandaman@gmail.com", Username: "gulandaman", Password: "change-me"},
					{FirstName: "Gul", Email: "gulandaman@gmail.com", Username: "gulandaman", Password: "change-me"},
				},
				err: nil,
			},
			Name: "Test simple select all",
		},
		{
			Con:        con,
			Data:       &[]models.UserTest{},
			Conditions: []int{ADDTABLE, PROJECT, SETLIMIT},
			PageNumber: 1,
			TestData: TestData{
				projections: map[string]bool{"firstname": true},
				tableName:   "users",
				pageValues:  1,
			},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: PaginateTestResults{
				data: &[]models.UserTest{{FirstName: "Zahid1"}},
				err:  nil,
			},
			Name: "Testing select with projection",
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
				case SETLIMIT:
					con.SetPageSize(tc.TestData.pageValues)
				}
			}
			results, err := con.Paginate(tc.Data, tc.PageNumber)

			log.Println("Pagination Results:", results)

			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
				return
			}

			// // Use type assertions to check the type of the result
			// switch v := tc.Data.(type) {
			// case map[string]interface{}:
			// 	// Check if it's a map[string]interface{}
			// 	// Add validation logic for maps if needed
			// default:

			if value, ok := interface{}(tc.Data).(*map[string]interface{}); ok {
				log.Println(value)
				// expectedMap := tc.ExpectedResult.data.(*map[string]interface{})
				// for key, value1 := range *expectedMap {
				// 	value2, exists := (*value)[key]

				// 	valueInterface1 := reflect.ValueOf(value1)
				// 	valueInterface2 := reflect.ValueOf(value2)

				// 	if !exists || !reflect.DeepEqual(valueInterface1.Interface(), valueInterface2.Interface()) {
				// 		// log.Println(key, value1, value2)
				// 		t.Errorf("Expected userObject %+v, but got %+v ", tc.ExpectedResult.data, tc.Data)
				// 	}
				// }
			} else if !reflect.DeepEqual(tc.Data, tc.ExpectedResult.data) {
				t.Errorf("Expected userObject %+v, but got %+v", tc.ExpectedResult.data, tc.Data)
			}
			// }
		})
	}
}
