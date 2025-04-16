package dtos

import "time"

type CreateUserDTO struct {
	NickName  string `json:"nick_name" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email,unique"`
	Password  string `json:"password" validate:"required,min=6"`
}

type UpdateUserDTO struct { 
	ID 	  uint   `json:"id" validate:"required"`
	NickName  string `json:"nickname" validate:"required,min=6"`
	FirstName string `json:"firstname" validate:"required"`
	LastName  string `json:"lastname" validate:"required"`
	AvatarURL    string `json:"avatar_url" validate:"required"`
	Email     string `json:"email" validate:"required,email,unique"`
	Password  string `json:"password" validate:"required,min=6"`
}

type UserResponseDTO struct {
	ID        uint      `json:"id"`
	NickName  string    `json:"nickname"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstname,omitempty"`
	LastName  string    `json:"lastname,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginDTO struct {
	NickName string `json:"nickname" binding:"required_without=Email"`
	Email    string `json:"email" binding:"required_without=NickName,omitempty,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}