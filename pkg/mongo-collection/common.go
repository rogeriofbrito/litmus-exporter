package mongocollection

type User struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"username"`
	Email    string `bson:"email"`
}

type Probe struct {
	FaultName  string   `bson:"fault_name"`
	ProbeNames []string `bson:"probe_names"`
}
