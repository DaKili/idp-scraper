package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func delProjects() {
	db_connection := getConnectionString()
	client, ctx, cancel := createClientAndContext(db_connection)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("projects").Collection("project_collection")
	collection.DeleteMany(ctx, bson.M{})
}

// Get a map of currently stored projects on the database.
func getProjects() []Project {
	// Get connection context and collection.
	fmt.Println("Getting DB connection")
	db_connection := getConnectionString()
	client, ctx, cancel := createClientAndContext(db_connection)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("projects").Collection("project_collection")

	// Get iteratable cursor.
	fmt.Printf("Retrieving project_collection")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Get and parse projects.
	var retrievedProjects []Project
	fmt.Println("Converting project_collection to []Project")
	for cursor.Next(ctx) {
		var project Project
		if err := cursor.Decode(&project); err != nil {
			log.Fatal(err)
		}
		retrievedProjects = append(retrievedProjects, project)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return retrievedProjects
}

// Insert a map of projects into the database.
func saveProjects(newProjects *Projects) {
	if len(*newProjects) == 0 {
		log.Println("No new projects to save.")
		return
	}

	// Connect.
	db_connection := getConnectionString()
	client, ctx, cancel := createClientAndContext(db_connection)
	defer cancel()
	defer client.Disconnect(ctx)
	collection := client.Database("projects").Collection("project_collection")

	// Add projects.
	interfaceProjects := getInterfacesFromProjects(newProjects)
	_, err := collection.InsertMany(ctx, interfaceProjects)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Inserted %v new projects.", len(*newProjects))
	}
}

// Read environment variables into a connection string.
func getConnectionString() string {
	// Read necessary environment variables and return connection string.
	// @seili lmk to give you access to mongodb - you should need an acc.
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

// Connect to db and return necessary objects to close when leaving the scope.
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

// Convert all projects of a map into an interface array for MongoDB.
func getInterfacesFromProjects(projects *Projects) []interface{} {
	interfaceProjects := make([]interface{}, len(*projects))
	i := 0
	for _, v := range *projects {
		interfaceProjects[i] = v
		i++
	}
	return interfaceProjects
}
