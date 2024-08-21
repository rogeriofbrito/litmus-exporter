package mongoextractor

import (
	"context"

	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
)

type IMongoExtractor interface {
	ChaosExperimentsExtractor(ctx context.Context) ([]mongocollection.ChaosExperiment, error)
	ChaosExperimentsRunsExtractor(ctx context.Context) ([]mongocollection.ChaosExperimentRun, error)
}
