package yamlfield

type ChaosEngine struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Metadata   ChaosEngineMetadata `yaml:"metadata"`
	Spec       ChaosEngineSpec     `yaml:"spec"`
}

type ChaosEngineMetadata struct {
	Namespace    string                         `yaml:"namespace"`
	Labels       ChaosEngineMetadataLabels      `yaml:"labels"`
	Annotations  ChaosEngineMetadataAnnotations `yaml:"annotations"`
	GenerateName string                         `yaml:"generateName"`
}

type ChaosEngineMetadataLabels struct {
	WorkflowRunID string `yaml:"workflow_run_id"`
	WorkflowName  string `yaml:"workflow_name"`
}

type ChaosEngineMetadataAnnotations struct {
	ProbeRef string `yaml:"probeRef"`
}

type ChaosEngineSpec struct {
	EngineState         string                      `yaml:"engineState"`
	Appinfo             ChaosEngineSpecAppInfo      `yaml:"appinfo"`
	ChaosServiceAccount string                      `yaml:"chaosServiceAccount"`
	Experiments         []ChaosEngineSpecExperiment `yaml:"experiments"`
}

type ChaosEngineSpecAppInfo struct {
	Appns    string `yaml:"appns"`
	Applabel string `yaml:"applabel"`
	Appkind  string `yaml:"appkind"`
}

type ChaosEngineSpecExperiment struct {
	Name string                        `yaml:"name"`
	Spec ChaosEngineSpecExperimentSpec `yaml:"spec"`
}

type ChaosEngineSpecExperimentSpec struct {
	Components ChaosEngineSpecExperimentSpecCompoments `yaml:"components"`
}

type ChaosEngineSpecExperimentSpecCompoments struct {
	Env []ChaosEngineSpecExperimentSpecCompomentsEnv `yaml:"env"`
}

type ChaosEngineSpecExperimentSpecCompomentsEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
