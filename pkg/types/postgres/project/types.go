package typespostgresproject

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID  `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	UpdatedAt *time.Time `gorm:"column:updated_at"`
	CreatedAt *time.Time `gorm:"column:created_at"`
	//CreatedBy User             `gorm:"column:created_by"`
	//UpdatedBy User             `gorm:"column:updated_by"`
	IsRemoved bool             `gorm:"column:is_removed"`
	Name      string           `gorm:"column:name"`
	Members   []ProjectMembers `gorm:"foreignKey:project_id"`
	State     *string          `gorm:"column:state"`
}

type ProjectMembers struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;default:uuid_generate_v4();primaryKey"`
	ProjectID uuid.UUID `gorm:"column:project_id"`
	//UserID     string    `gorm:"column:user_id"`
	//Username   string    `gorm:"column:username"`
	//Email      string    `gorm:"column:email"`
	Role       string     `gorm:"column:role"`
	Invitation string     `gorm:"column:invitation"`
	JoinedAt   *time.Time `gorm:"column:joined_at"`
}
