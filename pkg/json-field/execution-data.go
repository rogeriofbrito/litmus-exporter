package jsonfield

type ChaosExperimentRunExecutionData struct {
	ExperimentType    string                            `json:"experimentType"`
	RevisionID        string                            `json:"revisionID"`
	NotifyID          string                            `json:"notifyID"`
	ExperimentID      string                            `json:"experimentID"`
	EventType         string                            `json:"eventType"`
	UID               string                            `json:"uid"`
	Namespace         string                            `json:"namespace"`
	Name              string                            `json:"name"`
	CreationTimestamp string                            `json:"creationTimestamp"`
	Phase             string                            `json:"phase"`
	Message           string                            `json:"message"`
	StartedAt         string                            `json:"startedAt"`
	FinishedAt        string                            `json:"finishedAt"`
	Nodes             map[string]ChaosExperimentRunNode `json:"nodes"`
}

type ChaosExperimentRunNode struct {
	Name       string                      `json:"name"`
	Phase      string                      `json:"phase"`
	Message    string                      `json:"message"`
	StartedAt  string                      `json:"startedAt"`
	FinishedAt string                      `json:"finishedAt"`
	Children   []string                    `json:"children"`
	Type       string                      `json:"type"`
	ChaosData  ChaosExperimentRunChaosData `json:"chaosData"`
}

type ChaosExperimentRunChaosData struct {
	EngineUID              string                        `json:"engineUID"`
	EngineContext          string                        `json:"engineContext"`
	EngineName             string                        `json:"engineName"`
	Namespace              string                        `json:"namespace"`
	ExperimentName         string                        `json:"experimentName"`
	ExperimentStatus       string                        `json:"experimentStatus"`
	LastUpdatedAt          string                        `json:"lastUpdatedAt"`
	ExperimentVerdict      string                        `json:"experimentVerdict"`
	ExperimentPod          string                        `json:"experimentPod"`
	RunnerPod              string                        `json:"runnerPod"`
	ProbeSuccessPercentage string                        `json:"probeSuccessPercentage"`
	FailStep               string                        `json:"failStep"`
	ChaosResult            ChaosExperimentRunChaosResult `json:"chaosResult"`
}

type ChaosExperimentRunChaosResult struct {
	Metadata ChaosExperimentRunMetadata `json:"metadata"`
	Spec     ChaosExperimentRunSpec     `json:"spec"`
	Status   ChaosExperimentRunStatus   `json:"status"`
}

type ChaosExperimentRunMetadata struct {
	Name              string                   `json:"name"`
	Namespace         string                   `json:"namespace"`
	UID               string                   `json:"uid"`
	ResourceVersion   string                   `json:"resourceVersion"`
	Generation        int                      `json:"generation"`
	CreationTimestamp string                   `json:"creationTimestamp"`
	Labels            ChaosExperimentRunLabels `json:"labels"`
}

type ChaosExperimentRunSpec struct {
	Engine     string `json:"engine"`
	Experiment string `json:"experiment"`
}

type ChaosExperimentRunStatus struct {
	ExperimentStatus ChaosExperimentRunExperimentStatus `json:"experimentStatus"`
	ProbeStatuses    []ChaosExperimentRunProbeStatuses  `json:"probeStatuses"`
	History          ChaosExperimentRunHistory          `json:"history"`
}

type ChaosExperimentRunLabels struct {
	AppKubernetesIoComponent       string `json:"app.kubernetes.io/component"`
	AppKubernetesIoPartOf          string `json:"app.kubernetes.io/part-of"`
	AppKubernetesIoVersion         string `json:"app.kubernetes.io/version"`
	BatchKubernetesIoControllerUID string `json:"batch.kubernetes.io/controller-uid"`
	BatchKubernetesIoJobName       string `json:"batch.kubernetes.io/job-name"`
	ChaosUID                       string `json:"chaosUID"`
	ControllerUID                  string `json:"controller-uid"`
	InfraID                        string `json:"infra_id"`
	JobName                        string `json:"job-name"`
	Name                           string `json:"name"`
	StepPodName                    string `json:"step_pod_name"`
	WorkflowName                   string `json:"workflow_name"`
	WorkflowRunID                  string `json:"workflow_run_id"`
}

type ChaosExperimentRunExperimentStatus struct {
	Phase                  string `json:"phase"`
	Verdict                string `json:"verdict"`
	ProbeSuccessPercentage string `json:"probeSuccessPercentage"`
}

type ChaosExperimentRunProbeStatuses struct {
	Name   string                                `json:"name"`
	Type   string                                `json:"type"`
	Mode   string                                `json:"mode"`
	Status ChaosExperimentRunProbeStatusesStatus `json:"status"`
}

type ChaosExperimentRunProbeStatusesStatus struct {
	Verdict     string `json:"verdict"`
	Description string `json:"description"`
}

type ChaosExperimentRunHistory struct {
	PassedRuns  int                               `json:"passedRuns"`
	FailedRuns  int                               `json:"failedRuns"`
	StoppedRuns int                               `json:"stoppedRuns"`
	Targets     []ChaosExperimentRunHistoryTarget `json:"targets"`
}

type ChaosExperimentRunHistoryTarget struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	ChaosStatus string `json:"chaosStatus"`
}
