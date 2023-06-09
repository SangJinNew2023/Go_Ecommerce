package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017")) //ApplyURI parses the given URI and sets options accordingly
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.TODO(), nil) // 기한,값이 없고 취소되지 않는 빈 context 반환, Send a ping to confirm a successful connection
	if err != nil {
		log.Println("failed to connect to mongodb:(")
		return nil
	}
	fmt.Println("Successfully connect to  mongodb")
	return client
}

var Client *mongo.client = DBset()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Cllection(collectionName)
	return collection

}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productCollection *mongo.Collection = client.Database("Ecommerce").Cllection(collectionName)
	return productCollection
}
