package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/dbconnections"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoConditions(t *testing.T) {

	condition1 := conditions.MongoCondition{
		ConditionType: conditions.NOT_EQUAL,
		Key:           "firstname",
		Value:         "aafaq",
	}

	condition2 := conditions.MongoCondition{
		ConditionType: conditions.NOT_EQUAL,
		Key:           "email",
		Value:         "aafaqzahid9@gmail.com",
	}

	condition3 := conditions.MongoCondition{
		ConditionType: conditions.NOT_EQUAL,
		Key:           "email",
		Value:         "aafaqzahid9@gmail.com",
	}

	condition4 := conditions.MongoCondition{
		ConditionType: conditions.NOT_EQUAL,
		Key:           "firstname",
		Value:         "aafaq",
	}

	testCases := []struct {
		Conditions     []conditions.MongoCondition
		ExpectedResult primitive.M
		Name           string
		Type           int
	}{
		{
			Conditions:     []conditions.MongoCondition{condition1, condition2},
			ExpectedResult: primitive.M{"firstname": "aafaq", "email": "aafaqzahid9@gmail.com"},
			Name:           "Checking simple query creation",
			Type:           1,
		},
		{
			Conditions:     []conditions.MongoCondition{condition1, condition2},
			ExpectedResult: primitive.M{"firstname": "aafaq", "email": "aafaqzahid9@gmail.com"},
			Name:           "Test Complex inner queries",
			Type:           2,
		},
		{
			Conditions:     []conditions.MongoCondition{condition1, condition2},
			ExpectedResult: primitive.M{"firstname": "aafaq", "email": "aafaqzahid9@gmail.com"},
			Name:           "Test Complex inner queries",
			Type:           3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			switch tc.Type {
			case 1:
				con := &dbconnections.MongoConnection{}
				con.GetQuery(con.Add(condition1, condition2))
				// log.Println(query)
			case 2:
				{
					con := &dbconnections.MongoConnection{}
					con.GetQuery(
						con.Add(
							condition1,
							con.And(
								&condition3,
								con.Or(
									&conditions.MongoCondition{ConditionType: conditions.NOT_EQUAL, Key: "lastname", Value: "zahid"},
									&conditions.MongoCondition{ConditionType: conditions.EQUAL, Key: "lastname", Value: "khan"}),
								&condition4),
							condition3))
					// log.Println(query)
				}
			case 3:
				{
					con := dbconnections.MongoConnection{}
					con.GetQuery(con.Add(
						&conditions.MongoCondition{ConditionType: conditions.EQUAL, Key: "lastname", Value: "zahid"},
						&conditions.MongoCondition{ConditionType: conditions.EQUAL, Key: "lastname", Value: "khan"},
					))
					// log.Println(query)
				}
			}
		})
	}
}
