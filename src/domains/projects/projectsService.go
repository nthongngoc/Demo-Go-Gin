package projects

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/khoa5773/go-server/src/domains/repositories"
	"github.com/khoa5773/go-server/src/domains/users"
	"github.com/khoa5773/go-server/src/helpers"
	"github.com/khoa5773/go-server/src/shared"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

type Result struct {
	ValidDataLength int           `json:"length"`
	ValidData       []interface{} `json:"data"`
}

func ValidateAccessToSingleForProjects(data Project, credentials shared.Credentials) (Project, error) {
	if &data == nil {
		return Project{}, errors.New("invalid Data")
	}

	if (data.CreatedBy == credentials.Id) && credentials.HasPersonalScopes {
		return data, nil
	}

	projectUsers := append(data.MemberIDs, data.ManagerIDs...)

	_, ok := helpers.CheckItemExists(projectUsers, credentials.Id)
	if ok && credentials.HasProjectScopes {
		return data, nil
	}

	return Project{}, errors.New("validation error")

}

func ValidateAccessToListForProjects(data []Project, credentials shared.Credentials) (Result, error) {
	if len(data) == 0 {
		return Result{}, errors.New("inValid Data")
	}

	validDataLength := 0

	var validData []interface{}

	for _, v := range data {
		if (v.CreatedBy == credentials.Id) && credentials.HasPersonalScopes {
			validDataLength++
			validData = append(validData, v)
		}

		projectUsers := append(v.MemberIDs, v.ManagerIDs...)

		_, ok := helpers.CheckItemExists(projectUsers, credentials.Id)
		if ok && credentials.HasProjectScopes {
			validDataLength++
			validData = append(validData, v)
		}

	}

	return Result{ValidData: validData, ValidDataLength: validDataLength}, nil
}

func InitRootRepositoryForProject(projectName string) repositories.CreateOneRepositoryDto {
	if projectName == "" {
		return repositories.CreateOneRepositoryDto{}
	}

	return repositories.CreateOneRepositoryDto{
		Name:               projectName,
		Description:        fmt.Sprintf("Root repository for %s", projectName),
		ParentRepositoryID: nil,
	}
}

func FindManyProjects(projectIDs []bson.ObjectId, credentials shared.Credentials) (Result, error) {
	var projects []Project
	ProjectsModel := shared.MongoSession.C("projects")

	err := ProjectsModel.Find(bson.M{"_id": bson.M{"$in": projectIDs}}).All(&projects)
	if err != nil {
		return Result{}, err
	}

	data, err := ValidateAccessToListForProjects(projects, credentials)
	if err != nil {
		return Result{}, err
	}

	return data, nil
}

func FindOneProject(findOneProjectDto *FindOneProjectDto, credentials shared.Credentials) (Project, error) {
	var project Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return Project{}, err
	}

	data, err := ValidateAccessToSingleForProjects(project, credentials)

	if err != nil {
		return Project{}, err
	}

	return data, nil
}

func CreateOneProject(createOneProjectDto *CreateOneProjectDto, credentials shared.Credentials) (bool, error) {
	createOneProjectDto.CreatedBy = credentials.Id

	createOneProjectDto.CreatedAt = time.Now()
	createOneProjectDto.UpdatedAt = time.Now()

	project := &Project{}
	test := (structs.Map(createOneProjectDto))
	err := mapstructure.Decode(test, project)
	if err != nil {
		return false, err
	}

	_, err = ValidateAccessToSingleForProjects(*project, credentials)
	if err != nil {
		return false, err
	}

	ProjectsModel := shared.MongoSession.C("projects")
	err = ProjectsModel.Insert(createOneProjectDto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateProject(findOneProjectDto *FindOneProjectDto, updateProjectDto *UpdateProjectDto, credentials shared.Credentials) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")

	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)
	if err != nil {
		return false, err
	}

	_ = mapstructure.Decode(data, project)

	updateProjectDto.UpdatedAt = time.Now()

	err = ProjectsModel.Update(findOneProjectDto, bson.M{"$set": updateProjectDto})
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteProject(deleteProjectDto *DeleteProjectDto, credentials shared.Credentials) (bool, []string, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(deleteProjectDto).One(&project)
	if err != nil {
		return false, []string{}, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)
	if err != nil {
		return false, []string{}, err
	}

	_ = mapstructure.Decode(data, project)

	userIDs := append(append(project.ManagerIDs, project.MemberIDs...), project.CreatedBy)

	err = ProjectsModel.Remove(deleteProjectDto)
	if err != nil {
		return false, []string{}, err
	}

	return true, userIDs, nil

}

func AddProjectMembers(findOneProjectDto *FindOneProjectDto, addProjectMembersDto *AddProjectMembersDto, credentials shared.Credentials) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)

	_ = mapstructure.Decode(data, project)

	newMemberIDs := addProjectMembersDto.MemberIDs

	for _, v := range newMemberIDs {
		err = ProjectsModel.Update(findOneProjectDto, bson.M{"$addToSet": bson.M{"memberIDs": v}})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func RemoveProjectMembers(findOneProjectDto *FindOneProjectDto, removeProjectMembersDto *RemoveProjectMembersDto, credentials shared.Credentials) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)

	_ = mapstructure.Decode(data, project)

	memberIDsNeededRemove := removeProjectMembersDto.MemberIDs

	for _, v := range memberIDsNeededRemove {
		err = ProjectsModel.Update(findOneProjectDto, bson.M{"$pull": bson.M{"memberIDs": v}})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func AddProjectManagers(findOneProjectDto *FindOneProjectDto, addProjectManagersDto *AddProjectManagersDto, credentials shared.Credentials) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")

	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)

	_ = mapstructure.Decode(data, project)

	newManagerIDs := addProjectManagersDto.ManagerIDs

	for _, v := range newManagerIDs {
		_, ok := helpers.CheckItemExists(project.MemberIDs, v)
		if !ok {
			return false, errors.New("invalid user")
		}
	}

	err = ProjectsModel.Update(findOneProjectDto, bson.M{
		"$addToSet": bson.M{"managerIDs": bson.M{"$each": newManagerIDs}},
		"$pullAll":  bson.M{"memberIDs": newManagerIDs},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveProjectManagers(findOneProjectDto *FindOneProjectDto, removeProjectManagersDto *RemoveProjectManagersDto, credentials shared.Credentials) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	data, err := ValidateAccessToSingleForProjects(*project, credentials)

	_ = mapstructure.Decode(data, project)

	managerIDsNeededRemove := removeProjectManagersDto.ManagerIDs

	for _, v := range managerIDsNeededRemove {
		err = ProjectsModel.Update(findOneProjectDto, bson.M{"$pull": bson.M{"managerIDs": v}})
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func CheckValidUsers(findOneProjectDto *FindOneProjectDto, users []*users.User) (bool, error) {
	var project *Project
	ProjectsModel := shared.MongoSession.C("projects")
	err := ProjectsModel.Find(findOneProjectDto).One(&project)
	if err != nil {
		return false, err
	}

	memberIDs := project.MemberIDs
	managerIds := project.ManagerIDs
	createdBy := project.CreatedBy

	projectUsers := append(append(memberIDs, managerIds...), createdBy)

	for _, v := range users {
		_, ok := helpers.CheckItemExists(projectUsers, v.ID)
		if ok {
			return false, errors.New("user already existed in Project")
		}
	}

	return true, nil

}
