package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Gigfinder-io/util.gigfinder.io/log"
)

var (
	Address string = "localhost"
	User    string = ""
	Pass    string = ""
)

type Query map[string]interface{}
type Fields bson.D

var gigfinderDB *mongo.Database

func Connect() error {
	// Set client options
	uri := ""
	var clientOptions *options.ClientOptions
	if User == "" {
		uri = fmt.Sprintf("mongodb://%v", Address)
		clientOptions = options.Client().ApplyURI(uri)
	} else {
		uri = fmt.Sprintf("mongodb+srv://%v:%v@%v/BandMatch", User, Pass, Address)
		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions = options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)
	}
	log.Msgf(0, "mongo srv: [%v]", uri)

	// Connect to MongoDB
	var err error
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return fmt.Errorf("could not connect to mongodb: %v", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return fmt.Errorf("could not ping mongodb: %v", err)
	}

	log.Msg(log.V, "Connected to MongoDB!")

	gigfinderDB = client.Database("BandMatch")

	return nil
}

func Disconnect() {
	gigfinderDB.Client().Disconnect(context.TODO())
}
