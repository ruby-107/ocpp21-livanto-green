package transport

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	conn    *websocket.Conn
	writeMu sync.Mutex
}

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"ocpp2.1"},
	CheckOrigin: func(r *http.Request) bool {
		// Local learning demo only.
		// Production mein origin ko allowlist se validate karein.
		return true
	},
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &Conn{conn: conn}, nil
}

func Dial(url string) (*Conn, error) {
	dialer := websocket.Dialer{
		Subprotocols: []string{"ocpp2.1"},
	}

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return &Conn{conn: conn}, nil
}

func (c *Conn) Subprotocol() string {
	return c.conn.Subprotocol()
}

func (c *Conn) ReadMessage() ([]byte, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *Conn) WriteMessage(data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
