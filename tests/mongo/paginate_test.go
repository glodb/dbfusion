package mongotest

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

func TestMongoPaginate(t *testing.T) {
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

	// textQuery := query.QMap{"$text": query.QMap{"$search": "some search"}}
	testCases := []struct {
		Con            connections.MongoConnection
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
					{FirstName: "Aafaq", Email: "aafaqzahid9@gmail.com", Username: "", Password: "0f0bf2567ec111697671d2fd76af0d6c"},
					{FirstName: "Aafaq", Email: "aafaqzahid9@gmail.com", Username: "", Password: "0f0bf2567ec111697671d2fd76af0d6c"},
				},
				err: nil,
			},
			Name: "Test simple select all",
		},
		{
			Con:        con,
			Data:       &[]models.UserTest{},
			Conditions: []int{ADDTABLE, SETLIMIT},
			PageNumber: 1,
			TestData: TestData{
				tableName:  "users",
				pageValues: 1,
			},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: PaginateTestResults{
				data: &[]models.UserTest{{FirstName: "Aafaq", Email: "aafaqzahid9@gmail.com", Username: "", Password: "0f0bf2567ec111697671d2fd76af0d6c"}},
				err:  nil,
			},
			Name: "Testing select with projection",
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
				case SETLIMIT:
					con.SetPageSize(tc.TestData.pageValues)
				}
			}
			paginationResults, err := con.Paginate(tc.Data, tc.PageNumber)

			log.Println("Pagination Result:", paginationResults)

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
