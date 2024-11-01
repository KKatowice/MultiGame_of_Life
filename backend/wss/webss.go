package main

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"slices"
	"strconv"
	"sync"
	"time"

	databases "backend/database"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

/*
	 type IncomingPositions struct {
		UID int     `json:"uid"`
		X   float64 `json:"x"`
		Y   float64 `json:"y"`
	}
*/
type EnemyPosition struct {
	RID    int
	X      float64
	Y      float64
	IsEnem bool
}
type EnemyJson struct {
	RID       int
	Positions []EnemyPosition
}

type ConnWrapper struct {
	conn    *websocket.Conn
	writeCh chan []byte // Dedicated channel for writing to this connection
}
type ConnPool struct {
	conn      []*ConnWrapper
	enemyPos  *chan EnemyJson //forse useless?
	playerPos *chan []byte    // all positions pl+enem
}

func calculateGrid(allPositions *[][]EnemyPosition, ww *int32, hh *int32, stDens *int) {
	w := float64(*ww)
	h := float64(*hh)
	var enemySize int32 = 15
	flotES := float64(enemySize)
	var spacing int32 = 10
	var cols int32 = *ww / (enemySize + spacing)
	var rows int32 = *hh / (enemySize + spacing)
	paddingX := (w - (float64(cols * enemySize))) / (float64(cols) + 1) // Spaziatura orizzontale
	paddingY := (h - (float64(rows * enemySize))) / (float64(rows) + 1) // Spaziatura verticale

	//var allPositions [][][2]float64 //[[row]]
	for r := 1; r <= int(rows); r++ {
		var tempRow []EnemyPosition
		for c := 1; c <= int(cols); c++ {
			/* var tempCoord [3]float64
			tempCoord[0] = float64(c) * (flotES + paddingX)
			tempCoord[1] = float64(r) * (flotES + paddingY) */

			isEnem := rand.IntN(100)
			enem := false
			if isEnem <= *stDens {
				enem = true
			}
			tempCoord := EnemyPosition{
				X:      float64(c) * (flotES + paddingX),
				Y:      float64(r) * (flotES + paddingY),
				IsEnem: enem,
			}
			//tempCoord[2] = float64(enem) ///0 dead - 1 alive
			tempRow = append(tempRow, tempCoord)
		}
		*allPositions = append(*allPositions, tempRow)
	}

}

