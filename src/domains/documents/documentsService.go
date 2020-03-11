package documents

import (
	"context"
	"time"

	"github.com/fatih/structs"
	cloudstorage "github.com/khoa5773/go-server/src/domains/cloudStorage"
	"github.com/khoa5773/go-server/src/domains/repositories"
	"github.com/khoa5773/go-server/src/shared"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/mgo.v2/bson"
)

func FindManyDocumentsInRepository(findManyDocumentDto FindManyDocumentsDto, credentials shared.Credentials) (shared.Result, error) {
	var documents []map[string]interface{}
	DocumentsModel := shared.MongoSession.C("documents")

	err := DocumentsModel.Find(findManyDocumentDto).All(&documents)
	if err != nil {
		return shared.Result{}, err
	}

	data, err := shared.ValidateAccessToList(documents, credentials)
	if err != nil {
		return shared.Result{}, err
	}

	for _, val := range data.ValidData {
		val["path"], err = cloudstorage.GenerateSignedUrl(val["_id"].(bson.ObjectId).Hex(), val["projectID"].(bson.ObjectId).Hex(), val["mimeType"].(string))
		if err != nil {
			return shared.Result{}, err
		}
	}

	return *data, nil
}

func FindOneDocument(findOneDocumentDto *FindOneDocumentDto, credentials shared.Credentials) (Document, error) {
	var document Document
	DocumentsModel := shared.MongoSession.C("documents")
	err := DocumentsModel.Find(findOneDocumentDto).One(&document)
	if err != nil {
		return Document{}, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(document), credentials)
	if err != nil {
		return Document{}, err
	}

	err = mapstructure.Decode(data, &document)
	if err != nil {
		return Document{}, err
	}

	document.Path, err = cloudstorage.GenerateSignedUrl(document.ID.Hex(), document.ProjectID.Hex(), document.MimeType)
	if err != nil {
		return Document{}, err
	}

	return document, nil
}

func CreateOneDocument(createOneDocumentDto *CreateOneDocumentDto, credentials shared.Credentials) (bool, error) {
	createOneDocumentDto.CreatedBy = credentials.Id
	createOneDocumentDto.UpdatedBy = credentials.Id
	createOneDocumentDto.CreatedAt = time.Now()
	createOneDocumentDto.UpdatedAt = time.Now()

	_, err := repositories.CheckRepositoriesExists([]bson.ObjectId{createOneDocumentDto.RepositoryID}, createOneDocumentDto.ProjectID)
	if err != nil {
		return false, err
	}

	document := &Document{}
	err = mapstructure.Decode(structs.Map(createOneDocumentDto), document)
	if err != nil {
		return false, err
	}

	_, err = shared.ValidateAccessToSingle(structs.Map(document), credentials)
	if err != nil {
		return false, err
	}

	DocumentsModel := shared.MongoSession.C("documents")
	err = DocumentsModel.Insert(createOneDocumentDto)
	if err != nil {
		return false, err
	}

	credentials.IsAdmin = true
	_, err = repositories.AddDocumentIDs(createOneDocumentDto.RepositoryID, createOneDocumentDto.ProjectID, []bson.ObjectId{createOneDocumentDto.ID}, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateDocument(findOneDocumentDto *FindOneDocumentDto, updateDocumentDto *UpdateDocumentDto, credentials shared.Credentials) (bool, error) {
	var document *Document
	DocumentsModel := shared.MongoSession.C("documents")

	_, err := repositories.CheckRepositoriesExists([]bson.ObjectId{updateDocumentDto.RepositoryID}, findOneDocumentDto.ProjectID)
	if err != nil {
		return false, err
	}

	err = DocumentsModel.Find(findOneDocumentDto).One(&document)
	if err != nil {
		return false, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(*document), credentials)
	if err != nil {
		return false, err
	}

	_ = mapstructure.Decode(data, document)

	updateDocumentDto.UpdatedAt = time.Now()
	updateDocumentDto.UpdatedBy = credentials.Id

	err = DocumentsModel.Update(findOneDocumentDto, bson.M{"$set": updateDocumentDto})
	if err != nil {
		return false, err
	}

	if updateDocumentDto.RepositoryID == "" {
		return true, nil
	}

	_, err = repositories.RemoveDocumentIDs(document.RepositoryID, document.ProjectID, []bson.ObjectId{document.ID}, credentials)
	if err != nil {
		return false, err
	}

	_, err = repositories.AddDocumentIDs(updateDocumentDto.RepositoryID, document.ProjectID, []bson.ObjectId{document.ID}, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteDocument(deleteDocumentDto *DeleteDocumentDto, credentials shared.Credentials) (bool, error) {
	var document *Document
	DocumentsModel := shared.MongoSession.C("documents")
	err := DocumentsModel.Find(deleteDocumentDto).One(&document)
	if err != nil {
		return false, err
	}

	data, err := shared.ValidateAccessToSingle(structs.Map(*document), credentials)
	if err != nil {
		return false, err
	}

	err = mapstructure.Decode(data, document)
	if err != nil {
		return false, err
	}

	err = DocumentsModel.Remove(deleteDocumentDto)
	if err != nil {
		return false, err
	}

	_, err = cloudstorage.RemoveFileFromCloudStorage(context.Background(), document.ID.Hex(), document.ProjectID.Hex(), document.MimeType)
	if err != nil {
		return false, err
	}

	_, err = repositories.RemoveDocumentIDs(document.RepositoryID, document.ProjectID, []bson.ObjectId{document.ID}, credentials)
	if err != nil {
		return false, err
	}

	return true, nil
}
