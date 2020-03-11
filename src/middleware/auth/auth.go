package auth

import (
	oidc "github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"github.com/khoa5773/go-server/src/configs"
	"gopkg.in/mgo.v2/bson"

	"github.com/khoa5773/go-server/src/shared"
)

type Credentials struct {
	ID         string          `json:"_id" bson:"_id"`
	ProjectIDs []bson.ObjectId `json:"projectIDs" bson:"projectIDs"`
	Scopes     [][]string      `json:"scopes" bson:"scopes"`
}

func JWTRequired(c *gin.Context) {
	token := c.GetHeader("Authorization")[7:]

	authenticator, err := shared.NewAuthenticator("login")

	if err != nil {
		_ = c.Error(err)
		return
	}

	oidcConfig := &oidc.Config{
		ClientID: configs.ConfigsService.Auth0ID,
	}

	idToken, err := authenticator.Provider.Verifier(oidcConfig).Verify(c, token)

	if err != nil {
		_ = c.Error(err)
		return
	}

	var profile map[string]interface{}
	if err = idToken.Claims(&profile); err != nil {
		_ = c.Error(err)
		return
	}

	c.Set("userID", profile["sub"].(string))
	userID := profile["sub"].(string)
	var userData Credentials
	err = GetCurrentUserCredentials(userID, &userData)
	if err != nil {
		_ = c.Error(err)
		return
	}
	c.Set("userProjectIDs", userData.ProjectIDs)
	c.Set("userScopes", userData.Scopes)
	c.Next()
}

func GetCurrentUserCredentials(id string, credentials *Credentials) error {
	UsersModel := shared.MongoSession.C("users")
	err := UsersModel.Pipe([]bson.M{
		{"$match": bson.M{"_id": id}},
		{"$unwind": bson.M{
			"path":                       "$access",
			"preserveNullAndEmptyArrays": true,
		}},
		{
			"$lookup": bson.M{
				"from":         "roles",
				"localField":   "access.roleID",
				"foreignField": "_id",
				"as":           "access.roleID",
			},
		},
		{"$unwind": bson.M{
			"path":                       "$access.roleID",
			"preserveNullAndEmptyArrays": true,
		}},
		{"$unwind": bson.M{
			"path":                       "$access",
			"preserveNullAndEmptyArrays": true,
		}},
		{
			"$group": bson.M{
				"_id":        "$_id",
				"projectIDs": bson.M{"$first": "$projectIDs"},
				"scopes": bson.M{
					"$push": "$access.roleID.scopes",
				},
			},
		},
	}).One(&credentials)

	if err != nil {
		return err
	}

	return nil
}
