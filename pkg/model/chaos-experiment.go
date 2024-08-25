package model

import (
	"github.com/google/uuid"
)

type ChaosExperiment struct {
	ID          uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MongoID     string    `gorm:"column:mongo_id"`
	Name        string    `gorm:"column:name"`
	Description string    `gorm:"column:description"`
	Tags        string    `gorm:"column:tags"`
	UpdatedAt   int64     `gorm:"column:updated_at"`
	CreatedAt   int64     `gorm:"column:created_at"`
	//UpdatedBy                  User                              `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:id"`
	IsRemoved                  bool                        `gorm:"column:is_removed"`
	ProjectID                  string                      `gorm:"column:project_id"`
	ExperimentID               string                      `gorm:"column:experiment_id"`
	CronSyntax                 string                      `gorm:"column:cron_syntax"`
	InfraID                    string                      `gorm:"column:infra_id"`
	ExperimentType             string                      `gorm:"column:experiment_type"`
	Revision                   []Revision                  `gorm:"foreignKey:experiment_id"`
	IsCustomExperiment         bool                        `gorm:"column:is_custom_experiment"`
	RecentExperimentRunDetails []RecentExperimentRunDetail `gorm:"foreignKey:experiment_id"`
	TotalExperimentRuns        int                         `gorm:"column:total_experiment_runs"`
}

type User struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID   string    `gorm:"column:user_id"`
	UserName string    `gorm:"column:username"`
	Email    string    `gorm:"column:email"`
}

type Revision struct {
	ID                 uuid.UUID          `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID       uuid.UUID          `gorm:"column:experiment_id"`
	RevisionID         string             `gorm:"column:revision_id"`
	ExperimentManifest ExperimentManifest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_revision_id"`
}

type ExperimentManifest struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentRevisionID uuid.UUID `gorm:"column:experiment_revision_id"`
	Kind                 string    `gorm:"column:kind"`
	APIVersion           string    `gorm:"column:api_version"`
	Metadata             Metadata  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
	Spec                 Spec      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
	Status               Status    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
}

type Metadata struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID `gorm:"column:experiment_manifest_id"`
	Name                 string    `gorm:"column:name"`
	CreationTimestamp    int64     `gorm:"column:creation_timestamp"`
	Labels               Labels    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_metadata_id"`
}

type Labels struct {
	ID                           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestMetadataID uuid.UUID `gorm:"column:experiment_manifest_metadata_id"`
	InfraID                      string    `gorm:"column:infra_id"`
	RevisionID                   string    `gorm:"column:revision_id"`
	WorkflowID                   string    `gorm:"column:workflow_id"`
	ControllerInstanceID         string    `gorm:"column:controller_instance_id"`
}

type Spec struct {
	ID                   uuid.UUID       `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID       `gorm:"column:experiment_manifest_id"`
	Templates            []Template      `gorm:"foreignKey:experiment_manifest_spec_id"`
	Entrypoint           string          `gorm:"column:entrypoint"`
	Arguments            Arguments       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
	ServiceAccountName   string          `gorm:"column:serviceAccountName"`
	PodGC                PodGC           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
	SecurityContext      SecurityContext `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
}

type Template struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID `gorm:"column:experiment_manifest_spec_id"`
	Name                     string    `gorm:"column:name"`
	Steps                    Steps     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_template_id"`
	Container                Container `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_template_id"`
}

type Steps struct {
	ID                               uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecTemplateID uuid.UUID `gorm:"column:experiment_manifest_spec_template_id"`
	Name                             string    `gorm:"column:name"`
	Template                         string    `gorm:"column:template"`
}

type Container struct {
	ID                               uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecTemplateID uuid.UUID `gorm:"column:experiment_manifest_spec_template_id"`
	Name                             string    `gorm:"column:name"`
	Image                            string    `gorm:"column:image"`
	Command                          string    `gorm:"column:command"`
	Args                             string    `gorm:"column:args"`
}

type Arguments struct {
	ID                       uuid.UUID   `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID   `gorm:"column:experiment_manifest_spec_id"`
	Parameters               []Parameter `gorm:"foreignKey:experiment_manifest_spec_arguments_id"`
}

type Parameter struct {
	ID                                uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecArgumentsID uuid.UUID `gorm:"column:experiment_manifest_spec_arguments_id"`
	Name                              string    `gorm:"column:name"`
	Value                             string    `gorm:"column:value"`
}

type PodGC struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID `gorm:"column:experiment_manifest_spec_id"`
	Strategy                 string    `gorm:"column:strategy"`
}

type SecurityContext struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID `gorm:"column:experiment_manifest_spec_id"`
	RunAsUser                int       `gorm:"column:run_as_user"`
	RunAsNonRoot             bool      `gorm:"column:run_as_non_root"`
}

type Status struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID `gorm:"column:experiment_manifest_id"`
	StartedAt            int64     `gorm:"column:started_at"`
	FinishedAt           int64     `gorm:"column:finished_at"`
}

type RecentExperimentRunDetail struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID uuid.UUID `gorm:"column:experiment_id"`
	UpdatedAt    int64     `gorm:"column:updated_at"`
	CreatedAt    int64     `gorm:"column:created_at"`
	//CreatedBy       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:id"`
	//UpdatedBy       User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:id"`
	IsRemoved       bool    `gorm:"column:is_removed"`
	ProjectID       string  `gorm:"column:project_id"`
	ExperimentRunID string  `gorm:"column:experiment_run_id"`
	Phase           string  `gorm:"column:phase"`
	NotifyID        string  `gorm:"column:notify_id"`
	Completed       bool    `gorm:"column:completed"`
	RunSequence     int     `gorm:"column:run_sequence"`
	Probes          []Probe `gorm:"foreignKey:experiment_recent_run_details_id"`
	ResiliencyScore int     `gorm:"column:resiliency_score"`
}

type Probe struct {
	ID                           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentRecentRunDetailsID uuid.UUID `gorm:"column:experiment_recent_run_details_id"`
	FaultName                    string    `gorm:"column:fault_name"`
	ProbeNames                   string    `gorm:"column:probe_names"`
}
