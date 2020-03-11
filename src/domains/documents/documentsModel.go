package documents

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Document model
type Document struct {
	ID           bson.ObjectId `json:"_id" bson:"_id"`
	Name         string        `json:"name" bson:"name"`
	Description  string        `json:"description" bson:"description"`
	MimeType     string        `json:"mimeType" bson:"mimeType"`
	ProjectID    bson.ObjectId `json:"projectID" bson:"projectID"`
	RepositoryID bson.ObjectId `json:"repositoryID" bson:"repositoryID"`
	Type         string        `json:"type" bson:"type"`
	Path         string        `json:"path" bson:"path"`
	CreatedBy    string        `json:"createdBy" bson:"createdBy"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
	UpdatedBy    string        `json:"updatedBy" bson:"updatedBy"`
}
