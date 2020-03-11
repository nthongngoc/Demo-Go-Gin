package repositories

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type FindManyRepositoriesDto struct {
	ProjectID          bson.ObjectId `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
	ParentRepositoryID bson.ObjectId `json:"parentRepositoryID,omitempty" bson:"parentRepositoryID,omitempty" binding:"mongoid"`
}

type FindOneRepositoryDto struct {
	ID                 bson.ObjectId  `json:"_id,omitempty" bson:"_id,omitempty" uri:"repositoryID" binding:"required,mongoid"`
	ProjectID          bson.ObjectId  `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
	ParentRepositoryID *bson.ObjectId `bson:"parentRepositoryID,omitempty"`
}

type CreateOneRepositoryDto struct {
	ID                 bson.ObjectId  `bson:"_id,omitempty"`
	Name               string         `json:"name,omitempty" bson:"name,omitempty" binding:"required"`
	Description        string         `json:"description,omitempty" bson:"description,omitempty"`
	ParentRepositoryID *bson.ObjectId `json:"parentRepositoryID" bson:"parentRepositoryID" binding:"required"`
	ProjectID          bson.ObjectId  `bson:"projectID"`
	CreatedBy          string         `bson:"createdBy"`
	UpdatedBy          string         `bson:"updatedBy"`
	CreatedAt          time.Time      `bson:"createdAt"`
	UpdatedAt          time.Time      `bson:"updatedAt"`
}

type UpdateRepositoryDto struct {
	Name               string        `json:"name" bson:"name,omitempty"`
	Description        string        `json:"description" bson:"description,omitempty"`
	ParentRepositoryID bson.ObjectId `json:"parentRepositoryID,omitempty" bson:"parentRepositoryID,omitempty" binding:"mongoid"`
	UpdatedBy          string        `bson:"updatedBy,omitempty"`
	UpdatedAt          time.Time     `bson:"updatedAt,omitempty"`
}

type DeleteRepositoryDto struct {
	ID        bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" uri:"repositoryID" binding:"required,mongoid"`
	ProjectID bson.ObjectId `json:"projectID,omitempty" bson:"projectID,omitempty" uri:"projectID" binding:"required,mongoid"`
}
