package mongoextractor

import (
	"context"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/project"
)

type IMongoExtractor interface {
	ProjectsExtractor(ctx context.Context) ([]project.Project, error)
	ChaosExperimentsExtractor(ctx context.Context) ([]chaos_experiment.ChaosExperimentRequest, error)
	ChaosExperimentsRunsExtractor(ctx context.Context) ([]chaos_experiment_run.ChaosExperimentRun, error)
}
