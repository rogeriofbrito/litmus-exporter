package main

import (
	"context"
	"log"
	"os"
	"time"

	litmusextractor "github.com/rogeriofbrito/litmus-exporter/pkg/litmus-extractor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	le := litmusextractor.NewLitmusExtractorDefault(client)

	_, err = le.ChaosExperimentsExtractor(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = le.ChaosExperimentsRunsExtractor(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
