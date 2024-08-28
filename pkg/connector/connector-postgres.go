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
	"github.com/rogeriofbrito/litmus-exporter/pkg/util"
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
		&model.ChaosExperimentRevision{},
		&model.ChaosExperimentManifest{},
		&model.ChaosExperimentMetadata{},
		&model.ChaosExperimentLabels{},
		&model.ChaosExperimentSpec{},
		&model.ChaosExperimentTemplate{},
		&model.ChaosExperimentSteps{},
		&model.ChaosExperimentContainer{},
		&model.ChaosExperimentArguments{},
		&model.ChaosExperimentParameter{},
		&model.ChaosExperimentPodGC{},
		&model.ChaosExperimentSecurityContext{},
		&model.ChaosExperimentStatus{},
		&model.ChaosExperimentRecentExperimentRunDetail{},
		&model.ChaosExperimentProbe{},
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

	cems := util.SliceMap(ces, func(ce mongocollection.ChaosExperiment) model.ChaosExperiment {
		return model.ChaosExperiment{
			MongoID:     ce.ID.String(),
			Name:        ce.Name,
			Description: ce.Description,
			Tags:        strings.Join(ce.Tags, ","),
			UpdatedAt:   pc.getTime(ce.UpdatedAt),
			CreatedAt:   pc.getTime(ce.CreatedAt),
			/*UpdatedBy: model.User{
				UserID:   ce.UpdatedBy.UserID,
				UserName: ce.UpdatedBy.UserName,
				Email:    ce.UpdatedBy.Email,
			},*/
			IsRemoved:      ce.IsRemoved,
			ProjectID:      ce.ProjectID,
			ExperimentID:   ce.ExperimentID,
			CronSyntax:     ce.CronSyntax,
			InfraID:        ce.InfraID,
			ExperimentType: ce.ExperimentType,
			Revision: util.SliceMap(ce.Revision, func(rev mongocollection.Revision) model.ChaosExperimentRevision {
				return model.ChaosExperimentRevision{
					RevisionID: rev.RevisionId,
					ExperimentManifest: model.ChaosExperimentManifest{
						Kind:       rev.ExperimentManifest.Kind,
						APIVersion: rev.ExperimentManifest.APIVersion,
						Metadata: model.ChaosExperimentMetadata{
							Name:              rev.ExperimentManifest.Metadata.Name,
							CreationTimestamp: pc.getTime(rev.ExperimentManifest.Metadata.CreationTimestamp),
							Labels: model.ChaosExperimentLabels{
								InfraID:              rev.ExperimentManifest.Metadata.Labels.InfraID,
								RevisionID:           rev.ExperimentManifest.Metadata.Labels.RevisionID,
								WorkflowID:           rev.ExperimentManifest.Metadata.Labels.WorkflowID,
								ControllerInstanceID: rev.ExperimentManifest.Metadata.Labels.WorkflowsArgoprojIoControllerInstanceid,
							},
						},
						Spec: model.ChaosExperimentSpec{
							Templates: util.SliceMap(rev.ExperimentManifest.Spec.Templates, func(temp jsonfield.Template) model.ChaosExperimentTemplate {
								return model.ChaosExperimentTemplate{
									Name: temp.Name,
									Steps: model.ChaosExperimentSteps{
										Name: func(steps jsonfield.Steps) string {
											if len(steps) == 0 {
												return ""
											}
											return steps[0][0].Name
										}(temp.Steps),
										Template: func(steps jsonfield.Steps) string {
											if len(steps) == 0 {
												return ""
											}
											return steps[0][0].Template
										}(temp.Steps),
									},
									Container: model.ChaosExperimentContainer{
										Name:    temp.Container.Name,
										Image:   temp.Container.Image,
										Command: strings.Join(temp.Container.Command, ","),
										Args:    strings.Join(temp.Container.Args, ","),
									},
								}
							}),
							Entrypoint: rev.ExperimentManifest.Spec.Entrypoint,
							Arguments: model.ChaosExperimentArguments{
								Parameters: util.SliceMap(rev.ExperimentManifest.Spec.Arguments.Parameters, func(param jsonfield.Parameter) model.ChaosExperimentParameter {
									return model.ChaosExperimentParameter{
										Name:  param.Name,
										Value: param.Value,
									}
								}),
							},
							ServiceAccountName: rev.ExperimentManifest.Spec.ServiceAccountName,
							PodGC: model.ChaosExperimentPodGC{
								Strategy: rev.ExperimentManifest.Spec.PodGC.Strategy,
							},
							SecurityContext: model.ChaosExperimentSecurityContext{
								RunAsUser:    rev.ExperimentManifest.Spec.SecurityContext.RunAsUser,
								RunAsNonRoot: rev.ExperimentManifest.Spec.SecurityContext.RunAsNonRoot,
							},
						},
						Status: model.ChaosExperimentStatus{
							StartedAt:  pc.getTime(rev.ExperimentManifest.Status.StartedAt),
							FinishedAt: pc.getTime(rev.ExperimentManifest.Status.FinishedAt),
						},
					},
					ChaosExperimentYamls: pc.getChaosExperimentYamls(rev),
					ChaosEngineYamls:     pc.getChaosEngineYamls(rev),
				}
			}),
			IsCustomExperiment: ce.IsCustomExperiment,
			RecentExperimentRunDetails: util.SliceMap(ce.RecentExperimentRunDetails, func(detail mongocollection.RecentExperimentRunDetail) model.ChaosExperimentRecentExperimentRunDetail {
				return model.ChaosExperimentRecentExperimentRunDetail{
					UpdatedAt: pc.getTime(detail.UpdatedAt),
					CreatedAt: pc.getTime(detail.CreatedAt),
					/*CreatedBy: model.User{
						UserID:   ci.CreatedBy.UserID,
						UserName: ci.CreatedBy.UserName,
						Email:    ci.CreatedBy.Email,
					},
					UpdatedBy: model.User{
						UserID:   ci.UpdatedBy.UserID,
						UserName: ci.UpdatedBy.UserName,
						Email:    ci.UpdatedBy.Email,
					},*/
					IsRemoved:       detail.IsRemoved,
					ProjectID:       detail.ProjectID,
					ExperimentRunID: detail.ExperimentRunID,
					Phase:           detail.Phase,
					NotifyID:        detail.NotifyID,
					Completed:       detail.Completed,
					RunSequence:     detail.RunSequence,
					Probes: util.SliceMap(detail.Probes, func(probe mongocollection.Probe) model.ChaosExperimentProbe {
						return model.ChaosExperimentProbe{
							FaultName:  probe.FaultName,
							ProbeNames: strings.Join(probe.ProbeNames, ","),
						}
					}),
					ResiliencyScore: detail.ResiliencyScore,
				}
			}),
			TotalExperimentRuns: ce.TotalExperimentRuns,
		}
	})

	for _, cem := range cems {
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

func (pc PostgresConnector) getTime(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	time := time.Unix(t/1000, t%1000)
	return &time
}

func (pc PostgresConnector) getChaosExperimentYamls(ce mongocollection.Revision) []model_chaos_experiment_yaml.ChaosExperimentYaml {
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
							Scope: ce.Spec.Definition.Scope,
							Permissions: util.SliceMap(ce.Spec.Definition.Permissions, func(perm yamlfield.Permission) model_chaos_experiment_yaml.ChaosExperimentYamlPermission {
								return model_chaos_experiment_yaml.ChaosExperimentYamlPermission{
									APIGroups: strings.Join(perm.APIGroups, ","),
									Resources: strings.Join(perm.Resources, ","),
									Verbs:     strings.Join(perm.Verbs, ","),
								}
							}),
							Image:           ce.Spec.Definition.Image,
							ImagePullPolicy: ce.Spec.Definition.ImagePullPolicy,
							Args:            strings.Join(ce.Spec.Definition.Args, ","),
							Command:         strings.Join(ce.Spec.Definition.Command, ","),
							Env: util.SliceMap(ce.Spec.Definition.Env, func(env yamlfield.ChaosExperimentEnv) model_chaos_experiment_yaml.ChaosExperimentYamlEnv {
								return model_chaos_experiment_yaml.ChaosExperimentYamlEnv{
									Name:  env.Name,
									Value: env.Value,
								}
							}),
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

func (pc PostgresConnector) getChaosEngineYamls(ce mongocollection.Revision) []model_chaos_engine_yaml.ChaosEngineYaml {
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
						Experiments: util.SliceMap(ce.Spec.Experiments, func(exp yamlfield.ChaosEngineExperiment) model_chaos_engine_yaml.ChaosEngineYamlExperiment {
							return model_chaos_engine_yaml.ChaosEngineYamlExperiment{
								Name: exp.Name,
								Spec: model_chaos_engine_yaml.ChaosEngineYamlExperimentSpec{
									Components: model_chaos_engine_yaml.ChaosEngineYamlCompoments{
										Env: util.SliceMap(exp.Spec.Components.Env, func(env yamlfield.ChaosEngineEnv) model_chaos_engine_yaml.ChaosEngineYamlEnv {
											return model_chaos_engine_yaml.ChaosEngineYamlEnv{
												Name:  env.Name,
												Value: env.Value,
											}
										}),
									},
								},
							}
						}),
					},
				})
			}
		}
	}
	return mces
}
