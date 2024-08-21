package litmusextractor

import (
	"context"

	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLitmusExtractorDefault(mongoClient *mongo.Client) *LitmusExtractorDefault {
	return &LitmusExtractorDefault{
		mongoClient: mongoClient,
	}
}

type LitmusExtractorDefault struct {
	mongoClient *mongo.Client
}

func (le LitmusExtractorDefault) ChaosExperimentsExtractor(ctx context.Context) ([]mongocollection.ChaosExperiment, error) {
	mongoCollection := le.mongoClient.Database("litmus").Collection("chaosExperiments")
	cur, err := mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var collections []mongocollection.ChaosExperiment

	for cur.Next(ctx) {
		var collection mongocollection.ChaosExperiment
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		err = collection.ParseExperimentManifests()
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (le LitmusExtractorDefault) ChaosExperimentsRunsExtractor(ctx context.Context) ([]mongocollection.ChaosExperimentRun, error) {
	mongoCollection := le.mongoClient.Database("litmus").Collection("chaosExperimentRuns")
	cur, err := mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var collections []mongocollection.ChaosExperimentRun

	for cur.Next(ctx) {
		var collection mongocollection.ChaosExperimentRun
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		err = collection.ParseExecutionData()
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
