package main

import (
	"fmt"
	"log"
	"os"

	"github.com/upstash/vector-go"
)

func main() {
	url := os.Getenv("UPSTASH_VECTOR_REST_URL")
	token := os.Getenv("UPSTASH_VECTOR_REST_TOKEN")

	fmt.Printf("URL: %s\n", url)
	fmt.Printf("Token: %s\n", token)

	if url == "" || token == "" {
		log.Fatal("URL or token is empty")
	}

	// Try to create index
	index := vector.NewIndex(url, token)
	fmt.Printf("Index created: %v\n", index != nil)

	// Try to list namespaces (this should work)
	namespaces, err := index.ListNamespaces()
	if err != nil {
		log.Printf("Error listing namespaces: %v", err)
		return
	}

	fmt.Printf("Namespaces: %v\n", namespaces)

	// Try to upsert some data
	err = index.UpsertData(vector.UpsertData{
		Id:   "test",
		Data: "test data",
	})

	if err != nil {
		log.Printf("Error upserting data: %v", err)
		return
	}

	fmt.Println("Data upserted successfully")
}
