package util

import (
	"encoding/json"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
)

func ParseExperimentManifests(rev chaos_experiment.ExperimentRevision) (*jsonfield.ExperimentManifest, error) {
	em := &jsonfield.ExperimentManifest{}
	if err := json.Unmarshal([]byte(rev.ExperimentManifest), em); err != nil {
		return nil, err
	}
	return em, nil
}

func ParseExecutionData(cer chaos_experiment_run.ChaosExperimentRun) (*jsonfield.ChaosExperimentRunExecutionData, error) {
	ed := &jsonfield.ChaosExperimentRunExecutionData{}
	if err := json.Unmarshal([]byte(cer.ExecutionData), ed); err != nil {
		return nil, err
	}

	return ed, nil
}
