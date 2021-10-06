package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/gofiber/fiber/v2"
)

var casbinEnforcer *casbin.Enforcer

func main() {
	// setup casbin auth rules
	var err error
	casbinEnforcer, err = casbin.NewEnforcer("./auth_model.conf", "./policy.csv")
	casbinEnforcer.EnableLog(true)
	casbinEnforcer.AddFunction("condition_match", ConditionMatchFunc)
	if err != nil {
		log.Fatalf("Failed to initialize casbin: %s", err.Error())
	}

	app := setup()

	if err := app.Listen(":5000"); err != nil {
		log.Fatalf("Failed to start: %s", err.Error())
	}
}

func setup() *fiber.App {
	app := fiber.New()
	app.Get("/zone/:zoneId", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		zoneId, err := strconv.Atoi(c.Params("zoneId"))
		if err != nil {
			return c.SendStatus(400)
		}
		return authorizeZoneRoute(c, zoneId, "view")
	})

	app.Delete("/zone/:zoneId", func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		zoneId, err := strconv.Atoi(c.Params("zoneId"))
		if err != nil {
			return c.SendStatus(400)
		}
		return authorizeZoneRoute(c, zoneId, "delete")
	})
	return app
}

func authorizeZoneRoute(c *fiber.Ctx, zoneId int, action string) error {
	apiKey := c.Get("x-api-key", "")
	u, err := GetCurrentUser(apiKey)
	if err != nil {
		return c.Status(401).SendString(fmt.Sprintf("<h1>Whoops!<h1><p>%s</p>", err.Error()))
	}

	z, err := GetZoneById(zoneId)
	if err != nil {
		return c.Status(404).SendString("<h1>Whoops!</h1><p>That zone was not found</p>")
	}
	ok, err := casbinEnforcer.Enforce(u.Name, string(c.Request().URI().Path()), action, z)
	if err != nil {
		fmt.Printf("Err: %s\n", err)
		return c.Status(404).SendString("<h1>Whoops!</h1><p>That zone was not found</p>")
	}
	if !ok {
		return c.Status(403).SendString("<h1>Not Authorized</h1><p>User is unauthorized to view that zone</p>")
	}

	return c.Status(200).SendString(fmt.Sprintf("<h1>A zone<h1><p>Welcome %s to zone %s</p>", u.Name, z.Name))
}