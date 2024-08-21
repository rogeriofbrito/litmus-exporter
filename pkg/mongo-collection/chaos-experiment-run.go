package mongocollection

import (
	"encoding/json"

	jsonfield "github.com/rogeriofbrito/litmus-exporter/pkg/json-field"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChaosExperimentRun struct {
	ID               primitive.ObjectID                        `bson:"_id"`
	ExecutionDataStr string                                    `bson:"execution_data"`
	ExecutionData    jsonfield.ChaosExperimentRunExecutionData `bson:"-"`
}

func (cer *ChaosExperimentRun) ParseExecutionData() error {
	ed := jsonfield.ChaosExperimentRunExecutionData{}
	if err := json.Unmarshal([]byte(cer.ExecutionDataStr), &ed); err != nil {
		return err
	}
	cer.ExecutionData = ed

	return nil
}
