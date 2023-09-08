package dbconnections

import (
	"context"

	"github.com/glodb/dbfusion/query"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection struct {
	DBCommon
	currentDB string
	client    *mongo.Client
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

func (mc *MongoConnection) Insert(data interface{}) error {

	precreatedata, err := mc.PreInsert(data)

	if err != nil {
		return err
	}

	_, err = mc.client.Database(mc.currentDB).Collection(precreatedata.EntityName).InsertOne(context.TODO(), data)

	if err == nil {
		err = mc.PostInsert(mc.cache, precreatedata.Data, precreatedata.mData, mc.currentDB, precreatedata.EntityName)
	}
	return err
}

func (mc *MongoConnection) Find(interface{}) error {
	return nil
}

func (mc *MongoConnection) Update(interface{}) error {
	return nil
}

func (mc *MongoConnection) Delete(interface{}) error {
	return nil
}

func (mc *MongoConnection) ChangeDatabase(dbName string) error {
	mc.currentDB = dbName
	return nil
}

func (mc *MongoConnection) DisConnect() {

}

func (mc *MongoConnection) Filter(qmap query.QMap) {

}

func (mc *MongoConnection) Sort(order interface{}, args ...interface{}) {

}

func (mc *MongoConnection) Paginate(qmap query.QMap) {

}
func (mc *MongoConnection) Distinct(field string) {
}

func (mc *MongoConnection) RegisterSchema() {}

// New method for specifying query conditions.
func (mc *MongoConnection) Where(condition string, args ...interface{}) {}

// New methods for grouping and ordering.
func (mc *MongoConnection) GroupBy(keys string)                            {}
func (mc *MongoConnection) OrderBy(order interface{}, args ...interface{}) {}

// New methods for bulk operations.
func (mc *MongoConnection) CreateMany([]interface{}) {

}
func (mc *MongoConnection) UpdateMany([]interface{}) {

}
func (mc *MongoConnection) DeleteMany(qmap query.QMap) {

}

func (mc *MongoConnection) Skip(skip int64)             {}
func (mc *MongoConnection) Limit(limit int64)           {}
func (mc *MongoConnection) CreateTable(ifNotExist bool) {}
