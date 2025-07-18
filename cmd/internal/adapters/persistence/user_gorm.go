package persistence

import (
	"gorm.io/gorm"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/config"
	"swapp-go/cmd/internal/domain"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() ports.UserRepository {
	return &GormUserRepository{
		db: config.DB,
	}
}

type UserModel struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
}

func (UserModel) TableName() string {
	return "users"
}

func toUserModel(user *domain.User) *UserModel {
	return &UserModel{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
}

func toDomainUser(model *UserModel) *domain.User {
	return &domain.User{
		ID:       model.ID,
		Username: model.Username,
		Password: model.Password,
		Email:    model.Email,
	}
}

func (gur *GormUserRepository) CreateUser(user *domain.User) error {
	model := toUserModel(user)

	result := gur.db.Create(model)
	if result.Error != nil {
		return result.Error
	}

	user.ID = model.ID

	return nil
}

func (gur *GormUserRepository) GetUserByID(id uint) (*domain.User, error) {
	var usermodel UserModel

	if err := gur.db.First(&usermodel, id).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}

func (gur *GormUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var usermodel UserModel

	if err := gur.db.Where("username = ?", username).First(&usermodel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}
