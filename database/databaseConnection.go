package database

import (
	"com.aharakitchen/app/config"
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

func ConnectToDB() (*Connection, error) {
	p := config.Config("DB_PORT")
	n := config.Config("DB_NAME")
	h := config.Config("DB_HOST")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(n+h+p))
	if err != nil {
		return nil, err
	}

	// create database
	db := client.Database("blog-service")

	// create collection
	postsCollection := db.Collection("posts")
	tagsCollection := db.Collection("tags")
	blackListCollection := db.Collection("blacklist")

	dbConnection := &Connection{client, postsCollection, tagsCollection,blackListCollection,db}
	return dbConnection, nil
}
