package repositories

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khoa5773/go-server/src/middleware/auth"
	"github.com/khoa5773/go-server/src/middleware/authz"
	"github.com/khoa5773/go-server/src/shared"
	"gopkg.in/mgo.v2/bson"
)

func ApplyRoutes(r *gin.Engine) {
	RepositoriesControllers := r.Group("projects/:projectID/repositories")
	RepositoriesControllers.GET("", auth.JWTRequired, authz.Scopes([]string{"repositories:read", "proj:repositories:read"}), findManyRepositoriesController)
	RepositoriesControllers.POST("", auth.JWTRequired, authz.Scopes([]string{"repositories:create", "proj:repositories:create"}), createRepositoryController)
	RepositoriesControllers.GET("/:repositoryID", auth.JWTRequired, authz.Scopes([]string{"repositories:read", "proj:repositories:read"}), findOneRepositoryController)
	RepositoriesControllers.PUT("/:repositoryID", auth.JWTRequired, authz.Scopes([]string{"repositories:update", "proj:repositories:update"}), updateRepositoryController)
	RepositoriesControllers.DELETE("/:repositoryID", auth.JWTRequired, authz.Scopes([]string{"repositories:delete", "proj:repositories:delete"}), deleteRepositoryController)
}

func findManyRepositoriesController(c *gin.Context) {
	var findManyRepositoriesDto FindManyRepositoriesDto
	err := c.ShouldBindUri(&findManyRepositoriesDto)
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

	repositories, err := FindManyRepositories(findManyRepositoriesDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"repositories": repositories})
}

func findOneRepositoryController(c *gin.Context) {
	var findOneRepositoryDto FindOneRepositoryDto
	err := c.ShouldBindUri(&findOneRepositoryDto)
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

	repository, err := FindOneRepository(&findOneRepositoryDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"repository": repository})
}

func createRepositoryController(c *gin.Context) {
	var createOneRepositoryDto CreateOneRepositoryDto
	err := c.ShouldBindJSON(&createOneRepositoryDto)
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
	createOneRepositoryDto.ProjectID = bson.ObjectIdHex(c.Param("projectID"))

	isSuccess, err := CreateOneRepository(&createOneRepositoryDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": isSuccess})
}

func updateRepositoryController(c *gin.Context) {
	var findOneRepositoryDto FindOneRepositoryDto
	var updateRepositoryDto UpdateRepositoryDto

	err := c.ShouldBindUri(&findOneRepositoryDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&updateRepositoryDto)
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

	isSuccess, err := UpdateRepository(&findOneRepositoryDto, &updateRepositoryDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func deleteRepositoryController(c *gin.Context) {
	var deleteRepositoryDto DeleteRepositoryDto

	err := c.ShouldBindUri(&deleteRepositoryDto)
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

	isSuccess, err := DeleteRepository(&deleteRepositoryDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}
