package implementations

import (
	"context"

	"github.com/glodb/dbfusion/conditions"
	"github.com/glodb/dbfusion/connections"
	"github.com/glodb/dbfusion/ftypes"
	"github.com/glodb/dbfusion/queryoptions"
	"github.com/glodb/dbfusion/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	DBCommon
	client *mongo.Client
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

func (mc *MongoConnection) UpdateOne(interface{}) error {
	defer mc.refreshValues()
	return nil
}

func (mc *MongoConnection) DeleteOne(interface{}) error {
	defer mc.refreshValues()
	return nil
}

func (mc *MongoConnection) DisConnect() {

}

func (mc *MongoConnection) Paginate(interface{}, ...queryoptions.FindOptions) {

}
func (mc *MongoConnection) Distinct(field string) {
}

func (mc *MongoConnection) RegisterSchema() {}

// New methods for bulk operations.
func (mc *MongoConnection) CreateMany([]interface{}) {

}
func (mc *MongoConnection) UpdateMany([]interface{}) {

}
func (mc *MongoConnection) DeleteMany(qmap ftypes.QMap) {

}

func (mc *MongoConnection) CreateTable(ifNotExist bool) {}

func (mc *MongoConnection) Skip(skip int64) connections.MongoConnection {
	mc.skip = skip
	return mc
}
func (mc *MongoConnection) Limit(limit int64) connections.MongoConnection {
	mc.limit = limit
	return mc
}
func (mc *MongoConnection) Project(keys map[string]bool) connections.MongoConnection {
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
