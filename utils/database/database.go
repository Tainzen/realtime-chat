package database

import (
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Db setup connection with mongodb and returns mongo client
func Db() *mongo.Client {

	dbDriver := os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dburl := dbHost + ":" + dbPort

	//setting credentials for mongodb user
	credential := options.Credential{
		AuthSource: dbName,
		Username:   dbUser,
		Password:   dbPass,
	}

	clientOptions := options.Client().ApplyURI(dbDriver + "://" + dburl + "/").SetAuth(credential)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Panic().Msgf("Error connecting to mongodb")
	}

	// // Check the connection to MongoDB
	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Panic().Msgf("Error while pinging to mongodb")
	// }

	log.Info().Msgf("Successfully established connection to mongodb")
	return client
}
