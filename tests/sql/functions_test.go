package sqltest

import (
	"log"
	"reflect"
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/dbfusionErrors"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/joins"
	"github.com/glodb/dbfusion/queryoptions"
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

type FindTestResults struct {
	data interface{}
	err  error
}

const (
	WHERE = iota
	PROJECT
	ADDTABLE
	GROUPBY
	HAVING
	SORT
	LIMIT
	SKIP
	JOIN
)

type SortVal struct {
	key      string
	sortdesc bool
}

type TestData struct {
	whereConditions  interface{}
	projections      map[string]bool
	tableName        string
	groupByFields    []string
	havingConditions interface{}
	sortValues       []SortVal
	limitValues      int
	skipValues       int
	joinValues       []joins.Join
}

func TestSqlFind(t *testing.T) {
	validDBName := "testDBFusion"
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
		TestData       TestData
		Conditions     []int
		Options        queryoptions.FindOptions
		ExpectedResult FindTestResults
		Name           string
	}{
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{},
			Options:    queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
					Email:     "aafaqzahid9@gmail.com",
					Password:  "0f0bf2567ec111697671d2fd76af0d6c",
					UpdatedAt: 0,
					CreatedAt: 1694159585,
				},
				err: nil,
			},
			Name: "Test simple select all",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT},
			TestData: TestData{
				projections: map[string]bool{"firstname": true},
			},
			Options: queryoptions.FindOptions{ForceDB: false},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing select with projection",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE},
			TestData: TestData{
				projections: map[string]bool{"firstname": true},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"}),
				),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with select with single condition",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE},
			TestData: TestData{
				projections: map[string]bool{"firstname": true},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with select with multiple conditions",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true},
				groupByFields: []string{"firstname"},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with single group By",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY, HAVING},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true, "password": true},
				groupByFields: []string{"firstname", "password"},
				havingConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "password", ConditionType: conditions.EQUAL, Value: "0f0bf2567ec111697671d2fd76af0d6c"})),
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
					Password:  "0f0bf2567ec111697671d2fd76af0d6c",
				},
				err: nil,
			},
			Name: "Testing having with conditions object",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY, HAVING},
			TestData: TestData{
				projections:      map[string]bool{"firstname": true, "password": true},
				groupByFields:    []string{"firstname", "password"},
				havingConditions: ftypes.DMap{{Key: "password = ?", Value: "0f0bf2567ec111697671d2fd76af0d6c"}},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
					Password:  "0f0bf2567ec111697671d2fd76af0d6c",
				},
				err: nil,
			},
			Name: "Testing having with QDMap",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY, SORT},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true},
				groupByFields: []string{"firstname", "email"},
				sortValues:    []SortVal{{key: "firstname", sortdesc: true}},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with multiple group By single sort",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY, SORT, SORT},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true},
				groupByFields: []string{"firstname", "password"},
				sortValues:    []SortVal{{key: "firstname", sortdesc: true}, {key: "password", sortdesc: false}},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with multiple group By multiple sort",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true},
				groupByFields: []string{"firstname", "email"},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with single sort",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, GROUPBY, LIMIT},
			TestData: TestData{
				projections:   map[string]bool{"firstname": true},
				limitValues:   2,
				groupByFields: []string{"firstname", "email"},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"},
						con.And(conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_EQUAL, Value: "gulandaman@gmail.com"}))),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with limit as 2",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, GROUPBY, LIMIT, SKIP},
			TestData: TestData{
				projections:   map[string]bool{"createdAt": true},
				limitValues:   2,
				skipValues:    1,
				groupByFields: []string{"createdAt"},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"})),
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{},
				err:  nil,
			},
			Name: "Testing with limit and skip",
		},
		{
			Con:        con,
			Conditions: []int{ADDTABLE, WHERE},
			TestData:   TestData{tableName: "users", whereConditions: ftypes.DMap{{Key: "firstname = ?", Value: "Aafaq"}, {Key: " AND email <> ?", Value: "gulandaman@gmail.com"}}},
			Data:       &models.UserTest{},
			Options:    queryoptions.FindOptions{ForceDB: true},
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
			Name: "Testing select with DMap coditions",
		},
		{
			Con:        con,
			Data:       &models.UserTest{},
			Conditions: []int{PROJECT, WHERE, JOIN},
			TestData: TestData{
				projections: map[string]bool{"users.firstname": true},
				whereConditions: con.GetQuery(
					con.Add(conditions.SqlCondition{Key: "b.firstname", ConditionType: conditions.EQUAL, Value: "Aafaq"})),
				joinValues: []joins.Join{{Operator: joins.INNER_JOIN, TableName: "users b", Condition: "users.email=b.email"}},
			},
			Options: queryoptions.FindOptions{ForceDB: true},
			ExpectedResult: FindTestResults{
				data: &models.UserTest{
					FirstName: "Aafaq",
				},
				err: nil,
			},
			Name: "Testing with join",
		},
	}

	// { $text: { $search: "search query" } } TODO: Test this query also

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
			err := con.FindOne(tc.Data, tc.Options)

			// log.Println("Result:", tc.Data)

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
