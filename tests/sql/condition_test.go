package sqltest

import (
	"testing"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/implementations"
)

func compareInterfaceArrays(value []interface{}, expected []interface{}) bool {
	if len(value) != len(expected) {
		return false
	}
	for id, val := range value {
		if expected[id] != val {
			return false
		}
	}
	return true
}

func TestMongoConditions(t *testing.T) {
	con := &implementations.MySql{}

	testCases := []struct {
		Conditions     conditions.ConditionBuilder
		ExpectedResult conditions.SqlData
		Name           string
	}{
		{
			Conditions: con.Add(
				con.And(conditions.SqlCondition{Key: "users.id", ConditionType: conditions.EQUAL, Value: 20},
					conditions.SqlCondition{Key: "age", ConditionType: conditions.GREATER_THAN, Value: 50},
				),
				con.Or(conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"}),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "users.id = ? AND age > ? OR name <> ?",
				Values: []interface{}{20, 50, "aafaq"}},
			Name: "Testing multiple And, single OR ",
		},
		{
			Conditions: con.Add(
				con.Or(conditions.SqlCondition{Key: "users.id", ConditionType: conditions.EQUAL, Value: 20},
					conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"},
				),
				con.And(conditions.SqlCondition{Key: "age", ConditionType: conditions.GREATER_THAN, Value: 50}),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "users.id = ? OR name <> ? AND age > ?",
				Values: []interface{}{20, "aafaq", 50}},
			Name: "Testing Multiple OR, Single And",
		},
		{
			Conditions: con.Add(
				con.Or(conditions.SqlCondition{Key: "users.id", ConditionType: conditions.EQUAL, Value: 20},
					conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "users.id = ? OR name <> ?",
				Values: []interface{}{20, "aafaq"}},
			Name: "Testing OR only",
		},
		{
			Conditions: con.Add(
				con.And(conditions.SqlCondition{Key: "users.id", ConditionType: conditions.EQUAL, Value: 20},
					conditions.SqlCondition{Key: "age", ConditionType: conditions.GREATER_THAN, Value: 50},
				),
				con.Or(conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"}),
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "score", ConditionType: conditions.GREATER_THAN, Value: 60},
						con.Or(conditions.SqlCondition{Key: "name", ConditionType: conditions.LIKE, Value: "%aa"}),
					),
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "users.id = ? AND age > ? OR name <> ? AND (score > ? OR name LIKE ?)",
				Values: []interface{}{20, 50, "aafaq", 60, "%aa"}},
			Name: "Testing Group With Like",
		},
		{
			Conditions: con.Add(
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "users.id", ConditionType: conditions.NOT_EQUAL, Value: 60},
						conditions.SqlCondition{Key: "age", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 60},
					),
				),
				con.Or(
					conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"},
				),
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "score", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 60},
						con.Or(conditions.SqlCondition{Key: "name", ConditionType: conditions.LIKE, Value: "%aa"}),
					),
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "(users.id <> ? AND age <= ?) OR name <> ? AND (score <= ? OR name LIKE ?)",
				Values: []interface{}{60, 60, "aafaq", 60, "%aa"}},
			Name: "Testing top level brackets",
		},
		{
			Conditions: con.Add(
				con.GroupConditions(
					con.Or(
						conditions.SqlCondition{Key: "age", ConditionType: conditions.EQUAL, Value: 50},
						conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_EQUAL, Value: "aafaq"},
					),
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "(age = ? OR name <> ?)",
				Values: []interface{}{50, "aafaq"}},
			Name: "Testing only brackets ",
		},
		{
			Conditions: con.Add(
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "users.id", ConditionType: conditions.NOT_EQUAL, Value: 50},
						conditions.SqlCondition{Key: "age", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 50},
					),
				),
				con.And(
					conditions.SqlCondition{Key: "score", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 50},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "(users.id <> ? AND age <= ?) AND score <= ?",
				Values: []interface{}{50, 50, 50}},
			Name: "Testing brackets with AND ",
		},
		{
			Conditions: con.Add(
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "users.id", ConditionType: conditions.NOT_EQUAL, Value: 50},
						conditions.SqlCondition{Key: "age", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 50},
					),
				),
				con.Or(
					conditions.SqlCondition{Key: "score", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 50},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "(users.id <> ? AND age <= ?) OR score <= ?",
				Values: []interface{}{50, 50, 50}},
			Name: "Testing brackets with OR ",
		},
		{
			Conditions: con.Add(
				con.GroupConditions(
					con.And(
						conditions.SqlCondition{Key: "users.id", ConditionType: conditions.NOT_EQUAL, Value: 50},
						conditions.SqlCondition{Key: "name", ConditionType: conditions.NOT_LIKE, Value: "%aa"},
					),
				),
				con.Or(
					conditions.SqlCondition{Key: "score", ConditionType: conditions.LESSER_THAN_EQUAL, Value: 50},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "(users.id <> ? AND name NOT LIKE ?) OR score <= ?",
				Values: []interface{}{50, "%aa", 50}},
			Name: "Testing brackets with OR ",
		},
		{
			Conditions: con.Add(
				con.And(
					conditions.SqlCondition{Key: "age", ConditionType: conditions.IN, Value: []interface{}{50, 30, 60}},
					conditions.SqlCondition{Key: "email", ConditionType: conditions.NOT_IN, Value: []interface{}{10, 20}},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "age IN (?, ?, ?) AND email NOT IN (?, ?)",
				Values: []interface{}{50, 30, 60, 10, 20}},
			Name: "Testing brackets with OR ",
		},
		{
			Conditions: con.Add(
				con.Or(conditions.SqlCondition{Key: "email", ConditionType: conditions.EQUAL, Value: "aafaqzahid9@gmail.com"},
					conditions.SqlCondition{Key: "name", ConditionType: conditions.EQUAL, Value: "aafaq"},
				),
			),
			ExpectedResult: conditions.SqlData{
				Query:  "email = ? OR name = ?",
				Values: []interface{}{"aafaqzahid9@gmail.com", "aafaq"}},
			Name: "Testing Cache Query",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			data := con.GetQuery(tc.Conditions).(conditions.SqlData)
			if data.Query != tc.ExpectedResult.Query {
				t.Errorf("Expected Query %s, But Got %s", tc.ExpectedResult.Query, data.Query)
			}

			if !compareInterfaceArrays(data.Values, tc.ExpectedResult.Values) {
				t.Errorf("Data arrays are not same expected %v but got %v", tc.ExpectedResult.Values, data.Values)
			}

			// log.Println(data.QueryDefaultCache, data.CacheValues)
		})
	}
}
