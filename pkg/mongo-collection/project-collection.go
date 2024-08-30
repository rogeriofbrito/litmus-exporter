package mongocollection

import "go.mongodb.org/mongo-driver/bson/primitive"

type Project struct {
	ID        primitive.ObjectID `bson:"_id"`
	UpdatedAt string             `bson:"updated_at"`
	CreatedAt string             `bson:"created_at"`
	CreatedBy User               `bson:"created_by"`
	UpdatedBy User               `bson:"updated_by"`
	IsRemoved bool               `bson:"is_removed"`
	Name      string             `bson:"name"`
	Members   []ProjectMembers   `bson:"members"`
	State     string             `bson:"state"`
}

type ProjectMembers struct {
	UserID     string `bson:"user_id"`
	Username   string `bson:"username"`
	Email      string `bson:"email"`
	Name       string `bson:"name"`
	Role       string `bson:"role"`
	Invitation string `bson:"invitation"`
	JoinedAt   string `bson:"joined_at"`
}
