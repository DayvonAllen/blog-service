package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Connection struct {
	*mongo.Client
	PostCollection     *mongo.Collection
	TagCollection    *mongo.Collection
	BlackListCollection    *mongo.Collection
	*mongo.Database
}

var MongoConn *Connection

func ConnectToDB() {
	//p := config.Config("DB_PORT")
	//n := config.Config("DB_NAME")
	//h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	//socketTimeOut := time.Second *4
	dbOptions := options.ClientOptions{
		//SocketTimeout: &socketTimeOut,
	}

	client, err := mongo.Connect(ctx, dbOptions.ApplyURI("mongodb://backend-mongo-srv:27017/blog"))
	if err != nil {
		panic(err)
	}

	// create database
	db := client.Database("blog")

	// create collection
	postsCollection := db.Collection("posts")
	tagsCollection := db.Collection("tags")
	blackListCollection := db.Collection("blacklist")


	dbConnection := &Connection{client, postsCollection, tagsCollection,blackListCollection,db}
	MongoConn = dbConnection
}
