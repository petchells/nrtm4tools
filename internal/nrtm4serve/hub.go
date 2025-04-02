package nrtm4serve

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/petchells/nrtm4tools/internal/nrtm4/util"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 32 * 1024
)

var (
	newline = []byte{'\n'}
)

type message struct {
	ID      string
	Content map[string]any
}

type Client struct {
	ID   string
	conn *websocket.Conn
	send chan message
	hub  *Hub
}

func (c *Client) readPump(wg *sync.WaitGroup) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		wg.Done()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, mbytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Info("Websocket closed unexpectedly", "error", err)
			}
			logger.Warn("Cannot read from websocket", "error", err)
			break
		}
		var msg message
		if err = json.Unmarshal(mbytes, &msg); err != nil {
			logger.Warn("Unmarshal message failed", "mbytes", mbytes, "error", err)
		}
		c.hub.send <- msg
	}
}

func (c *Client) writePump(wg *sync.WaitGroup) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		wg.Done()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warn("SetWriteDeadline failed", "error", err)
			}
			if !ok {
				logger.Info("Hub closed the channel")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				logger.Warn("NextWriter returned an error", "error", err)
				return
			}
			var b []byte
			if b, err = json.Marshal(msg); err != nil {
				logger.Warn("Marshal message failed", "message", msg)
				return
			}
			w.Write(b)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for range n {
				w.Write(newline)
				if b, err = json.Marshal(<-c.send); err != nil {
					logger.Warn("Marshal <-c.send failed", "message", msg)
					return
				}
				w.Write(b)
			}

			if err := w.Close(); err != nil {
				logger.Warn("Websocket closed with error", "error", err)
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				logger.Warn("SetWriteDeadline returned error", "error", err)
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Warn("Failed to write ping to websocket", "error", err)
				return
			}
		}
	}
}

type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Inbound messages from the clients.
	send chan message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	//connections map[string]*websocket.Conn
}

func newHub() *Hub {
	return &Hub{
		send:       make(chan message, 20),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}
}

var broadcastChannels = util.NewSet("logs")

func (h *Hub) run() {
	sendMessage := func(client *Client, msg message) {
		select {
		case client.send <- msg:
		default:
			close(client.send)
			delete(h.clients, client.ID)
		}
	}
	for {
		select {
		case client := <-h.register:
			h.clients[client.ID] = client
			logger.Debug("Registered client", "client.ID", client.ID)
		case client := <-h.unregister:
			logger.Debug("Unregistering client", "client.ID", client.ID)
			if _, ok := h.clients[client.ID]; ok {
				delete(h.clients, client.ID)
				close(client.send)
			}
		case msg := <-h.send:
			if broadcastChannels.Contains(msg.ID) {
				for _, client := range h.clients {
					sendMessage(client, msg)
				}
			} else if client, ok := h.clients[msg.ID]; ok {
				sendMessage(client, msg)
			}
		}
	}
}
