package database

import (
	"fmt"
	"sync"
)

var MongoConnectionPool = sync.Pool{
	// function to execute when no instance of a buffer is not found
	New: func() interface{} {
		fmt.Println("Created connection...")
		dbConnection, err := ConnectToDB()
		if err != nil {
			panic("connection to DB failed")
		}
		return dbConnection
	},
}
