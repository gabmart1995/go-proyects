package song

import (
	"api-music/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func Get(c *fiber.Ctx) error {
	idSong, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	for _, song := range models.Songs {
		if song.Id == idSong {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"ok":      true,
				"status":  http.StatusOK,
				"data":    song,
				"message": nil,
			})
		}
	}

	return c.Status(http.StatusBadRequest).JSON(fiber.Map{
		"ok":      false,
		"status":  http.StatusBadRequest,
		"data":    nil,
		"message": "Song not found",
	})
}

func GetAll(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    models.Songs,
		"message": nil,
	})
}

func Create(c *fiber.Ctx) error {
	song := models.Song{
		Id:  len(models.Songs) + 1,
		Url: "http://localhost:3000/music_" + strconv.Itoa(len(models.Songs)+1),
	}

	if err := c.BodyParser(&song); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in parsing JSON",
		})
	}

	if errors := song.ValidateStruct(); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"message": "Error in entries fields",
			"data":    errors,
		})
	}

	models.Songs = append(models.Songs, song)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusCreated,
		"data":    nil,
		"message": "Song created successfully",
	})
}

func Update(c *fiber.Ctx) error {
	idSong, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	songModel := models.Song{
		Id:  idSong,
		Url: "http://localhost:3000/music_" + strconv.Itoa(idSong),
	}

	if err := c.BodyParser(&songModel); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in parsing JSON",
		})
	}

	if errors := songModel.ValidateStruct(); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"message": "Error in entries fields",
			"data":    errors,
		})
	}

	models.Songs = songModel.MapSlice(
		models.Songs,
		func(song models.Song, index int, songs []models.Song) models.Song {

			if song.Id == idSong {
				return songModel
			}

			return song
		},
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"message": "Song update succesffully",
	})
}

func Delete(c *fiber.Ctx) error {
	idSong, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in get id song",
		})
	}

	song := models.Song{}

	models.Songs = song.FilterSlice(
		models.Songs,
		func(item models.Song, index int, slice []models.Song) bool {
			return item.Id != idSong
		},
	)

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    nil,
		"message": "Song delete successfully",
	})
}
