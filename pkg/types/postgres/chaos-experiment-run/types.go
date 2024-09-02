package typespostgreschaosexperimentrun

import (
	"time"

	"github.com/google/uuid"
)

type ChaosExperimentRun struct {
	ID               uuid.UUID                       `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID        string                          `gorm:"column:project_id"`
	UpdatedAt        *time.Time                      `gorm:"column:updated_at"`
	CreatedAt        *time.Time                      `gorm:"column:created_at"`
	CreatedBy        string                          `gorm:"column:created_by"`
	UpdatedBy        string                          `gorm:"column:updated_by"`
	IsRemoved        bool                            `gorm:"column:is_removed"`
	InfraID          string                          `gorm:"column:infra_id"`
	ExperimentRunID  string                          `gorm:"column:experiment_run_id"`
	ExperimentID     string                          `gorm:"column:experiment_id"`
	ExperimentName   string                          `gorm:"column:experiment_name"`
	Phase            string                          `gorm:"column:phase"`
	Probes           []ChaosExperimentRunProbe       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_experiment_run_id"`
	ExecutionDataStr string                          `gorm:"column:execution_data"`
	ExecutionData    ChaosExperimentRunExecutionData `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_experiment_run_id"`
	RevisionID       string                          `gorm:"column:revision_id"`
	NotifyID         *string                         `gorm:"column:notify_id"`
	ResiliencyScore  *float64                        `gorm:"column:resiliency_score"`
	RunSequence      int                             `gorm:"column:run_sequence"`
	Completed        bool                            `gorm:"column:completed"`
	FaultsAwaited    *int                            `gorm:"column:faults_awaited"`
	FaultsFailed     *int                            `gorm:"column:faults_failed"`
	FaultsNa         *int                            `gorm:"column:faults_na"`
	FaultsPassed     *int                            `gorm:"column:faults_passed"`
	FaultsStopped    *int                            `gorm:"column:faults_stopped"`
	TotalFaults      *int                            `gorm:"column:total_faults"`
}

type ChaosExperimentRunProbe struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosExperimentRunID uuid.UUID `gorm:"column:chaos_experiment_run_id"`
	FaultName            string    `gorm:"column:fault_name"`
	ProbeNames           string    `gorm:"column:probe_names"`
}

type ChaosExperimentRunExecutionData struct {
	ID                   uuid.UUID                `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosExperimentRunID uuid.UUID                `gorm:"column:chaos_experiment_run_id"`
	ExperimentType       string                   `gorm:"column:experiment_type"`
	RevisionID           string                   `gorm:"column:revision_id"`
	ExperimentID         string                   `gorm:"column:experiment_id"`
	EventType            string                   `gorm:"column:event_type"`
	UID                  string                   `gorm:"column:uid"`
	Namespace            string                   `gorm:"column:namespace"`
	Name                 string                   `gorm:"column:name"`
	CreationTimestamp    *time.Time               `gorm:"column:creation_timestamp"`
	Phase                string                   `gorm:"column:phase"`
	Message              string                   `gorm:"column:message"`
	StartedAt            *time.Time               `gorm:"column:started_at"`
	FinishedAt           *time.Time               `gorm:"column:finished_at"`
	Nodes                []ChaosExperimentRunNode `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:execution_data_id"`
}

type ChaosExperimentRunNode struct {
	ID                                uuid.UUID                    `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosExperimentRunExecutionDataID uuid.UUID                    `gorm:"column:execution_data_id"`
	NodeName                          string                       `gorm:"column:node_name"`
	Name                              string                       `gorm:"column:name"`
	Phase                             string                       `gorm:"column:phase"`
	Message                           string                       `gorm:"column:message"`
	StartedAt                         *time.Time                   `gorm:"column:started_at"`
	FinishedAt                        *time.Time                   `gorm:"column:finished_at"`
	Children                          string                       `gorm:"column:children"`
	Type                              string                       `gorm:"column:type"`
	ChaosData                         *ChaosExperimentRunChaosData `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:node_id"`
}

type ChaosExperimentRunChaosData struct {
	ID                     uuid.UUID                     `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	NodeID                 uuid.UUID                     `gorm:"column:node_id"`
	EngineUID              string                        `gorm:"column:engine_uid"`
	EngineContext          string                        `gorm:"column:engine_context"`
	EngineName             string                        `gorm:"column:engine_name"`
	Namespace              string                        `gorm:"column:namespace"`
	ExperimentName         string                        `gorm:"column:experiment_name"`
	ExperimentStatus       string                        `gorm:"column:experiment_status"`
	LastUpdatedAt          string                        `gorm:"column:last_updated_at"`
	ExperimentVerdict      string                        `gorm:"column:experiment_verdict"`
	ExperimentPod          string                        `gorm:"column:experiment_pod"`
	RunnerPod              string                        `gorm:"column:runner_pod"`
	ProbeSuccessPercentage string                        `gorm:"column:probe_success_percentage"`
	FailStep               string                        `gorm:"column:fail_step"`
	ChaosResult            ChaosExperimentRunChaosResult `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_data_id"`
}

