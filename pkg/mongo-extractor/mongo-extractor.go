package mongoextractor

import (
	"context"

	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
)

type IMongoExtractor interface {
	ProjectsExtractor(ctx context.Context) ([]mongocollection.Project, error)
	ChaosExperimentsExtractor(ctx context.Context) ([]mongocollection.ChaosExperiment, error)
	ChaosExperimentsRunsExtractor(ctx context.Context) ([]mongocollection.ChaosExperimentRun, error)
}
