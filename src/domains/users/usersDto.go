package users

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type CreateUserDto struct {
	ID               string                   `json:"_id" bson:"_id" binding:"required"`
	Tutorial         bool                     `json:"tutorial" bson:"tutorial" binding:"required"`
	LastAccess       time.Time                `json:"lastAccess" bson:"lastAccess" binding:"required"`
	Access           []map[string]interface{} `json:"access" bson:"access"`
	ProjectIDs       []bson.ObjectId          `json:"projectIDs" bson:"projectIDs"`
	SettingSelection map[string]interface{}   `json:"settingSelection" bson:"settingSelection"`
	SettingMethod    map[string]interface{}   `json:"settingMethod" bson:"settingMethod"`
	Name             string                   `json:"name" bson:"name" binding:"required"`
	Picture          string                   `json:"picture" bson:"picture" binding:"required"`
	Email            string                   `json:"email" bson:"email" binding:"required,email"`
}

type FindOneUserDto struct {
	ID string `json:"_id,omitempty" bson:"_id,omitempty" uri:"_id" binding:"required"`
}

type UpdateUserDto struct {
	Tutorial         bool                   `bson:"tutorial,omitempty"`
	LastAccess       time.Time              `bson:"lastAccess,omitempty"`
	Access           map[string]interface{} `bson:"access,omitempty"`
	ProjectIDs       []bson.ObjectId        `bson:"projectIDs,omitempty"`
	SettingSelection map[string]interface{} `json:"settingSelection" bson:"settingSelection,omitempty"`
	SettingMethod    map[string]interface{} `json:"settingMethod" bson:"settingMethod,omitempty"`
	Name             string                 `json:"name" bson:"name,omitempty"`
	Picture          string                 `bson:"picture,omitempty"`
}

type ProjectIDNeededAddDto struct {
	ProjectID bson.ObjectId `json:"projectID" bson:"projectID" binding:"require,mongoid"`
}

type ProjectIDNeededRemoveDto struct {
	ProjectID bson.ObjectId `json:"projectID" bson:"projectID" binding:"require,mongoid"`
}

type FindManyDto struct {
	UserIDs []string `json:"userIDs" bson:"userIDs" binding:"require,unique"`
}
