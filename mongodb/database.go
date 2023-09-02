package mongodb

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
)

func C(name string) *mongo.Collection {
	return DB(resolveDatabase()).Collection(name)
}

func DB(name string) *mongo.Database {
	return GetClient().Database(name)
}

func GetClient() *mongo.Client {
	if client == nil {
		connect()
	} else {
		if err := client.Ping(context.Background(), nil); err != nil {
			connect()
		}
	}
	return client
}

func connect() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environmental variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
}

func resolveDatabase() string {
	db := os.Getenv("MONGODB_DATABASE")
	if db == "" {
		panic("cant resolve mongo dsn")
	}
	return db
}

type M bson.M
type D bson.D
type A bson.A
type E bson.E
type Raw bson.Raw
type RawElement bson.RawElement
