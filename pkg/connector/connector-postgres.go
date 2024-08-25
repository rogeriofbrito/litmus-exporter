package connector

import (
	"context"
	"fmt"
	"os"
	"strings"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	"github.com/rogeriofbrito/litmus-exporter/pkg/model"
	mongocollection "github.com/rogeriofbrito/litmus-exporter/pkg/mongo-collection"
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
		&model.Metadata{},
		&model.Labels{},
		&model.Spec{},
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

	revisionConv := func(c []mongocollection.Revision) []model.Revision {
		var m []model.Revision
		for _, ci := range c {
			m = append(m, model.Revision{
				RevisionID: ci.RevisionId,
				ExperimentManifest: model.ExperimentManifest{
					Kind:       ci.ExperimentManifest.Kind,
					APIVersion: ci.ExperimentManifest.APIVersion,
					Metadata: model.Metadata{
						Name:              ci.ExperimentManifest.Metadata.Name,
						CreationTimestamp: ci.ExperimentManifest.Metadata.CreationTimestamp,
						Labels: model.Labels{
							InfraID:              ci.ExperimentManifest.Metadata.Labels.InfraID,
							RevisionID:           ci.ExperimentManifest.Metadata.Labels.RevisionID,
							WorkflowID:           ci.ExperimentManifest.Metadata.Labels.WorkflowID,
							ControllerInstanceID: ci.ExperimentManifest.Metadata.Labels.WorkflowsArgoprojIoControllerInstanceid,
						},
					},
					Spec: model.Spec{
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
						StartedAt:  ci.ExperimentManifest.Status.StartedAt,
						FinishedAt: ci.ExperimentManifest.Status.FinishedAt,
					},
				},
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
				UpdatedAt: ci.UpdatedAt,
				CreatedAt: ci.CreatedAt,
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
			UpdatedAt:   ce.UpdatedAt,
			CreatedAt:   ce.CreatedAt,
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