func handleGameOfLife(rmId int, chPool *chan EnemyJson, wg *sync.WaitGroup, speed int, startDensity int) {
	//TODO later as param
	speed = 500       // millisec between cycles = 0.5s
	startDensity = 15 // %
	//positions := make([][]int, 0) //size = calc con w e h
	var allPositions [][]EnemyPosition //[[row]] -- keep track all positions for logic
	w, h, err := databases.Get_wh(rmId)
	calculateGrid(&allPositions, &w, &h, &startDensity)
	if err {
		fmt.Println("error get wh: ", err)
	}

	//get w,h from db
	defer wg.Done()

	for {
		var newPos []EnemyPosition /// only enemy positions
		///logica

		for row := 0; row < len(allPositions); row++ {
			checkAbove := false ////row=0
			checkBelow := true  ////row=0
			checkRight := true  //elem==0
			checkLeft := false  ////elem==0
			if row > 0 {
				checkAbove = true
				checkBelow = true
			}
			if row == len(allPositions)-1 {
				checkBelow = false
				checkAbove = true
			}
			singRow := allPositions[row]
			for elem := 0; elem < len(singRow); elem++ {
				vicini := 0

				if elem == 0 {
					checkRight = true //elem==0
					checkLeft = false ////elem==0
				}
				if elem > 0 {
					checkLeft = true
					checkRight = true
				}
				if elem == len(singRow)-1 {
					checkRight = false
					checkLeft = true
				}

				if checkAbove {
					sopra := allPositions[row-1][elem].IsEnem
					if sopra {
						vicini++
					}
				}
				if checkBelow {
					sotto := allPositions[row+1][elem].IsEnem
					if sotto {
						vicini++
					}
				}
				if checkLeft {
					sx := allPositions[row][elem-1].IsEnem
					if sx {
						vicini++
					}
					if checkAbove {
						tsx := allPositions[row-1][elem-1].IsEnem
						if tsx {
							vicini++
						}
					}
					if checkBelow {
						bsx := allPositions[row+1][elem-1].IsEnem
						if bsx {
							vicini++
						}
					}
				}
				if checkRight {
					dx := allPositions[row][elem+1].IsEnem
					if dx {
						vicini++
					}
					if checkAbove {
						tdx := allPositions[row-1][elem+1].IsEnem
						if tdx {
							vicini++
						}
					}
					if checkBelow {
						bdx := allPositions[row+1][elem+1].IsEnem
						if bdx {
							vicini++
						}
					}
				}
				if allPositions[row][elem].IsEnem {
					if vicini < 2 {
						allPositions[row][elem].IsEnem = false
					} else if vicini > 3 {
						allPositions[row][elem].IsEnem = false
					}
					newPos = append(newPos, allPositions[row][elem])
				} else {
					if vicini == 3 {
						allPositions[row][elem].IsEnem = true
						newPos = append(newPos, allPositions[row][elem])
					}
				}

			}
		}

		///

		/* en, err := json.Marshal(EnemyJson{
			RID:       rmId,
			Positions: newPos,
		})
		if err != nil {
			fmt.Println("json error")
			panic(err)
		} */
		a := EnemyJson{
			RID:       rmId,
			Positions: newPos,
		}
		*chPool <- a
		time.Sleep(time.Duration(speed) * time.Millisecond)
	}
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
	fmt.Println("cazzo dici grandiocane?")
	defer wg.Done()
	for {
		select {
		case pos := <-*poolrm.playerPos:
			fmt.Println("entro/?", pos)
			for _, v := range poolrm.conn {
				fmt.Println("che probnlemi??")

				v.writeCh <- pos

			}
		case posE := <-*poolrm.enemyPos:
			for _, v := range poolrm.conn {
				fmt.Println("che cazzo di probnlemi??")
				byteList, err := json.Marshal(posE)
				if err != nil {
					panic(err)
				}
				v.writeCh <- byteList

			}

		}
	}
}
func handleConnectionWrites(cw *ConnWrapper, wg *sync.WaitGroup, lcw *map[int]*ConnPool, rmid int) {
	fmt.Println("cazzo dici gr?")
	frapls := (*lcw)[rmid]
	defer wg.Done()
	for msg := range cw.writeCh {
		err := cw.conn.WriteMessage(1, msg)
		if err != nil {
			log.Println("Write error:", err)
			RemoveIndex(&frapls.conn, cw)
			cw.conn.Close() //// PORCODIOOOOOOOOOOOOOOOOO PERCHEEEEEEEEEEEEEE
			if len(frapls.conn) == 0 {
				delete(*lcw, rmid)
			}

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

		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))            // true
		log.Println("roomID: ", c.Params("roomId")) // 123
		log.Println(c.Cookies("session"))           // ""

		cw := &ConnWrapper{
			conn:    c,
			writeCh: make(chan []byte, 10), // Buffered channel for this connection
		}
		rmId, errA := strconv.Atoi(c.Params("roomId"))
		if errA != nil {
			return ///bho
		}

		/* defer func() {
			///ps. mesa funz a parte co api o something bho ??
			//TODO manda userid, metti listaconnessione a struct co id:c
			//cancella utente e se len0 delete(poolz, rmId)
			RemoveIndex(&poolz[rmId].conn, cw)
			if len(poolz[rmId].conn) == 0 {
				delete(poolz, rmId)
			}
		}() */
		/* socketio.On(socketio.EventDisconnect, func(ep *socketio.EventPayload) {
			fmt.Println("FRAAAAAAAAAA SE STA A DISCONNETTE DIOIOMERDA")

			// Remove the user from the local clients
			RemoveIndex(&poolz[rmId].conn, cw)
			if len(poolz[rmId].conn) == 0 {
				delete(poolz, rmId)
			}
			fmt.Printf("Disconnection event - User")
		}) */

		_, ok := poolz[rmId] //if _, ok := poolz[rmId]; !ok {...
		if !ok {
			incomPosz := make(chan []byte)
			enemPosz := make(chan EnemyJson)
			poolz[rmId] = &ConnPool{
				conn:      make([]*ConnWrapper, 0),
				enemyPos:  &enemPosz,
				playerPos: &incomPosz,
			}
			wg.Add(1)
			go handleGameOfLife(rmId, poolz[rmId].enemyPos, &wg, 0, 0)
		}
		poolz[rmId].conn = append(poolz[rmId].conn, cw)

		wg.Add(3)
		go handleMessages(c, poolz[rmId].playerPos, &wg)
		go handleSendMessages(poolz[rmId], &wg)
		go handleConnectionWrites(cw, &wg, &poolz, rmId)

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

func RemoveIndex(s *[]*ConnWrapper, user *ConnWrapper) {
	var index int
	for i, v := range *s {
		if v == user {
			index = i
			v.conn.Close()
			break
		}
	}
	*s = slices.Delete(*s, index, index)

}
