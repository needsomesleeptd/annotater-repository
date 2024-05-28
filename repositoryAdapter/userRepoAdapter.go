package repo_adapter

import (
	"syscall"

	"github.com/needsomesleeptd/annotater-core/models"
	repository "github.com/needsomesleeptd/annotater-core/repositoryPorts"
	models_da "github.com/needsomesleeptd/annotater-repository/modelsDA"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepositoryAdapter struct {
	db *gorm.DB
}

func NewUserRepositoryAdapter(srcDB *gorm.DB) repository.IUserRepository {
	return &UserRepositoryAdapter{
		db: srcDB,
	}
}

func (repo *UserRepositoryAdapter) GetUserByID(id uint64) (*models.User, error) {
	var user_da models_da.User
	user_da.ID = id
	tx := repo.db.First(&user_da)
	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return nil, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}
	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := models_da.FromDaUser(&user_da)
	return &user, nil
}

func (repo *UserRepositoryAdapter) GetUserByLogin(login string) (*models.User, error) {
	var user_da models_da.User
	tx := repo.db.Where("login = ?", login).First(&user_da)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return nil, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}
	if tx.Error == gorm.ErrRecordNotFound {
		return nil, models.ErrNotFound
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error getting user by ID")
	}
	user := models_da.FromDaUser(&user_da)
	return &user, nil
}

func (repo *UserRepositoryAdapter) UpdateUserByLogin(login string, user *models.User) error {
	userDA := models_da.ToDaUser(*user)
	tx := repo.db.Model(&models_da.User{}).Where("login = ?", login).Updates(userDA)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) DeleteUserByLogin(login string) error {
	tx := repo.db.Where("login = ?", login).Delete(models_da.User{}) // specifically for gorm

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "Error in updating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) CreateUser(user *models.User) error {
	tx := repo.db.Create(models_da.ToDaUser(*user))

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error == gorm.ErrDuplicatedKey {
		return models.ErrDuplicateuserData
	}

	if tx.Error != nil {
		return errors.Wrap(tx.Error, "error in creating user")
	}
	return nil
}

func (repo *UserRepositoryAdapter) GetAllUsers() ([]models.User, error) {
	var usersDA []models_da.User
	tx := repo.db.Find(&usersDA)

	if errors.Is(tx.Error, syscall.ECONNREFUSED) {
		return nil, errors.Wrap(models.ErrDatabaseConnection, tx.Error.Error())
	}

	if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "Error in getting all users")
	}
	users := models_da.FromDaUserSlice(usersDA)
	return users, nil
}
