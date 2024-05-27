package storage

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/needsomesleeptd/annotater-core/models"
	repository "github.com/needsomesleeptd/annotater-core/repositoryPorts"
	"github.com/pkg/errors"
)

type DocumentDataRepositoryAdapter struct {
	root          string
	fileExtension string //it is optional
}

func NewReportRepositoryAdapter(rootSrc string, ext string) repository.IDocumentDataRepository {
	return &DocumentDataRepositoryAdapter{
		root:          rootSrc,
		fileExtension: ext,
	}
}

func (repo *DocumentDataRepositoryAdapter) MakeDir(dir string) error {
	dirPath := fmt.Sprintf("%s/%s", repo.root, dir) + repo.fileExtension
	return os.MkdirAll(dirPath, 0755)
}

func (repo *DocumentDataRepositoryAdapter) Exists(path string) bool {
	fullPath := fmt.Sprintf("%s/%s", repo.root, path) + repo.fileExtension
	_, err := os.Stat(fullPath)

	return !os.IsNotExist(err)
}

func (repo *DocumentDataRepositoryAdapter) AddDocument(doc *models.DocumentData) error {
	if !repo.Exists(repo.root) {
		err := repo.MakeDir(repo.root)
		if err != nil {
			return errors.Wrap(err, "error in saving document data")
		}
	}

	filePath := fmt.Sprintf("%s/%s", repo.root, doc.ID) + repo.fileExtension

	err := os.WriteFile(filePath, doc.DocumentBytes, 0644)
	if err != nil {
		return errors.Wrap(err, "error in saving document data")
	}

	return nil
}

func (repo *DocumentDataRepositoryAdapter) DeleteDocumentByID(id uuid.UUID) error {
	filePath := fmt.Sprintf("%s/%s", repo.root, id) + repo.fileExtension
	err := os.Remove(filePath)
	if err != nil {
		return errors.Wrap(err, "error in deleting document data")
	}

	return nil
}

func (repo *DocumentDataRepositoryAdapter) GetDocumentByID(id uuid.UUID) (*models.DocumentData, error) {
	filePath := fmt.Sprintf("%s/%s", repo.root, id) + repo.fileExtension
	fileBytes, err := os.ReadFile(filePath)

	if err == os.ErrNotExist {
		return nil, models.ErrNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "error getting file")
	}

	documentData := &models.DocumentData{
		DocumentBytes: fileBytes,
		ID:            id,
	}
	return documentData, nil
}
