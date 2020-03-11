package users

import (
	"errors"
	"strings"

	"github.com/fatih/structs"
	cloudstorage "github.com/khoa5773/go-server/src/domains/cloudStorage"
	"github.com/khoa5773/go-server/src/domains/roles"
	"github.com/khoa5773/go-server/src/shared"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

func CreateUser(createUserDto *CreateUserDto) (bool, error) {
	UsersModel := shared.MongoSession.C("users")
	err := UsersModel.Insert(createUserDto)

	if err != nil {
		return false, err
	}

	return true, nil
}

func FindOneUser(findOneUserDto *FindOneUserDto, credentials shared.Credentials) (*User, error) {
	var user *User
	UsersModel := shared.MongoSession.C("users")
	err := UsersModel.Find(findOneUserDto).One(&user)

	if err != nil {
		return nil, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(*user), credentials)

	if err != nil {
		return nil, err
	}

	_ = mapstructure.Decode(data, user)
	if strings.Contains(user.Picture, "storage.googleapis.com") {
		user.Picture, err = cloudstorage.GenerateSignedUrl(user.ID, "picture", "jpg")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(userID string, updateUserDto *UpdateUserDto, credentials shared.Credentials) (bool, error) {
	UsersModel := shared.MongoSession.C("users")

	var user *User
	err := UsersModel.Find(&FindOneUserDto{ID: userID}).One(&user)

	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(*user), credentials)

	if err != nil {
		return false, err
	}

	err = UsersModel.UpdateId(userID, bson.M{"$set": updateUserDto})

	if err != nil {
		return false, err
	}

	return true, nil
}

func AddProjectForUsers(userIDs []string, projectIDNeededAdd bson.ObjectId, roleName roles.RoleName, credentials shared.Credentials) (bool, error) {
	UsersModel := shared.MongoSession.C("users")
	RolesModel := shared.MongoSession.C("roles")

	var users []User
	var role *roles.Role

	err := UsersModel.Find(bson.M{"_id": bson.M{"$in": userIDs}}).All(&users)
	if err != nil {
		return false, err
	}

	err = RolesModel.Find(bson.M{"name": roleName}).One(&role)
	if err != nil {
		return false, err
	}

	var usersData []map[string]interface{}
	for _, v := range users {
		usersData = append(usersData, structs.Map(v))
	}

	_, err = shared.ValidateAccessToList(usersData, credentials)
	if err != nil {
		return false, err
	}

	findingQuery := bson.M{"_id": bson.M{"$in": userIDs}}
	updatingQuery := bson.M{"$addToSet": bson.M{
		"projectIDs": projectIDNeededAdd,
		"access": bson.M{
			"projectID": projectIDNeededAdd,
			"roleID":    role.ID,
		},
	}}

	if role.Name == roles.MANAGER {
		findingQuery = bson.M{"_id": bson.M{"$in": userIDs}, "access": bson.M{"$elemMatch": bson.M{"projectID": projectIDNeededAdd}}}
		updatingQuery = bson.M{"$set": bson.M{"access.$.roleID": role.ID}}
	}

	_, err = UsersModel.UpdateAll(findingQuery, updatingQuery)
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveProjectFromUsers(userIDs []string, projectIDNeededRemove bson.ObjectId, isManager bool, credentials shared.Credentials) (bool, error) {
	UsersModel := shared.MongoSession.C("users")

	var users []*User

	err := UsersModel.Find(bson.M{"_id": bson.M{"$in": userIDs}}).All(&users)
	if err != nil {
		return false, err
	}

	var usersData []map[string]interface{}
	for _, v := range users {
		usersData = append(usersData, structs.Map(v))
	}

	_, err = shared.ValidateAccessToList(usersData, credentials)
	if err != nil {
		return false, err
	}

	if isManager {
		RolesModel := shared.MongoSession.C("roles")
		var role *roles.Role

		err = RolesModel.Find(bson.M{"name": roles.MEMBER}).One(&role)
		if err != nil {
			return false, err
		}

		_, err = UsersModel.UpdateAll(
			bson.M{"_id": bson.M{"$in": userIDs}, "access": bson.M{"$elemMatch": bson.M{"projectID": projectIDNeededRemove}}},
			bson.M{"$set": bson.M{"access.$.roleID": role.ID}},
		)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	_, err = UsersModel.UpdateAll(bson.M{"_id": bson.M{"$in": userIDs}}, bson.M{"$pull": bson.M{
		"projectIDs": projectIDNeededRemove,
		"access": bson.M{
			"projectID": projectIDNeededRemove,
		}}})
	if err != nil {
		return false, err
	}

	return true, nil
}

func FindMany(findManyDto *FindManyDto) ([]*User, error) {
	UsersModel := shared.MongoSession.C("users")
	var users []*User
	err := UsersModel.Find(bson.M{"_id": bson.M{"$in": findManyDto.UserIDs}}).All(&users)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if len(findManyDto.UserIDs) != len(users) {
		return nil, errors.New("at least one user is invalid")
	}

	return users, nil
}
