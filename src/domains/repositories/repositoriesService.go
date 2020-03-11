package repositories

import (
	"errors"
	"time"

	"github.com/fatih/structs"
	"github.com/khoa5773/go-server/src/shared"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

func FindManyRepositories(findManyRepositoriesDto FindManyRepositoriesDto, credentials shared.Credentials) (shared.Result, error) {
	var repositories []map[string]interface{}
	RepositoriesModel := shared.MongoSession.C("repositories")

	err := RepositoriesModel.Find(findManyRepositoriesDto).All(&repositories)
	if err != nil {
		return shared.Result{}, err
	}

	data, err := shared.ValidateAccessToList(repositories, credentials)
	if err != nil {
		return shared.Result{}, err
	}

	return *data, nil
}

func FindOneRepository(findOneRepositoryDto *FindOneRepositoryDto, credentials shared.Credentials) (Repository, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(findOneRepositoryDto).One(&repository)
	if err != nil {
		return Repository{}, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return Repository{}, err
	}

	err = mapstructure.Decode(data, &repository)
	if err != nil {
		return Repository{}, err
	}

	return repository, nil
}

func CreateOneRepository(createOneRepositoryDto *CreateOneRepositoryDto, credentials shared.Credentials) (bool, error) {
	createOneRepositoryDto.CreatedBy = credentials.Id
	createOneRepositoryDto.UpdatedBy = credentials.Id
	createOneRepositoryDto.CreatedAt = time.Now()
	createOneRepositoryDto.UpdatedAt = time.Now()

	var parentRepositoryID bson.ObjectId
	if createOneRepositoryDto.ParentRepositoryID != nil {
		parentRepositoryID = *createOneRepositoryDto.ParentRepositoryID
	}

	_, err := CheckRepositoriesExists([]bson.ObjectId{parentRepositoryID}, createOneRepositoryDto.ProjectID)
	if err != nil {
		return false, err
	}

	repository := &Repository{}
	err = mapstructure.Decode(structs.Map(createOneRepositoryDto), repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	createOneRepositoryDto.ID = bson.NewObjectId()
	RepositoriesModel := shared.MongoSession.C("repositories")
	err = RepositoriesModel.Insert(createOneRepositoryDto)
	if err != nil {
		return false, err
	}

	credentials.IsAdmin = true
	if createOneRepositoryDto.ParentRepositoryID != nil {
		_, err = AddChildRepositoryID(*createOneRepositoryDto.ParentRepositoryID, createOneRepositoryDto.ID, repository.ProjectID, credentials)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func UpdateRepository(findOneRepositoryDto *FindOneRepositoryDto, updateRepositoryDto *UpdateRepositoryDto, credentials shared.Credentials) (bool, error) {
	var repository *Repository
	RepositoriesModel := shared.MongoSession.C("repositories")

	_, err := CheckRepositoriesExists([]bson.ObjectId{updateRepositoryDto.ParentRepositoryID}, findOneRepositoryDto.ProjectID)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Find(findOneRepositoryDto).One(&repository)
	if err != nil {
		return false, err
	}

	if repository.ParentRepositoryID == "" && updateRepositoryDto.ParentRepositoryID != "" {
		return false, errors.New("you can't move the root repository")
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(*repository), credentials)
	if err != nil {
		return false, err
	}

	_ = mapstructure.Decode(data, repository)

	updateRepositoryDto.UpdatedAt = time.Now()
	updateRepositoryDto.UpdatedBy = credentials.Id

	err = RepositoriesModel.Update(findOneRepositoryDto, bson.M{"$set": updateRepositoryDto})
	if err != nil {
		return false, err
	}

	if updateRepositoryDto.ParentRepositoryID == "" {
		return true, nil
	}

	if repository.ParentRepositoryID == "" {
		_, err = AddChildRepositoryID(updateRepositoryDto.ParentRepositoryID, repository.ID, findOneRepositoryDto.ProjectID, credentials)
		if err != nil {
			return false, err
		}
	}

	_, err = RemoveChildRepositoryID(repository.ParentRepositoryID, repository.ID, findOneRepositoryDto.ProjectID, credentials)
	if err != nil {
		return false, err
	}

	_, err = AddChildRepositoryID(updateRepositoryDto.ParentRepositoryID, repository.ID, findOneRepositoryDto.ProjectID, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteRepository(deleteRepositoryDto *DeleteRepositoryDto, credentials shared.Credentials) (bool, error) {
	var repository *Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(deleteRepositoryDto).One(&repository)
	if err != nil {
		return false, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(*repository), credentials)
	if err != nil {
		return false, err
	}

	err = mapstructure.Decode(data, repository)
	if err != nil {
		return false, err
	}

	if repository.ParentRepositoryID == "" {
		return false, errors.New("you can not delete the root repository")
	}

	_, err = DeleteManyRepositoriesAndDocumentInRepository(repository.ID, repository.ProjectID, credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Remove(deleteRepositoryDto)
	if err != nil {
		return false, err
	}

	_, err = RemoveChildRepositoryID(repository.ParentRepositoryID, repository.ID, repository.ProjectID, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CheckRepositoriesExists(repositoryIDs []bson.ObjectId, projectID bson.ObjectId) (bool, error) {
	if repositoryIDs == nil || len(repositoryIDs) == 0 || repositoryIDs[0] == "" {
		return true, nil
	}

	var repositories []Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(bson.M{"_id": bson.M{"$in": repositoryIDs}, "projectID": projectID}).All(&repositories)
	if err != nil {
		return false, err
	}

	if len(repositories) != len(repositoryIDs) {
		return false, errors.New("invalid repositories")
	}

	return true, nil
}

func AddChildRepositoryID(repositoryID, childRepositoryID, projectID bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}).One(&repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Update(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}, bson.M{
		"$addToSet": bson.M{
			"childRepositoryIDs": childRepositoryID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": credentials.Id,
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveChildRepositoryID(repositoryID, childRepositoryID, projectID bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}).One(&repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Update(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}, bson.M{
		"$pull": bson.M{
			"childRepositoryIDs": childRepositoryID,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": credentials.Id,
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteManyRepositoriesAndDocumentInRepository(repositoryID, projectID bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repositories []Repository
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")

	err := RepositoriesModel.Find(FindManyRepositoriesDto{ParentRepositoryID: repositoryID, ProjectID: projectID}).All(&repositories)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Find(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}).One(&repository)
	if err != nil {
		return false, err
	}

	repositoriesMap := []map[string]interface{}{}
	for _, val := range repositories {
		repositoriesMap = append(repositoriesMap, structs.Map(val))
	}

	_, err = shared.ValidateAccessToList(repositoriesMap, credentials)
	if err != nil {
		return false, err
	}

	for _, repo := range repositories {
		_, err = DeleteManyRepositoriesAndDocumentInRepository(repo.ID, projectID, credentials)
		if err != nil {
			return false, err
		}
	}

	_, err = RepositoriesModel.RemoveAll(FindManyRepositoriesDto{ParentRepositoryID: repositoryID, ProjectID: projectID})
	if err != nil {
		return false, err
	}

	_, err = DeleteManyDocuments(repository.DocumentIDs, repositoryID, projectID, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteRootRepositoryOfProject(projectID bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	findOneRepositoryDto := FindOneRepositoryDto{ProjectID: projectID, ParentRepositoryID: nil}
	err := RepositoriesModel.Find(findOneRepositoryDto).One(&repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	_, err = DeleteManyRepositoriesAndDocumentInRepository(repository.ID, projectID, credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Remove(findOneRepositoryDto)
	if err != nil {
		return false, err
	}

	return true, nil
}

func AddDocumentIDs(repositoryID, projectID bson.ObjectId, documentIDs []bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}).One(&repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Update(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}, bson.M{
		"$addToSet": bson.M{
			"documentIDs": bson.M{
				"$each": documentIDs,
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": credentials.Id,
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveDocumentIDs(repositoryID, projectID bson.ObjectId, documentIDs []bson.ObjectId, credentials shared.Credentials) (bool, error) {
	var repository Repository
	RepositoriesModel := shared.MongoSession.C("repositories")
	err := RepositoriesModel.Find(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}).One(&repository)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(repository), credentials)
	if err != nil {
		return false, err
	}

	err = RepositoriesModel.Update(FindOneRepositoryDto{ID: repositoryID, ProjectID: projectID}, bson.M{
		"$pullAll": bson.M{
			"documentIDs": documentIDs,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
			"updatedBy": credentials.Id,
		},
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteManyDocuments(documentIDs []bson.ObjectId, repositoryID, projectID bson.ObjectId, credentials shared.Credentials) (bool, error) {
	DocumentsModel := shared.MongoSession.C("documents")
	_, err := DocumentsModel.RemoveAll(bson.M{
		"_id": bson.M{
			"$in": documentIDs,
		},
		"repositoryID": repositoryID,
		"projectID":    projectID,
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
