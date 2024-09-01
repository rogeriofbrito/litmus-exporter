package argoworkflowstypes

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Workflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              Spec   `json:"spec"`
	Status            Status `json:"status"`
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
	Name      string           `json:"name"`
	Inputs    Inputs           `json:"inputs"`
	Outputs   Outputs          `json:"outputs"`
	Metadata  TemplateMetadata `json:"metadata"`
	Steps     Steps            `json:"steps,omitempty"`
	Container Container        `json:"container,omitempty"`
}

type Inputs struct {
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Raw  Raw    `json:"raw"`
}

type Raw struct {
	Data string `json:"data"`
}

type Outputs struct {
}

type TemplateMetadata struct {
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
