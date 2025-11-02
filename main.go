package main

import (
	"AntColonyProject/server"
	"log"
	"net/http"
	"os"
)

func main() {
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)

	// Sirve la aplicación de React
	http.Handle("/", http.FileServer(http.Dir("./client/build")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal((err))
	}

}
