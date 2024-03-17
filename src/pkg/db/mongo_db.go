package db

import (
	"context"
	"log"

	"github.com/UrlShortener/src/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShortenUrlCollection struct {
	_id         int
	Url         string
	Hash        string
	CreatedDate string
	Clicks      int
}

var mongoClient *mongo.Client

func connect() {

	if mongoClient != nil {
		log.Panicf("Already connected to mongo db\n")
		return
	}

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.AppConfig.MongoConnectionString).SetServerAPIOptions(serverApi)
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Panicf("Error while connecting err %e", err)
		panic(err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	// Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }

	mongoClient = client

}

func (sh *ShortenUrlCollection) Save() {
	connect()
	res, err := mongoClient.Database(config.AppConfig.MongoDBName).Collection(config.AppConfig.MongoCollectionName).InsertOne(context.Background(), sh)

	if err != nil {
		log.Panicf("Error while inserting entry to mongo %e\n", err)
		panic(err)
	}

	log.Printf("Entry saved %s\n", res.InsertedID)
}

func FindOneByHash(hash string) ShortenUrlCollection {
	connect()
	var op ShortenUrlCollection
	res := mongoClient.Database(config.AppConfig.MongoDBName).Collection(config.AppConfig.MongoCollectionName).FindOne(context.TODO(), bson.M{"hash": hash})

	err := res.Decode(&op)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("No matching entry for hash %s\n", hash)
		} else {
			log.Panicf("Error while fetch entry %e\n", err)
		}
	}
	return op
}
