package jsonfield

type ChaosExperimentManifest struct {
	Kind       string                          `json:"kind"`
	APIVersion string                          `json:"apiVersion"`
	Metadata   ChaosExperimentManifestMetadata `json:"metadata"`
	Spec       ChaosExperimentManifestSpec     `json:"spec"`
	Status     ChaosExperimentManifestStatus   `json:"status"`
}

type ChaosExperimentManifestMetadata struct {
	Name              string                                `json:"name"`
	Namespace         string                                `json:"namespace"`
	CreationTimestamp int64                                 `json:"creationTimestamp"`
	Labels            ChaosExperimentManifestMetadataLabels `json:"labels"`
}

type ChaosExperimentManifestMetadataLabels struct {
	InfraID                                 string `json:"infra_id"`
	RevisionID                              string `json:"revision_id"`
	WorkflowID                              string `json:"workflow_id"`
	WorkflowsArgoprojIoControllerInstanceid string `json:"workflows.argoproj.io/controller-instanceid"`
}

type ChaosExperimentManifestSpec struct {
	Templates          []ChaosExperimentManifestSpecTemplate      `json:"templates"`
	Entrypoint         string                                     `json:"entrypoint"`
	Arguments          ChaosExperimentManifestSpecArguments       `json:"arguments"`
	ServiceAccountName string                                     `json:"serviceAccountName"`
	PodGC              ChaosExperimentManifestSpecPodGC           `json:"podGC"`
	SecurityContext    ChaosExperimentManifestSpecSecurityContext `json:"securityContext"`
}

type ChaosExperimentManifestSpecTemplate struct {
	Name      string                                   `json:"name"`
	Inputs    struct{}                                 `json:"inputs"`
	Outputs   struct{}                                 `json:"outputs"`
	Metadata  struct{}                                 `json:"metadata"`
	Steps     ChaosExperimentManifestSpecTemplateSteps `json:"steps,omitempty"`
	Container ChaosExperimentManifestSpecContainer     `json:"container,omitempty"`
}

type ChaosExperimentManifestStatus struct {
	StartedAt  interface{} `json:"startedAt"`
	FinishedAt interface{} `json:"finishedAt"`
}

type ChaosExperimentManifestSpecTemplateSteps [][]struct {
	Name      string   `json:"name"`
	Template  string   `json:"template"`
	Arguments struct{} `json:"arguments"`
}

type ChaosExperimentManifestSpecContainer struct {
	Name      string   `json:"name"`
	Image     string   `json:"image"`
	Command   []string `json:"command"`
	Args      []string `json:"args"`
	Resources struct{} `json:"resources"`
}

type ChaosExperimentManifestSpecArguments struct {
	Parameters []ChaosExperimentManifestSpecArgumentsParameter `json:"parameters"`
}

type ChaosExperimentManifestSpecArgumentsParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ChaosExperimentManifestSpecPodGC struct {
	Strategy string `json:"strategy"`
}

type ChaosExperimentManifestSpecSecurityContext struct {
	RunAsUser    int  `json:"runAsUser"`
	RunAsNonRoot bool `json:"runAsNonRoot"`
}
