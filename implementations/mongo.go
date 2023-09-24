package implementations

import (
	"context"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/hooks"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	DBCommon
	client         *mongo.Client
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

//TODO: add the communication with certificate
//TODO: add the options for cache enabling and disabling
//TODO: add thhe
func (mc *MongoConnection) ConnectWithCertificate(uri string, filePath string) error {
	return nil
}

func (mc *MongoConnection) Connect(uri string) error {
	var err error

	clientOptions := options.Client().ApplyURI(uri)

	// Connect to the MongoDB server
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	mc.client = client

	// if mc.certSet {
	// 	connectionURI := fmt.Sprintf(connectionStringMain, connectionConfigMap["dbUserMongo"].(string), connectionConfigMap["dbPasswordMongo"].(string), connectionConfigMap["dbHostMongo"].(string)+":"+strconv.Itoa(int(connectionConfigMap["dbPortMongo"].(float64))), connectionConfigMap["dbNameMongo"].(string), readPreference)
	// 	tlsConfig, err := u.getCustomTLSConfig(connectionConfigMap["dbCAFileMongo"].(string))
	// 	if err != nil {
	// 		return nil, errors.New("Unable to get tls config")
	// 	}
	// 	client, err = mongo.NewClient(options.Client().ApplyURI(connectionURI).SetTLSConfig(tlsConfig))
	// 	if err != nil {
	// 		return nil, errors.New("failed to create client")
	// 	}
	// } else {
	// 	connectionURI := fmt.Sprintf(connectionStringDev, connectionConfigMap["dbHostMongo"].(string)+":"+strconv.Itoa(int(connectionConfigMap["dbPortMongo"].(float64))), connectionConfigMap["dbNameMongo"].(string))
	// 	client, err = mongo.NewClient(options.Client().ApplyURI(connectionURI))

	// 	if err != nil {
	// 		return nil, errors.New("failed to create client")
	// 	}
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	// defer cancel()

	// err = client.Connect()
	// if err != nil {
	// 	log.Fatalf("Failed to connect to cluster: %v", err)
	// }

	// // Force a connection to verify our connection string
	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatalf("Failed to ping cluster: %v", err)
	// }
	// u.client = client
	// u.dbName = dbName
	return nil
}

func (mc *MongoConnection) Table(tablename string) connections.MongoConnection {
	mc.setTable(tablename)
	// mc.tableName = tablename
	return mc
}

func (mc *MongoConnection) InsertOne(data interface{}) error {
	defer mc.refreshValues()
	preCreateData, err := mc.preInsert(data)

	if err != nil {
		return err
	}
	_, err = mc.client.Database(mc.currentDB).Collection(preCreateData.entityName).InsertOne(context.TODO(), preCreateData.mData)

	if err == nil {
		err = mc.postInsert(mc.cache, preCreateData.Data, preCreateData.mData, mc.currentDB, preCreateData.entityName)
	}
	return err
}

func (mc *MongoConnection) FindOne(result interface{}, dbFusionOptions ...queryoptions.FindOptions) error {
	defer mc.refreshValues()
	if mc.whereQuery != nil {
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		mc.whereQuery = query
	} else {
		mc.whereQuery = &conditions.MongoData{}
	}
	prefindReturn, err := mc.preFind(mc.cache, result, dbFusionOptions...)
	if err != nil {
		return err
	}
	if prefindReturn.queryDatabase {
		opts := options.FindOneOptions{}
		if mc.projection != nil {
			opts.SetProjection(mc.projection)
		}

		if mc.skip != 0 {
			opts.SetSkip(mc.skip)
		}

		if mc.sort != nil {
			opts.SetSort(mc.sort)
		}

		err = mc.client.Database(mc.currentDB).Collection(prefindReturn.entityName).FindOne(context.TODO(), prefindReturn.query).Decode(result)
		if err != nil {
			return err
		}
	}

	err = mc.postFind(mc.cache, result, prefindReturn.entityName, dbFusionOptions...)
	return err
}

func (mc *MongoConnection) UpdateAndFindOne(data interface{}, result interface{}, upsert bool) error {
	defer mc.refreshValues()

	var fusionQuery conditions.DBFusionData
	if mc.whereQuery != nil {
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		fusionQuery = query
	} else {
		fusionQuery = &conditions.MongoData{}
	}

	preUpdateReturn, err := mc.preUpdate(data, connections.MONGO)
	if err != nil {
		return err
	}
	opts := options.FindOneAndUpdateOptions{}
	if mc.projection != nil {
		opts.SetProjection(mc.projection)
	}

	if mc.sort != nil {
		opts.SetSort(mc.sort)
	}

	opts.SetUpsert(upsert)

	opts.SetReturnDocument(options.After)

	updateCache := false
	oldKeys := []string{}
	newKeys := []string{}
	var cacheHook hooks.CacheHook
	if value, ok := interface{}(result).(hooks.CacheHook); ok {
		err = mc.client.Database(mc.currentDB).Collection(preUpdateReturn.entityName).FindOne(context.TODO(), fusionQuery.GetQuery().(primitive.D)).Decode(result)
		if err != nil {
			return err
		}
		tagMapValue, err := mc.createTagValueMap(result)
		if err == nil {
			oldKeys = mc.getAllCacheValues(value, tagMapValue, preUpdateReturn.entityName)
			updateCache = true
			cacheHook = value
		}
	}
	err = mc.client.Database(mc.currentDB).Collection(preUpdateReturn.entityName).FindOneAndUpdate(context.TODO(), fusionQuery.GetQuery().(primitive.D), preUpdateReturn.queryData.(primitive.D), &opts).Decode(result)
	if err != nil {
		return err
	}
	if updateCache {
		tagMapValue, _ := mc.createTagValueMap(result)
		newKeys = mc.getAllCacheValues(cacheHook, tagMapValue, preUpdateReturn.entityName)
	}
	err = mc.postUpdate(mc.cache, result, preUpdateReturn.entityName, oldKeys, newKeys)
	return err
}

func (mc *MongoConnection) DeleteOne(sliceData ...interface{}) error {
	defer mc.refreshValues()
	var data interface{}
	if len(sliceData) != 0 {
		data = sliceData[0]
	}
	if mc.whereQuery != nil {
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return err
		}
		mc.whereQuery = query
	} else {
		mc.whereQuery = &conditions.MongoData{}
	}
	preDeleteData, err := mc.preDelete(data)

	if err != nil {
		return err
	}
	var deleteQuery primitive.D
	var results primitive.M
	// tagValueMap := make(map[string]interface{})
	if data != nil { //Need to delete from a struct
		deleteQuery = mc.buildMongoData(preDeleteData.dataType, preDeleteData.dataValue)

		err = mc.client.Database(mc.currentDB).Collection(preDeleteData.entityName).FindOneAndDelete(context.TODO(), deleteQuery).Decode(&results)

	} else { //need to delete from whereClause

		//Simple Delete the data here as cann't check the cache for this
		_, err = mc.client.Database(mc.currentDB).Collection(preDeleteData.entityName).DeleteOne(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery())
	}

	if err != nil {
		return err
	}

	err = mc.postDelete(mc.cache, data, preDeleteData.entityName, results)

	return err
}

