package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khoa5773/go-server/src/domains/repositories"
	"github.com/khoa5773/go-server/src/domains/roles"
	"github.com/khoa5773/go-server/src/domains/users"
	"github.com/khoa5773/go-server/src/middleware/auth"
	"github.com/khoa5773/go-server/src/middleware/authz"
	"github.com/khoa5773/go-server/src/shared"
	"gopkg.in/mgo.v2/bson"
)

func ApplyRoutes(r *gin.Engine) {
	ProjectsController := r.Group("/projects")
	ProjectsController.GET("", auth.JWTRequired, authz.Scopes([]string{"projects:read", "proj:projects:read"}), findManyProjectsController)
	ProjectsController.GET("/:projectID", auth.JWTRequired, authz.Scopes([]string{"projects:read", "proj:projects:read"}), findOneProjectController)
	ProjectsController.POST("", auth.JWTRequired, authz.Scopes([]string{"projects:create", "proj:projects:create"}), createProjectController)
	ProjectsController.PUT("/:projectID", auth.JWTRequired, authz.Scopes([]string{"projects:update", "proj:projects:update"}), updateProjectController)
	ProjectsController.DELETE("/:projectID", auth.JWTRequired, authz.Scopes([]string{"projects:delete", "proj:projects:delete"}), deleteProjectController)
	ProjectsController.PUT("/:projectID/members/add", auth.JWTRequired, authz.Scopes([]string{"projects:update", "proj:projects:update"}), addProjectMembersController)
	ProjectsController.PUT("/:projectID/members/remove", auth.JWTRequired, authz.Scopes([]string{"projects:update", "proj:projects:update"}), removeProjectMembersController)
	ProjectsController.PUT("/:projectID/managers/add", auth.JWTRequired, authz.Scopes([]string{"projects:update", "proj:projects:update"}), addProjectManagersController)
	ProjectsController.PUT("/:projectID/managers/remove", auth.JWTRequired, authz.Scopes([]string{"projects:update", "proj:projects:update"}), removeProjectManagersController)
}

func findManyProjectsController(c *gin.Context) {
	projectIDs := c.MustGet("userProjectIDs").([]bson.ObjectId)

	credentials := shared.Credentials{
		Id:                c.MustGet("userID").(string),
		HasPersonalScopes: c.MustGet("hasPersonalScopes").(bool),
		HasProjectScopes:  c.MustGet("hasProjectScopes").(bool),
		ProjectIDs:        c.MustGet("userProjectIDs").([]bson.ObjectId),
	}

	projects, err := FindManyProjects(projectIDs, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

func findOneProjectController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	err := c.ShouldBindUri(&findOneProjectDto)
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

	project, err := FindOneProject(&findOneProjectDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

func createProjectController(c *gin.Context) {
	var createOneProjectDto CreateOneProjectDto
	err := c.ShouldBindJSON(&createOneProjectDto)
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

	projectID := bson.NewObjectId()
	createOneProjectDto.ID = projectID

	isSuccess, err := CreateOneProject(&createOneProjectDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	ownerRole := roles.OWNER
	credentials.IsAdmin = true

	_, err = users.AddProjectForUsers([]string{credentials.Id}, projectID, ownerRole, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	createOneRepositoryDto := InitRootRepositoryForProject(createOneProjectDto.Name)
	createOneRepositoryDto.ProjectID = projectID
	_, err = repositories.CreateOneRepository(&createOneRepositoryDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": isSuccess})
}

func updateProjectController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	var updateProjectDto UpdateProjectDto

	err := c.ShouldBindUri(&findOneProjectDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&updateProjectDto)
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

	isSuccess, err := UpdateProject(&findOneProjectDto, &updateProjectDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func deleteProjectController(c *gin.Context) {
	var deleteProjectDto DeleteProjectDto

	err := c.ShouldBindUri(&deleteProjectDto)
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

	isSuccess, userIDs, err := DeleteProject(&deleteProjectDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials.IsAdmin = true
	_, err = users.RemoveProjectFromUsers(userIDs, deleteProjectDto.ID, false, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	_, err = repositories.DeleteRootRepositoryOfProject(deleteProjectDto.ID, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})

}

func addProjectMembersController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	var addProjectMembersDto AddProjectMembersDto

	err := c.ShouldBindUri(&findOneProjectDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&addProjectMembersDto)
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

	usersData, err := users.FindMany(&users.FindManyDto{UserIDs: addProjectMembersDto.MemberIDs})
	if err != nil {
		_ = c.Error(err)
		return
	}

	_, err = CheckValidUsers(&findOneProjectDto, usersData)
	if err != nil {
		_ = c.Error(err)
		return
	}

	isSuccess, err := AddProjectMembers(&findOneProjectDto, &addProjectMembersDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	var userIDs []string
	for _, v := range usersData {
		userIDs = append(userIDs, v.ID)
	}

	memberRole := roles.MEMBER
	credentials.IsAdmin = true

	_, err = users.AddProjectForUsers(userIDs, findOneProjectDto.ID, memberRole, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func removeProjectMembersController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	var removeProjectMembersDto RemoveProjectMembersDto

	err := c.ShouldBindUri(&findOneProjectDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&removeProjectMembersDto)
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

	isSuccess, err := RemoveProjectMembers(&findOneProjectDto, &removeProjectMembersDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials.IsAdmin = true

	_, err = users.RemoveProjectFromUsers(removeProjectMembersDto.MemberIDs, findOneProjectDto.ID, false, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func addProjectManagersController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	var addProjectManagersDto AddProjectManagersDto

	err := c.ShouldBindUri(&findOneProjectDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&addProjectManagersDto)
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

	isSuccess, err := AddProjectManagers(&findOneProjectDto, &addProjectManagersDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	managerRole := roles.MANAGER
	credentials.IsAdmin = true

	isSuccess, err = users.AddProjectForUsers(addProjectManagersDto.ManagerIDs, findOneProjectDto.ID, managerRole, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}

func removeProjectManagersController(c *gin.Context) {
	var findOneProjectDto FindOneProjectDto
	var removeProjectManagersDto RemoveProjectManagersDto

	err := c.ShouldBindUri(&findOneProjectDto)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = c.ShouldBindJSON(&removeProjectManagersDto)
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

	isSuccess, err := RemoveProjectManagers(&findOneProjectDto, &removeProjectManagersDto, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	isSuccess, err = AddProjectMembers(&findOneProjectDto, &AddProjectMembersDto{MemberIDs: removeProjectManagersDto.ManagerIDs}, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	credentials.IsAdmin = true

	isSuccess, err = users.RemoveProjectFromUsers(removeProjectManagersDto.ManagerIDs, findOneProjectDto.ID, true, credentials)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": isSuccess})
}
