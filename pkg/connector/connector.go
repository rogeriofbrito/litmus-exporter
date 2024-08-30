package connector

import (
	"context"

	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
)

type IConnector interface {
	Init(ctx context.Context) error
	SaveProjects(ctx context.Context, projs []mongocollection.Project) error
	SaveChaosExperiments(ctx context.Context, ces []mongocollection.ChaosExperiment) error
	SaveChaosExperimentRuns(ctx context.Context, cers []mongocollection.ChaosExperimentRun) error
}
