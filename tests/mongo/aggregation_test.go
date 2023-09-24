package mongotest

import (
	"testing"

	"github.com/glodb/dbfusion"
	"github.com/glodb/dbfusion/caches"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/tests/models"
)

type AggregationResults struct {
	data interface{}
	err  error
}

const (
	MATCH = iota
	COUNT
	BUCKET
	BUCKETAUTO
	ADDFIELDS
	GEONEAR
	GROUP
	LIMIT
	SKIP
	SORT
	PROJECT
	SORTCOUNT
	UNSET
	REPLACEWITH
	MERGE
	OUT
	REPLACEROOT
	FACET
	COLLSTATS
	INDEXSTATS
	PLANCACHESTATS
	REDACT
	REPLACECOUNT
	SAMPLE
	SET
	UNWIND
	LOOKUP
	GRAPHLOOKUP
)

type AggregationTestData struct {
	match          interface{}
	count          interface{}
	bucket         interface{}
	bucketAuto     interface{}
	addFields      interface{}
	geoNear        interface{}
	group          interface{}
	limitAggregate int
	skipAggregate  int
	sortAggregate  interface{}
	project        interface{}
	sortCount      interface{}
	unset          interface{}
	replaceWith    interface{}
	merge          interface{}
	out            interface{}
	replaceRoot    interface{}
	facet          interface{}
	collStats      interface{}
	indexStats     interface{}
	planCacheStats interface{}
	redact         interface{}
	replaceCount   interface{}
	sample         interface{}
	set            interface{}
	unwind         interface{}
	lookup         interface{}
	graphLookup    interface{}
}

func TestMongoAggregation(t *testing.T) {
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
	testCases := []struct {
		Con            connections.MongoConnection
		Conditions     []int
		TestData       AggregationTestData
		Data           interface{}
		TableName      string
		ExpectedResult AggregationResults
		Name           string
	}{
		{
			Con:        con,
			Conditions: []int{MATCH},
			TestData:   AggregationTestData{match: ftypes.QMap{"firstname": "Aafaq"}},
			Data:       &[]models.UserTest{},
			TableName:  "users",
			Name:       "Testing aggregate with Match",
		},
		{
			Con:        con,
			Conditions: []int{MATCH, PROJECT},
			TestData: AggregationTestData{match: ftypes.QMap{"firstname": "Aafaq"},
				project: ftypes.QMap{"firstname": 1}},
			Data:      &[]models.UserTest{},
			TableName: "users",
			Name:      "Testing aggregate with Match",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			con.Table(tc.TableName)
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
			err := con.Aggregate(tc.Data)

			// log.Println("Aggregation Result:", tc.Data)

			if err != tc.ExpectedResult.err {
				t.Errorf("Expected %v, got %v", tc.ExpectedResult, err)
				return
			}

		})
	}
}
