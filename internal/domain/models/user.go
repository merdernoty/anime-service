package models

import (
	"emperror.dev/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrPasswordEmpty = errors.New("password is empty")
	ErrPasswordHashing = errors.New("failed to hash password")
)

type User struct {
	gorm.Model           
	Nickname  string `gorm:"unique;not null" json:"nickname"`
	Firstname string `json:"firstname,omitempty"`
	Lastname  string `json:"lastname,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"` 
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
}

func (u *User) HashPassword() error {
	if u.Password == "" {
		return ErrPasswordEmpty
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashing
	}

	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	if u.Password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

