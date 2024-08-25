package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rogeriofbrito/litmus-exporter/pkg/connector"
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

	var c connector.IConnector
	ct := os.Getenv("CONNECTOR_TYPE")
	switch ct {
	case "POSTGRES":
		c = connector.NewPostgresConnector()
	default:
		log.Fatalf("%s connector type not supported", ct)
	}

	err = c.Init(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx = context.Background()

	ce, err := me.ChaosExperimentsExtractor(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = c.SaveChaosExperiments(ctx, ce)
	if err != nil {
		log.Fatal(err)
	}

	cer, err := me.ChaosExperimentsRunsExtractor(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = c.SaveChaosExperimentRuns(ctx, cer)
	if err != nil {
		log.Fatal(err)
	}
}
