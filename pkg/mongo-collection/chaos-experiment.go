package mongocollection

import (
	"encoding/json"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChaosExperiment struct {
	ID                         primitive.ObjectID          `bson:"_id"`
	Name                       string                      `bson:"name"`
	Description                string                      `bson:"description"`
	Tags                       []string                    `bson:"tags"`
	UpdatedAt                  int64                       `bson:"updated_at"`
	CreatedAt                  int64                       `bson:"created_at"`
	CreatedBy                  User                        `bson:"created_by"`
	UpdatedBy                  User                        `bson:"updated_by"`
	IsRemoved                  bool                        `bson:"is_removed"`
	ProjectID                  string                      `bson:"project_id"`
	ExperimentID               string                      `bson:"experiment_id"`
	CronSyntax                 string                      `bson:"cron_syntax"`
	InfraID                    string                      `bson:"infra_id"`
	ExperimentType             string                      `bson:"experiment_type"`
	Revision                   []Revision                  `bson:"revision"`
	IsCustomExperiment         bool                        `bson:"is_custom_experiment"`
	RecentExperimentRunDetails []RecentExperimentRunDetail `bson:"recent_experiment_run_details"`
	TotalExperimentRuns        int                         `bson:"total_experiment_runs"`
}

type Revision struct {
	RevisionId            string                       `bson:"revision_id"`
	ExperimentManifestStr string                       `bson:"experiment_manifest"`
	ExperimentManifest    jsonfield.ExperimentManifest `bson:"-"`
}

type RecentExperimentRunDetail struct {
	UpdatedAt       int64   `bson:"updated_at"`
	CreatedAt       int64   `bson:"created_at"`
	CreatedBy       User    `bson:"created_by"`
	UpdatedBy       User    `bson:"updated_by"`
	IsRemoved       bool    `bson:"is_removed"`
	ProjectID       string  `bson:"project_id"`
	ExperimentRunID string  `bson:"experiment_run_id"`
	Phase           string  `bson:"phase"`
	NotifyID        string  `bson:"notify_id"`
	Completed       bool    `bson:"completed"`
	RunSequence     int     `bson:"run_sequence"`
	Probes          []Probe `bson:"probes"`
	ResiliencyScore float64 `bson:"resiliency_score"`
}

func (ce ChaosExperiment) ParseExperimentManifests() error {
	for i := range ce.Revision {
		em := jsonfield.ExperimentManifest{}
		if err := json.Unmarshal([]byte(ce.Revision[i].ExperimentManifestStr), &em); err != nil {
			return err
		}
		ce.Revision[i].ExperimentManifest = em
	}

	return nil
}
