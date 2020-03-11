package documents

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type DocumentType string

const (
	GENOTYPE  DocumentType = "GENOTYPE"
	PHENOTYPE DocumentType = "PHENOTYPE"
)

type FindManyDocumentsDto struct {
	RepositoryID bson.ObjectId `json:"repositoryID,omitempty" bson:"repositoryID,omitempty" uri:"repositoryID" binding:"required,mongoid"`
	ProjectID    bson.ObjectId `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
}

type FindOneDocumentDto struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" uri:"documentID" binding:"required,mongoid"`
	ProjectID    bson.ObjectId `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
	RepositoryID bson.ObjectId `json:"repositoryID,omitempty" bson:"repositoryID,omitempty" uri:"repositoryID" binding:"required,mongoid"`
}

type CreateOneDocumentDto struct {
	ID           bson.ObjectId `bson:"_id"`
	Name         string        `json:"name" bson:"name" form:"name" binding:"required"`
	Description  string        `json:"description" bson:"description" form:"description"`
	Type         DocumentType  `json:"type" bson:"type" form:"type" binding:"required,oneof=GENOTYPE PHENOTYPE"`
	MimeType     string        `bson:"mimeType"`
	RepositoryID bson.ObjectId `bson:"repositoryID"`
	ProjectID    bson.ObjectId `bson:"projectID"`
	Path         string        `bson:"path"`
	CreatedBy    string        `bson:"createdBy"`
	UpdatedBy    string        `bson:"updatedBy"`
	CreatedAt    time.Time     `bson:"createdAt"`
	UpdatedAt    time.Time     `bson:"updatedAt"`
}

type UpdateDocumentDto struct {
	Name         string        `json:"name" bson:"name,omitempty"`
	Description  string        `json:"description" bson:"description,omitempty"`
	Type         string        `json:"type" bson:"type,omitempty" binding:"required,oneof=GENOTYPE PHENOTYPE"`
	RepositoryID bson.ObjectId `json:"repositoryID" bson:"repositoryID,omitempty"`
	CreatedBy    string        `bson:"createdBy,omitempty"`
	UpdatedBy    string        `bson:"updatedBy,omitempty"`
	CreatedAt    time.Time     `bson:"createdAt,omitempty"`
	UpdatedAt    time.Time     `bson:"updatedAt,omitempty"`
}

type DeleteDocumentDto struct {
	ID           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" uri:"documentID" binding:"required,mongoid"`
	ProjectID    bson.ObjectId `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
	RepositoryID bson.ObjectId `json:"repositoryID,omitempty" bson:"repositoryID,omitempty" uri:"repositoryID" binding:"required,mongoid"`
}
