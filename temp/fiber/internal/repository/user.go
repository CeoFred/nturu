package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/nturu/microservice-template/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (a *UserRepository) UsersAMonthAgo() ([]*models.User, error) {
	var users []*models.User
	query := `
		SELECT *
		FROM users
		WHERE created_at >= $1
	`
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	err := a.database.Raw(query, oneMonthAgo).Scan(&users).Error
	return users, err
}

func (a *UserRepository) AllUsers() ([]*models.User, error) {
	var users []*models.User
	err := a.database.Raw("select * from users").Scan(&users).Error
	return users, err
}

func (a *UserRepository) FindRecordsByCondition(condition, value string) ([]*models.User, error) {
	var users []*models.User
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM users WHERE %s = ?`, condition), value).Scan(&users).Error
	return users, err

}

func (a *UserRepository) FindUserByCondition(condition, value string) (*models.User, bool, error) {
	var user *models.User
	err := a.database.Raw(fmt.Sprintf(`SELECT * FROM users WHERE %s = ?`, condition), value).Scan(&user).Error
	if err != nil {
		return nil, false, err
	}
	if user != nil {
		return user, true, nil
	}
	return nil, false, nil
}

func (a *UserRepository) FindByAccountType(value string) ([]*models.User, bool, error) {
	var user []*models.User
	err := a.database.Raw(`SELECT * FROM users WHERE account_type = ?`, value).Scan(&user).Error
	if err != nil {
		return nil, false, err
	}
	if user != nil {
		return user, true, nil
	}
	return nil, false, nil
}

func (a *UserRepository) CreateUser(user *models.User) error {
	return a.database.Model(&models.User{}).Create(user).Error
}

func (a *UserRepository) UpdateUserByCondition(condition, value string, update *models.User) (*models.User, error) {
	user := &models.User{}
	rows := a.database.Model(user).Where(fmt.Sprintf(`%s = ?`, condition), value).Updates(&update).First(user)
	if rows.RowsAffected == 0 {
		return nil, errors.New("no record updated")
	}
	return user, nil
}
