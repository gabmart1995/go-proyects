package user

import (
	"api-music/helpers"
	"api-music/middlewares"
	"api-music/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Get(c *fiber.Ctx, db *gorm.DB) error {

	idUser, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	// busca un registro en particular
	user := models.User{Id: idUser}
	result := db.First(&user)

	// valida si llega mas de un registro
	if result.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"data":    nil,
			"message": "User not found",
		})
	}

	if result.Error != nil {
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    user,
		"message": nil,
	})
}

func GetAll(c *fiber.Ctx, db *gorm.DB) error {

	var users []models.User
	result := db.Scopes(helpers.Paginate(c, "start")).Find(&users)

	if result.RowsAffected == 0 {
		return c.Status(http.StatusOK).JSON(fiber.Map{
			"ok":      true,
			"status":  http.StatusOK,
			"data":    users,
			"message": nil,
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    users,
		"message": nil,
	})
}

func Create(c *fiber.Ctx, db *gorm.DB) error {

	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in parsing JSON",
		})
	}

	if errors := user.ValidateStruct(); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"message": "Error in entries fields",
			"data":    errors,
		})
	}

	// ciframos la password
	cipherPassword, err := helpers.HashPassword(user.Password, 10)

	if err != nil {
		return err
	}

	user.Password = cipherPassword

	// inserta en la Base de datos
	result := db.Create(&user)

	if result.Error != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error inserting DB",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusCreated,
		"data":    user,
		"message": "User created successfully",
	})
}

func Update(c *fiber.Ctx, db *gorm.DB) error {

	idUser, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return err
	}

	user := models.User{Id: idUser, UpdatedAt: time.Now()}
	userDB := models.User{Id: idUser}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in parsing JSON",
		})
	}

	if errors := user.ValidateStruct(); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"message": "Error in entries fields",
			"data":    errors,
		})
	}

	// busca al usuario
	result := db.First(&userDB)

	if result.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"data":    nil,
			"message": "User not found",
		})
	}

	if result.Error != nil {
		return result.Error
	}

	// actualizamos la contraseña si el campo no llega vacío
	if len(user.Password) == 0 {
		user.Password = userDB.Password
	}

	// actualizamos los datos
	result = db.Model(&models.User{Id: idUser}).
		Omit("is_active", "created_at"). // omite campos para actualizar
		Updates(&user)

	if result.Error != nil {
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    user,
		"message": "User update successfully",
	})
}

func Delete(c *fiber.Ctx, db *gorm.DB) error {

	idUser, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in get id user",
		})
	}

	// borra el registro de usuario
	result := db.Delete(&models.User{Id: idUser})

	if result.Error != nil {
		return result.Error
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"data":    nil,
		"message": "User delete successfully",
	})
}

func Login(c *fiber.Ctx, db *gorm.DB) error {

	var (
		user   models.User
		userDB models.User
	)

	// obtenemos los datos
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusInternalServerError,
			"data":    nil,
			"message": "Error in parsing JSON",
		})
	}

	// validamos los datos
	if errors := user.ValidateStructLogin(); errors != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"status":  http.StatusBadRequest,
			"message": "Error in entries fields",
			"data":    errors,
		})
	}

	// verificamos si el usuario existe
	result := db.Where(&models.User{Email: user.Email, IsActive: true}).Find(&userDB)

	if result.RowsAffected == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": "Invalid credentials",
			"status":  http.StatusBadRequest,
			"data":    nil,
		})
	}

	if result.Error != nil {
		return result.Error
	}

	// comparamos la contraseñas
	if err := bcrypt.CompareHashAndPassword(
		[]byte(userDB.Password),
		[]byte(user.Password),
	); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"ok":      false,
			"message": "Invalid credentials",
			"status":  http.StatusBadRequest,
			"data":    nil,
		})
	}

	// generamos y firmamos el token JWT
	token, err := middlewares.GenerateJWT(userDB.Name + " " + userDB.Surname)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Problemas internos del servidor",
			"status":  http.StatusInternalServerError,
			"ok":      false,
			"data":    nil,
		})
	}

	// limpiamos la contraseña
	userDB.Password = ""

	// retornamos los datos al usuario
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"ok":      true,
		"status":  http.StatusOK,
		"message": nil,
		"data": fiber.Map{
			"user":  userDB,
			"token": token, // token jwt
		},
	})
}
