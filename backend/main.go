package main

import (
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
		panic("Error loading .env file")
	} /*
		go Apiz()
		go Webss() */
}
