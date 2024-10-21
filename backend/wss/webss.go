package main

import (
	//"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	//databases "backend/database"
)

type IncomingPositions struct {
	UID int     `json:"uid"`
	X   float64 `json:"x"`
	Y   float64 `json:"y"`
}
type EnemyPosition struct {
	RID int
	X   float64
	Y   float64
}
type ConnPool struct {
	conn      []*websocket.Conn
	enemyPos  *chan EnemyPosition
	playerPos *chan []byte
}

func handleMessages(c *websocket.Conn, chPool *chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("readerr:", err)
			//incomingPos <- mt
			break
		}
		fmt.Println("read:", string(msg), mt)

		//var tstz IncomingPositions
		//_ = json.Unmarshal(msg, &tstz)
		*chPool <- msg
	}
}
func handleSendMessages(poolrm *ConnPool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case pos := <-*poolrm.playerPos:
			fmt.Println("entro/?", pos)
			for _, v := range poolrm.conn {

				err1 := v.WriteMessage(1, pos)
				if err1 != nil {
					log.Println("write:", err1)
					break
				}
			}
		case posE := <-*poolrm.enemyPos:
			fmt.Println("TODO", posE)

		}
	}
}

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Error loading .env file")
	}
	app := fiber.New()

	poolz := make(map[int]*ConnPool)
	var wg sync.WaitGroup

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:roomId", websocket.New(func(c *websocket.Conn) {
		//TODO map room:poolqueue  --map usid pool?

		/* defer func() {
			///ps. mesa funz a parte co api o something bho ??
			//TODO manda userid, metti listaconnessione a struct co id:c
			//cancella utente e se len0 delete(poolz, rmId)
			c.Close()
		}() */

		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))            // true
		log.Println("roomID: ", c.Params("roomId")) // 123
		log.Println(c.Cookies("session"))           // ""

		rmId, errA := strconv.Atoi(c.Params("roomId"))
		if errA != nil {
			return ///bho
		}

		_, ok := poolz[rmId] //if _, ok := poolz[rmId]; !ok {...
		if !ok {
			incomPosz := make(chan []byte)
			enemPosz := make(chan EnemyPosition)
			poolz[rmId] = &ConnPool{
				conn:      make([]*websocket.Conn, 0),
				enemyPos:  &enemPosz,
				playerPos: &incomPosz,
			}
		}
		poolz[rmId].conn = append(poolz[rmId].conn, c)

		wg.Add(2)
		go handleMessages(c, poolz[rmId].playerPos, &wg)
		go handleSendMessages(poolz[rmId], &wg)
		wg.Wait()
		//var tstz IncomingPositions
		//_ = json.Unmarshal(msg, &tstz)
		//*poolz[rmId].playerPos <- msg
		//log.Println("recv: ", msg)

		///write

	}))

	log.Fatal(app.Listen(":3001"))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}
