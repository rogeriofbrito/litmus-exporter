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
	typesargoworkflows "github.com/rogeriofbrito/litmus-exporter/pkg/types/argo-workflows"
	typespostgreschaosengineyaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-engine-yaml"
	typespostgreschaosexperiment "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment"
	typespostgreschaosexperimentrun "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment-run"
	typespostgreschaosexperimentyaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment-yaml"
	typespostgresproject "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/project"
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
		&typespostgresproject.Project{},
		&typespostgresproject.ProjectMembers{},
		//ChaosExperiment
		&typespostgreschaosexperiment.ChaosExperiment{},
		&typespostgreschaosexperiment.ChaosExperimentRevision{},
		&typespostgreschaosexperiment.ChaosExperimentManifest{},
		&typespostgreschaosexperiment.ChaosExperimentMetadata{},
		&typespostgreschaosexperiment.ChaosExperimentLabels{},
		&typespostgreschaosexperiment.ChaosExperimentSpec{},
		&typespostgreschaosexperiment.ChaosExperimentTemplate{},
		&typespostgreschaosexperiment.ChaosExperimentSteps{},
		&typespostgreschaosexperiment.ChaosExperimentContainer{},
		&typespostgreschaosexperiment.ChaosExperimentArguments{},
		&typespostgreschaosexperiment.ChaosExperimentParameter{},
		&typespostgreschaosexperiment.ChaosExperimentPodGC{},
		&typespostgreschaosexperiment.ChaosExperimentSecurityContext{},
		&typespostgreschaosexperiment.ChaosExperimentStatus{},
		&typespostgreschaosexperiment.ChaosExperimentRecentExperimentRunDetail{},
		&typespostgreschaosexperiment.ChaosExperimentProbe{},
		//ChaosExperimentYaml
		&typespostgreschaosexperimentyaml.ChaosExperimentYaml{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlMetadata{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlLabels{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlSpec{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlDefinition{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlPermission{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlEnv{},
		&typespostgreschaosexperimentyaml.ChaosExperimentYamlDefinitionLabels{},
		//ChaosEngineYaml
		&typespostgreschaosengineyaml.ChaosEngineYaml{},
		&typespostgreschaosengineyaml.ChaosEngineYamlMetadata{},
		&typespostgreschaosengineyaml.ChaosEngineYamlSpec{},
		&typespostgreschaosengineyaml.ChaosEngineYamlLabels{},
		&typespostgreschaosengineyaml.ChaosEngineYamlAnnotations{},
		&typespostgreschaosengineyaml.ChaosEngineYamlAppInfo{},
		&typespostgreschaosengineyaml.ChaosEngineYamlExperiment{},
		&typespostgreschaosengineyaml.ChaosEngineYamlExperimentSpec{},
		&typespostgreschaosengineyaml.ChaosEngineYamlCompoments{},
		&typespostgreschaosengineyaml.ChaosEngineYamlEnv{},
		//ChaosExperimentRun
		&typespostgreschaosexperimentrun.ChaosExperimentRun{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunProbe{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunExecutionData{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunNode{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunChaosData{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunChaosResult{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunMetadata{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunSpec{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunStatus{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunLabels{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunExperimentStatus{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatus{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunHistory{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatusesStatus{},
		&typespostgreschaosexperimentrun.ChaosExperimentRunHistoryTarget{},
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

	pms := util.SliceMap(projs, func(p project.Project) typespostgresproject.Project {
		return typespostgresproject.Project{
			UpdatedAt: pc.getTimeFromMiliSecInt64(p.UpdatedAt),
			CreatedAt: pc.getTimeFromMiliSecInt64(p.CreatedAt),
			IsRemoved: p.IsRemoved,
			Name:      p.Name,
			Members: util.SliceMap(p.Members, func(m *project.Member) typespostgresproject.ProjectMembers {
				return typespostgresproject.ProjectMembers{
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

	cems := util.SliceMap(ces, func(ce chaos_experiment.ChaosExperimentRequest) typespostgreschaosexperiment.ChaosExperiment {
		return typespostgreschaosexperiment.ChaosExperiment{
			Name:           ce.Name,
			Description:    ce.Description,
			Tags:           strings.Join(ce.Tags, ","),
			UpdatedAt:      pc.getTimeFromMiliSecInt64(ce.UpdatedAt),
			CreatedAt:      pc.getTimeFromMiliSecInt64(ce.CreatedAt),
			CreatedBy:      ce.CreatedBy.Username,
			UpdatedBy:      ce.UpdatedBy.Username,
			IsRemoved:      ce.IsRemoved,
			ProjectID:      ce.ProjectID,
			ExperimentID:   ce.ExperimentID,
			CronSyntax:     ce.CronSyntax,
			InfraID:        ce.InfraID,
			ExperimentType: string(ce.ExperimentType),
			Revision: util.SliceMap(ce.Revision, func(rev chaos_experiment.ExperimentRevision) typespostgreschaosexperiment.ChaosExperimentRevision {
				w, err := util.ParseExperimentManifests(rev)
				if err != nil {
					panic(err)
				}
				return typespostgreschaosexperiment.ChaosExperimentRevision{
					RevisionID: rev.RevisionID,
					ExperimentManifest: typespostgreschaosexperiment.ChaosExperimentManifest{
						Kind:       w.Kind,
						APIVersion: w.APIVersion,
						Metadata: typespostgreschaosexperiment.ChaosExperimentMetadata{
							Name:              w.ObjectMeta.Name,
							CreationTimestamp: &w.ObjectMeta.CreationTimestamp.Time,
							Labels: typespostgreschaosexperiment.ChaosExperimentLabels{
								InfraID:              w.ObjectMeta.Labels["infra_id"],
								RevisionID:           w.ObjectMeta.Labels["revision_id"],
								WorkflowID:           w.ObjectMeta.Labels["workflow_id"],
								ControllerInstanceID: w.ObjectMeta.Labels["controller_instance_id"],
							},
						},
						Spec: typespostgreschaosexperiment.ChaosExperimentSpec{
							Templates: util.SliceMap(w.Spec.Templates, func(temp typesargoworkflows.Template) typespostgreschaosexperiment.ChaosExperimentTemplate {
								return typespostgreschaosexperiment.ChaosExperimentTemplate{
									Name: temp.Name,
									Steps: typespostgreschaosexperiment.ChaosExperimentSteps{
										Name: func(steps typesargoworkflows.Steps) string {
											if len(steps) == 0 {
												return ""
											}
											return steps[0][0].Name
										}(temp.Steps),
										Template: func(steps typesargoworkflows.Steps) string {
											if len(steps) == 0 {
												return ""
											}
											return steps[0][0].Template
										}(temp.Steps),
									},
									Container: typespostgreschaosexperiment.ChaosExperimentContainer{
										Name:    temp.Container.Name,
										Image:   temp.Container.Image,
										Command: strings.Join(temp.Container.Command, ","),
										Args:    strings.Join(temp.Container.Args, ","),
									},
								}
							}),
							Entrypoint: w.Spec.Entrypoint,
							Arguments: typespostgreschaosexperiment.ChaosExperimentArguments{
								Parameters: util.SliceMap(w.Spec.Arguments.Parameters, func(param typesargoworkflows.Parameter) typespostgreschaosexperiment.ChaosExperimentParameter {
									return typespostgreschaosexperiment.ChaosExperimentParameter{
										Name:  param.Name,
										Value: param.Value,
									}
								}),
							},
							ServiceAccountName: w.Spec.ServiceAccountName,
							PodGC: typespostgreschaosexperiment.ChaosExperimentPodGC{
								Strategy: w.Spec.PodGC.Strategy,
							},
							SecurityContext: typespostgreschaosexperiment.ChaosExperimentSecurityContext{
								RunAsUser:    w.Spec.SecurityContext.RunAsUser,
								RunAsNonRoot: w.Spec.SecurityContext.RunAsNonRoot,
							},
						},
						Status: typespostgreschaosexperiment.ChaosExperimentStatus{
							StartedAt:  pc.getTimeFromMiliSecInt64(w.Status.StartedAt),
							FinishedAt: pc.getTimeFromMiliSecInt64(w.Status.FinishedAt),
						},
					},
					ChaosExperimentYamls: pc.getChaosExperimentYamls(w),
					ChaosEngineYamls:     pc.getChaosEngineYamls(w),
				}
			}),
			IsCustomExperiment: ce.IsCustomExperiment,
			RecentExperimentRunDetails: util.SliceMap(ce.RecentExperimentRunDetails, func(detail chaos_experiment.ExperimentRunDetail) typespostgreschaosexperiment.ChaosExperimentRecentExperimentRunDetail {
				return typespostgreschaosexperiment.ChaosExperimentRecentExperimentRunDetail{
					UpdatedAt:       pc.getTimeFromMiliSecInt64(detail.UpdatedAt),
					CreatedAt:       pc.getTimeFromMiliSecInt64(detail.CreatedAt),
					CreatedBy:       detail.CreatedBy.Username,
					UpdatedBy:       detail.UpdatedBy.Username,
					IsRemoved:       detail.IsRemoved,
					ProjectID:       detail.ProjectID,
					ExperimentRunID: detail.ExperimentRunID,
					Phase:           detail.Phase,
					NotifyID:        detail.NotifyID,
					Completed:       detail.Completed,
					RunSequence:     detail.RunSequence,
					Probes: util.SliceMap(detail.Probe, func(probe chaos_experiment.Probes) typespostgreschaosexperiment.ChaosExperimentProbe {
						return typespostgreschaosexperiment.ChaosExperimentProbe{
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

	cerms := util.SliceMap(cers, func(cer chaos_experiment_run.ChaosExperimentRun) typespostgreschaosexperimentrun.ChaosExperimentRun {
		ed, err := util.ParseExecutionData(cer)
		if err != nil {
			panic(err)
		}

		return typespostgreschaosexperimentrun.ChaosExperimentRun{
			ProjectID:       cer.ProjectID,
			UpdatedAt:       pc.getTimeFromMiliSecInt64(cer.UpdatedAt),
			CreatedAt:       pc.getTimeFromMiliSecInt64(cer.CreatedAt),
			CreatedBy:       cer.CreatedBy.Username,
			UpdatedBy:       cer.UpdatedBy.Username,
			IsRemoved:       cer.IsRemoved,
			InfraID:         cer.InfraID,
			ExperimentRunID: cer.ExperimentRunID,
			ExperimentID:    cer.ExperimentID,
			ExperimentName:  cer.ExperimentName,
			Phase:           cer.Phase,
			Probes: util.SliceMap(cer.Probes, func(probe chaos_experiment_run.Probes) typespostgreschaosexperimentrun.ChaosExperimentRunProbe {
				return typespostgreschaosexperimentrun.ChaosExperimentRunProbe{
					FaultName:  probe.FaultName,
					ProbeNames: strings.Join(probe.ProbeNames, ","),
				}
			}),
			ExecutionData: typespostgreschaosexperimentrun.ChaosExperimentRunExecutionData{
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

func (pc PostgresConnector) getChaosExperimentYamls(w *typesargoworkflows.Workflow) []typespostgreschaosexperimentyaml.ChaosExperimentYaml {
	var mces []typespostgreschaosexperimentyaml.ChaosExperimentYaml
	for _, t := range w.Spec.Templates {
		if t.Name == "install-chaos-faults" {
			for _, a := range t.Inputs.Artifacts {
				ce, err := util.ParseChaosExperimentYaml(a.Raw.Data)
				if err != nil {
					panic(err)
				}

				mces = append(mces, typespostgreschaosexperimentyaml.ChaosExperimentYaml{
					APIVersion: ce.APIVersion,
					Kind:       ce.Kind,
					Metadata: typespostgreschaosexperimentyaml.ChaosExperimentYamlMetadata{
						Name: ce.Name,
						Labels: typespostgreschaosexperimentyaml.ChaosExperimentYamlLabels{
							AppKubernetesIoPartOf:    ce.ObjectMeta.Labels["app.kubernetes.io/part-of"],
							AppKubernetesIoComponent: ce.ObjectMeta.Labels["app.kubernetes.io/component"],
							AppKubernetesIoVersion:   ce.ObjectMeta.Labels["app.kubernetes.io/version"],
						},
					},
					Spec: typespostgreschaosexperimentyaml.ChaosExperimentYamlSpec{
						Definition: typespostgreschaosexperimentyaml.ChaosExperimentYamlDefinition{
							Scope: ce.Spec.Definition.Scope,
							Permissions: util.SliceMap(ce.Spec.Definition.Permissions, func(perm rbacV1.PolicyRule) typespostgreschaosexperimentyaml.ChaosExperimentYamlPermission {
								return typespostgreschaosexperimentyaml.ChaosExperimentYamlPermission{
									APIGroups: strings.Join(perm.APIGroups, ","),
									Resources: strings.Join(perm.Resources, ","),
									Verbs:     strings.Join(perm.Verbs, ","),
								}
							}),
							Image:           ce.Spec.Definition.Image,
							ImagePullPolicy: string(ce.Spec.Definition.ImagePullPolicy),
							Args:            strings.Join(ce.Spec.Definition.Args, ","),
							Command:         strings.Join(ce.Spec.Definition.Command, ","),
							Env: util.SliceMap(ce.Spec.Definition.ENVList, func(env corev1.EnvVar) typespostgreschaosexperimentyaml.ChaosExperimentYamlEnv {
								return typespostgreschaosexperimentyaml.ChaosExperimentYamlEnv{
									Name:  env.Name,
									Value: env.Value,
								}
							}),
							Labels: typespostgreschaosexperimentyaml.ChaosExperimentYamlDefinitionLabels{
								Name:                           ce.Spec.Definition.Labels["name"],
								AppKubernetesIoPartOf:          ce.Spec.Definition.Labels["app.kubernetes.io/part-of"],
								AppKubernetesIoComponent:       ce.Spec.Definition.Labels["app.kubernetes.io/component"],
								AppKubernetesIoRuntimeAPIUsage: ce.Spec.Definition.Labels["app.kubernetes.io/runtime-api-usage"],
								AppKubernetesIoVersion:         ce.Spec.Definition.Labels["app.kubernetes.io/version"],
							},
						},
					},
				})
			}
		}
	}
	return mces
}

func (pc PostgresConnector) getChaosEngineYamls(w *typesargoworkflows.Workflow) []typespostgreschaosengineyaml.ChaosEngineYaml {
	var mces []typespostgreschaosengineyaml.ChaosEngineYaml
	for _, t := range w.Spec.Templates {
		if strings.Contains(t.Container.Image, "litmus-checker") {
			for _, a := range t.Inputs.Artifacts {
				ce, err := util.ParseChaosEngineYaml(a.Raw.Data)
				if err != nil {
					panic(err)
				}

				mces = append(mces, typespostgreschaosengineyaml.ChaosEngineYaml{
					APIVersion: ce.APIVersion,
					Kind:       ce.Kind,
					Metadata: typespostgreschaosengineyaml.ChaosEngineYamlMetadata{
						Namespace: ce.Namespace,
						Labels: typespostgreschaosengineyaml.ChaosEngineYamlLabels{
							WorkflowRunID: ce.Labels["workflow_run_id"],
							WorkflowName:  ce.Labels["workflow_name"],
						},
						Annotations: typespostgreschaosengineyaml.ChaosEngineYamlAnnotations{
							ProbeRef: ce.ObjectMeta.Annotations["probeRef"],
						},
						GenerateName: ce.ObjectMeta.GenerateName,
					},
					Spec: typespostgreschaosengineyaml.ChaosEngineYamlSpec{
						EngineState: string(ce.Spec.EngineState),
						Appinfo: typespostgreschaosengineyaml.ChaosEngineYamlAppInfo{
							Appns:    ce.Spec.Appinfo.Appns,
							Applabel: ce.Spec.Appinfo.Applabel,
							Appkind:  ce.Spec.Appinfo.AppKind,
						},
						ChaosServiceAccount: ce.Spec.ChaosServiceAccount,
						Experiments: util.SliceMap(ce.Spec.Experiments, func(exp v1alpha1.ExperimentList) typespostgreschaosengineyaml.ChaosEngineYamlExperiment {
							return typespostgreschaosengineyaml.ChaosEngineYamlExperiment{
								Name: exp.Name,
								Spec: typespostgreschaosengineyaml.ChaosEngineYamlExperimentSpec{
									Components: typespostgreschaosengineyaml.ChaosEngineYamlCompoments{
										Env: util.SliceMap(exp.Spec.Components.ENV, func(env corev1.EnvVar) typespostgreschaosengineyaml.ChaosEngineYamlEnv {
											return typespostgreschaosengineyaml.ChaosEngineYamlEnv{
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

func (pc PostgresConnector) getNodes(cerns map[string]litmus_chaos_experiment_run.Node) []typespostgreschaosexperimentrun.ChaosExperimentRunNode {
	var mcerns []typespostgreschaosexperimentrun.ChaosExperimentRunNode
	for name, cern := range cerns {
		mcerns = append(mcerns, typespostgreschaosexperimentrun.ChaosExperimentRunNode{
			NodeName:   name,
			Name:       cern.Name,
			Phase:      cern.Phase,
			Message:    cern.Message,
			StartedAt:  pc.getTimeFromSecString(cern.StartedAt),
			FinishedAt: pc.getTimeFromSecString(cern.FinishedAt),
			Children:   strings.Join(cern.Children, ","),
			Type:       cern.Type,
			ChaosData: func(cern litmus_chaos_experiment_run.Node) *typespostgreschaosexperimentrun.ChaosExperimentRunChaosData {
				if cern.ChaosExp == nil {
					return nil
				}
				return &typespostgreschaosexperimentrun.ChaosExperimentRunChaosData{
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
					ChaosResult: typespostgreschaosexperimentrun.ChaosExperimentRunChaosResult{
						Metadata: typespostgreschaosexperimentrun.ChaosExperimentRunMetadata{
							Name:              cern.ChaosExp.ChaosResult.ObjectMeta.Name,
							Namespace:         cern.ChaosExp.ChaosResult.ObjectMeta.Namespace,
							UID:               string(cern.ChaosExp.ChaosResult.ObjectMeta.UID),
							ResourceVersion:   cern.ChaosExp.ChaosResult.ObjectMeta.ResourceVersion,
							Generation:        cern.ChaosExp.ChaosResult.ObjectMeta.Generation,
							CreationTimestamp: &cern.ChaosExp.ChaosResult.ObjectMeta.CreationTimestamp.Time,
							Labels: typespostgreschaosexperimentrun.ChaosExperimentRunLabels{
								AppKubernetesIoComponent:       cern.ChaosExp.ChaosResult.ObjectMeta.Labels["app.kubernetes.io/component"],
								AppKubernetesIoPartOf:          cern.ChaosExp.ChaosResult.ObjectMeta.Labels["app.kubernetes.io/part-of"],
								AppKubernetesIoVersion:         cern.ChaosExp.ChaosResult.ObjectMeta.Labels["app.kubernetes.io/version"],
								BatchKubernetesIoControllerUID: cern.ChaosExp.ChaosResult.ObjectMeta.Labels["batch.kubernetes.io/controller-uid"],
								BatchKubernetesIoJobName:       cern.ChaosExp.ChaosResult.ObjectMeta.Labels["batch.kubernetes.io/job-name"],
								ChaosUID:                       cern.ChaosExp.ChaosResult.ObjectMeta.Labels["chaosUID"],
								ControllerUID:                  cern.ChaosExp.ChaosResult.ObjectMeta.Labels["controller-uid"],
								InfraID:                        cern.ChaosExp.ChaosResult.ObjectMeta.Labels["infra_id"],
								JobName:                        cern.ChaosExp.ChaosResult.ObjectMeta.Labels["job-name"],
								Name:                           cern.ChaosExp.ChaosResult.ObjectMeta.Labels["name"],
								StepPodName:                    cern.ChaosExp.ChaosResult.ObjectMeta.Labels["step_pod_name"],
								WorkflowName:                   cern.ChaosExp.ChaosResult.ObjectMeta.Labels["workflow_name"],
								WorkflowRunID:                  cern.ChaosExp.ChaosResult.ObjectMeta.Labels["workflow_run_id"],
							},
						},
						Spec: typespostgreschaosexperimentrun.ChaosExperimentRunSpec{
							EngineName:     cern.ChaosExp.ChaosResult.Spec.EngineName,
							ExperimentName: cern.ChaosExp.ChaosResult.Spec.ExperimentName,
						},
						Status: typespostgreschaosexperimentrun.ChaosExperimentRunStatus{
							ExperimentStatus: typespostgreschaosexperimentrun.ChaosExperimentRunExperimentStatus{
								Phase:                  string(cern.ChaosExp.ChaosResult.Status.ExperimentStatus.Phase),
								Verdict:                string(cern.ChaosExp.ChaosResult.Status.ExperimentStatus.Verdict),
								ProbeSuccessPercentage: cern.ChaosExp.ChaosResult.Status.ExperimentStatus.ProbeSuccessPercentage,
							},
							ProbeStatuses: util.SliceMap(cern.ChaosExp.ChaosResult.Status.ProbeStatuses, func(probeStatus v1alpha1.ProbeStatuses) typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatus {
								return typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatus{
									Name: probeStatus.Name,
									Type: probeStatus.Type,
									Mode: probeStatus.Mode,
									Status: typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatusesStatus{
										Verdict:     string(probeStatus.Status.Verdict),
										Description: probeStatus.Status.Description,
									},
								}
							}),
							History: typespostgreschaosexperimentrun.ChaosExperimentRunHistory{
								PassedRuns:  cern.ChaosExp.ChaosResult.Status.History.PassedRuns,
								FailedRuns:  cern.ChaosExp.ChaosResult.Status.History.FailedRuns,
								StoppedRuns: cern.ChaosExp.ChaosResult.Status.History.StoppedRuns,
								Targets: util.SliceMap(cern.ChaosExp.ChaosResult.Status.History.Targets, func(target v1alpha1.TargetDetails) typespostgreschaosexperimentrun.ChaosExperimentRunHistoryTarget {
									return typespostgreschaosexperimentrun.ChaosExperimentRunHistoryTarget{
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
