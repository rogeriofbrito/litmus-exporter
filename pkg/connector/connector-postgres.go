package connector

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	model_chaos_engine_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-engine-yaml"
	model_chaos_experiment "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment"
	model_chaos_experiment_run "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment-run"
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
		//ChaosExperiment
		&model_chaos_experiment.ChaosExperiment{},
		&model_chaos_experiment.User{},
		&model_chaos_experiment.ChaosExperimentRevision{},
		&model_chaos_experiment.ChaosExperimentManifest{},
		&model_chaos_experiment.ChaosExperimentMetadata{},
		&model_chaos_experiment.ChaosExperimentLabels{},
		&model_chaos_experiment.ChaosExperimentSpec{},
		&model_chaos_experiment.ChaosExperimentTemplate{},
		&model_chaos_experiment.ChaosExperimentSteps{},
		&model_chaos_experiment.ChaosExperimentContainer{},
		&model_chaos_experiment.ChaosExperimentArguments{},
		&model_chaos_experiment.ChaosExperimentParameter{},
		&model_chaos_experiment.ChaosExperimentPodGC{},
		&model_chaos_experiment.ChaosExperimentSecurityContext{},
		&model_chaos_experiment.ChaosExperimentStatus{},
		&model_chaos_experiment.ChaosExperimentRecentExperimentRunDetail{},
		&model_chaos_experiment.ChaosExperimentProbe{},
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
		//ChaosExperimentRun
		&model_chaos_experiment_run.ChaosExperimentRun{},
		&model_chaos_experiment_run.ChaosExperimentRunProbe{},
		&model_chaos_experiment_run.ChaosExperimentRunExecutionData{},
		&model_chaos_experiment_run.ChaosExperimentRunNode{},
		&model_chaos_experiment_run.ChaosExperimentRunChaosData{},
		&model_chaos_experiment_run.ChaosExperimentRunChaosResult{},
		&model_chaos_experiment_run.ChaosExperimentRunMetadata{},
		&model_chaos_experiment_run.ChaosExperimentRunSpec{},
		&model_chaos_experiment_run.ChaosExperimentRunStatus{},
		&model_chaos_experiment_run.ChaosExperimentRunLabels{},
		&model_chaos_experiment_run.ChaosExperimentRunExperimentStatus{},
		&model_chaos_experiment_run.ChaosExperimentRunProbeStatus{},
		&model_chaos_experiment_run.ChaosExperimentRunHistory{},
		&model_chaos_experiment_run.ChaosExperimentRunProbeStatusesStatus{},
		&model_chaos_experiment_run.ChaosExperimentRunHistoryTarget{},
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

	cems := util.SliceMap(ces, func(ce mongocollection.ChaosExperiment) model_chaos_experiment.ChaosExperiment {
		return model_chaos_experiment.ChaosExperiment{
			MongoID:     ce.ID.String(),
			Name:        ce.Name,
			Description: ce.Description,
			Tags:        strings.Join(ce.Tags, ","),
			UpdatedAt:   pc.getTimeFromMiliSecInt64(ce.UpdatedAt),
			CreatedAt:   pc.getTimeFromMiliSecInt64(ce.CreatedAt),
			/*UpdatedBy: model_chaos_experiment.User{
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
			Revision: util.SliceMap(ce.Revision, func(rev mongocollection.Revision) model_chaos_experiment.ChaosExperimentRevision {
				return model_chaos_experiment.ChaosExperimentRevision{
					RevisionID: rev.RevisionId,
					ExperimentManifest: model_chaos_experiment.ChaosExperimentManifest{
						Kind:       rev.ExperimentManifest.Kind,
						APIVersion: rev.ExperimentManifest.APIVersion,
						Metadata: model_chaos_experiment.ChaosExperimentMetadata{
							Name:              rev.ExperimentManifest.Metadata.Name,
							CreationTimestamp: pc.getTimeFromMiliSecInt64(rev.ExperimentManifest.Metadata.CreationTimestamp),
							Labels: model_chaos_experiment.ChaosExperimentLabels{
								InfraID:              rev.ExperimentManifest.Metadata.Labels.InfraID,
								RevisionID:           rev.ExperimentManifest.Metadata.Labels.RevisionID,
								WorkflowID:           rev.ExperimentManifest.Metadata.Labels.WorkflowID,
								ControllerInstanceID: rev.ExperimentManifest.Metadata.Labels.WorkflowsArgoprojIoControllerInstanceid,
							},
						},
						Spec: model_chaos_experiment.ChaosExperimentSpec{
							Templates: util.SliceMap(rev.ExperimentManifest.Spec.Templates, func(temp jsonfield.Template) model_chaos_experiment.ChaosExperimentTemplate {
								return model_chaos_experiment.ChaosExperimentTemplate{
									Name: temp.Name,
									Steps: model_chaos_experiment.ChaosExperimentSteps{
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
									Container: model_chaos_experiment.ChaosExperimentContainer{
										Name:    temp.Container.Name,
										Image:   temp.Container.Image,
										Command: strings.Join(temp.Container.Command, ","),
										Args:    strings.Join(temp.Container.Args, ","),
									},
								}
							}),
							Entrypoint: rev.ExperimentManifest.Spec.Entrypoint,
							Arguments: model_chaos_experiment.ChaosExperimentArguments{
								Parameters: util.SliceMap(rev.ExperimentManifest.Spec.Arguments.Parameters, func(param jsonfield.Parameter) model_chaos_experiment.ChaosExperimentParameter {
									return model_chaos_experiment.ChaosExperimentParameter{
										Name:  param.Name,
										Value: param.Value,
									}
								}),
							},
							ServiceAccountName: rev.ExperimentManifest.Spec.ServiceAccountName,
							PodGC: model_chaos_experiment.ChaosExperimentPodGC{
								Strategy: rev.ExperimentManifest.Spec.PodGC.Strategy,
							},
							SecurityContext: model_chaos_experiment.ChaosExperimentSecurityContext{
								RunAsUser:    rev.ExperimentManifest.Spec.SecurityContext.RunAsUser,
								RunAsNonRoot: rev.ExperimentManifest.Spec.SecurityContext.RunAsNonRoot,
							},
						},
						Status: model_chaos_experiment.ChaosExperimentStatus{
							StartedAt:  pc.getTimeFromMiliSecInt64(rev.ExperimentManifest.Status.StartedAt),
							FinishedAt: pc.getTimeFromMiliSecInt64(rev.ExperimentManifest.Status.FinishedAt),
						},
					},
					ChaosExperimentYamls: pc.getChaosExperimentYamls(rev),
					ChaosEngineYamls:     pc.getChaosEngineYamls(rev),
				}
			}),
			IsCustomExperiment: ce.IsCustomExperiment,
			RecentExperimentRunDetails: util.SliceMap(ce.RecentExperimentRunDetails, func(detail mongocollection.RecentExperimentRunDetail) model_chaos_experiment.ChaosExperimentRecentExperimentRunDetail {
				return model_chaos_experiment.ChaosExperimentRecentExperimentRunDetail{
					UpdatedAt: pc.getTimeFromMiliSecInt64(detail.UpdatedAt),
					CreatedAt: pc.getTimeFromMiliSecInt64(detail.CreatedAt),
					/*CreatedBy: model_chaos_experiment.User{
						UserID:   ci.CreatedBy.UserID,
						UserName: ci.CreatedBy.UserName,
						Email:    ci.CreatedBy.Email,
					},
					UpdatedBy: model_chaos_experiment.User{
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
					Probes: util.SliceMap(detail.Probes, func(probe mongocollection.Probe) model_chaos_experiment.ChaosExperimentProbe {
						return model_chaos_experiment.ChaosExperimentProbe{
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
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	cerms := util.SliceMap(cers, func(cer mongocollection.ChaosExperimentRun) model_chaos_experiment_run.ChaosExperimentRun {
		return model_chaos_experiment_run.ChaosExperimentRun{
			MongoID:   cer.ID.String(),
			ProjectID: cer.ProjectID,
			UpdatedAt: pc.getTimeFromMiliSecInt64(cer.UpdatedAt),
			CreatedAt: pc.getTimeFromMiliSecInt64(cer.CreatedAt),
			/*CreatedBy: model_chaos_experiment.User{
				UserID:   ci.CreatedBy.UserID,
				UserName: ci.CreatedBy.UserName,
				Email:    ci.CreatedBy.Email,
			},
			UpdatedBy: model_chaos_experiment.User{
				UserID:   ci.UpdatedBy.UserID,
				UserName: ci.UpdatedBy.UserName,
				Email:    ci.UpdatedBy.Email,
			},*/
			IsRemoved:       cer.IsRemoved,
			InfraID:         cer.InfraID,
			ExperimentRunID: cer.ExperimentRunID,
			ExperimentID:    cer.ExperimentID,
			ExperimentName:  cer.ExperimentName,
			Phase:           cer.Phase,
			Probes: util.SliceMap(cer.Probes, func(probe mongocollection.Probe) model_chaos_experiment_run.ChaosExperimentRunProbe {
				return model_chaos_experiment_run.ChaosExperimentRunProbe{
					FaultName:  probe.FaultName,
					ProbeNames: strings.Join(probe.ProbeNames, ","),
				}
			}),
			ExecutionData: model_chaos_experiment_run.ChaosExperimentRunExecutionData{
				ExperimentType:    cer.ExecutionData.ExperimentType,
				RevisionID:        cer.ExecutionData.RevisionID,
				NotifyID:          cer.ExecutionData.NotifyID,
				ExperimentID:      cer.ExecutionData.ExperimentID,
				EventType:         cer.ExecutionData.EventType,
				UID:               cer.ExecutionData.UID,
				Namespace:         cer.ExecutionData.Namespace,
				Name:              cer.ExecutionData.Name,
				CreationTimestamp: pc.getTimeFromSecString(cer.ExecutionData.CreationTimestamp),
				Phase:             cer.ExecutionData.Phase,
				Message:           cer.ExecutionData.Message,
				StartedAt:         pc.getTimeFromSecString(cer.ExecutionData.StartedAt),
				FinishedAt:        pc.getTimeFromSecString(cer.ExecutionData.FinishedAt),
				Nodes:             pc.getNodes(cer.ExecutionData.Nodes),
			},
			RevisionID:      cer.RevisionID,
			NotifyID:        cer.NotifyID,
			ResiliencyScore: cer.ResiliencyScore,
			RunSequence:     cer.RunSequence,
			Completed:       cer.Completed,
			FaultsAwaited:   cer.FaultsAwaited,
			FaultsFailed:    cer.FaultsFailed,
			FaultsNa:        cer.FaultsNa,
			FaultsPassed:    cer.FaultsPassed,
			FaultsStopped:   cer.FaultsStopped,
			TotalFaults:     cer.TotalFaults,
		}
	})

	for _, cerm := range cerms {
		db.Save(&cerm)
	}

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

