package Databases

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Driver       string `json:"driver"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	Replica      string `json:"replica"`
	DatabaseName string `json:"name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Connection   *mongo.Client
	Database     *mongo.Database
}

func (mdb *MongoDB) IsEmpty() bool {
	return mdb.DatabaseName == ""
}

func (mdb *MongoDB) Initial() error {
	if !mdb.IsEmpty() {
		var err error
		var client *mongo.Client
		var clientOptions *options.ClientOptions
		if mdb.User != "" {
			// Below clientOptions is a correct way to connect to MongoDB using MongoDB Atlas (mongo in da cloud)
			// If you are using standalone instance or your own on premises cluster, you must use this format instead:
			// clientOptions = options.Client().ApplyURI("mongodb://" + dbConfigurator.User + ":" + dbConfigurator.Password + "@" + dbConfigurator.Host + ":" + dbConfigurator.Port)
			// We have database.type config parameter in dev.env.json
			// So I suggest to as additional check:
			// If database.type = "atlas" then use below connection string
			// Else use connection string like in this comment (not using 'mongodb+srv' and appending port number at the end)
			clientOptions = options.Client().ApplyURI("mongodb://" + mdb.User + ":" + mdb.Password + "@" + mdb.Host + ":" + mdb.Port + mdb.Replica)
		} else {
			clientOptions = options.Client().ApplyURI("mongodb://" + mdb.Host + ":" + mdb.Port)
		}
		client, err = mongo.Connect(context.TODO(), clientOptions)

		if err == nil {
			mdb.Connection = client
			mdb.Database = client.Database(mdb.DatabaseName)
			err = mdb.Connection.Ping(context.TODO(), nil)

			return err
		}
		return err
	}
	return nil
}

func (mdb *MongoDB) GetConnection() (*mongo.Client, error) {
	return mdb.Connection, nil
}

func (mdb *MongoDB) GetDatabase() *mongo.Database {
	return mdb.Database
}

func (mdb *MongoDB) Close(ctx context.Context) error {
	return mdb.Connection.Disconnect(ctx)
}
