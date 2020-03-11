package repositories

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Repository model
type Repository struct {
	ID                 bson.ObjectId   `json:"_id" bson:"_id"`
	Name               string          `json:"name" bson:"name"`
	Description        string          `json:"description" bson:"description"`
	ProjectID          bson.ObjectId   `json:"projectID" bson:"projectID"`
	ParentRepositoryID bson.ObjectId   `json:"parentRepositoryID" bson:"parentRepositoryID"`
	ChildRepositoryIDs []bson.ObjectId `json:"childRepositoryIDs" bson:"childRepositoryIDs"`
	DocumentIDs        []bson.ObjectId `json:"documentIDs" bson:"documentIDs"`
	CreatedBy          string          `json:"createdBy" bson:"createdBy"`
	CreatedAt          time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy          string          `json:"updatedBy" bson:"updatedBy"`
}
