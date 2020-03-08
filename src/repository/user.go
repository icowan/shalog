package repository

import (
	"github.com/icowan/shalom/src/repository/types"
	"github.com/jinzhu/gorm"
	"time"
)

type UserRepository interface {
	Find(username string) (res *types.User, err error)
	FindAndPwd(username, password string) (res *types.User, err error)
	FindById(id int64) (res types.User, err error)
	Update(user *types.User) (err error)
	Create(username, password, email string) error
}

type user struct {
	db *gorm.DB
}

func (c *user) Create(username, password, email string) error {
	now := time.Now()
	return c.db.Save(&types.User{
		Username:            username,
		UsernameCanonical:   username,
		Email:               email,
		EmailCanonical:      email,
		Enabled:             1,
		Salt:                "",
		Password:            password,
		ConfirmationToken:   "",
		LastLogin:           now,
		ExpiresAt:           now,
		PasswordRequestedAt: now,
		CredentialsExpireAt: now,
	}).Error
}

func (c *user) Update(user *types.User) (err error) {
	return c.db.Model(&types.User{}).Where("id = ?", user.ID).Update(user).Error
}

func (c *user) FindById(id int64) (res types.User, err error) {
	err = c.db.Select("id, username, username_canonical, email, email_canonical, created_at, updated_at, last_login ").First(&res, "id = ?", id).Error
	return
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &user{db: db}
}

func (c *user) Find(username string) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "username_canonical = ?", username).Error
	return &rs, err
}

func (c *user) FindAndPwd(username, password string) (res *types.User, err error) {
	var rs types.User
	err = c.db.First(&rs, "username_canonical = ? AND password = ?", username, password).Error
	return &rs, err
}
