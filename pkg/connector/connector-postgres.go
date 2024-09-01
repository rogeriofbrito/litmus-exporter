package connector

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1" //TODO: rename import
	litmus_chaos_experiment_run "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/chaos_experiment_run"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"     //TODO: rename import
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run" //TODO: rename import
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/project"              //TODO: rename import
	argoworkflowstypes "github.com/rogeriofbrito/litmus-exporter/pkg/argo-workflows-types"
	model_chaos_engine_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-engine-yaml"
	model_chaos_experiment "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment"
	model_chaos_experiment_run "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment-run"
	model_chaos_experiment_yaml "github.com/rogeriofbrito/litmus-exporter/pkg/model/chaos-experiment-yaml"
	model_project "github.com/rogeriofbrito/litmus-exporter/pkg/model/project"
	"github.com/rogeriofbrito/litmus-exporter/pkg/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
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
		//Project
		&model_project.Project{},
		&model_project.ProjectMembers{},
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

func (pc PostgresConnector) SaveProjects(ctx context.Context, projs []project.Project) error {
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	pms := util.SliceMap(projs, func(p project.Project) model_project.Project {
		return model_project.Project{
			UpdatedAt: pc.getTimeFromMiliSecInt64(p.UpdatedAt),
			CreatedAt: pc.getTimeFromMiliSecInt64(p.CreatedAt),
			IsRemoved: p.IsRemoved,
			Name:      p.Name,
			Members: util.SliceMap(p.Members, func(m *project.Member) model_project.ProjectMembers {
				return model_project.ProjectMembers{
					Role:       string(m.Role),
					Invitation: string(m.Invitation),
					JoinedAt:   pc.getTimeFromMiliSecInt64(m.JoinedAt),
				}
			}),
			State: p.State,
		}
	})

	for _, cem := range pms {
		db.Save(&cem)
	}

	return nil
}

