package connector

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		&model.ChaosExperimentRevision{},
		&model.ChaosExperimentManifest{},
		&model.ChaosExperimentRecentRunDetails{},
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

	experimentRevisionConv := func(c []mongocollection.ChaosExperimentRevision) []model.ChaosExperimentRevision {
		var m []model.ChaosExperimentRevision
		for _, ci := range c {
			m = append(m, model.ChaosExperimentRevision{
				RevisionID: ci.RevisionId,
				ExperimentManifest: model.ChaosExperimentManifest{
					Kind: ci.ExperimentManifest.Kind,
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

	recentExperimentRunDetailsConv := func(c []mongocollection.ChaosExperimentRecentRunDetails) []model.ChaosExperimentRecentRunDetails {
		var m []model.ChaosExperimentRecentRunDetails
		for _, ci := range c {
			m = append(m, model.ChaosExperimentRecentRunDetails{
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
			Revision:                   experimentRevisionConv(ce.Revision),
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
