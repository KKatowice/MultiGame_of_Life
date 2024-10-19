package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	//databases "backend/database"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Error loading .env file")
	}
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	type IncomingPositions struct {
		UID int     `json:"uid"`
		X   float64 `json:"x"`
		Y   float64 `json:"y"`
	}
	app.Use("/ws/:roomId", websocket.New(func(c *websocket.Conn) {
		//TODO map room:poolqueue  --map usid pool?
		defer func() {
			c.Close()
		}()
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))            // true
		log.Println("roomID: ", c.Params("roomId")) // 123
		log.Println(c.Cookies("session"))           // ""
		//incomingPos := make(chan IncomingPositions)
		//enemyPos := make(chan [][]int32)

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			mt, msg, err = c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				//incomingPos <- mt
				break
			}
			var tstz IncomingPositions
			_ = json.Unmarshal(msg, &tstz)
			log.Println("recv: ", tstz)

			err1 := c.WriteMessage(mt, msg)
			if err1 != nil {
				log.Println("write:", err)
				break
			}
		}

	}))

	log.Fatal(app.Listen(":3001"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}
