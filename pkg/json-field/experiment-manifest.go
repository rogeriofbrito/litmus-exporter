package jsonfield

type ExperimentManifest struct {
	Kind       string   `json:"kind"`
	APIVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}

type Metadata struct {
	Name              string `json:"name"`
	Namespace         string `json:"namespace"`
	CreationTimestamp int64  `json:"creationTimestamp"`
	Labels            Labels `json:"labels"`
}

type Labels struct {
	InfraID                                 string `json:"infra_id"`
	RevisionID                              string `json:"revision_id"`
	WorkflowID                              string `json:"workflow_id"`
	WorkflowsArgoprojIoControllerInstanceid string `json:"workflows.argoproj.io/controller-instanceid"`
}

type Spec struct {
	Templates          []Template      `json:"templates"`
	Entrypoint         string          `json:"entrypoint"`
	Arguments          Arguments       `json:"arguments"`
	ServiceAccountName string          `json:"serviceAccountName"`
	PodGC              PodGC           `json:"podGC"`
	SecurityContext    SecurityContext `json:"securityContext"`
}

type Template struct {
	Name      string    `json:"name"`
	Inputs    struct{}  `json:"inputs"`
	Outputs   struct{}  `json:"outputs"`
	Metadata  struct{}  `json:"metadata"`
	Steps     Steps     `json:"steps,omitempty"`
	Container Container `json:"container,omitempty"`
}

type Status struct {
	StartedAt  int64 `json:"startedAt"`
	FinishedAt int64 `json:"finishedAt"`
}

type Steps [][]struct {
	Name      string   `json:"name"`
	Template  string   `json:"template"`
	Arguments struct{} `json:"arguments"`
}

type Container struct {
	Name      string   `json:"name"`
	Image     string   `json:"image"`
	Command   []string `json:"command"`
	Args      []string `json:"args"`
	Resources struct{} `json:"resources"`
}

type Arguments struct {
	Parameters []Parameter `json:"parameters"`
}

type Parameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PodGC struct {
	Strategy string `json:"strategy"`
}

type SecurityContext struct {
	RunAsUser    int  `json:"runAsUser"`
	RunAsNonRoot bool `json:"runAsNonRoot"`
}