type ChaosExperimentRunChaosResult struct {
	ID          uuid.UUID                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosDataID uuid.UUID                  `gorm:"column:chaos_data_id"`
	Metadata    ChaosExperimentRunMetadata `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_result_id"`
	Spec        ChaosExperimentRunSpec     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_result_id"`
	Status      ChaosExperimentRunStatus   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_result_id"`
}

type ChaosExperimentRunMetadata struct {
	ID                uuid.UUID                `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosResultID     uuid.UUID                `gorm:"column:chaos_result_id"`
	Name              string                   `gorm:"column:name"`
	Namespace         string                   `gorm:"column:namespace"`
	UID               string                   `gorm:"column:uid"`
	ResourceVersion   string                   `gorm:"column:resource_version"`
	Generation        int64                    `gorm:"column:generation"`
	CreationTimestamp *time.Time               `gorm:"column:creation_timestamp"`
	Labels            ChaosExperimentRunLabels `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:metadata_id"`
}

type ChaosExperimentRunSpec struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosResultID  uuid.UUID `gorm:"column:chaos_result_id"`
	EngineName     string    `gorm:"column:engine_name"`
	ExperimentName string    `gorm:"column:experiment_name"`
}

type ChaosExperimentRunStatus struct {
	ID               uuid.UUID                          `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosResultID    uuid.UUID                          `gorm:"column:chaos_result_id"`
	ExperimentStatus ChaosExperimentRunExperimentStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:status_id"`
	ProbeStatuses    []ChaosExperimentRunProbeStatus    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:status_id"`
	History          ChaosExperimentRunHistory          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:status_id"`
}

type ChaosExperimentRunLabels struct {
	ID                             uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MetadataID                     uuid.UUID `gorm:"column:metadata_id"`
	AppKubernetesIoComponent       string    `gorm:"column:component"`
	AppKubernetesIoPartOf          string    `gorm:"column:part_of"`
	AppKubernetesIoVersion         string    `gorm:"column:version"`
	BatchKubernetesIoControllerUID string    `gorm:"column:controller_uid"`
	BatchKubernetesIoJobName       string    `gorm:"column:job_name"`
	ChaosUID                       string    `gorm:"column:chaos_uid"`
	ControllerUID                  string    `gorm:"column:controller_uid"`
	InfraID                        string    `gorm:"column:infra_id"`
	JobName                        string    `gorm:"column:job_name"`
	Name                           string    `gorm:"column:name"`
	StepPodName                    string    `gorm:"column:step_pod_name"`
	WorkflowName                   string    `gorm:"column:workflow_name"`
	WorkflowRunID                  string    `gorm:"column:workflow_run_id"`
}

type ChaosExperimentRunExperimentStatus struct {
	ID                     uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	StatusID               uuid.UUID `gorm:"column:status_id"`
	Phase                  string    `gorm:"column:phase"`
	Verdict                string    `gorm:"column:verdict"`
	ProbeSuccessPercentage string    `gorm:"column:probe_success_percentage"`
}

type ChaosExperimentRunProbeStatus struct {
	ID       uuid.UUID                             `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	StatusID uuid.UUID                             `gorm:"column:status_id"`
	Name     string                                `gorm:"column:name"`
	Type     string                                `gorm:"column:type"`
	Mode     string                                `gorm:"column:mode"`
	Status   ChaosExperimentRunProbeStatusesStatus `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:probe_status_id"`
}

type ChaosExperimentRunHistory struct {
	ID          uuid.UUID                         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	StatusID    uuid.UUID                         `gorm:"column:status_id"`
	PassedRuns  int                               `gorm:"column:passedRuns"`
	FailedRuns  int                               `gorm:"column:failedRuns"`
	StoppedRuns int                               `gorm:"column:stoppedRuns"`
	Targets     []ChaosExperimentRunHistoryTarget `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:history_id"`
}

type ChaosExperimentRunProbeStatusesStatus struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ProbeStatusID uuid.UUID `gorm:"column:probe_status_id"`
	Verdict       string    `gorm:"column:verdict"`
	Description   string    `gorm:"column:description"`
}

type ChaosExperimentRunHistoryTarget struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	HistoryID   uuid.UUID `gorm:"column:history_id"`
	Name        string    `gorm:"column:name"`
	Kind        string    `gorm:"column:kind"`
	ChaosStatus string    `gorm:"column:chaos_status"`
}
