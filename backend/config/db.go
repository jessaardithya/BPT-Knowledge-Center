package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/joho/godotenv"
)

var Cluster *gocb.Cluster

// ConnectDB establishes the connection to Couchbase using environment variables
func ConnectDB() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment variables")
	}

	// Retrieve credentials
	connString := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	if connString == "" || username == "" || password == "" {
		log.Fatal("Error: Database credentials (DB_HOST, DB_USERNAME, DB_PASSWORD) are missing in .env")
	}

	// Configure Cluster Options
	opts := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		// Recommended for multi-node production setups to handle topology changes
		SecurityConfig: gocb.SecurityConfig{
			TLSSkipVerify: true, // Set to false if you are using verified certificates
		},
	}

	// Initialize Connection
	var err error
	Cluster, err = gocb.Connect(connString, opts)
	if err != nil {
		log.Fatalf("Critical: Could not initialize Couchbase connection: %v", err)
	}

	// Verify Connection (Wait up to 10 seconds for remote clusters)
	if err = Cluster.WaitUntilReady(10*time.Second, nil); err != nil {
		log.Fatalf("Critical: Couchbase cluster is not reachable: %v", err)
	}

	fmt.Printf("Connected to Couchbase Cluster at %s\n", connString)
}

// GetCollection returns a handle to the specific collection
func GetCollection(bucketName, scopeName, collectionName string) *gocb.Collection {
	if Cluster == nil {
		log.Fatal("Error: Cluster is not initialized. Call ConnectDB() first.")
	}

	bucket := Cluster.Bucket(bucketName)
	err := bucket.WaitUntilReady(5*time.Second, nil)
	if err != nil {
		log.Printf("Warning: Bucket '%s' might not be ready: %v", bucketName, err)
	}

	return bucket.Scope(scopeName).Collection(collectionName)
}
