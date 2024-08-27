package yamlfield

type ChaosEngine struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Metadata   ChaosEngineMetadata `yaml:"metadata"`
	Spec       ChaosEngineSpec     `yaml:"spec"`
}

type ChaosEngineMetadata struct {
	Namespace    string                 `yaml:"namespace"`
	Labels       ChaosEngineLabels      `yaml:"labels"`
	Annotations  ChaosEngineAnnotations `yaml:"annotations"`
	GenerateName string                 `yaml:"generateName"`
}

type ChaosEngineLabels struct {
	WorkflowRunID string `yaml:"workflow_run_id"`
	WorkflowName  string `yaml:"workflow_name"`
}

type ChaosEngineAnnotations struct {
	ProbeRef string `yaml:"probeRef"`
}

type ChaosEngineSpec struct {
	EngineState         string                  `yaml:"engineState"`
	Appinfo             ChaosEngineAppInfo      `yaml:"appinfo"`
	ChaosServiceAccount string                  `yaml:"chaosServiceAccount"`
	Experiments         []ChaosEngineExperiment `yaml:"experiments"`
}

type ChaosEngineAppInfo struct {
	Appns    string `yaml:"appns"`
	Applabel string `yaml:"applabel"`
	Appkind  string `yaml:"appkind"`
}

type ChaosEngineExperiment struct {
	Name string                    `yaml:"name"`
	Spec ChaosEngineExperimentSpec `yaml:"spec"`
}

type ChaosEngineExperimentSpec struct {
	Components ChaosEngineCompoments `yaml:"components"`
}

type ChaosEngineCompoments struct {
	Env []ChaosEngineEnv `yaml:"env"`
}

type ChaosEngineEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}
