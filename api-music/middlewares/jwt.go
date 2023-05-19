package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type customClaims struct {
	Name string
	jwt.RegisteredClaims
}

var secret []byte

func GenerateJWT(nameUser string) (string, error) {

	secret = []byte("my_secret_phase")

	claims := customClaims{
		nameUser,
		jwt.RegisteredClaims{
			// testing: tiempo de expiracion de un 1 minuto
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			Issuer:    nameUser,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(c *fiber.Ctx) error {

	rawToken := c.Get("Authorization")

	// validamos si llega el token por los headers
	if len(rawToken) == 0 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"ok":      false,
			"message": "Usuario no autorizado",
			"status":  http.StatusUnauthorized,
		})
	}

	// obtenemos el token y lo parseamos
	tokenString := strings.Split(rawToken, " ")[1]
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// si el token esta mal formado
	if errors.Is(err, jwt.ErrTokenMalformed) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"ok":      false,
			"message": "Usuario no autorizado (token malformado)",
			"status":  http.StatusUnauthorized,
		})
	}

	// si el token expiro o antes de ser valido
	if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {

		// generamos el nuevo token y validamos el error
		tokenString, err := GenerateJWT("token_refresh")

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"ok":      false,
				"message": "Problemas internos con el servidor",
				"status":  http.StatusInternalServerError,
			})
		}

		// se lo pasamos al 401 para que frontend pueda actualizarlo
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"ok":                false,
			"message":           "Usuario no autorizado (token vencido) se actualiza el token",
			"status":            http.StatusUnauthorized,
			"gabalfusers_token": tokenString, // se manda el token actualizado
		})
	}

	// si el token no es valido
	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"ok":      false,
			"message": "Usuario no autorizado (token no valido)",
			"status":  http.StatusUnauthorized,
		})
	}

	// en este punto el token es valido y funciona
	return c.Next()
}
