package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getProjects() map[string]Project {
	db_connection := getConnectionString()
	client, ctx, cancel := createClientAndContext(db_connection)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("projects").Collection("project_collection")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	retrievedProjects := make(map[string]Project)

	for cursor.Next(ctx) {
		var project Project
		if err := cursor.Decode(&project); err != nil {
			log.Fatal(err)
		}

		retrievedProjects[project.Title] = project
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return retrievedProjects
}

func saveProjects(projects map[string]Project) {
	// connect
	db_connection := getConnectionString()
	client, ctx, cancel := createClientAndContext(db_connection)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("projects").Collection("project_collection")

	// add projects
	interfaceProjects := getInterfacesFromProjects(projects)
	_, err := collection.InsertMany(ctx, interfaceProjects)
	if err != nil {
		log.Fatal(err)
	}
}

func getConnectionString() string {
	host, exists := os.LookupEnv("IDP_SCRAPER_HOST")
	if !exists {
		log.Fatal("No host environment variable found.")
	}
	user := os.Getenv("IDP_SCRAPER_USER")
	if !exists {
		log.Fatal("No user environment variable found.")
	}
	pass := os.Getenv("IDP_SCRAPER_PASSWORD")
	if !exists {
		log.Fatal("No password environment variable found.")
	}

	db_connection := "mongodb+srv://" + user + ":" + pass + "@" + host + "/?retryWrites=true&w=majority"
	return db_connection
}

func createClientAndContext(db_connection string) (*mongo.Client, context.Context, context.CancelFunc) {
	client, err := mongo.NewClient(options.Client().ApplyURI(db_connection))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, ctx, cancel
}

func getInterfacesFromProjects(projects map[string]Project) []interface{} {
	interfaceProjects := make([]interface{}, len(projects))
	i := 0
	for _, v := range projects {
		interfaceProjects[i] = v
		i++
	}
	return interfaceProjects
}
