package main

import (
	"log"

	"github.com/akshaysangma/go-chat/server/db"
	"github.com/akshaysangma/go-chat/server/internal/user"
	"github.com/akshaysangma/go-chat/server/internal/ws"
	"github.com/akshaysangma/go-chat/server/router"
)

func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database initialize failed: %s", err)
	}

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:8080")
}
