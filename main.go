package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/moazfarrukh/go-realtime-chat/chat"
)

func main() {
	Hub := chat.NewHub()
	serverAddress := "localhost:1234"
	http.Handle("/", http.FileServer(http.Dir("frontend")))
	http.HandleFunc("/ws", Hub.ServeWS)
	fmt.Printf("server live at http://%s\n", serverAddress)
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		log.Println(err)
	}
	
	_, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

}
