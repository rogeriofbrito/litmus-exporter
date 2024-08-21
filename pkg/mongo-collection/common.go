package mongocollection

type User struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"username"`
	Email    string `bson:"email"`
}
