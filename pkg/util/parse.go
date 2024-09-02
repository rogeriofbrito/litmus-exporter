package util

import (
	"encoding/json"

	typeslitmusk8s "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	typeslitmuschaosexperimentrun "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/chaos_experiment_run"
	typesmongodbchaosexperiment "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	typesmongodbchaosexperimentrun "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	typesargoworkflows "github.com/rogeriofbrito/litmus-exporter/pkg/types/argo-workflows"
	"gopkg.in/yaml.v3"
)

func ParseExperimentManifests(rev typesmongodbchaosexperiment.ExperimentRevision) (*typesargoworkflows.Workflow, error) {
	w := &typesargoworkflows.Workflow{}
	if err := json.Unmarshal([]byte(rev.ExperimentManifest), w); err != nil {
		return nil, err
	}
	return w, nil
}

func ParseExecutionData(cer typesmongodbchaosexperimentrun.ChaosExperimentRun) (*typeslitmuschaosexperimentrun.ExecutionData, error) {
	ed := &typeslitmuschaosexperimentrun.ExecutionData{}
	if err := json.Unmarshal([]byte(cer.ExecutionData), ed); err != nil {
		return nil, err
	}

	return ed, nil
}

func ParseChaosExperimentYaml(yamlStr string) (*typeslitmusk8s.ChaosExperiment, error) {
	yamlData := []byte(yamlStr)
	jsonData, err := yamlToJson(yamlData)
	if err != nil {
		return nil, err
	}
	var ce typeslitmusk8s.ChaosExperiment
	err = json.Unmarshal(jsonData, &ce)
	if err != nil {
		panic(err)
	}
	return &ce, nil
}

func ParseChaosEngineYaml(yamlStr string) (*typeslitmusk8s.ChaosEngine, error) {
	yamlData := []byte(yamlStr)
	jsonData, err := yamlToJson(yamlData)
	if err != nil {
		return nil, err
	}
	var ce typeslitmusk8s.ChaosEngine
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
