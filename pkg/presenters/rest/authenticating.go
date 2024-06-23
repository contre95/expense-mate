package rest

import (
	"expenses-app/pkg/app/authenticating"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func login(a authenticating.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// Throws Unauthorized error
		req := authenticating.LoginReq{
			Username: c.FormValue("user"),
			Password: c.FormValue("pass"),
		}
		resp, err := a.UserAuthenticator.Authenticate(req)
		if err != nil {
			return c.JSON(&fiber.Map{
				"err": c.SendStatus(fiber.StatusUnauthorized),
				"msg": "Unexistent username or wrong password.",
			})
		}
		// Create the Claims
		claims := jwt.MapClaims{
			"name": resp.UserID,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}
		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_SEED")))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.JSON(fiber.Map{"token": t})
	}
}
