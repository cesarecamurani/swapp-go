package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) ports.UserRepository {
	return &UserGormRepository{db: db}
}

func toUserModel(user *domain.User) *models.UserModel {
	id := user.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &models.UserModel{
		ID:       id,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Phone:    user.Phone,
		Address:  user.Address,
	}
}

func toDomainUser(model *models.UserModel) *domain.User {
	return &domain.User{
		ID:       model.ID,
		Username: model.Username,
		Password: model.Password,
		Email:    model.Email,
		Phone:    model.Phone,
		Address:  model.Address,
	}
}

func (userGorm *UserGormRepository) CreateUser(user *domain.User) error {
	model := toUserModel(user)

	if result := userGorm.db.Create(model); result.Error != nil {
		return result.Error
	}

	user.ID = model.ID

	return nil
}

func (userGorm *UserGormRepository) UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	if err := userGorm.db.Model(&models.UserModel{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return nil, err
	}

	var updatedUserModel models.UserModel
	if err := userGorm.db.Where("id = ?", id).First(&updatedUserModel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&updatedUserModel), nil
}

func (userGorm *UserGormRepository) DeleteUser(id uuid.UUID) error {
	return userGorm.db.Delete(&models.UserModel{}, id).Error
}

func (userGorm *UserGormRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var usermodel models.UserModel

	if err := userGorm.db.First(&usermodel, id).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}

func (userGorm *UserGormRepository) GetUserByUsername(username string) (*domain.User, error) {
	var usermodel models.UserModel

	if err := userGorm.db.Where("username = ?", username).First(&usermodel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}

func (userGorm *UserGormRepository) GetUserByEmail(email string) (*domain.User, error) {
	var usermodel models.UserModel

	if err := userGorm.db.Where("email = ?", email).First(&usermodel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}
