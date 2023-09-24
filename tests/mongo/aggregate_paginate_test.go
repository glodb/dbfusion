package mongotest

import (
	"log"
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/tests/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMongoAggregatePagination(t *testing.T) {
	// Define valid database name and URI for MongoDB connection.
	validDBName := "testDBFusion"
	validUri := "mongodb://localhost:27017"

	// Create a Redis cache and connect to it.
	cache := caches.RedisCache{}
	err := cache.ConnectCache("localhost:6379")
	if err != nil {
		t.Errorf("Error in Redis connection: %v", err)
	}

	// Configure options for the DBFusion connection, including database name, URI, and cache.
	options := dbfusion.Options{
		DbName: &validDBName,
		Uri:    &validUri,
		Cache:  &cache,
	}

	// Get a MongoDB connection instance using DBFusion with the specified options.
	con, err := dbfusion.GetInstance().GetMongoConnection(options)
	if err != nil {
		t.Errorf("DBConnection failed with error: %v", err)
	}

	// Define test cases for aggregation with pagination.
	testCases := []struct {
		Con            connections.MongoConnection
		Conditions     []int
		TestData       AggregationTestData
		Data           interface{}
		TableName      string
		ExpectedResult AggregationResults
		PageNumber     int
		Name           string
	}{
		{
			Con:        con,
			Conditions: []int{MATCH},
			TestData:   AggregationTestData{match: primitive.M{"firstname": "Aafaq"}},
			Data:       &[]models.UserTest{},
			TableName:  "users",
			Name:       "Testing aggregate with Match",
			PageNumber: 2,
		},
	}

	// Iterate through the test cases.
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Set the MongoDB collection to operate on.
			con.Table(tc.TableName)

			// Apply aggregation stages based on test data.
			for _, condition := range tc.Conditions {
				switch condition {
				case MATCH:
					con.Match(tc.TestData.match)
				case COUNT:
					con.Count(tc.TestData.count)
				case BUCKET:
					con.Bucket(tc.TestData.bucket)
				case BUCKETAUTO:
					con.BucketsAuto(tc.TestData.bucketAuto)
				case ADDFIELDS:
					con.AddFields(tc.TestData.addFields)
				case GEONEAR:
					con.GeoNear(tc.TestData.geoNear)
				case GROUP:
					con.Group(tc.TestData.group)
				case LIMIT:
					con.LimitAggregate(tc.TestData.limitAggregate)
				case SKIP:
					con.SkipAggregate(tc.TestData.skipAggregate)
				case SORT:
					con.SortAggregate(tc.TestData.sortAggregate)
				case PROJECT:
					con.Project(tc.TestData.project)
				case SORTCOUNT:
					con.SortByCount(tc.TestData.sortCount)
				case UNSET:
					con.Unset(tc.TestData.unset)
				case REPLACEWITH:
					con.ReplaceWith(tc.TestData.replaceWith)
				case MERGE:
					con.Merge(tc.TestData.merge)
				case OUT:
					con.Out(tc.TestData.out)
				case REPLACEROOT:
					con.ReplaceRoot(tc.TestData.replaceRoot)
				case FACET:
					con.Facet(tc.TestData.facet)
				case COLLSTATS:
					con.CollStats(tc.TestData.collStats)
				case INDEXSTATS:
					con.IndexStats(tc.TestData.indexStats)
				case PLANCACHESTATS:
					con.PlanCacheStats(tc.TestData.planCacheStats)
				case REDACT:
					con.Redact(tc.TestData.redact)
				case REPLACECOUNT:
					con.ReplaceCount(tc.TestData.replaceCount)
				case SAMPLE:
					con.Sample(tc.TestData.sample)
				case SET:
					con.Set(tc.TestData.set)
				case UNWIND:
					con.Unwind(tc.TestData.unwind)
				case LOOKUP:
					con.Lookup(tc.TestData.lookup)
				case GRAPHLOOKUP:
					con.GraphLookup(tc.TestData.graphLookup)
				}
			}

			// Perform aggregation with pagination and capture results.
			pageinationResults, err := con.AggregatePaginate(tc.Data, tc.PageNumber)

			// Log the aggregation paginate result.
			log.Println("Aggregation Paginate Result:", pageinationResults)

			// Check if the actual error matches the expected error.
			if err != tc.ExpectedResult.err {
				t.Errorf("Expected error: %v, got error: %v", tc.ExpectedResult.err, err)
				return
			}

			// Check if the pagination returned the expected number of results.
			if len(*tc.Data.(*[]models.UserTest)) != 10 {
				t.Errorf("Expected %d results in pagination, but got %d", 10, len(*tc.Data.(*[]models.UserTest)))
				return
			}
		})
	}
}
