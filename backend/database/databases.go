package databases

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	jsoniter "github.com/json-iterator/go"

	_ "github.com/go-sql-driver/mysql"

	"backend/database/sqlc"

	"math/rand"
)

var (
	dbse    *sql.DB
	queries *sqlc.Queries
	once    sync.Once
	json    jsoniter.API
)

func initDB() {
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	var err error
	envurl := os.Getenv("DB_URL")
	fmt.Println("asd????", envurl)
	dbse, err = sql.Open("mysql", envurl)

	dbse.SetMaxOpenConns(1000)
	if err != nil {
		panic(err)
	}
	queries = sqlc.New(dbse)

}

func GetUser(name string) (bool, int32) {
	once.Do(initDB)
	ctx := context.Background()
	nme := sql.NullString{
		String: name,
		Valid:  true,
	}
	fmt.Println("nome?", name)

	urs, err := queries.ListUsers(ctx, nme)
	fmt.Println("porcodio", urs, err)
	/* if err != nil {
		fmt.Println("ao listuser crash")
		panic(err)
	} */
	if urs != 0 {
		return true, urs
	}
	fmt.Println("exxist? ", urs, name)
	uid, err1 := CreateUser(name)
	if err1 != nil {
		fmt.Println("crashone")
		panic(err1)
	}
	return false, int32(uid)
	/* jsonResponse, err := json.Marshal(&urs)
	if err != nil {
		fmt.Println("jsonResponse error convert")
	}
	fmt.Println("dai?", jsonResponse) */

}

func RemoveUser(uid int) error {
	once.Do(initDB)
	ctx := context.Background()
	err := queries.DeleteUser(ctx, int32(uid))
	return err
}

/*
	 func randString(length int) string {
		var random = rand.New(rand.NewSource(1))
		b := make([]byte, length)
		for i := 0; i < length; i++ {
			b[i] = byte(random.Int63() & 0xff)
		}
		return string(b)
	}
*/
func createRoomId() int32 {
	///check che non esiste
	return rand.Int31n(918273645)
}

func CreateGame(heig int, widt int) (int32, error) {
	once.Do(initDB)
	ctx := context.Background()
	rid := createRoomId()
	parm := sqlc.CreateGameParams{
		Roomid: rid,
		Hei: sql.NullInt32{
			Int32: int32(heig),
			Valid: true,
		},
		Wid: sql.NullInt32{
			Int32: int32(widt),
			Valid: true,
		},
	}
	err := queries.CreateGame(ctx, parm)
	if err != nil {
		fmt.Println("Error", err)
		return 0, err
	}
	/* jsonResponse, err := json.Marshal(&urs)
	if err != nil {
		fmt.Println("jsonResponse error convert")
	} */

	return rid, nil
}

func CreateLobby(uid int32, roId int32) error {
	once.Do(initDB)
	ctx := context.Background()
	parm := sqlc.CreateLobbyParams{
		Roomid: sql.NullInt32{
			Int32: roId,
			Valid: true,
		},
		Userid: sql.NullInt32{
			Int32: uid,
			Valid: true,
		},
	}
	err := queries.CreateLobby(ctx, parm)
	if err != nil {
		fmt.Println("Error", err)
		return err
	}

	return nil
}
func JoinLobby(uid int32, roId int32) (int32, int32, error) {
	once.Do(initDB)
	ctx := context.Background()
	//TODO check if space & if non sta gia in
	parm := sqlc.CreateLobbyParams{
		Roomid: sql.NullInt32{
			Int32: roId,
			Valid: true,
		},
		Userid: sql.NullInt32{
			Int32: uid,
			Valid: true,
		},
	}
	fmt.Println("param i ndb?", roId, uid)

	err := queries.CreateLobby(ctx, parm)
	if err != nil {
		fmt.Println("Error", err)
		return 0, 0, err
	}
	wh, err2 := queries.Get_wh(ctx, roId)
	if err2 != nil {
		fmt.Println("Error", err)
		return 0, 0, err
	}
	fmt.Println("wh????", wh)

	return wh.Wid.Int32, wh.Hei.Int32, nil
}

func CreateUser(name string) (uid int32, err error) {
	once.Do(initDB)
	ctx := context.Background()
	nme := sql.NullString{
		String: name,
		Valid:  true,
	}
	parm := sqlc.CreateUserParams{
		Name: nme,
	}
	err = queries.CreateUser(ctx, parm)
	if err != nil {
		fmt.Println("Error", err)
		return 0, err
	}
	urs, err := queries.ListUsers(ctx, nme)
	if err != nil {
		fmt.Println("Error", err)
		return 0, err
	}
	/* stm, err := dbse.Prepare("INSERT INTO users (userId, name) VALUES (?,?);")
	res, err := stm.Exec(name)
	uidi, err := res.LastInsertId() */

	return urs, nil
}
