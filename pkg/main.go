package main

import (
	"context"
	"log"
	"os"
	"time"

	mongoextractor "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-extractor"
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

	var me mongoextractor.IMongoExtractor
	mes := os.Getenv("MONGO_EXTRACTOR_STRATEGY")
	switch mes {
	case "FULL":
		me = mongoextractor.NewFullMongoExtractor(client)
	default:
		log.Fatalf("%s strategy not supported", mes)
	}

	_, err = me.ChaosExperimentsExtractor(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	_, err = me.ChaosExperimentsRunsExtractor(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
