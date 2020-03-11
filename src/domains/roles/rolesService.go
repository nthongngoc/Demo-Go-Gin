package roles

import (
	"log"

	"github.com/khoa5773/go-server/src/shared"
	"gopkg.in/mgo.v2/bson"
)

func findRoleByName(roleName string) (*Role, error) {
	RolesModel := shared.MongoSession.C("roles")
	var role *Role
	err := RolesModel.Find(bson.M{"name": roleName}).One(&role)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return role, nil
}

func findRoleByID(roleID bson.ObjectId) (*Role, error) {
	RolesModel := shared.MongoSession.C("roles")
	var role *Role
	err := RolesModel.FindId(roleID).One(&role)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return role, nil
}
