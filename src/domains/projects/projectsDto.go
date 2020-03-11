package projects

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type FindOneProjectDto struct {
	ID bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" uri:"projectID" binding:"required,mongoid"`
}

type CreateOneProjectDto struct {
	ID          bson.ObjectId `bson:"_id"`
	Name        string        `json:"name" bson:"name" binding:"required"`
	Description string        `json:"description" bson:"description"`
	CreatedBy   string        `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time     `bson:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt"`
	MemberIDs   []string      `bson:"memberIDs" binding:"unique"`
	ManagerIDs  []string      `bson:"managerIDs" binding:"unique"`
}

type UpdateProjectDto struct {
	Name        string    `json:"name" bson:"name,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	CreatedBy   string    `bson:"createdBy,omitempty"`
	CreatedAt   time.Time `bson:"createdAt,omitempty"`
	UpdatedAt   time.Time `bson:"updatedAt,omitempty"`
}

type DeleteProjectDto struct {
	ID bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty" uri:"projectID" binding:"required,mongoid"`
}

type AddProjectsDto struct {
	ProjectIDs []bson.ObjectId `json:"projectIDs" bson:"projectIDs" binding:"required,unique"`
}

type RemoveProjectsDto struct {
	ProjectIDs []bson.ObjectId `json:"projectIDs" bson:"projectIDs" binding:"required,unique"`
}

type AddProjectMembersDto struct {
	MemberIDs []string `json:"memberIDs" bson:"memberIDs" binding:"required,unique"`
}

type RemoveProjectMembersDto struct {
	MemberIDs []string `json:"memberIDs" bson:"memberIDs" binding:"required,unique"`
}

type AddProjectManagersDto struct {
	ManagerIDs []string `json:"managerIDs" bson:"managerIDs" binding:"required,unique"`
}

type RemoveProjectManagersDto struct {
	ManagerIDs []string `json:"managerIDs" bson:"managerIDs" binding:"required,unique"`
}
