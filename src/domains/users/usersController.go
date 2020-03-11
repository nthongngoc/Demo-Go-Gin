package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cloudstorage "github.com/khoa5773/go-server/src/domains/cloudStorage"
	"github.com/khoa5773/go-server/src/helpers"
	"github.com/khoa5773/go-server/src/middleware/auth"
	"github.com/khoa5773/go-server/src/middleware/authz"
	"github.com/khoa5773/go-server/src/shared"
	"gopkg.in/mgo.v2/bson"
)

func ApplyRoutes(r *gin.Engine) {
	UsersController := r.Group("/users")
	UsersController.GET("/:_id", auth.JWTRequired, authz.Scopes([]string{"users:read", "proj:users:read"}), findOneUserController)

	MyInfoController := r.Group("/me")
	MyInfoController.GET("", auth.JWTRequired, authz.Scopes([]string{"users:read", "proj:users:read"}), userInfoController)
	MyInfoController.PUT("", auth.JWTRequired, authz.Scopes([]string{"users:update", "proj:users:update"}), updateUserController)
	MyInfoController.PUT("/picture", auth.JWTRequired, authz.Scopes([]string{"users:update", "proj:users:update"}),
		updateUserPictureController)
}

func findOneUserController(c *gin.Context) {
	var findOneUserDto FindOneUserDto
	err := c.ShouldBindUri(&findOneUserDto)

	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials := shared.Credentials{
		Id:                c.MustGet("userID").(string),
		HasPersonalScopes: c.MustGet("hasPersonalScopes").(bool),
		HasProjectScopes:  c.MustGet("hasProjectScopes").(bool),
		ProjectIDs:        c.MustGet("userProjectIDs").([]bson.ObjectId),
	}

	user, err := FindOneUser(&findOneUserDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func userInfoController(c *gin.Context) {
	findOneUserDto := FindOneUserDto{ID: c.MustGet("userID").(string)}

	credentials := shared.Credentials{
		Id:                c.MustGet("userID").(string),
		HasPersonalScopes: c.MustGet("hasPersonalScopes").(bool),
		HasProjectScopes:  c.MustGet("hasProjectScopes").(bool),
		ProjectIDs:        c.MustGet("userProjectIDs").([]bson.ObjectId),
	}

	user, err := FindOneUser(&findOneUserDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func updateUserController(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var updateUserDto UpdateUserDto
	err := c.ShouldBindJSON(&updateUserDto)

	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials := shared.Credentials{
		Id:                c.MustGet("userID").(string),
		HasPersonalScopes: c.MustGet("hasPersonalScopes").(bool),
		HasProjectScopes:  c.MustGet("hasProjectScopes").(bool),
		ProjectIDs:        c.MustGet("userProjectIDs").([]bson.ObjectId),
	}

	isSuccess, err := UpdateUser(userID, &updateUserDto, credentials)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func updateUserPictureController(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	picture, err := c.FormFile("picture")
	if err != nil {
		_ = c.Error(err)
		return
	}

	mimeType, err := helpers.GetMimeTypeFromFilename(picture.Filename)
	if err != nil {
		_ = c.Error(err)
		return
	}

	url, err := cloudstorage.UploadFileToCloudStorage(c, picture, userID, "picture", mimeType)
	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials := shared.Credentials{
		Id:                c.MustGet("userID").(string),
		HasPersonalScopes: c.MustGet("hasPersonalScopes").(bool),
		HasProjectScopes:  c.MustGet("hasProjectScopes").(bool),
		ProjectIDs:        c.MustGet("userProjectIDs").([]bson.ObjectId),
	}

	isSuccess, err := UpdateUser(userID, &UpdateUserDto{Picture: url}, credentials)

	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}
