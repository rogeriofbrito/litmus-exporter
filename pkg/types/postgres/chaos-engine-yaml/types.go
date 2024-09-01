package typespostgreschaosengineyaml

import "github.com/google/uuid"

type ChaosEngineYaml struct {
	ID                   uuid.UUID               `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentRevisionID uuid.UUID               `gorm:"column:revision_id"`
	APIVersion           string                  `gorm:"column:api_version"`
	Kind                 string                  `gorm:"column:kind"`
	Metadata             ChaosEngineYamlMetadata `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_engine_yaml_id"`
	Spec                 ChaosEngineYamlSpec     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_engine_yaml_id"`
}

type ChaosEngineYamlMetadata struct {
	ID                uuid.UUID                  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosEngineYamlID uuid.UUID                  `gorm:"column:chaos_engine_yaml_id"`
	Namespace         string                     `gorm:"column:namespace"`
	Labels            ChaosEngineYamlLabels      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:metadata_id"`
	Annotations       ChaosEngineYamlAnnotations `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:metadata_id"`
	GenerateName      string                     `gorm:"column:generate_name"`
}

type ChaosEngineYamlSpec struct {
	ID                  uuid.UUID                   `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosEngineYamlID   uuid.UUID                   `gorm:"column:chaos_engine_yaml_id"`
	EngineState         string                      `gorm:"column:engine_state"`
	Appinfo             ChaosEngineYamlAppInfo      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:spec_id"`
	ChaosServiceAccount string                      `gorm:"column:chaos_service_account"`
	Experiments         []ChaosEngineYamlExperiment `gorm:"foreignKey:spec_id"`
}

type ChaosEngineYamlLabels struct {
	ID            uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MetadataID    uuid.UUID `gorm:"column:metadata_id"`
	WorkflowRunID string    `gorm:"column:workflow_run_id"`
	WorkflowName  string    `gorm:"column:workflow_name"`
}

type ChaosEngineYamlAnnotations struct {
	ID         uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MetadataID uuid.UUID `gorm:"column:metadata_id"`
	ProbeRef   string    `gorm:"column:probe_ref"`
}

type ChaosEngineYamlAppInfo struct {
	ID       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID   uuid.UUID `gorm:"column:spec_id"`
	Appns    string    `gorm:"column:app_ns"`
	Applabel string    `gorm:"column:app_label"`
	Appkind  string    `gorm:"column:app_kind"`
}

type ChaosEngineYamlExperiment struct {
	ID     uuid.UUID                     `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID uuid.UUID                     `gorm:"column:spec_id"`
	Name   string                        `gorm:"column:name"`
	Spec   ChaosEngineYamlExperimentSpec `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_id"`
}

type ChaosEngineYamlExperimentSpec struct {
	ID           uuid.UUID                 `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentID uuid.UUID                 `gorm:"column:experiment_id"`
	Components   ChaosEngineYamlCompoments `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:experiment_spec_id"`
}

type ChaosEngineYamlCompoments struct {
	ID               uuid.UUID            `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ExperimentSpecID uuid.UUID            `gorm:"column:experiment_spec_id"`
	Env              []ChaosEngineYamlEnv `gorm:"foreignKey:components_id"`
}

type ChaosEngineYamlEnv struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	CompomentsID uuid.UUID `gorm:"column:components_id"`
	Name         string    `gorm:"column:name"`
	Value        string    `gorm:"column:value"`
}