func (pc PostgresConnector) getTimeFromMiliSecInt64(t int64) *time.Time {
	if t == 0 {
		return nil
	}
	time := time.Unix(t/1000, t%1000)
	return &time
}

func (pc PostgresConnector) getTimeFromSecString(ts string) *time.Time {
	tn, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil
	}
	time := time.Unix(tn, 0)
	return &time
}

func (pc PostgresConnector) getTimeFromIso8601String(iso8601Date string) *time.Time {
	if iso8601Date == "" {
		return nil
	}
	parsedTime, err := time.Parse(time.RFC3339, iso8601Date)
	if err != nil {
		return nil
	}
	return &parsedTime
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

func (pc PostgresConnector) getNodes(cerns map[string]jsonfield.ChaosExperimentRunNode) []model_chaos_experiment_run.ChaosExperimentRunNode {
	var mcerns []model_chaos_experiment_run.ChaosExperimentRunNode
	for name, cern := range cerns {
		mcerns = append(mcerns, model_chaos_experiment_run.ChaosExperimentRunNode{
			NodeName:   name,
			Name:       cern.Name,
			Phase:      cern.Phase,
			Message:    cern.Message,
			StartedAt:  pc.getTimeFromSecString(cern.StartedAt),
			FinishedAt: pc.getTimeFromSecString(cern.FinishedAt),
			Children:   strings.Join(cern.Children, ","),
			Type:       cern.Type,
			ChaosData: model_chaos_experiment_run.ChaosExperimentRunChaosData{
				EngineUID:              cern.ChaosData.EngineUID,
				EngineContext:          cern.ChaosData.EngineContext,
				EngineName:             cern.ChaosData.EngineName,
				Namespace:              cern.ChaosData.Namespace,
				ExperimentName:         cern.ChaosData.ExperimentName,
				ExperimentStatus:       cern.ChaosData.ExperimentStatus,
				LastUpdatedAt:          cern.ChaosData.LastUpdatedAt,
				ExperimentVerdict:      cern.ChaosData.ExperimentVerdict,
				ExperimentPod:          cern.ChaosData.ExperimentPod,
				RunnerPod:              cern.ChaosData.RunnerPod,
				ProbeSuccessPercentage: cern.ChaosData.ProbeSuccessPercentage,
				FailStep:               cern.ChaosData.FailStep,
				ChaosResult: model_chaos_experiment_run.ChaosExperimentRunChaosResult{
					Metadata: model_chaos_experiment_run.ChaosExperimentRunMetadata{
						Name:              cern.ChaosData.ChaosResult.Metadata.Name,
						Namespace:         cern.ChaosData.ChaosResult.Metadata.Namespace,
						UID:               cern.ChaosData.ChaosResult.Metadata.UID,
						ResourceVersion:   cern.ChaosData.ChaosResult.Metadata.ResourceVersion,
						Generation:        cern.ChaosData.ChaosResult.Metadata.Generation,
						CreationTimestamp: pc.getTimeFromIso8601String(cern.ChaosData.ChaosResult.Metadata.CreationTimestamp),
						Labels: model_chaos_experiment_run.ChaosExperimentRunLabels{
							AppKubernetesIoComponent:       cern.ChaosData.ChaosResult.Metadata.Labels.AppKubernetesIoComponent,
							AppKubernetesIoPartOf:          cern.ChaosData.ChaosResult.Metadata.Labels.AppKubernetesIoPartOf,
							AppKubernetesIoVersion:         cern.ChaosData.ChaosResult.Metadata.Labels.AppKubernetesIoVersion,
							BatchKubernetesIoControllerUID: cern.ChaosData.ChaosResult.Metadata.Labels.BatchKubernetesIoControllerUID,
							BatchKubernetesIoJobName:       cern.ChaosData.ChaosResult.Metadata.Labels.BatchKubernetesIoJobName,
							ChaosUID:                       cern.ChaosData.ChaosResult.Metadata.Labels.ChaosUID,
							ControllerUID:                  cern.ChaosData.ChaosResult.Metadata.Labels.ControllerUID,
							InfraID:                        cern.ChaosData.ChaosResult.Metadata.Labels.InfraID,
							JobName:                        cern.ChaosData.ChaosResult.Metadata.Labels.JobName,
							Name:                           cern.ChaosData.ChaosResult.Metadata.Labels.Name,
							StepPodName:                    cern.ChaosData.ChaosResult.Metadata.Labels.StepPodName,
							WorkflowName:                   cern.ChaosData.ChaosResult.Metadata.Labels.WorkflowName,
							WorkflowRunID:                  cern.ChaosData.ChaosResult.Metadata.Labels.WorkflowRunID,
						},
					},
					Spec: model_chaos_experiment_run.ChaosExperimentRunSpec{
						Engine:     cern.ChaosData.ChaosResult.Spec.Engine,
						Experiment: cern.ChaosData.ChaosResult.Spec.Experiment,
					},
					Status: model_chaos_experiment_run.ChaosExperimentRunStatus{
						ExperimentStatus: model_chaos_experiment_run.ChaosExperimentRunExperimentStatus{
							Phase:                  cern.ChaosData.ChaosResult.Status.ExperimentStatus.Phase,
							Verdict:                cern.ChaosData.ChaosResult.Status.ExperimentStatus.Verdict,
							ProbeSuccessPercentage: cern.ChaosData.ChaosResult.Status.ExperimentStatus.ProbeSuccessPercentage,
						},
						ProbeStatuses: util.SliceMap(cern.ChaosData.ChaosResult.Status.ProbeStatuses, func(probeStatus jsonfield.ChaosExperimentRunProbeStatuses) model_chaos_experiment_run.ChaosExperimentRunProbeStatus {
							return model_chaos_experiment_run.ChaosExperimentRunProbeStatus{
								Name: probeStatus.Name,
								Type: probeStatus.Type,
								Mode: probeStatus.Mode,
								Status: model_chaos_experiment_run.ChaosExperimentRunProbeStatusesStatus{
									Verdict:     probeStatus.Status.Verdict,
									Description: probeStatus.Status.Description,
								},
							}
						}),
						History: model_chaos_experiment_run.ChaosExperimentRunHistory{
							PassedRuns:  cern.ChaosData.ChaosResult.Status.History.PassedRuns,
							FailedRuns:  cern.ChaosData.ChaosResult.Status.History.FailedRuns,
							StoppedRuns: cern.ChaosData.ChaosResult.Status.History.StoppedRuns,
							Targets: util.SliceMap(cern.ChaosData.ChaosResult.Status.History.Targets, func(target jsonfield.ChaosExperimentRunHistoryTarget) model_chaos_experiment_run.ChaosExperimentRunHistoryTarget {
								return model_chaos_experiment_run.ChaosExperimentRunHistoryTarget{
									Name:        target.Name,
									Kind:        target.Kind,
									ChaosStatus: target.ChaosStatus,
								}
							}),
						},
					},
				},
			},
		})
	}
	return mcerns
}
