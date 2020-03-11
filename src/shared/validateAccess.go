package shared

import (
	"errors"

	"github.com/khoa5773/go-server/src/helpers"
	"gopkg.in/mgo.v2/bson"
)

type Credentials struct {
	Id                string
	HasPersonalScopes bool
	HasProjectScopes  bool
	ProjectIDs        []bson.ObjectId
	IsAdmin           bool
}

type ValidateAccess struct {
	data        interface{}
	credentials Credentials
}

type Result struct {
	ValidDataLength int                      `json:"length"`
	ValidData       []map[string]interface{} `json:"data"`
}

func ValidateAccessToSingle(data map[string]interface{}, credentials Credentials) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	if credentials.IsAdmin {
		return data, nil
	}

	if (data["CreatedBy"] == credentials.Id || data["ID"] == credentials.Id) && credentials.HasPersonalScopes {
		return data, nil
	}

	projectID := data["ProjectID"].(bson.ObjectId).Hex()

	_, result := helpers.CheckObjectIdExists(credentials.ProjectIDs, projectID)

	if result && credentials.HasProjectScopes {
		return data, nil
	}

	return nil, errors.New("validation error")
}

func ValidateAccessToList(data []map[string]interface{}, credentials Credentials) (*Result, error) {
	if len(data) == 0 {
		return &Result{ValidData: []map[string]interface{}{}, ValidDataLength: 0}, nil
	}

	if credentials.IsAdmin {
		return &Result{ValidData: data, ValidDataLength: len(data)}, nil
	}

	validDataLength := 0
	var validData []map[string]interface{}

	for _, v := range data {
		if (v["CreatedBy"] == credentials.Id || v["ID"] == credentials.Id) && credentials.HasPersonalScopes {
			validDataLength++
			validData = append(validData, v)
			break
		}

		projectID := v["projectID"].(bson.ObjectId).Hex()
		_, result := helpers.CheckObjectIdExists(credentials.ProjectIDs, projectID)

		if result && credentials.HasProjectScopes {
			validDataLength++
			validData = append(validData, v)
		}
	}
	return &Result{ValidData: validData, ValidDataLength: validDataLength}, nil
}
