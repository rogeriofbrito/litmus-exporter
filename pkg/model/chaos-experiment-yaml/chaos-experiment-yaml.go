package model_chaos_experiment_yaml

import "github.com/google/uuid"

type ChaosExperimentYaml struct {
	ID         uuid.UUID                   `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	RevisionID uuid.UUID                   `gorm:"column:revision_id"`
	APIVersion string                      `gorm:"column:apiVersion"`
	Kind       string                      `gorm:"column:kind"`
	Metadata   ChaosExperimentYamlMetadata `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_experiment_yaml_id"`
	Spec       ChaosExperimentYamlSpec     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:chaos_experiment_yaml_id"`
}

type ChaosExperimentYamlMetadata struct {
	ID                    uuid.UUID                 `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosExperimentYamlID uuid.UUID                 `gorm:"column:chaos_experiment_yaml_id"`
	Name                  string                    `gorm:"column:name"`
	Labels                ChaosExperimentYamlLabels `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:metadata_id"`
}

type ChaosExperimentYamlSpec struct {
	ID                    uuid.UUID                     `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ChaosExperimentYamlID uuid.UUID                     `gorm:"column:chaos_experiment_yaml_id"`
	Definition            ChaosExperimentYamlDefinition `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:spec_id"`
}

type ChaosExperimentYamlLabels struct {
	ID                       uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	MetadataID               uuid.UUID `gorm:"column:metadata_id"`
	Name                     string    `gorm:"column:name"`
	AppKubernetesIoPartOf    string    `gorm:"column:part_of"`
	AppKubernetesIoComponent string    `gorm:"column:component"`
	AppKubernetesIoVersion   string    `gorm:"column:version"`
}

type ChaosExperimentYamlDefinition struct {
	ID              uuid.UUID                           `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	SpecID          uuid.UUID                           `gorm:"column:spec_id"`
	Scope           string                              `gorm:"column:scope"`
	Permissions     []ChaosExperimentYamlPermission     `gorm:"foreignKey:definition_id"`
	Image           string                              `gorm:"column:image"`
	ImagePullPolicy string                              `gorm:"column:image_pull_policy"`
	Args            string                              `gorm:"column:args"`
	Command         string                              `gorm:"column:command"`
	Env             []ChaosExperimentYamlEnv            `gorm:"foreignKey:definition_id"`
	Labels          ChaosExperimentYamlDefinitionLabels `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:definition_id"`
}

type ChaosExperimentYamlPermission struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	DefinitionID uuid.UUID `gorm:"column:definition_id"`
	APIGroups    string    `gorm:"column:api_groups"`
	Resources    string    `gorm:"column:resources"`
	Verbs        string    `gorm:"column:verbs"`
}

type ChaosExperimentYamlEnv struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	DefinitionID uuid.UUID `gorm:"definition_id"`
	Name         string    `gorm:"column:name"`
	Value        string    `gorm:"column:value"`
}

type ChaosExperimentYamlDefinitionLabels struct {
	ID                             uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	DefinitionID                   uuid.UUID `gorm:"column:definition_id"`
	Name                           string    `gorm:"column:name"`
	AppKubernetesIoPartOf          string    `gorm:"column:part_of"`
	AppKubernetesIoComponent       string    `gorm:"column:component"`
	AppKubernetesIoRuntimeAPIUsage string    `gorm:"column:runtime_api_usage"`
	AppKubernetesIoVersion         string    `gorm:"column:version"`
}
