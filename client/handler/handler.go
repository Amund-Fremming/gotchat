package handler

import (
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

func ConnectWebSocket() error {
	url := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/chat"}
	conn, _, err := websocket.DefaultDialer.Dial(url.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("[Client] Ping")
	err = conn.WriteMessage(websocket.TextMessage, []byte("Ping"))
	if err != nil {
		return err
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	log.Println("[SERVER]", string(msg))

	return nil
}

func Run() {
}

func tunnel() {
	//
}

func dispatcher() {
	//
}

/*
main
	spawner og holder select

tunnel
	tar imot meldinger fra server og logger
	snakker ikke til dispatcher
	mottar commands fra dispatcher, og handler videre

dispatcher
	lytter på input
	ved input sender melding til tunnel
*/
