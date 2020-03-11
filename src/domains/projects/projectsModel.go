package projects

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Project model
type Project struct {
	ID          bson.ObjectId `json:"_id" bson:"_id"`
	Name        string        `json:"name" bson:"name"`
	Description string        `json:"description" bson:"description"`
	CreatedBy   string        `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
	MemberIDs   []string      `json:"memberIDs" bson:"memberIDs"`
	ManagerIDs  []string      `json:"managerIDs" bson:"managerIDs"`
}
