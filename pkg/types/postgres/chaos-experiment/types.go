package typespostgreschaosexperiment

import (
	"time"

	"github.com/google/uuid"
	model_chaos_engine_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-engine-yaml"
	model_chaos_experiment_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment-yaml"
)

type ChaosExperiment struct {
	ID                         uuid.UUID                                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	Name                       string                                     `gorm:"column:name"`
	Description                string                                     `gorm:"column:description"`
	Tags                       string                                     `gorm:"column:tags"`
	UpdatedAt                  *time.Time                                 `gorm:"column:updated_at"`
	CreatedAt                  *time.Time                                 `gorm:"column:created_at"`
	CreatedBy                  string                                     `gorm:"column:created_by"`
	UpdatedBy                  string                                     `gorm:"column:updated_by"`
	IsRemoved                  bool                                       `gorm:"column:is_removed"`
	ProjectID                  string                                     `gorm:"column:project_id"`
	ExperimentID               string                                     `gorm:"column:experiment_id"`
	CronSyntax                 string                                     `gorm:"column:cron_syntax"`
	InfraID                    string                                     `gorm:"column:infra_id"`
	ExperimentType             string                                     `gorm:"column:experiment_type"`
	Revision                   []ChaosExperimentRevision                  `gorm:"foreignKey:experiment_id"`
	IsCustomExperiment         bool                                       `gorm:"column:is_custom_experiment"`
	RecentExperimentRunDetails []ChaosExperimentRecentExperimentRunDetail `gorm:"foreignKey:experiment_id"`
	TotalExperimentRuns        int                                        `gorm:"column:total_experiment_runs"`
}

type ChaosExperimentRevision struct {
	ID                   uuid.UUID                                         `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID         uuid.UUID                                         `gorm:"column:experiment_id"`
	RevisionID           string                                            `gorm:"column:revision_id"`
	ExperimentManifest   ChaosExperimentManifest                           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:revision_id"`
	ChaosExperimentYamls []model_chaos_experiment_yaml.ChaosExperimentYaml `gorm:"foreignKey:revision_id"`
	ChaosEngineYamls     []model_chaos_engine_yaml.ChaosEngineYaml         `gorm:"foreignKey:revision_id"`
}

type ChaosExperimentManifest struct {
	ID         uuid.UUID               `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	RevisionID uuid.UUID               `gorm:"column:revision_id"`
	Kind       string                  `gorm:"column:kind"`
	APIVersion string                  `gorm:"column:api_version"`
	Metadata   ChaosExperimentMetadata `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:manifest_id"`
	Spec       ChaosExperimentSpec     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:manifest_id"`
	Status     ChaosExperimentStatus   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:manifest_id"`
}

type ChaosExperimentMetadata struct {
	ID                uuid.UUID             `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ManifestID        uuid.UUID             `gorm:"column:manifest_id"`
	Name              string                `gorm:"column:name"`
	CreationTimestamp *time.Time            `gorm:"column:creation_timestamp"`
	Labels            ChaosExperimentLabels `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:metadata_id"`
}

type ChaosExperimentLabels struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MetadataID           uuid.UUID `gorm:"column:metadata_id"`
	InfraID              string    `gorm:"column:infra_id"`
	RevisionID           string    `gorm:"column:revision_id"`
	WorkflowID           string    `gorm:"column:workflow_id"`
	ControllerInstanceID string    `gorm:"column:controller_instance_id"`
}

type ChaosExperimentSpec struct {
	ID                 uuid.UUID                      `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ManifestID         uuid.UUID                      `gorm:"column:manifest_id"`
	Templates          []ChaosExperimentTemplate      `gorm:"foreignKey:spec_id"`
	Entrypoint         string                         `gorm:"column:entrypoint"`
	Arguments          ChaosExperimentArguments       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:spec_id"`
	ServiceAccountName string                         `gorm:"column:serviceAccountName"`
	PodGC              ChaosExperimentPodGC           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:spec_id"`
	SecurityContext    ChaosExperimentSecurityContext `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:spec_id"`
}

type ChaosExperimentTemplate struct {
	ID        uuid.UUID                `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID    uuid.UUID                `gorm:"column:spec_id"`
	Name      string                   `gorm:"column:name"`
	Steps     ChaosExperimentSteps     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:template_id"`
	Container ChaosExperimentContainer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:template_id"`
}

type ChaosExperimentSteps struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	TemplateID uuid.UUID `gorm:"column:template_id"`
	Name       string    `gorm:"column:name"`
	Template   string    `gorm:"column:template"`
}

type ChaosExperimentContainer struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	TemplateID uuid.UUID `gorm:"column:template_id"`
	Name       string    `gorm:"column:name"`
	Image      string    `gorm:"column:image"`
	Command    string    `gorm:"column:command"`
	Args       string    `gorm:"column:args"`
}

type ChaosExperimentArguments struct {
	ID         uuid.UUID                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID     uuid.UUID                  `gorm:"column:spec_id"`
	Parameters []ChaosExperimentParameter `gorm:"foreignKey:arguments_id"`
}

type ChaosExperimentParameter struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ArgumentsID uuid.UUID `gorm:"column:arguments_id"`
	Name        string    `gorm:"column:name"`
	Value       string    `gorm:"column:value"`
}

type ChaosExperimentPodGC struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID   uuid.UUID `gorm:"column:spec_id"`
	Strategy string    `gorm:"column:strategy"`
}

type ChaosExperimentSecurityContext struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID       uuid.UUID `gorm:"column:spec_id"`
	RunAsUser    int       `gorm:"column:run_as_user"`
	RunAsNonRoot bool      `gorm:"column:run_as_non_root"`
}

type ChaosExperimentStatus struct {
	ID         uuid.UUID  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ManifestID uuid.UUID  `gorm:"column:manifest_id"`
	StartedAt  *time.Time `gorm:"column:started_at"`
	FinishedAt *time.Time `gorm:"column:finished_at"`
}

type ChaosExperimentRecentExperimentRunDetail struct {
	ID              uuid.UUID              `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID    uuid.UUID              `gorm:"column:experiment_id"`
	UpdatedAt       *time.Time             `gorm:"column:updated_at"`
	CreatedAt       *time.Time             `gorm:"column:created_at"`
	CreatedBy       string                 `gorm:"column:created_by"`
	UpdatedBy       string                 `gorm:"column:updated_by"`
	IsRemoved       bool                   `gorm:"column:is_removed"`
	ProjectID       string                 `gorm:"column:project_id"`
	ExperimentRunID string                 `gorm:"column:experiment_run_id"`
	Phase           string                 `gorm:"column:phase"`
	NotifyID        *string                `gorm:"column:notify_id"`
	Completed       bool                   `gorm:"column:completed"`
	RunSequence     int                    `gorm:"column:run_sequence"`
	Probes          []ChaosExperimentProbe `gorm:"foreignKey:recent_run_details_id"`
	ResiliencyScore *float64               `gorm:"column:resiliency_score"`
}

type ChaosExperimentProbe struct {
	ID                 uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	RecentRunDetailsID uuid.UUID `gorm:"column:recent_run_details_id"`
	FaultName          string    `gorm:"column:fault_name"`
	ProbeNames         string    `gorm:"column:probe_names"`
}
