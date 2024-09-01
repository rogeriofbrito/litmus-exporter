package util

import (
	"encoding/json"

	litmus_v1alpha1 "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	litmus_chaos_experiment_run "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/chaos_experiment_run"
	mongodb_chaos_experiment "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	mongodb_chaos_experiment_run "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	argoworkflowstypes "github.com/rogeriofbrito/litmus-exporter/pkg/argo-workflows-types"
	"gopkg.in/yaml.v3"
)

func ParseExperimentManifests(rev mongodb_chaos_experiment.ExperimentRevision) (*argoworkflowstypes.Workflow, error) {
	w := &argoworkflowstypes.Workflow{}
	if err := json.Unmarshal([]byte(rev.ExperimentManifest), w); err != nil {
		return nil, err
	}
	return w, nil
}

func ParseExecutionData(cer mongodb_chaos_experiment_run.ChaosExperimentRun) (*litmus_chaos_experiment_run.ExecutionData, error) {
	ed := &litmus_chaos_experiment_run.ExecutionData{}
	if err := json.Unmarshal([]byte(cer.ExecutionData), ed); err != nil {
		return nil, err
	}

	return ed, nil
}

func ParseChaosExperimentYaml(yamlStr string) (*litmus_v1alpha1.ChaosExperiment, error) {
	yamlData := []byte(yamlStr)
	jsonData, err := yamlToJson(yamlData)
	if err != nil {
		return nil, err
	}
	var ce litmus_v1alpha1.ChaosExperiment
	err = json.Unmarshal(jsonData, &ce)
	if err != nil {
		panic(err)
	}
	return &ce, nil
}

func ParseChaosEngineYaml(yamlStr string) (*litmus_v1alpha1.ChaosEngine, error) {
	yamlData := []byte(yamlStr)
	jsonData, err := yamlToJson(yamlData)
	if err != nil {
		return nil, err
	}
	var ce litmus_v1alpha1.ChaosEngine
	err = json.Unmarshal(jsonData, &ce)
	if err != nil {
		panic(err)
	}
	return &ce, nil
}

func yamlToJson(yamlData []byte) ([]byte, error) {
	var data map[string]interface{}
	err := yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
