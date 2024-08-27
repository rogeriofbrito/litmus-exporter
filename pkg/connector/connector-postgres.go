package connector

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	"github.com/rogeriofbrito/litmus-exporter/pkg/model"
	model_chaos_engine_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-engine-yaml"
	model_chaos_experiment_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment-yaml"
	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
	yamlfield "github.com/rogeriofbrito/litmus-exporter/pkg/yaml-field"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresConnector() *PostgresConnector {
	return &PostgresConnector{}
}

type PostgresConnector struct{}

func (pc PostgresConnector) Init(ctx context.Context) error {
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	err = db.AutoMigrate(
		&model.ChaosExperiment{},
		&model.User{},
		&model.Revision{},
		&model.ExperimentManifest{},
		&model.ManifestMetadata{},
		&model.Labels{},
		&model.ManifestSpec{},
		&model.Template{},
		&model.Steps{},
		&model.Container{},
		&model.Arguments{},
		&model.Parameter{},
		&model.PodGC{},
		&model.SecurityContext{},
		&model.Status{},
		&model.RecentExperimentRunDetail{},
		&model.Probe{},
		//ChaosExperimentYaml
		&model_chaos_experiment_yaml.ChaosExperimentYaml{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlDescription{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlMetadata{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlLabels{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlSpec{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlDefinition{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlPermission{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlEnv{},
		&model_chaos_experiment_yaml.ChaosExperimentYamlDefinitionLabels{},
		//ChaosEngineYaml
		&model_chaos_engine_yaml.ChaosEngineYaml{},
		&model_chaos_engine_yaml.ChaosEngineYamlMetadata{},
		&model_chaos_engine_yaml.ChaosEngineYamlSpec{},
		&model_chaos_engine_yaml.ChaosEngineYamlLabels{},
		&model_chaos_engine_yaml.ChaosEngineYamlAnnotations{},
		&model_chaos_engine_yaml.ChaosEngineYamlAppInfo{},
		&model_chaos_engine_yaml.ChaosEngineYamlExperiment{},
		&model_chaos_engine_yaml.ChaosEngineYamlExperimentSpec{},
		&model_chaos_engine_yaml.ChaosEngineYamlCompoments{},
		&model_chaos_engine_yaml.ChaosEngineYamlEnv{},
	)
	if err != nil {
		return err
	}

	return nil
}

func (pc PostgresConnector) SaveChaosExperiments(ctx context.Context, ces []mongocollection.ChaosExperiment) error {
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	timeConv := func(t int64) *time.Time {
		if t == 0 {
			return nil
		}
		time := time.Unix(t/1000, t%1000)
		return &time
	}

	parametersConv := func(c []jsonfield.Parameter) []model.Parameter {
		var m []model.Parameter
		for _, ci := range c {
			m = append(m, model.Parameter{
				Name:  ci.Name,
				Value: ci.Value,
			})
		}
		return m
	}

	templatesConv := func(c []jsonfield.Template) []model.Template {
		getStepName := func(steps jsonfield.Steps) string {
			if len(steps) == 0 {
				return ""
			}
			return steps[0][0].Name
		}

		getStepTemplate := func(steps jsonfield.Steps) string {
			if len(steps) == 0 {
				return ""
			}
			return steps[0][0].Template
		}

		var m []model.Template
		for _, ci := range c {
			m = append(m, model.Template{
				Name: ci.Name,
				Steps: model.Steps{
					Name:     getStepName(ci.Steps),
					Template: getStepTemplate(ci.Steps),
				},
				Container: model.Container{
					Name:    ci.Container.Name,
					Image:   ci.Container.Image,
					Command: strings.Join(ci.Container.Command, ","),
					Args:    strings.Join(ci.Container.Args, ","),
				},
			})
		}
		return m
	}

	chaosExperimentYamlsConv := func(ce mongocollection.Revision) []model_chaos_experiment_yaml.ChaosExperimentYaml {
		permissionsConv := func(permissions []yamlfield.Permission) []model_chaos_experiment_yaml.ChaosExperimentYamlPermission {
			var m []model_chaos_experiment_yaml.ChaosExperimentYamlPermission
			for _, p := range permissions {
				m = append(m, model_chaos_experiment_yaml.ChaosExperimentYamlPermission{
					APIGroups: strings.Join(p.APIGroups, ","),
					Resources: strings.Join(p.Resources, ","),
					Verbs:     strings.Join(p.Verbs, ","),
				})
			}
			return m
		}

		envConv := func(envs []yamlfield.ChaosExperimentEnv) []model_chaos_experiment_yaml.ChaosExperimentYamlEnv {
			var m []model_chaos_experiment_yaml.ChaosExperimentYamlEnv
			for _, e := range envs {
				m = append(m, model_chaos_experiment_yaml.ChaosExperimentYamlEnv{
					Name:  e.Name,
					Value: e.Value,
				})
			}
			return m
		}

		var mces []model_chaos_experiment_yaml.ChaosExperimentYaml
		for _, t := range ce.ExperimentManifest.Spec.Templates {
			if t.Name == "install-chaos-faults" {
				for _, a := range t.Inputs.Artifacts {
					var ce yamlfield.ChaosExperiment
					err := yaml.Unmarshal([]byte(a.Raw.Data), &ce)
					if err != nil {
						panic(err)
					}
					mces = append(mces, model_chaos_experiment_yaml.ChaosExperimentYaml{
						APIVersion: ce.APIVersion,
						Description: model_chaos_experiment_yaml.ChaosExperimentYamlDescription{
							Message: ce.Description.Message,
						},
						Kind: ce.Kind,
						Metadata: model_chaos_experiment_yaml.ChaosExperimentYamlMetadata{
							Name: ce.Metadata.Name,
							Labels: model_chaos_experiment_yaml.ChaosExperimentYamlLabels{
								Name:                     ce.Metadata.Labels.Name,
								AppKubernetesIoPartOf:    ce.Metadata.Labels.AppKubernetesIoPartOf,
								AppKubernetesIoComponent: ce.Metadata.Labels.AppKubernetesIoComponent,
								AppKubernetesIoVersion:   ce.Metadata.Labels.AppKubernetesIoVersion,
							},
						},
						Spec: model_chaos_experiment_yaml.ChaosExperimentYamlSpec{
							Definition: model_chaos_experiment_yaml.ChaosExperimentYamlDefinition{
								Scope:           ce.Spec.Definition.Scope,
								Permissions:     permissionsConv(ce.Spec.Definition.Permissions),
								Image:           ce.Spec.Definition.Image,
								ImagePullPolicy: ce.Spec.Definition.ImagePullPolicy,
								Args:            strings.Join(ce.Spec.Definition.Args, ","),
								Command:         strings.Join(ce.Spec.Definition.Command, ","),
								Env:             envConv(ce.Spec.Definition.Env),
								Labels: model_chaos_experiment_yaml.ChaosExperimentYamlDefinitionLabels{
									Name:                           ce.Spec.Definition.Labels.Name,
									AppKubernetesIoPartOf:          ce.Spec.Definition.Labels.AppKubernetesIoPartOf,
									AppKubernetesIoComponent:       ce.Spec.Definition.Labels.AppKubernetesIoComponent,
									AppKubernetesIoRuntimeAPIUsage: ce.Spec.Definition.Labels.AppKubernetesIoRuntimeAPIUsage,
									AppKubernetesIoVersion:         ce.Spec.Definition.Labels.AppKubernetesIoVersion,
								},
							},
						},
					})
				}
			}
		}
		return mces
	}

	chaosEngineYamlsConv := func(ce mongocollection.Revision) []model_chaos_engine_yaml.ChaosEngineYaml {
		envConv := func(envs []yamlfield.ChaosEngineEnv) []model_chaos_engine_yaml.ChaosEngineYamlEnv {
			var m []model_chaos_engine_yaml.ChaosEngineYamlEnv
			for _, env := range envs {
				m = append(m, model_chaos_engine_yaml.ChaosEngineYamlEnv{
					Name:  env.Name,
					Value: env.Value,
				})
			}
			return m
		}

		experimentsConv := func(exps []yamlfield.ChaosEngineExperiment) []model_chaos_engine_yaml.ChaosEngineYamlExperiment {
			var m []model_chaos_engine_yaml.ChaosEngineYamlExperiment
			for _, exp := range exps {
				m = append(m, model_chaos_engine_yaml.ChaosEngineYamlExperiment{
					Name: exp.Name,
					Spec: model_chaos_engine_yaml.ChaosEngineYamlExperimentSpec{
						Components: model_chaos_engine_yaml.ChaosEngineYamlCompoments{
							Env: envConv(exp.Spec.Components.Env),
						},
					},
				})
			}
			return m
		}

		var mces []model_chaos_engine_yaml.ChaosEngineYaml
		for _, t := range ce.ExperimentManifest.Spec.Templates {
			if strings.Contains(t.Container.Image, "litmus-checker") {
				for _, a := range t.Inputs.Artifacts {
					var ce yamlfield.ChaosEngine
					err := yaml.Unmarshal([]byte(a.Raw.Data), &ce)
					if err != nil {
						panic(err)
					}
					mces = append(mces, model_chaos_engine_yaml.ChaosEngineYaml{
						APIVersion: ce.APIVersion,
						Kind:       ce.Kind,
						Metadata: model_chaos_engine_yaml.ChaosEngineYamlMetadata{
							Namespace: ce.Metadata.Namespace,
							Labels: model_chaos_engine_yaml.ChaosEngineYamlLabels{
								WorkflowRunID: ce.Metadata.Labels.WorkflowRunID,
								WorkflowName:  ce.Metadata.Labels.WorkflowName,
							},
							Annotations: model_chaos_engine_yaml.ChaosEngineYamlAnnotations{
								ProbeRef: ce.Metadata.Annotations.ProbeRef,
							},
							GenerateName: ce.Metadata.GenerateName,
						},
						Spec: model_chaos_engine_yaml.ChaosEngineYamlSpec{
							EngineState: ce.Spec.EngineState,
							Appinfo: model_chaos_engine_yaml.ChaosEngineYamlAppInfo{
								Appns:    ce.Spec.Appinfo.Appns,
								Applabel: ce.Spec.Appinfo.Applabel,
								Appkind:  ce.Spec.Appinfo.Appkind,
							},
							ChaosServiceAccount: ce.Spec.ChaosServiceAccount,
							Experiments:         experimentsConv(ce.Spec.Experiments),
						},
					})
				}

			}
		}
		return mces
	}

	revisionConv := func(c []mongocollection.Revision) []model.Revision {
		var m []model.Revision
		for _, ci := range c {
			m = append(m, model.Revision{
				RevisionID: ci.RevisionId,
				ExperimentManifest: model.ExperimentManifest{
					Kind:       ci.ExperimentManifest.Kind,
					APIVersion: ci.ExperimentManifest.APIVersion,
					Metadata: model.ManifestMetadata{
						Name:              ci.ExperimentManifest.Metadata.Name,
						CreationTimestamp: timeConv(ci.ExperimentManifest.Metadata.CreationTimestamp),
						Labels: model.Labels{
							InfraID:              ci.ExperimentManifest.Metadata.Labels.InfraID,
							RevisionID:           ci.ExperimentManifest.Metadata.Labels.RevisionID,
							WorkflowID:           ci.ExperimentManifest.Metadata.Labels.WorkflowID,
							ControllerInstanceID: ci.ExperimentManifest.Metadata.Labels.WorkflowsArgoprojIoControllerInstanceid,
						},
					},
					Spec: model.ManifestSpec{
						Templates:  templatesConv(ci.ExperimentManifest.Spec.Templates),
						Entrypoint: ci.ExperimentManifest.Spec.Entrypoint,
						Arguments: model.Arguments{
							Parameters: parametersConv(ci.ExperimentManifest.Spec.Arguments.Parameters),
						},
						ServiceAccountName: ci.ExperimentManifest.Spec.ServiceAccountName,
						PodGC: model.PodGC{
							Strategy: ci.ExperimentManifest.Spec.PodGC.Strategy,
						},
						SecurityContext: model.SecurityContext{
							RunAsUser:    ci.ExperimentManifest.Spec.SecurityContext.RunAsUser,
							RunAsNonRoot: ci.ExperimentManifest.Spec.SecurityContext.RunAsNonRoot,
						},
					},
					Status: model.Status{
						StartedAt:  timeConv(ci.ExperimentManifest.Status.StartedAt),
						FinishedAt: timeConv(ci.ExperimentManifest.Status.FinishedAt),
					},
				},
				ChaosExperimentYamls: chaosExperimentYamlsConv(ci),
				ChaosEngineYamls:     chaosEngineYamlsConv(ci),
			})
		}
		return m
	}

	probesConv := func(c []mongocollection.Probe) []model.Probe {
		var m []model.Probe
		for _, ci := range c {
			m = append(m, model.Probe{
				FaultName:  ci.FaultName,
				ProbeNames: strings.Join(ci.ProbeNames, ","),
			})
		}
		return m
	}

	recentExperimentRunDetailsConv := func(c []mongocollection.RecentExperimentRunDetail) []model.RecentExperimentRunDetail {
		var m []model.RecentExperimentRunDetail
		for _, ci := range c {
			m = append(m, model.RecentExperimentRunDetail{
				UpdatedAt: timeConv(ci.UpdatedAt),
				CreatedAt: timeConv(ci.CreatedAt),
				/*
					CreatedBy: model.User{
						UserID:   ci.CreatedBy.UserID,
						UserName: ci.CreatedBy.UserName,
						Email:    ci.CreatedBy.Email,
					},
					UpdatedBy: model.User{
						UserID:   ci.UpdatedBy.UserID,
						UserName: ci.UpdatedBy.UserName,
						Email:    ci.UpdatedBy.Email,
					},
				*/
				IsRemoved:       ci.IsRemoved,
				ProjectID:       ci.ProjectID,
				ExperimentRunID: ci.ExperimentRunID,
				Phase:           ci.Phase,
				NotifyID:        ci.NotifyID,
				Completed:       ci.Completed,
				RunSequence:     ci.RunSequence,
				Probes:          probesConv(ci.Probes),
				ResiliencyScore: ci.ResiliencyScore,
			})
		}
		return m
	}

	for _, ce := range ces {
		cem := model.ChaosExperiment{
			MongoID:     ce.ID.String(),
			Name:        ce.Name,
			Description: ce.Description,
			Tags:        strings.Join(ce.Tags, ","),
			UpdatedAt:   timeConv(ce.UpdatedAt),
			CreatedAt:   timeConv(ce.CreatedAt),
			/*
				UpdatedBy: model.User{
					UserID:   ce.UpdatedBy.UserID,
					UserName: ce.UpdatedBy.UserName,
					Email:    ce.UpdatedBy.Email,
				},
			*/
			IsRemoved:                  ce.IsRemoved,
			ProjectID:                  ce.ProjectID,
			ExperimentID:               ce.ExperimentID,
			CronSyntax:                 ce.CronSyntax,
			InfraID:                    ce.InfraID,
			ExperimentType:             ce.ExperimentType,
			Revision:                   revisionConv(ce.Revision),
			IsCustomExperiment:         ce.IsCustomExperiment,
			RecentExperimentRunDetails: recentExperimentRunDetailsConv(ce.RecentExperimentRunDetails),
			TotalExperimentRuns:        ce.TotalExperimentRuns,
		}
		db.Save(&cem)
	}

	return nil
}

func (pc PostgresConnector) SaveChaosExperimentRuns(ctx context.Context, cers []mongocollection.ChaosExperimentRun) error {
	return nil
}

func (pc PostgresConnector) getGormDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DATABASE_NAME"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSL_MODE"))
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
