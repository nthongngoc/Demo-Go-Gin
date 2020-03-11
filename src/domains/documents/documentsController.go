package documents

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
	DocumentsControllers := r.Group("projects/:projectID/repositories/:repositoryID/documents")
	DocumentsControllers.GET("", auth.JWTRequired, authz.Scopes([]string{"documents:read", "proj:documents:read"}), findManyDocumentsController)
	DocumentsControllers.POST("", auth.JWTRequired, authz.Scopes([]string{"documents:create", "proj:documents:create"}), createDocumentController)
	DocumentsControllers.GET("/:documentID", auth.JWTRequired, authz.Scopes([]string{"documents:read", "proj:documents:read"}), findOneDocumentController)
	DocumentsControllers.PUT("/:documentID", auth.JWTRequired, authz.Scopes([]string{"documents:update", "proj:documents:update"}), updateDocumentController)
	DocumentsControllers.DELETE("/:documentID", auth.JWTRequired, authz.Scopes([]string{"documents:delete", "proj:documents:delete"}), deleteDocumentController)
}

func findManyDocumentsController(c *gin.Context) {
	var findManyDocumentsDto FindManyDocumentsDto
	err := c.ShouldBindUri(&findManyDocumentsDto)
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

	documents, err := FindManyDocumentsInRepository(findManyDocumentsDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func findOneDocumentController(c *gin.Context) {
	var findOneDocumentDto FindOneDocumentDto
	err := c.ShouldBindUri(&findOneDocumentDto)
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

	document, err := FindOneDocument(&findOneDocumentDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"document": document})
}

func createDocumentController(c *gin.Context) {
	var createOneDocumentDto CreateOneDocumentDto
	err := c.ShouldBind(&createOneDocumentDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		_ = c.Error(err)
		return
	}

	mimeType, err := helpers.GetMimeTypeFromFilename(file.Filename)
	if err != nil {
		_ = c.Error(err)
		return
	}

	projectID := bson.ObjectIdHex(c.Param("projectID"))
	documentID := bson.NewObjectId()

	url, err := cloudstorage.UploadFileToCloudStorage(c, file, documentID.Hex(), projectID.Hex(), mimeType)
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
	createOneDocumentDto.ProjectID = projectID
	createOneDocumentDto.RepositoryID = bson.ObjectIdHex(c.Param("repositoryID"))
	createOneDocumentDto.Path = url
	createOneDocumentDto.MimeType = mimeType
	createOneDocumentDto.ID = documentID

	isSuccess, err := CreateOneDocument(&createOneDocumentDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": isSuccess})
}

func updateDocumentController(c *gin.Context) {
	var findOneDocumentDto FindOneDocumentDto
	var updateDocumentDto UpdateDocumentDto

	err := c.ShouldBindUri(&findOneDocumentDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&updateDocumentDto)
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

	isSuccess, err := UpdateDocument(&findOneDocumentDto, &updateDocumentDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func deleteDocumentController(c *gin.Context) {
	var deleteDocumentDto DeleteDocumentDto

	err := c.ShouldBindUri(&deleteDocumentDto)
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

	isSuccess, err := DeleteDocument(&deleteDocumentDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}
