package main

import (
	"go-chess/server/config"
	"go-chess/server/dao"
	"go-chess/server/handler"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dao.InitDB()
	config.InitConfig()

	http.HandleFunc("/room", handler.HandleCreateRoom)
	http.HandleFunc("/join", handler.HandleJoinRoom)
	http.HandleFunc("/registry", handler.HandleRegistry)
	http.HandleFunc("/login", handler.HandleLogin)

	log.Println("Server started. Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
