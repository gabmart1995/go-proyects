package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// types
type User struct {
	Id        int       `json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `json:"name" validate:"required,min=3,max=30"`
	Surname   string    `json:"surname" validate:"required,min=3,max=30"`
	Email     string    `json:"email" validate:"required,email,min=3,max=30" gorm:"unique"`
	Password  string    `json:"password" validate:"omitempty,min=8"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// campo que representa la relaciones
	Album Album `gorm:"foreignKey:UserId"`
}

type mapCallback func(user User, index int, slice []User) User
type filterCallback func(user User, index int, slice []User) bool

// calback que mapea y genera un nuevo slice de usuarios
func (user *User) Map(users []User, callback mapCallback) []User {
	var result []User

	for index, user := range users {
		result = append(result, callback(user, index, users))
	}

	return result
}

// callback que filtra y genera un nuevo slice de usuarios
func (user *User) Filter(users []User, callback filterCallback) []User {
	var result []User

	for index, user := range users {
		if callback(user, index, users) {
			result = append(result, user)
		}
	}

	return result
}

func (user *User) ValidateStruct() []*ErrorResponse {

	var (
		errors []*ErrorResponse
		err    error
	)

	validate := validator.New()

	if err = validate.StructExcept(user, "User.Id"); err != nil {
		// cambia la interface del error por uno personalizado
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()

			errors = append(errors, &element)
		}
	}

	return errors
}
