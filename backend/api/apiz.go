package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	databases "backend/database"
)

func main() { //api pe checknome
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Error loading .env file")
	}
	app := fiber.New()
	app.Use(cors.New())

	app.Get("/api/host_game", func(c *fiber.Ctx) error {
		h, err1 := strconv.Atoi(c.Query("h"))
		w, err2 := strconv.Atoi(c.Query("w"))
		n, err3 := strconv.Atoi(c.Query("userId"))
		if err1 != nil && err2 != nil && err3 != nil {
			// ... handle error
			panic("error W-H-n atoi")
		}

		rmId, err := databases.CreateGame(h, w)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating game")
		}
		errz := databases.CreateLobby(int32(n), rmId)
		if errz != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error creating lobby")
		}
		struid := strconv.Itoa(int(rmId))
		return c.Status(fiber.StatusOK).SendString(struid)
	})
	app.Get("/api/temp_user", func(c *fiber.Ctx) error {
		nme := c.Query("nome")
		if len(nme) >= 10 {
			return c.Status(fiber.StatusInternalServerError).SendString("Name too long")
		}
		exxist, uid := databases.GetUser(nme)
		if exxist {
			return c.Status(fiber.StatusInternalServerError).SendString("Exist")
		} else {
			struid := strconv.Itoa(int(uid))
			return c.Status(fiber.StatusOK).SendString(struid)
		}

	})
	///join_game?userId=${userId}&roomID=${rmId}
	///TODO

	//start server
	log.Fatal(app.Listen(":3000"))
}
