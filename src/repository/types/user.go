/**
 * @Time : 2019-09-06 11:43
 * @Author : solacowa@gmail.com
 * @File : user
 * @Software: GoLand
 */

package types

import (
	"time"
)

type User struct {
	ID                  int       `gorm:"column:id;primary_key" json:"id"`
	Username            string    `gorm:"column:username" json:"username"`
	UsernameCanonical   string    `gorm:"column:username_canonical" json:"username_canonical"`
	Email               string    `gorm:"column:email" json:"email"`
	EmailCanonical      string    `gorm:"column:email_canonical" json:"email_canonical"`
	Enabled             int       `gorm:"column:enabled" json:"enabled"`
	Salt                string    `gorm:"column:salt" json:"salt"`
	Password            string    `gorm:"column:password" json:"password"`
	LastLogin           time.Time `gorm:"column:last_login" json:"last_login"`
	Locked              int       `gorm:"column:locked" json:"locked"`
	Expired             int       `gorm:"column:expired" json:"expired"`
	ExpiresAt           time.Time `gorm:"column:expires_at" json:"expires_at"`
	ConfirmationToken   string    `gorm:"column:confirmation_token" json:"confirmation_token"`
	PasswordRequestedAt time.Time `gorm:"column:password_requested_at" json:"password_requested_at"`
	Roles               string    `gorm:"column:roles" json:"roles"`
	CredentialsExpired  int       `gorm:"column:credentials_expired" json:"credentials_expired"`
	CredentialsExpireAt time.Time `gorm:"column:credentials_expire_at" json:"credentials_expire_at"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (p *User) TableName() string {
	return "users"
}
