package mongocollection

import (
	"encoding/json"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChaosExperimentRun struct {
	ID               primitive.ObjectID                        `bson:"_id"`
	ProjectID        string                                    `bson:"project_id"`
	UpdatedAt        int64                                     `bson:"updated_at"`
	CreatedAt        int64                                     `bson:"created_at"`
	CreatedBy        User                                      `bson:"created_by"`
	UpdatedBy        User                                      `bson:"updated_by"`
	IsRemoved        bool                                      `bson:"is_removed"`
	InfraID          string                                    `bson:"infra_id"`
	ExperimentRunID  string                                    `bson:"experiment_run_id"`
	ExperimentID     string                                    `bson:"experiment_id"`
	ExperimentName   string                                    `bson:"experiment_name"`
	Phase            string                                    `bson:"phase"`
	Probes           []Probe                                   `bson:"probes"`
	ExecutionDataStr string                                    `bson:"execution_data"`
	ExecutionData    jsonfield.ChaosExperimentRunExecutionData `bson:"-"`
	RevisionID       string                                    `bson:"revision_id"`
	NotifyID         string                                    `bson:"notify_id"`
	ResiliencyScore  float64                                   `bson:"resiliency_score"`
	RunSequence      int                                       `bson:"run_sequence"`
	Completed        bool                                      `bson:"completed"`
	FaultsAwaited    int                                       `bson:"faults_awaited"`
	FaultsFailed     int                                       `bson:"faults_failed"`
	FaultsNa         int                                       `bson:"faults_na"`
	FaultsPassed     int                                       `bson:"faults_passed"`
	FaultsStopped    int                                       `bson:"faults_stopped"`
	TotalFaults      int                                       `bson:"total_faults"`
}

func (cer *ChaosExperimentRun) ParseExecutionData() error {
	ed := jsonfield.ChaosExperimentRunExecutionData{}
	if err := json.Unmarshal([]byte(cer.ExecutionDataStr), &ed); err != nil {
		return err
	}
	cer.ExecutionData = ed

	return nil
}
