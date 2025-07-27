package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
	"time"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) ports.UserRepository {
	return &GormUserRepository{db: db}
}

type UserModel struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}

func toUserModel(user *domain.User) *UserModel {
	id := user.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &UserModel{
		ID:       id,
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

func (gormUser *GormUserRepository) CreateUser(user *domain.User) error {
	model := toUserModel(user)

	result := gormUser.db.Create(model)
	if result.Error != nil {
		return result.Error
	}

	user.ID = model.ID

	return nil
}

func (gormUser *GormUserRepository) UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	if err := gormUser.db.Model(&UserModel{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return nil, err
	}

	var updatedModel UserModel
	if err := gormUser.db.Where("id = ?", id).First(&updatedModel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&updatedModel), nil
}

func (gormUser *GormUserRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var usermodel UserModel

	if err := gormUser.db.First(&usermodel, id).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}

func (gormUser *GormUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var usermodel UserModel

	if err := gormUser.db.Where("username = ?", username).First(&usermodel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}

func (gormUser *GormUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var usermodel UserModel

	if err := gormUser.db.Where("email = ?", email).First(&usermodel).Error; err != nil {
		return nil, err
	}

	return toDomainUser(&usermodel), nil
}
