package mongoextractor

import (
	"context"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/project"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewFullMongoExtractor(mongoClient *mongo.Client) *FullMongoExtractor {
	return &FullMongoExtractor{
		mongoClient: mongoClient,
	}
}

type FullMongoExtractor struct {
	mongoClient *mongo.Client
}

func (fme FullMongoExtractor) ProjectsExtractor(ctx context.Context) ([]project.Project, error) {
	mongoCollection := fme.mongoClient.Database("auth").Collection("project")
	cur, err := mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var collections []project.Project

	for cur.Next(ctx) {
		var collection project.Project
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (fme FullMongoExtractor) ChaosExperimentsExtractor(ctx context.Context) ([]chaos_experiment.ChaosExperimentRequest, error) {
	mongoCollection := fme.mongoClient.Database("litmus").Collection("chaosExperiments")
	cur, err := mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var collections []chaos_experiment.ChaosExperimentRequest

	for cur.Next(ctx) {
		var collection chaos_experiment.ChaosExperimentRequest
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}

func (fme FullMongoExtractor) ChaosExperimentsRunsExtractor(ctx context.Context) ([]chaos_experiment_run.ChaosExperimentRun, error) {
	mongoCollection := fme.mongoClient.Database("litmus").Collection("chaosExperimentRuns")
	cur, err := mongoCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var collections []chaos_experiment_run.ChaosExperimentRun

	for cur.Next(ctx) {
		var collection chaos_experiment_run.ChaosExperimentRun
		err := cur.Decode(&collection)
		if err != nil {
			return nil, err
		}
		collections = append(collections, collection)
	}

	return collections, nil
}
