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
	IsRemoved                  bool                              `gorm:"column:is_removed"`
	ProjectID                  string                            `gorm:"column:project_id"`
	ExperimentID               string                            `gorm:"column:experiment_id"`
	CronSyntax                 string                            `gorm:"column:cron_syntax"`
	InfraID                    string                            `gorm:"column:infra_id"`
	ExperimentType             string                            `gorm:"column:experiment_type"`
	Revision                   []ChaosExperimentRevision         `gorm:"foreignKey:experiment_id"`
	IsCustomExperiment         bool                              `gorm:"column:is_custom_experiment"`
	RecentExperimentRunDetails []ChaosExperimentRecentRunDetails `gorm:"foreignKey:experiment_id"`
	TotalExperimentRuns        int                               `gorm:"column:total_experiment_runs"`
}

type User struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID   string    `gorm:"column:user_id"`
	UserName string    `gorm:"column:username"`
	Email    string    `gorm:"column:email"`
}

type ChaosExperimentRevision struct {
	ID                 uuid.UUID               `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID       uuid.UUID               `gorm:"column:experiment_id"`
	RevisionID         string                  `gorm:"column:revision_id"`
	ExperimentManifest ChaosExperimentManifest `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_revision_id"`
}

type ChaosExperimentManifest struct {
	ID                   uuid.UUID                       `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentRevisionID uuid.UUID                       `gorm:"column:experiment_revision_id"`
	Kind                 string                          `gorm:"column:kind"`
	APIVersion           string                          `gorm:"column:api_version"`
	Metadata             ChaosExperimentManifestMetadata `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
	Spec                 ChaosExperimentManifestSpec     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
	Status               ChaosExperimentManifestStatus   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_id"`
}

type ChaosExperimentManifestMetadata struct {
	ID                   uuid.UUID                             `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID                             `gorm:"column:experiment_manifest_id"`
	Name                 string                                `gorm:"column:name"`
	CreationTimestamp    int64                                 `gorm:"column:creation_timestamp"`
	Labels               ChaosExperimentManifestMetadataLabels `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_metadata_id"`
}

type ChaosExperimentManifestMetadataLabels struct {
	ID                           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestMetadataID uuid.UUID `gorm:"column:experiment_manifest_metadata_id"`
	InfraID                      string    `gorm:"column:infra_id"`
	RevisionID                   string    `gorm:"column:revision_id"`
	WorkflowID                   string    `gorm:"column:workflow_id"`
	ControllerInstanceID         string    `gorm:"column:controller_instance_id"`
}

type ChaosExperimentManifestSpec struct {
	ID                   uuid.UUID                                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID                                  `gorm:"column:experiment_manifest_id"`
	Templates            []ChaosExperimentManifestSpecTemplate      `gorm:"foreignKey:experiment_manifest_spec_id"`
	Entrypoint           string                                     `gorm:"column:entrypoint"`
	Arguments            ChaosExperimentManifestSpecArguments       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
	ServiceAccountName   string                                     `gorm:"column:serviceAccountName"`
	PodGC                ChaosExperimentManifestSpecPodGC           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
	SecurityContext      ChaosExperimentManifestSpecSecurityContext `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_id"`
}

type ChaosExperimentManifestSpecTemplate struct {
	ID                       uuid.UUID                                `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID                                `gorm:"column:experiment_manifest_spec_id"`
	Name                     string                                   `gorm:"column:name"`
	Steps                    ChaosExperimentManifestSpecTemplateSteps `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_template_id"`
	Container                ChaosExperimentManifestSpecContainer     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_manifest_spec_template_id"`
}

type ChaosExperimentManifestSpecTemplateSteps struct {
	ID                               uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecTemplateID uuid.UUID `gorm:"column:experiment_manifest_spec_template_id"`
	Name                             string    `gorm:"column:name"`
	Template                         string    `gorm:"column:template"`
}

type ChaosExperimentManifestSpecContainer struct {
	ID                               uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecTemplateID uuid.UUID `gorm:"column:experiment_manifest_spec_template_id"`
	Name                             string    `gorm:"column:name"`
	Image                            string    `gorm:"column:image"`
	Command                          string    `gorm:"column:command"`
	Args                             string    `gorm:"column:args"`
}

type ChaosExperimentManifestSpecArguments struct {
	ID                       uuid.UUID                                       `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID                                       `gorm:"column:experiment_manifest_spec_id"`
	Parameters               []ChaosExperimentManifestSpecArgumentsParameter `gorm:"foreignKey:experiment_manifest_spec_arguments_id"`
}

type ChaosExperimentManifestSpecArgumentsParameter struct {
	ID                                uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecArgumentsID uuid.UUID `gorm:"column:experiment_manifest_spec_arguments_id"`
	Name                              string    `gorm:"column:name"`
	Value                             string    `gorm:"column:value"`
}

type ChaosExperimentManifestSpecPodGC struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID `gorm:"column:experiment_manifest_spec_id"`
	Strategy                 string    `gorm:"column:strategy"`
}

type ChaosExperimentManifestSpecSecurityContext struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestSpecID uuid.UUID `gorm:"column:experiment_manifest_spec_id"`
	RunAsUser                int       `gorm:"column:run_as_user"`
	RunAsNonRoot             bool      `gorm:"column:run_as_non_root"`
}

type ChaosExperimentManifestStatus struct {
	ID                   uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentManifestID uuid.UUID `gorm:"column:experiment_manifest_id"`
	StartedAt            int64     `gorm:"column:started_at"`
	FinishedAt           int64     `gorm:"column:finished_at"`
}

type ChaosExperimentRecentRunDetails struct {
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
