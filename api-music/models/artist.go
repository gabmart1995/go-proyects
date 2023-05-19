package models

import "github.com/go-playground/validator/v10"

type Artist struct {
	Id          int `json:"id"`
	AlbumId     int
	Name        string `json:"name" validate:"required"`
	Alias       string `json:"alias" validate:"required"`
	MusicGender string `json:"music_gender" validate:"required"`
}

var Artists = []Artist{
	{
		Id:          1,
		Name:        "Hector Velazquez",
		Alias:       "Dimension Latina",
		MusicGender: "salsa",
	},
	{
		Id:          2,
		Name:        "Emilio Garante",
		Alias:       "Billos, Caracas Boys",
		MusicGender: "bolero",
	},
	{
		Id:          3,
		Name:        "Freddie Mercury",
		Alias:       "Queen",
		MusicGender: "rock",
	},
}

func (artist *Artist) Validate() []*ErrorResponse {
	var (
		errors []*ErrorResponse
		err    error
	)

	validate := validator.New()

	if err = validate.StructExcept(artist, "Artist.Id"); err != nil {
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
