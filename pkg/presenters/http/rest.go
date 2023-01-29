package http

import (
	"expenses-app/pkg/app/authenticating"
	"expenses-app/pkg/app/health"
	"expenses-app/pkg/app/importing"
	"expenses-app/pkg/app/managing"
	"expenses-app/pkg/app/querying"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

// MapRoutes is where http REST routes are mapped to functions
func MapRoutes(fi *fiber.App, he *health.Service, m *managing.Service, i *importing.Service, a *authenticating.Service, q *querying.Service) {
	// Unrestricted
	fi.Get("/ping", ping(*he))
	fi.Post("/login", login(*a))
	fi.Use(jwtware.New(jwtware.Config{SigningKey: []byte(os.Getenv("JWT_SECRET_SEED"))}))
	// Restricted
	fi.Get("/restricted", restricted)
	fi.Post("/users", createUsers(m.UserCreator))
	fi.Get("/expenses/categories", getCategories(q.CategoryGetter))
	fi.Post("/importers/:id", importExpenses(i.ImportExpenses))
}

func ping(h health.Service) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{
			"ping": h.Ping(),
		})
	}
}

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

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}