func (pc PostgresConnector) SaveChaosExperiments(ctx context.Context, ces []chaos_experiment.ChaosExperimentRequest) error {
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	cems := util.SliceMap(ces, func(ce chaos_experiment.ChaosExperimentRequest) model_chaos_experiment.ChaosExperiment {
		return model_chaos_experiment.ChaosExperiment{
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
			ExperimentType: string(ce.ExperimentType),
			Revision: util.SliceMap(ce.Revision, func(rev chaos_experiment.ExperimentRevision) model_chaos_experiment.ChaosExperimentRevision {
				em, err := util.ParseExperimentManifests(rev)
				if err != nil {
					panic(err)
				}
				return model_chaos_experiment.ChaosExperimentRevision{
					RevisionID: rev.RevisionID,
					ExperimentManifest: model_chaos_experiment.ChaosExperimentManifest{
						Kind:       em.Kind,
						APIVersion: em.APIVersion,
						/*
							Metadata: model_chaos_experiment.ChaosExperimentMetadata{
								Name:              em.Metadata.Name,
								CreationTimestamp: pc.getTimeFromMiliSecInt64(em.Metadata.CreationTimestamp),
								Labels: model_chaos_experiment.ChaosExperimentLabels{
									InfraID:              em.Metadata.Labels.InfraID,
									RevisionID:           em.Metadata.Labels.RevisionID,
									WorkflowID:           em.Metadata.Labels.WorkflowID,
									ControllerInstanceID: em.Metadata.Labels.WorkflowsArgoprojIoControllerInstanceid,
								},
							},
						*/
						Spec: model_chaos_experiment.ChaosExperimentSpec{
							Templates: util.SliceMap(em.Spec.Templates, func(temp argoworkflowstypes.Template) model_chaos_experiment.ChaosExperimentTemplate {
								return model_chaos_experiment.ChaosExperimentTemplate{
									Name: temp.Name,
									Steps: model_chaos_experiment.ChaosExperimentSteps{
										Name: func(steps argoworkflowstypes.Steps) string {
											if len(steps) == 0 {
												return ""
											}
											return steps[0][0].Name
										}(temp.Steps),
										Template: func(steps argoworkflowstypes.Steps) string {
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
							Entrypoint: em.Spec.Entrypoint,
							Arguments: model_chaos_experiment.ChaosExperimentArguments{
								Parameters: util.SliceMap(em.Spec.Arguments.Parameters, func(param argoworkflowstypes.Parameter) model_chaos_experiment.ChaosExperimentParameter {
									return model_chaos_experiment.ChaosExperimentParameter{
										Name:  param.Name,
										Value: param.Value,
									}
								}),
							},
							ServiceAccountName: em.Spec.ServiceAccountName,
							PodGC: model_chaos_experiment.ChaosExperimentPodGC{
								Strategy: em.Spec.PodGC.Strategy,
							},
							SecurityContext: model_chaos_experiment.ChaosExperimentSecurityContext{
								RunAsUser:    em.Spec.SecurityContext.RunAsUser,
								RunAsNonRoot: em.Spec.SecurityContext.RunAsNonRoot,
							},
						},
						Status: model_chaos_experiment.ChaosExperimentStatus{
							StartedAt:  pc.getTimeFromMiliSecInt64(em.Status.StartedAt),
							FinishedAt: pc.getTimeFromMiliSecInt64(em.Status.FinishedAt),
						},
					},
					ChaosExperimentYamls: pc.getChaosExperimentYamls(em),
					ChaosEngineYamls:     pc.getChaosEngineYamls(em),
				}
			}),
			IsCustomExperiment: ce.IsCustomExperiment,
			RecentExperimentRunDetails: util.SliceMap(ce.RecentExperimentRunDetails, func(detail chaos_experiment.ExperimentRunDetail) model_chaos_experiment.ChaosExperimentRecentExperimentRunDetail {
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
					Probes: util.SliceMap(detail.Probe, func(probe chaos_experiment.Probes) model_chaos_experiment.ChaosExperimentProbe {
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

func (pc PostgresConnector) SaveChaosExperimentRuns(ctx context.Context, cers []chaos_experiment_run.ChaosExperimentRun) error {
	db, err := pc.getGormDB()
	if err != nil {
		return err
	}

	cerms := util.SliceMap(cers, func(cer chaos_experiment_run.ChaosExperimentRun) model_chaos_experiment_run.ChaosExperimentRun {
		ed, err := util.ParseExecutionData(cer)
		if err != nil {
			panic(err)
		}

		return model_chaos_experiment_run.ChaosExperimentRun{
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
			Probes: util.SliceMap(cer.Probes, func(probe chaos_experiment_run.Probes) model_chaos_experiment_run.ChaosExperimentRunProbe {
				return model_chaos_experiment_run.ChaosExperimentRunProbe{
					FaultName:  probe.FaultName,
					ProbeNames: strings.Join(probe.ProbeNames, ","),
				}
			}),
			ExecutionData: model_chaos_experiment_run.ChaosExperimentRunExecutionData{
				ExperimentType:    ed.ExperimentType,
				RevisionID:        ed.RevisionID,
				ExperimentID:      ed.ExperimentID,
				EventType:         ed.EventType,
				UID:               ed.UID,
				Namespace:         ed.Namespace,
				Name:              ed.Name,
				CreationTimestamp: pc.getTimeFromSecString(ed.CreationTimestamp),
				Phase:             ed.Phase,
				Message:           ed.Message,
				StartedAt:         pc.getTimeFromSecString(ed.StartedAt),
				FinishedAt:        pc.getTimeFromSecString(ed.FinishedAt),
				Nodes:             pc.getNodes(ed.Nodes),
			},
			RevisionID:      cer.RevisionID,
			NotifyID:        cer.NotifyID,
			ResiliencyScore: cer.ResiliencyScore,
			RunSequence:     cer.RunSequence,
			Completed:       cer.Completed,
			FaultsAwaited:   cer.FaultsAwaited,
			FaultsFailed:    cer.FaultsFailed,
			FaultsNa:        cer.FaultsNA,
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

func (pc PostgresConnector) getChaosExperimentYamls(w *argoworkflowstypes.Workflow) []model_chaos_experiment_yaml.ChaosExperimentYaml {
	var mces []model_chaos_experiment_yaml.ChaosExperimentYaml
	for _, t := range w.Spec.Templates {
		if t.Name == "install-chaos-faults" {
			for _, a := range t.Inputs.Artifacts {
				ce, err := util.ParseChaosExperimentYaml(a.Raw.Data)
				if err != nil {
					panic(err)
				}

				mces = append(mces, model_chaos_experiment_yaml.ChaosExperimentYaml{
					APIVersion: ce.APIVersion,
					/*
						Description: model_chaos_experiment_yaml.ChaosExperimentYamlDescription{
							Message: ce.Description.Message,
						},
					*/
					Kind: ce.Kind,
					Metadata: model_chaos_experiment_yaml.ChaosExperimentYamlMetadata{
						Name: ce.Name,
						/*
							Labels: model_chaos_experiment_yaml.ChaosExperimentYamlLabels{
								Name:                     ce.GetObjectMeta().GetName(),
								AppKubernetesIoPartOf:    ce.Metadata.Labels.AppKubernetesIoPartOf,
								AppKubernetesIoComponent: ce.Metadata.Labels.AppKubernetesIoComponent,
								AppKubernetesIoVersion:   ce.Metadata.Labels.AppKubernetesIoVersion,
							},
						*/
					},
					Spec: model_chaos_experiment_yaml.ChaosExperimentYamlSpec{
						Definition: model_chaos_experiment_yaml.ChaosExperimentYamlDefinition{
							Scope: ce.Spec.Definition.Scope,
							Permissions: util.SliceMap(ce.Spec.Definition.Permissions, func(perm rbacV1.PolicyRule) model_chaos_experiment_yaml.ChaosExperimentYamlPermission {
								return model_chaos_experiment_yaml.ChaosExperimentYamlPermission{
									APIGroups: strings.Join(perm.APIGroups, ","),
									Resources: strings.Join(perm.Resources, ","),
									Verbs:     strings.Join(perm.Verbs, ","),
								}
							}),
							Image:           ce.Spec.Definition.Image,
							ImagePullPolicy: string(ce.Spec.Definition.ImagePullPolicy),
							Args:            strings.Join(ce.Spec.Definition.Args, ","),
							Command:         strings.Join(ce.Spec.Definition.Command, ","),
							Env: util.SliceMap(ce.Spec.Definition.ENVList, func(env corev1.EnvVar) model_chaos_experiment_yaml.ChaosExperimentYamlEnv {
								return model_chaos_experiment_yaml.ChaosExperimentYamlEnv{
									Name:  env.Name,
									Value: env.Value,
								}
							}),
							/*
								Labels: model_chaos_experiment_yaml.ChaosExperimentYamlDefinitionLabels{
									Name:                           ce.Spec.Definition.Labels.Name,
									AppKubernetesIoPartOf:          ce.Spec.Definition.Labels.AppKubernetesIoPartOf,
									AppKubernetesIoComponent:       ce.Spec.Definition.Labels.AppKubernetesIoComponent,
									AppKubernetesIoRuntimeAPIUsage: ce.Spec.Definition.Labels.AppKubernetesIoRuntimeAPIUsage,
									AppKubernetesIoVersion:         ce.Spec.Definition.Labels.AppKubernetesIoVersion,
								},
							*/
						},
					},
				})
			}
		}
	}
	return mces
}

func (pc PostgresConnector) getChaosEngineYamls(w *argoworkflowstypes.Workflow) []model_chaos_engine_yaml.ChaosEngineYaml {
	var mces []model_chaos_engine_yaml.ChaosEngineYaml
	for _, t := range w.Spec.Templates {
		if strings.Contains(t.Container.Image, "litmus-checker") {
			for _, a := range t.Inputs.Artifacts {
				ce, err := util.ParseChaosEngineYaml(a.Raw.Data)
				if err != nil {
					panic(err)
				}

				mces = append(mces, model_chaos_engine_yaml.ChaosEngineYaml{
					APIVersion: ce.APIVersion,
					Kind:       ce.Kind,
					Metadata: model_chaos_engine_yaml.ChaosEngineYamlMetadata{
						Namespace: ce.Namespace,
						/*
							Labels: model_chaos_engine_yaml.ChaosEngineYamlLabels{
								WorkflowRunID: ce.Metadata.Labels.WorkflowRunID,
								WorkflowName:  ce.Metadata.Labels.WorkflowName,
							},
						*/
						/*
							Annotations: model_chaos_engine_yaml.ChaosEngineYamlAnnotations{
								ProbeRef: ce.Metadata.Annotations.ProbeRef,
							},
						*/
						//GenerateName: ce.Metadata.GenerateName,
					},
					Spec: model_chaos_engine_yaml.ChaosEngineYamlSpec{
						EngineState: string(ce.Spec.EngineState),
						Appinfo: model_chaos_engine_yaml.ChaosEngineYamlAppInfo{
							Appns:    ce.Spec.Appinfo.Appns,
							Applabel: ce.Spec.Appinfo.Applabel,
							Appkind:  ce.Spec.Appinfo.AppKind,
						},
						ChaosServiceAccount: ce.Spec.ChaosServiceAccount,
						Experiments: util.SliceMap(ce.Spec.Experiments, func(exp v1alpha1.ExperimentList) model_chaos_engine_yaml.ChaosEngineYamlExperiment {
							return model_chaos_engine_yaml.ChaosEngineYamlExperiment{
								Name: exp.Name,
								Spec: model_chaos_engine_yaml.ChaosEngineYamlExperimentSpec{
									Components: model_chaos_engine_yaml.ChaosEngineYamlCompoments{
										Env: util.SliceMap(exp.Spec.Components.ENV, func(env corev1.EnvVar) model_chaos_engine_yaml.ChaosEngineYamlEnv {
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

func (pc PostgresConnector) getNodes(cerns map[string]litmus_chaos_experiment_run.Node) []model_chaos_experiment_run.ChaosExperimentRunNode {
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
			ChaosData: func(cern litmus_chaos_experiment_run.Node) *model_chaos_experiment_run.ChaosExperimentRunChaosData {
				if cern.ChaosExp == nil {
					return nil
				}
				return &model_chaos_experiment_run.ChaosExperimentRunChaosData{
					EngineUID:              cern.ChaosExp.EngineUID,
					EngineContext:          cern.ChaosExp.EngineContext,
					EngineName:             cern.ChaosExp.EngineName,
					Namespace:              cern.ChaosExp.Namespace,
					ExperimentName:         cern.ChaosExp.ExperimentName,
					ExperimentStatus:       cern.ChaosExp.ExperimentStatus,
					LastUpdatedAt:          cern.ChaosExp.LastUpdatedAt,
					ExperimentVerdict:      cern.ChaosExp.ExperimentVerdict,
					ExperimentPod:          cern.ChaosExp.ExperimentPod,
					RunnerPod:              cern.ChaosExp.RunnerPod,
					ProbeSuccessPercentage: cern.ChaosExp.ProbeSuccessPercentage,
					FailStep:               cern.ChaosExp.FailStep,
					ChaosResult: model_chaos_experiment_run.ChaosExperimentRunChaosResult{
						Metadata: model_chaos_experiment_run.ChaosExperimentRunMetadata{
							/*
								Name:              cern.ChaosExp.ChaosResult.Metadata.Name,
								Namespace:         cern.ChaosExp.ChaosResult.Metadata.Namespace,
								UID:               cern.ChaosExp.ChaosResult.Metadata.UID,
								ResourceVersion:   cern.ChaosExp.ChaosResult.Metadata.ResourceVersion,
								Generation:        cern.ChaosExp.ChaosResult.Metadata.Generation,
								CreationTimestamp: pc.getTimeFromIso8601String(cern.ChaosExp.ChaosResult.Metadata.CreationTimestamp),
							*/
							Labels: model_chaos_experiment_run.ChaosExperimentRunLabels{
								/*
									AppKubernetesIoComponent:       cern.ChaosExp.ChaosResult.Metadata.Labels.AppKubernetesIoComponent,
									AppKubernetesIoPartOf:          cern.ChaosExp.ChaosResult.Metadata.Labels.AppKubernetesIoPartOf,
									AppKubernetesIoVersion:         cern.ChaosExp.ChaosResult.Metadata.Labels.AppKubernetesIoVersion,
									BatchKubernetesIoControllerUID: cern.ChaosExp.ChaosResult.Metadata.Labels.BatchKubernetesIoControllerUID,
									BatchKubernetesIoJobName:       cern.ChaosExp.ChaosResult.Metadata.Labels.BatchKubernetesIoJobName,
									ChaosUID:                       cern.ChaosExp.ChaosResult.Metadata.Labels.ChaosUID,
									ControllerUID:                  cern.ChaosExp.ChaosResult.Metadata.Labels.ControllerUID,
									InfraID:                        cern.ChaosExp.ChaosResult.Metadata.Labels.InfraID,
									JobName:                        cern.ChaosExp.ChaosResult.Metadata.Labels.JobName,
									Name:                           cern.ChaosExp.ChaosResult.Metadata.Labels.Name,
									StepPodName:                    cern.ChaosExp.ChaosResult.Metadata.Labels.StepPodName,
									WorkflowName:                   cern.ChaosExp.ChaosResult.Metadata.Labels.WorkflowName,
									WorkflowRunID:                  cern.ChaosExp.ChaosResult.Metadata.Labels.WorkflowRunID,
								*/
							},
						},
						Spec: model_chaos_experiment_run.ChaosExperimentRunSpec{
							EngineName:     cern.ChaosExp.ChaosResult.Spec.EngineName,
							ExperimentName: cern.ChaosExp.ChaosResult.Spec.ExperimentName,
						},
						Status: model_chaos_experiment_run.ChaosExperimentRunStatus{
							ExperimentStatus: model_chaos_experiment_run.ChaosExperimentRunExperimentStatus{
								Phase:                  string(cern.ChaosExp.ChaosResult.Status.ExperimentStatus.Phase),
								Verdict:                string(cern.ChaosExp.ChaosResult.Status.ExperimentStatus.Verdict),
								ProbeSuccessPercentage: cern.ChaosExp.ChaosResult.Status.ExperimentStatus.ProbeSuccessPercentage,
							},
							ProbeStatuses: util.SliceMap(cern.ChaosExp.ChaosResult.Status.ProbeStatuses, func(probeStatus v1alpha1.ProbeStatuses) model_chaos_experiment_run.ChaosExperimentRunProbeStatus {
								return model_chaos_experiment_run.ChaosExperimentRunProbeStatus{
									Name: probeStatus.Name,
									Type: probeStatus.Type,
									Mode: probeStatus.Mode,
									Status: model_chaos_experiment_run.ChaosExperimentRunProbeStatusesStatus{
										Verdict:     string(probeStatus.Status.Verdict),
										Description: probeStatus.Status.Description,
									},
								}
							}),
							History: model_chaos_experiment_run.ChaosExperimentRunHistory{
								PassedRuns:  cern.ChaosExp.ChaosResult.Status.History.PassedRuns,
								FailedRuns:  cern.ChaosExp.ChaosResult.Status.History.FailedRuns,
								StoppedRuns: cern.ChaosExp.ChaosResult.Status.History.StoppedRuns,
								Targets: util.SliceMap(cern.ChaosExp.ChaosResult.Status.History.Targets, func(target v1alpha1.TargetDetails) model_chaos_experiment_run.ChaosExperimentRunHistoryTarget {
									return model_chaos_experiment_run.ChaosExperimentRunHistoryTarget{
										Name:        target.Name,
										Kind:        target.Kind,
										ChaosStatus: target.ChaosStatus,
									}
								}),
							},
						},
					},
				}
			}(cern),
		})
	}
	return mcerns
}
