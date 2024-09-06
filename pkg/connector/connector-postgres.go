package connector

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	typeslitmusk8s "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	typeslitmuschaosexperimentrun "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/chaos_experiment_run"
	typesmongodbchaosexperiment "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment"
	typesmongodbchaosexperimentrun "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/chaos_experiment_run"
	typesmongodbproject "github.com/litmuschaos/litmus/chaoscenter/graphql/server/pkg/database/mongodb/project"
	typesargoworkflows "github.com/rogeriofbrito/litmus-exporter/pkg/types/argo-workflows"
	typespostgreschaosengineyaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-engine-yaml"
	typespostgreschaosexperiment "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment"
	typespostgreschaosexperimentrun "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment-run"
	typespostgreschaosexperimentyaml "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/chaos-experiment-yaml"
	typespostgresproject "github.com/rogeriofbrito/litmus-exporter/pkg/types/postgres/project"
	"github.com/rogeriofbrito/litmus-exporter/pkg/util"
	"gorm.io/gorm"
	corev1 "k8s.io/api/core/v1"
	rbacV1 "k8s.io/api/rbac/v1"
)

type ctxValue string

func NewPostgresConnector(db *gorm.DB) *PostgresConnector {
	return &PostgresConnector{
		DB: db,
	}
}

type PostgresConnector struct {
	DB *gorm.DB
}

func (pc PostgresConnector) InitCtx(ctx context.Context) (context.Context, error) {
	return pc.getCtx(ctx)
}

func (pc PostgresConnector) Begin(ctx context.Context) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	err = tx.AutoMigrate(
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

	if tx.Where("1 = 1").Delete(&typespostgresproject.Project{}).Error != nil {
		return nil
	}

	if tx.Where("1 = 1").Delete(&typespostgreschaosexperiment.ChaosExperiment{}).Error != nil {
		return nil
	}

	if tx.Where("1 = 1").Delete(&typespostgreschaosexperimentrun.ChaosExperimentRun{}).Error != nil {
		return nil
	}

	return nil
}

func (pc PostgresConnector) Commit(ctx context.Context) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

func (pc PostgresConnector) Rollback(ctx context.Context) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	err = tx.Rollback().Error
	if err != nil {
		return err
	}

	return nil
}

