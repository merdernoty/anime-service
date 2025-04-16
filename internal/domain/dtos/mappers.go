package dtos

import "github.com/merdernoty/anime-service/internal/domain/models"

func ToUserResponse(user models.User) UserResponseDTO {
	return UserResponseDTO{
		ID: 	  user.ID,
		NickName: user.Nickname,
		Email:    user.Email,	
		FirstName: user.Firstname,
		LastName:  user.Lastname,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