func (mc *MongoConnection) DisConnect() error {
	return mc.client.Disconnect(context.TODO())
}

func (mc *MongoConnection) Paginate(results interface{}, pageNumber int) (connections.PaginationResults, error) {
	defer mc.refreshValues()
	if mc.whereQuery != nil {
		query, err := utils.GetInstance().GetMongoFusionData(mc.whereQuery)
		if err != nil {
			return connections.PaginationResults{}, err
		}
		mc.whereQuery = query
	} else {
		mc.whereQuery = &conditions.MongoData{}
	}
	var paginationResults connections.PaginationResults
	opts := options.FindOptions{}
	if mc.projection != nil {
		opts.SetProjection(mc.projection)
	}

	count, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).CountDocuments(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery())

	if err != nil {
		return connections.PaginationResults{}, err
	}
	paginationResults.TotalDocuments = count
	paginationResults.TotalPages = int64(math.Ceil((float64(count) / float64(mc.pageSize))))
	paginationResults.Limit = int64(mc.pageSize)
	paginationResults.CurrentPage = int64(pageNumber)

	mc.limit = int64(mc.pageSize)
	mc.skip = int64((pageNumber - 1) * mc.pageSize)
	opts.SetSkip(mc.skip)
	opts.SetSort(mc.sort)
	opts.SetLimit(mc.limit)

	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Find(context.TODO(), mc.whereQuery.(conditions.DBFusionData).GetQuery(), &opts)

	if err = cursor.All(context.TODO(), results); err != nil {
		return connections.PaginationResults{}, err
	}

	return paginationResults, nil
}