func (pc PostgresConnector) SaveProjects(ctx context.Context, projs []typesmongodbproject.Project) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	pms := util.SliceMap(projs, func(p typesmongodbproject.Project) typespostgresproject.Project {
		return typespostgresproject.Project{
			UpdatedAt: pc.getTimeFromMiliSecInt64(p.UpdatedAt),
			CreatedAt: pc.getTimeFromMiliSecInt64(p.CreatedAt),
			IsRemoved: p.IsRemoved,
			Name:      p.Name,
			Members: util.SliceMap(p.Members, func(m *typesmongodbproject.Member) typespostgresproject.ProjectMembers {
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
		if err := tx.Save(&cem).Error; err != nil {
			return err
		}
	}

	return nil
}

func (pc PostgresConnector) SaveChaosExperiments(ctx context.Context, ces []typesmongodbchaosexperiment.ChaosExperimentRequest) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	cems := util.SliceMap(ces, func(ce typesmongodbchaosexperiment.ChaosExperimentRequest) typespostgreschaosexperiment.ChaosExperiment {
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
			Revision: util.SliceMap(ce.Revision, func(rev typesmongodbchaosexperiment.ExperimentRevision) typespostgreschaosexperiment.ChaosExperimentRevision {
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
			RecentExperimentRunDetails: util.SliceMap(ce.RecentExperimentRunDetails, func(detail typesmongodbchaosexperiment.ExperimentRunDetail) typespostgreschaosexperiment.ChaosExperimentRecentExperimentRunDetail {
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
					Probes: util.SliceMap(detail.Probe, func(probe typesmongodbchaosexperiment.Probes) typespostgreschaosexperiment.ChaosExperimentProbe {
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
		if err := tx.Save(&cem).Error; err != nil {
			return err
		}
	}

	return nil
}

func (pc PostgresConnector) SaveChaosExperimentRuns(ctx context.Context, cers []typesmongodbchaosexperimentrun.ChaosExperimentRun) error {
	tx, err := pc.getTx(ctx)
	if err != nil {
		return err
	}

	cerms := util.SliceMap(cers, func(cer typesmongodbchaosexperimentrun.ChaosExperimentRun) typespostgreschaosexperimentrun.ChaosExperimentRun {
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
			Probes: util.SliceMap(cer.Probes, func(probe typesmongodbchaosexperimentrun.Probes) typespostgreschaosexperimentrun.ChaosExperimentRunProbe {
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
		if err := tx.Save(&cerm).Error; err != nil {
			return err
		}
	}

	return nil
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
						Experiments: util.SliceMap(ce.Spec.Experiments, func(exp typeslitmusk8s.ExperimentList) typespostgreschaosengineyaml.ChaosEngineYamlExperiment {
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

func (pc PostgresConnector) getNodes(cerns map[string]typeslitmuschaosexperimentrun.Node) []typespostgreschaosexperimentrun.ChaosExperimentRunNode {
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
			ChaosData: func(cern typeslitmuschaosexperimentrun.Node) *typespostgreschaosexperimentrun.ChaosExperimentRunChaosData {
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
					ChaosResult: func(chaosResult *typeslitmusk8s.ChaosResult) typespostgreschaosexperimentrun.ChaosExperimentRunChaosResult {
						if chaosResult == nil {
							return typespostgreschaosexperimentrun.ChaosExperimentRunChaosResult{}
						}

						return typespostgreschaosexperimentrun.ChaosExperimentRunChaosResult{
							Metadata: typespostgreschaosexperimentrun.ChaosExperimentRunMetadata{
								Name:              chaosResult.ObjectMeta.Name,
								Namespace:         chaosResult.ObjectMeta.Namespace,
								UID:               string(chaosResult.ObjectMeta.UID),
								ResourceVersion:   chaosResult.ObjectMeta.ResourceVersion,
								Generation:        chaosResult.ObjectMeta.Generation,
								CreationTimestamp: &chaosResult.ObjectMeta.CreationTimestamp.Time,
								Labels: typespostgreschaosexperimentrun.ChaosExperimentRunLabels{
									AppKubernetesIoComponent:       chaosResult.ObjectMeta.Labels["app.kubernetes.io/component"],
									AppKubernetesIoPartOf:          chaosResult.ObjectMeta.Labels["app.kubernetes.io/part-of"],
									AppKubernetesIoVersion:         chaosResult.ObjectMeta.Labels["app.kubernetes.io/version"],
									BatchKubernetesIoControllerUID: chaosResult.ObjectMeta.Labels["batch.kubernetes.io/controller-uid"],
									BatchKubernetesIoJobName:       chaosResult.ObjectMeta.Labels["batch.kubernetes.io/job-name"],
									ChaosUID:                       chaosResult.ObjectMeta.Labels["chaosUID"],
									ControllerUID:                  chaosResult.ObjectMeta.Labels["controller-uid"],
									InfraID:                        chaosResult.ObjectMeta.Labels["infra_id"],
									JobName:                        chaosResult.ObjectMeta.Labels["job-name"],
									Name:                           chaosResult.ObjectMeta.Labels["name"],
									StepPodName:                    chaosResult.ObjectMeta.Labels["step_pod_name"],
									WorkflowName:                   chaosResult.ObjectMeta.Labels["workflow_name"],
									WorkflowRunID:                  chaosResult.ObjectMeta.Labels["workflow_run_id"],
								},
							},
							Spec: typespostgreschaosexperimentrun.ChaosExperimentRunSpec{
								EngineName:     chaosResult.Spec.EngineName,
								ExperimentName: chaosResult.Spec.ExperimentName,
							},
							Status: typespostgreschaosexperimentrun.ChaosExperimentRunStatus{
								ExperimentStatus: typespostgreschaosexperimentrun.ChaosExperimentRunExperimentStatus{
									Phase:                  string(chaosResult.Status.ExperimentStatus.Phase),
									Verdict:                string(chaosResult.Status.ExperimentStatus.Verdict),
									ProbeSuccessPercentage: chaosResult.Status.ExperimentStatus.ProbeSuccessPercentage,
								},
								ProbeStatuses: util.SliceMap(chaosResult.Status.ProbeStatuses, func(probeStatus typeslitmusk8s.ProbeStatuses) typespostgreschaosexperimentrun.ChaosExperimentRunProbeStatus {
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
									PassedRuns:  chaosResult.Status.History.PassedRuns,
									FailedRuns:  chaosResult.Status.History.FailedRuns,
									StoppedRuns: chaosResult.Status.History.StoppedRuns,
									Targets: util.SliceMap(chaosResult.Status.History.Targets, func(target typeslitmusk8s.TargetDetails) typespostgreschaosexperimentrun.ChaosExperimentRunHistoryTarget {
										return typespostgreschaosexperimentrun.ChaosExperimentRunHistoryTarget{
											Name:        target.Name,
											Kind:        target.Kind,
											ChaosStatus: target.ChaosStatus,
										}
									}),
								},
							},
						}
					}(cern.ChaosExp.ChaosResult),
				}
			}(cern),
		})
	}
	return mcerns
}

func (pc PostgresConnector) getCtx(ctx context.Context) (context.Context, error) {
	tx := pc.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return context.WithValue(ctx, ctxValue("tx"), tx), nil
}

func (pc PostgresConnector) getTx(ctx context.Context) (*gorm.DB, error) {
	tx := ctx.Value(ctxValue("tx"))
	if tx == nil {
		return nil, fmt.Errorf("tx not found in context")
	}

	return tx.(*gorm.DB), nil
}
