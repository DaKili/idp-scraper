package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	projects := []Project{}
	ScanERI(&projects)

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
	client, err := mongo.NewClient(options.Client().ApplyURI(db_connection))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection := client.Database("projects").Collection("project_collection")
	interfaceProjects := make([]interface{}, len(projects))
	for i, p := range projects {
		interfaceProjects[i] = p
	}

	_, err = collection.InsertMany(ctx, interfaceProjects)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("iserted all.")
}