// New methods for bulk operations.
func (mc *MongoConnection) CreateMany([]interface{}) {

}
func (mc *MongoConnection) FindMany(interface{}) {
}
func (mc *MongoConnection) UpdateMany([]interface{}) {

}
func (mc *MongoConnection) DeleteMany(qmap ftypes.QMap) {

}

func (mc *MongoConnection) Skip(skip int64) connections.MongoConnection {
	mc.skip = skip
	return mc
}
func (mc *MongoConnection) Limit(limit int64) connections.MongoConnection {
	mc.limit = limit
	return mc
}
func (mc *MongoConnection) Select(keys map[string]bool) connections.MongoConnection {
	selectionKeys := make(map[string]int, 0)

	for key, val := range keys {
		if val {
			selectionKeys[key] = 1
		} else {
			selectionKeys[key] = 0
		}
	}
	mc.projection = selectionKeys
	return mc
}

func (mc *MongoConnection) Sort(sortKey string, sortdesc ...bool) connections.MongoConnection {
	sortString := sortKey
	sortVal := 1
	if len(sortdesc) > 0 {
		if !sortdesc[0] {
			sortVal = -1
		}
	}

	if mc.sort != nil {
		sortMap := mc.sort.(map[string]interface{})
		sortMap[sortString] = sortVal
		mc.sort = sortMap
	} else {
		sortMap := make(map[string]interface{})
		mc.sort = sortMap
	}
	return mc
}

func (mc *MongoConnection) Where(query interface{}) connections.MongoConnection {
	mc.whereQuery = query
	return mc
}

func (mc *MongoConnection) SetPageSize(limit int) {
	mc.pageSize = limit
}

func (mc *MongoConnection) parseSortableIndexes(indexes []string) []bson.D {
	mongoIndexes := []bson.D{}
	for _, index := range indexes {
		compounds := strings.Split(index, ",")
		localIndex := bson.D{}
		for _, compoundIndex := range compounds {
			indexValues := strings.Split(compoundIndex, ":")
			if len(indexValues) > 1 {
				intVal, _ := strconv.ParseInt(indexValues[1], 10, 64)
				singleIndex := bson.E{Key: indexValues[0], Value: intVal}
				localIndex = append(localIndex, singleIndex)
			}
		}
		mongoIndexes = append(mongoIndexes, localIndex)
	}
	return mongoIndexes
}

