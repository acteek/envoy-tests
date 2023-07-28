package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {

	port := os.Getenv("PORT")
	meet := os.Getenv("MEET")
	redirectHost := os.Getenv("REDIRECT")

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		log.Println("http ping ...")
		return
	})

	http.HandleFunc("/meetings", func(w http.ResponseWriter, r *http.Request) {

		log.Println("-------Request Headers------")
		for k, v := range r.Header {
			log.Printf("%s: %s", k, v[0])
		}
		log.Println("------------END--------------")
		meetingId := r.Header.Get("X-Meeting-Id")
		if meetingId == meet {
			r.Header.Set("Content-Type", "")
			http.Redirect(w, r, fmt.Sprintf("http://%s/redirect", redirectHost), 302)
			return
		}

		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			log.Printf("%s PONG: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	log.Println("Ws server starting...")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		panic(err)
	}
}
