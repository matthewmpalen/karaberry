package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type (
	Hub struct {
		clients    map[*HubClient]struct{}
		register   chan *HubClient
		unregister chan *HubClient
		broadcast  chan string
	}

	HubClient struct {
		hub  *Hub
		conn *websocket.Conn
		send chan []byte
	}
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var (
	hub = &Hub{
		clients:    make(map[*HubClient]struct{}),
		register:   make(chan *HubClient),
		unregister: make(chan *HubClient),
		broadcast:  make(chan string, 100),
	}
	upgrader = websocket.Upgrader{}
)

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = struct{}{}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- []byte(message):
					// Do nothing else
				default:
					// If the client's send buffer is full, assumes
					// that the client is dead or stuck
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (c *HubClient) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}

func (c *HubClient) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if closeErr := w.Close(); closeErr != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	context := map[string]interface{}{"songs": songList}
	RenderTemplate(w, "home", context)
}

func SongsHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]interface{}{
		"count": len(songList),
		"songs": songList,
	}
	JSON(w, resp, http.StatusOK)
}

func QueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResp := map[string]string{"error": "Requires POST"}
		JSON(w, errorResp, http.StatusMethodNotAllowed)
		return
	}

	id := r.PostFormValue("songId")
	songID, parseErr := strconv.Atoi(id)
	if parseErr != nil || songID < 0 || songID > len(songList)-1 {
		msg := fmt.Sprintf("Invalid song: %s", id)
		errorResp := map[string]string{"error": msg}
		JSON(w, errorResp, http.StatusBadRequest)
		return
	}

	song := songList[songID]
	karaoke.Queue <- song
	hub.broadcast <- fmt.Sprintf("[ADDED] %s", song.String())

	resp := map[string]interface{}{
		"count": len(karaoke.Queue),
		"added": song.String(),
	}
	JSON(w, resp, http.StatusOK)
}

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	songs := karaoke.History().Songs()
	resp := map[string]interface{}{
		"count": len(songs),
		"songs": songs,
	}
	JSON(w, resp, http.StatusOK)
}

func IdleHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, "idle", nil)
}

func SkipHandler(w http.ResponseWriter, r *http.Request) {
	program := Config.MediaPlayer
	// Override for omxplayer
	if program == "omxplayer" {
		program += ".bin"
	}
	exec.Command("killall", "-q", program).Run()
	msg := "SKIPPING CURRENT SONG"
	hub.broadcast <- msg
	log.Printf("%s\n", msg)
	resp := map[string]string{"message": msg}
	JSON(w, resp, http.StatusOK)
}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade: %v", err)
		return
	}

	client := &HubClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	client.hub.register <- client
	go client.writePump()
}