func (mc *MongoConnection) parseIndexes(indexes []string, indexName string) []bson.D {
	mongoIndexes := []bson.D{}
	for _, index := range indexes {
		compounds := strings.Split(index, ",")
		localIndex := bson.D{}
		for _, compoundIndex := range compounds {
			singleIndex := bson.E{Key: compoundIndex, Value: indexName}
			localIndex = append(localIndex, singleIndex)
		}
		mongoIndexes = append(mongoIndexes, localIndex)
	}
	return mongoIndexes
}
func (mc *MongoConnection) CreateIndexes(data interface{}) error {

	name, _ := mc.getEntityName(data)

	if val, ok := interface{}(data).(hooks.NormalIndexes); ok {
		indexes := val.GetNormalIndexes()
		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}

	}

	if val, ok := interface{}(data).(hooks.UniqueIndexes); ok {
		indexes := val.GetUniqueIndexes()
		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys:    index,
				Options: options.Index().SetUnique(true),
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	if val, ok := interface{}(data).(hooks.TextIndexes); ok {
		indexes := val.GetTextIndex()
		compundIndexes := mc.parseIndexes([]string{indexes}, "text")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	if val, ok := interface{}(data).(hooks.TwoDimensialIndexes); ok {
		indexes := val.Get2DIndexes()
		compundIndexes := mc.parseIndexes(indexes, "2d")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	if val, ok := interface{}(data).(hooks.TwoDimensialSpatialIndexes); ok {
		indexes := val.Get2DSpatialIndexes()
		compundIndexes := mc.parseIndexes(indexes, "2dsphere")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	if val, ok := interface{}(data).(hooks.HashedIndexes); ok {
		indexes := val.GetHashedIndexes()
		compundIndexes := mc.parseIndexes(indexes, "hashed")
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys: index,
			}
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			if err != nil {
				return err
			}
		}
	}

	if val, ok := interface{}(data).(hooks.SparseIndexes); ok {
		indexes := val.GetSparseIndexes()

		compundIndexes := mc.parseSortableIndexes(indexes)
		for _, index := range compundIndexes {
			indexModel := mongo.IndexModel{
				Keys:    index,
				Options: options.Index().SetSparse(true),
			}
			log.Println(mc.parseSortableIndexes(indexes))
			_, err := mc.client.Database(mc.currentDB).Collection(name.entityName).Indexes().CreateOne(context.TODO(), indexModel)
			log.Println(err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (mc *MongoConnection) Match(data interface{}) connections.MongoConnection {
	mc.match = data
	return mc
}
func (mc *MongoConnection) Bucket(data interface{}) connections.MongoConnection {
	mc.bucket = data
	return mc
}
func (mc *MongoConnection) BucketsAuto(data interface{}) connections.MongoConnection {
	mc.bucketAuto = data
	return mc
}
func (mc *MongoConnection) AddFields(data interface{}) connections.MongoConnection {
	mc.addFields = data
	return mc
}
func (mc *MongoConnection) GeoNear(data interface{}) connections.MongoConnection {
	mc.geoNear = data
	return mc
}
func (mc *MongoConnection) Group(data interface{}) connections.MongoConnection {
	mc.group = data
	return mc
}
func (mc *MongoConnection) LimitAggregate(data int) connections.MongoConnection {
	mc.limitAggregate = data
	return mc
}
func (mc *MongoConnection) SkipAggregate(data int) connections.MongoConnection {
	mc.skipAggregate = data
	return mc
}
func (mc *MongoConnection) SortAggregate(data interface{}) connections.MongoConnection {
	mc.sortAggregate = data
	return mc
}
func (mc *MongoConnection) SortByCount(data interface{}) connections.MongoConnection {
	mc.sortCount = data
	return mc
}
func (mc *MongoConnection) Project(data interface{}) connections.MongoConnection {
	mc.project = data
	return mc
}
func (mc *MongoConnection) Unset(data interface{}) connections.MongoConnection {
	mc.unset = data
	return mc
}
func (mc *MongoConnection) ReplaceWith(data interface{}) connections.MongoConnection {
	mc.replaceWith = data
	return mc
}
func (mc *MongoConnection) Merge(data interface{}) connections.MongoConnection {
	mc.merge = data
	return mc
}
func (mc *MongoConnection) Out(data interface{}) connections.MongoConnection {
	mc.out = data
	return mc
}
func (mc *MongoConnection) Facet(data interface{}) connections.MongoConnection {
	mc.facet = data
	return mc
}
func (mc *MongoConnection) CollStats(data interface{}) connections.MongoConnection {
	mc.collStats = data
	return mc
}
func (mc *MongoConnection) IndexStats(data interface{}) connections.MongoConnection {
	mc.indexStats = data
	return mc
}
func (mc *MongoConnection) PlanCacheStats(data interface{}) connections.MongoConnection {
	mc.planCacheStats = data
	return mc
}
func (mc *MongoConnection) Redact(data interface{}) connections.MongoConnection {
	mc.redact = data
	return mc
}
func (mc *MongoConnection) ReplaceRoot(data interface{}) connections.MongoConnection {
	mc.replaceRoot = data
	return mc
}
func (mc *MongoConnection) ReplaceCount(data interface{}) connections.MongoConnection {
	mc.replaceCount = data
	return mc
}
func (mc *MongoConnection) Sample(data interface{}) connections.MongoConnection {
	mc.sample = data
	return mc
}
func (mc *MongoConnection) Set(data interface{}) connections.MongoConnection {
	mc.set = data
	return mc
}
func (mc *MongoConnection) Unwind(data interface{}) connections.MongoConnection {
	mc.unwind = data
	return mc
}
func (mc *MongoConnection) Lookup(data interface{}) connections.MongoConnection {
	mc.lookup = data
	return mc
}
func (mc *MongoConnection) GraphLookup(data interface{}) connections.MongoConnection {
	mc.graphLookup = data
	return mc
}
func (mc *MongoConnection) Count(data interface{}) connections.MongoConnection {
	mc.count = data
	return mc
}

func (mc *MongoConnection) refreshAggregation() {
	mc.match = nil
	mc.count = nil
	mc.bucket = nil
	mc.bucketAuto = nil
	mc.addFields = nil
	mc.geoNear = nil
	mc.group = nil
	mc.limitAggregate = 0
	mc.skipAggregate = 0
	mc.sortAggregate = nil
	mc.project = nil
	mc.sortCount = nil
	mc.unset = nil
	mc.replaceWith = nil
	mc.merge = nil
	mc.out = nil
	mc.replaceRoot = nil
	mc.facet = nil
	mc.collStats = nil
	mc.indexStats = nil
	mc.planCacheStats = nil
	mc.redact = nil
	mc.replaceCount = nil
	mc.sample = nil
	mc.set = nil
	mc.unwind = nil
	mc.lookup = nil
	mc.graphLookup = nil
}
func (mc *MongoConnection) createAggregation() bson.A {
	pipelines := primitive.A{}

	if mc.match != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$match", Value: mc.match}})
	}
	if mc.count != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$count", Value: mc.count}})
	}
	if mc.bucket != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$bucket", Value: mc.bucket}})
	}
	if mc.bucketAuto != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$bucketAuto", Value: mc.bucketAuto}})
	}
	if mc.addFields != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$addFields", Value: mc.addFields}})
	}
	if mc.geoNear != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$geoNear", Value: mc.geoNear}})
	}
	if mc.group != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$group", Value: mc.group}})
	}
	if mc.skipAggregate != 0 {
		pipelines = append(pipelines, primitive.D{{Key: "$skip", Value: mc.skipAggregate}})
	}
	if mc.limitAggregate != 0 {
		pipelines = append(pipelines, primitive.D{{Key: "$limit", Value: mc.limitAggregate}})
	}
	if mc.sortAggregate != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sort", Value: mc.sortAggregate}})
	}
	if mc.project != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$project", Value: mc.project}})
	}
	if mc.sortCount != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sortCount", Value: mc.sortCount}})
	}
	if mc.unset != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$unset", Value: mc.unset}})
	}
	if mc.replaceWith != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceWith", Value: mc.replaceWith}})
	}
	if mc.merge != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$merge", Value: mc.merge}})
	}
	if mc.out != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$out", Value: mc.out}})
	}
	if mc.replaceRoot != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceRoot", Value: mc.replaceRoot}})
	}
	if mc.facet != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$facet", Value: mc.facet}})
	}
	if mc.collStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$collStats", Value: mc.collStats}})
	}
	if mc.indexStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$indexStats", Value: mc.indexStats}})
	}
	if mc.planCacheStats != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$planCacheStats", Value: mc.planCacheStats}})
	}
	if mc.redact != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$redact", Value: mc.redact}})
	}
	if mc.replaceCount != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$replaceCount", Value: mc.replaceCount}})
	}
	if mc.sample != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$sample", Value: mc.sample}})
	}
	if mc.set != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$set", Value: mc.set}})
	}
	if mc.unwind != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$unwind", Value: mc.unwind}})
	}
	if mc.lookup != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$lookup", Value: mc.lookup}})
	}
	if mc.graphLookup != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$graphLookup", Value: mc.graphLookup}})
	}
	return pipelines
}
func (mc *MongoConnection) Aggregate(data interface{}) error {
	defer mc.refreshAggregation()
	defer mc.refreshValues()
	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), mc.createAggregation())
	if err != nil {
		return err
	}
	if err = cursor.All(context.TODO(), data); err != nil {
		return err
	}
	return nil
}

