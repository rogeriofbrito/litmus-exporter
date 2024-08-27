package yamlfield

type ChaosExperiment struct {
	APIVersion  string                     `yaml:"apiVersion"`
	Description ChaosExperimentDescription `yaml:"description"`
	Kind        string                     `yaml:"kind"`
	Metadata    ChaosExperimentMetadata    `yaml:"metadata"`
	Spec        ChaosExperimentSpec        `yaml:"spec"`
}

type ChaosExperimentDescription struct {
	Message string `yaml:"message"`
}

type ChaosExperimentMetadata struct {
	Name   string                `yaml:"name"`
	Labels ChaosExperimentLabels `yaml:"labels"`
}

type ChaosExperimentLabels struct {
	Name                     string `yaml:"name"`
	AppKubernetesIoPartOf    string `yaml:"app.kubernetes.io/part-of"`
	AppKubernetesIoComponent string `yaml:"app.kubernetes.io/component"`
	AppKubernetesIoVersion   string `yaml:"app.kubernetes.io/version"`
}

type ChaosExperimentSpec struct {
	Definition ChaosExperimentDefinition `yaml:"definition"`
}

type Permission struct {
	APIGroups []string `yaml:"apiGroups"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}

type ChaosExperimentEnv struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type ChaosExperimentDefinition struct {
	Scope           string                          `yaml:"scope"`
	Permissions     []Permission                    `yaml:"permissions"`
	Image           string                          `yaml:"image"`
	ImagePullPolicy string                          `yaml:"imagePullPolicy"`
	Args            []string                        `yaml:"args"`
	Command         []string                        `yaml:"command"`
	Env             []ChaosExperimentEnv            `yaml:"env"`
	Labels          ChaosExperimentDefinitionLabels `yaml:"labels"`
}

type ChaosExperimentDefinitionLabels struct {
	Name                           string `yaml:"name"`
	AppKubernetesIoPartOf          string `yaml:"app.kubernetes.io/part-of"`
	AppKubernetesIoComponent       string `yaml:"app.kubernetes.io/component"`
	AppKubernetesIoRuntimeAPIUsage string `yaml:"app.kubernetes.io/runtime-api-usage"`
	AppKubernetesIoVersion         string `yaml:"app.kubernetes.io/version"`
}
