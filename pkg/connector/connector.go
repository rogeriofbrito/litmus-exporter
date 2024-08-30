package connector

import (
	"context"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/project"
)

type IConnector interface {
	Init(ctx context.Context) error
	SaveProjects(ctx context.Context, projs []project.Project) error
	SaveChaosExperiments(ctx context.Context, ces []chaos_experiment.ChaosExperimentRequest) error
	SaveChaosExperimentRuns(ctx context.Context, cers []chaos_experiment_run.ChaosExperimentRun) error
}
