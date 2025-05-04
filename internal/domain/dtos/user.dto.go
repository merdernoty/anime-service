package dtos

// UpdateUserDTO модель для обновления данных пользователя
// Все поля необязательные, обновляются только указанные поля
type UserUpdateDTO struct {
	NickName  *string `json:"nickname,omitempty" example:"johndoe123" validate:"omitempty,min=6"`
	Email     *string `json:"email,omitempty" example:"john.doe@example.com" validate:"omitempty,email"`
	Avatar    *string `json:"avatar,omitempty" example:"https://example.com/avatar.jpg" validate:"omitempty,url"`
	LastName  *string `json:"lastname,omitempty" example:"Doe" validate:"omitempty"`
	FirstName *string `json:"firstname,omitempty" example:"John" validate:"omitempty"`
}