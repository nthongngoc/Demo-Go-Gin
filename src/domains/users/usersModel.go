package users

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Access struct {
	ProjectID bson.ObjectId `json:"projectID" bson:"projectID"`
	RoleID    bson.ObjectId `json:"roleID" bson:"roleID"`
}

// User model
type User struct {
	ID               string                 `json:"_id" bson:"_id"`
	Tutorial         bool                   `json:"tutorial" bson:"tutorial"`
	LastAccess       time.Time              `json:"lastAccess" bson:"lastAccess"`
	Access           []Access               `json:"access" bson:"access"`
	ProjectIDs       []bson.ObjectId        `json:"projectIDs" bson:"projectIDs"`
	SettingSelection map[string]interface{} `json:"settingSelection" bson:"settingSelection"`
	SettingMethod    map[string]interface{} `json:"settingMethod" bson:"settingMethod"`
	Name             string                 `json:"name" bson:"name"`
	Picture          string                 `json:"picture" bson:"picture"`
	Email            string                 `json:"email" bson:"email"`
}