func (mc *MongoConnection) AggregatePaginate(data interface{}, pageNumber int) (paginationResults connections.PaginationResults, err error) {
	defer mc.refreshAggregation()
	defer mc.refreshValues()

	countGoupStage := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: nil}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

	pipelines := primitive.A{}

	if mc.match != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$match", Value: mc.match}})
	}

	if mc.group != nil {
		pipelines = append(pipelines, primitive.D{{Key: "$group", Value: mc.group}})
	}
	pipelines = append(pipelines, countGoupStage)

	countData := []ftypes.QMap{}
	cursor, err := mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), pipelines)
	if err != nil {
		return
	}
	if err = cursor.All(context.TODO(), &countData); err != nil {
		return
	}

	if len(countData) > 0 {
		count := int64(countData[0]["count"].(int32))
		paginationResults.TotalDocuments = count
		paginationResults.TotalPages = int64(math.Ceil((float64(count) / float64(mc.pageSize))))
		paginationResults.Limit = int64(mc.pageSize)
		paginationResults.CurrentPage = int64(pageNumber)

		mc.limitAggregate = mc.pageSize
		mc.skipAggregate = int((pageNumber - 1) * mc.pageSize)

		cursor, err = mc.client.Database(mc.currentDB).Collection(mc.tableName).Aggregate(context.TODO(), mc.createAggregation())
		if err != nil {
			return
		}
		if err = cursor.All(context.TODO(), data); err != nil {
			return
		}
	}
	return
}
