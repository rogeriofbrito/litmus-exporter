package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/rogeriofbrito/litmus-exporter/pkg/connector"
	mongoextractor "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-extractor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db, err := getGormDB()
	if err != nil {
		log.Fatal(err)
	}

	var c connector.IConnector
	ct := os.Getenv("CONNECTOR_TYPE")
	switch ct {
	case "POSTGRES":
		c = connector.NewPostgresConnector(db)
	default:
		log.Fatalf("%s connector type not supported", ct)
	}

	ctx = context.Background()

	ctx, err = c.InitCtx(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Begin(ctx)
	if err != nil {
		log.Fatal(err)
	}

	p, err := me.ProjectsExtractor(ctx)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}
	err = c.SaveProjects(ctx, p)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}

	ce, err := me.ChaosExperimentsExtractor(ctx)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}
	err = c.SaveChaosExperiments(ctx, ce)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}

	cer, err := me.ChaosExperimentsRunsExtractor(ctx)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}
	err = c.SaveChaosExperimentRuns(ctx, cer)
	if err != nil {
		_ = c.Rollback(ctx)
		log.Fatal(err)
	}

	err = c.Commit(ctx)
	if err != nil {
		log.Fatal(err)
	}

}

func getGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE_NAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSL_MODE"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
}
