package models

import "github.com/go-playground/validator/v10"

type Song struct {
	Id      int    `json:"id" gorm:"primaryKey,autoIncrement"`
	AlbumId int    // foreign key
	Name    string `json:"name" validate:"required,min=2,max=30"`
	Url     string `json:"url"`
}

type mapSongCallback func(song Song, index int, slice []Song) Song
type filterSongCallback func(song Song, index int, slice []Song) bool

var Songs = []Song{
	{
		Id:   1,
		Name: "I can't get GO (Remastered)",
		Url:  "http://localhost:3000/music_1",
	},
	{
		Id:   2,
		Name: "Whising on the star",
		Url:  "http://localhost:3000/music_2",
	},
	{
		Id:   3,
		Name: "Don't Say Goodbye",
		Url:  "http://localhost:3000/music_3",
	},
	{
		Id:   4,
		Name: "Aesthetic Sounds",
		Url:  "http://localhost:3000/music_4",
	},
}

// calback que mapea y genera un nuevo slice de usuarios
func (song *Song) MapSlice(songs []Song, callback mapSongCallback) []Song {
	var result []Song

	for index, song := range songs {
		result = append(result, callback(song, index, songs))
	}

	return result
}

// callback que filtra y genera un nuevo slice de usuarios
func (song *Song) FilterSlice(songs []Song, callback filterSongCallback) []Song {
	var result []Song

	for index, song := range Songs {
		if callback(song, index, Songs) {
			result = append(result, song)
		}
	}

	return result
}

func (song *Song) ValidateStruct() []*ErrorResponse {

	var (
		errors []*ErrorResponse
		err    error
	)

	validate := validator.New()

	if err = validate.StructExcept(song, "Song.Id", "Song.Url"); err != nil {
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
