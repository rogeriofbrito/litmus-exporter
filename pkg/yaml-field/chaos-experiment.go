package yamlfield

type ChaosExperiment struct {
	APIVersion  string      `yaml:"apiVersion"`
	Description Description `yaml:"description"`
	Kind        string      `yaml:"kind"`
	Metadata    Metadata    `yaml:"metadata"`
	Spec        Spec        `yaml:"spec"`
}

type Description struct {
	Message string `yaml:"message"`
}

type Metadata struct {
	Name   string         `yaml:"name"`
	Labels MetadataLabels `yaml:"labels"`
}

type MetadataLabels struct {
	Name                     string `yaml:"name"`
	AppKubernetesIoPartOf    string `yaml:"app.kubernetes.io/part-of"`
	AppKubernetesIoComponent string `yaml:"app.kubernetes.io/component"`
	AppKubernetesIoVersion   string `yaml:"app.kubernetes.io/version"`
}

type Spec struct {
	Definition Definition `yaml:"definition"`
}

type Permission struct {
	APIGroups []string `yaml:"apiGroups"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}

type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Definition struct {
	Scope           string           `yaml:"scope"`
	Permissions     []Permission     `yaml:"permissions"`
	Image           string           `yaml:"image"`
	ImagePullPolicy string           `yaml:"imagePullPolicy"`
	Args            []string         `yaml:"args"`
	Command         []string         `yaml:"command"`
	Env             []Env            `yaml:"env"`
	Labels          DefinitionLabels `yaml:"labels"`
}

type DefinitionLabels struct {
	Name                           string `yaml:"name"`
	AppKubernetesIoPartOf          string `yaml:"app.kubernetes.io/part-of"`
	AppKubernetesIoComponent       string `yaml:"app.kubernetes.io/component"`
	AppKubernetesIoRuntimeAPIUsage string `yaml:"app.kubernetes.io/runtime-api-usage"`
	AppKubernetesIoVersion         string `yaml:"app.kubernetes.io/version"`
}