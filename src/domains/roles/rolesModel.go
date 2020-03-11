package roles

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RoleName string

const (
	OWNER   RoleName = "EXPRESS_OWNER"
	MANAGER RoleName = "EXPRESS_MANAGER"
	MEMBER  RoleName = "EXPRESS_MEMBER"
)

// Role model
type Role struct {
	ID        bson.ObjectId `json:"_id" bson:"_id"`
	Name      RoleName      `json:"name" bson:"name"`
	Scopes    []string      `json:"scopes" bson:"scopes"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}
