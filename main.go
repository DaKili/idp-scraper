package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	projects := []Project{}
	ScanERI(&projects)
	client, err := mongo.NewClient(options.Client().ApplyURI(""))
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

	// Insert the projects into the MongoDB collection
	_, err = collection.InsertMany(ctx, interfaceProjects)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted successfully")
}
