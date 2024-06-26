package repo_adapter

import (
	"syscall"

	"github.com/google/uuid"
	"github.com/needsomesleeptd/annotater-core/models"
	repository "github.com/needsomesleeptd/annotater-core/repositoryPorts"
	models_da "github.com/needsomesleeptd/annotater-repository/modelsDA"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DocumentMetaDataRepositoryAdapter struct {
	db *gorm.DB
}

func NewDocumentRepositoryAdapter(srcDB *gorm.DB) repository.IDocumentMetaDataRepository {
	return &DocumentMetaDataRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *DocumentMetaDataRepositoryAdapter) AddDocument(doc *models.DocumentMetaData) error {
	tx := repo.db.Create(models_da.ToDaDocument(*doc))

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating document")
	}
	return nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentByID(id uuid.UUID) (*models.DocumentMetaData, error) {
	var documentDa models_da.Document
	documentDa.ID = id
	tx := repo.db.First(&documentDa)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return nil, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting document by ID")
	}
	document := models_da.FromDaDocument(&documentDa)
	return &document, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) DeleteDocumentByID(id uuid.UUID) error {
	tx := repo.db.Delete(models.DocumentMetaData{}, id) // specifically for gorm

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in deleting document")
	}
	return nil
}
func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentsByCreatorID(id uint64) ([]models.DocumentMetaData, error) {
	var documentsDA []models_da.Document
	tx := repo.db.Where("creator_id = ?", id).Find(&documentsDA)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return nil, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting documents by creator")
	}
	documents := models_da.FromDaDocumentSlice(documentsDA)
	return documents, nil
}

func (repo *DocumentMetaDataRepositoryAdapter) GetDocumentCountByCreator(id uint64) (int64, error) {
	var count int64
	tx := repo.db.Model(models_da.Document{}).Where("creator_id = ?", id).Count(&count)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return -1, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return -1, errors.Wrap(tx.Error, "Error in getting count by creator")
	}
	return count, nil
}
