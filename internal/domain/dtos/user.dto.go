package dtos

import "time"

type CreateUserDTO struct {
	NickName  string `json:"nick_name" validate:"required,min=6" example:"johndoe123" swaggertype:"string"`
	FirstName string `json:"first_name" validate:"required" example:"John" swaggertype:"string"`
	LastName  string `json:"last_name" validate:"required" example:"Doe" swaggertype:"string"`
	Email     string `json:"email" validate:"required,email,unique" example:"john.doe@example.com" swaggertype:"string"`
	Password  string `json:"password" validate:"required,min=6" example:"StrongPass123!" swaggertype:"string"`
}

type UpdateUserDTO struct { 
	ID        uint   `json:"id" validate:"required" example:"1" swaggertype:"integer"`
	NickName  string `json:"nickname" validate:"required,min=6" example:"johndoe123" swaggertype:"string"`
	FirstName string `json:"firstname" validate:"required" example:"John" swaggertype:"string"`
	LastName  string `json:"lastname" validate:"required" example:"Doe" swaggertype:"string"`
	AvatarURL string `json:"avatar_url" validate:"required" example:"https://example.com/avatar.jpg" swaggertype:"string"`
	Email     string `json:"email" validate:"required,email,unique" example:"john.doe@example.com" swaggertype:"string"`
	Password  string `json:"password" validate:"required,min=6" example:"StrongPass123!" swaggertype:"string"`
}

type UserResponseDTO struct {
	ID        uint      `json:"id" example:"1" swaggertype:"integer"`
	NickName  string    `json:"nickname" example:"johndoe123" swaggertype:"string"`
	Email     string    `json:"email" example:"john.doe@example.com" swaggertype:"string"`
	FirstName string    `json:"firstname,omitempty" example:"John" swaggertype:"string"`
	LastName  string    `json:"lastname,omitempty" example:"Doe" swaggertype:"string"`
	AvatarURL string    `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg" swaggertype:"string"`
	CreatedAt time.Time `json:"created_at" example:"2024-04-28T10:30:00Z" swaggertype:"string"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-04-28T10:30:00Z" swaggertype:"string"`
}

type LoginDTO struct {
	NickName string `json:"nickname" validate:"required_without=Email" example:"johndoe123" swaggertype:"string"`
	Email    string `json:"email" validate:"required_without=NickName,omitempty,email" example:"john.doe@example.com" swaggertype:"string"`
	Password string `json:"password" validate:"required" example:"StrongPass123!" swaggertype:"string"`
}

type TokenResponseDTO struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." swaggertype:"string"`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." swaggertype:"string"`
	ExpiresIn    int64  `json:"expires_in" example:"3600" swaggertype:"integer"`
	TokenType    string `json:"token_type" example:"Bearer" swaggertype:"string"`
}


type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." swaggertype:"string"`
}


type ErrorResponse struct {
	Code    int    `json:"code" example:"400" swaggertype:"integer"`
	Message string `json:"message" example:"Invalid input" swaggertype:"string"`
}