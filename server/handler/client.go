package handler

import (
	"encoding/json"
	"errors"
	ws "github.com/gorilla/websocket"
	"log"
	"net"
	"time"
)

const (
	writeWait = 60 * time.Second

	pongWait = 360 * time.Second

	pingPeriod = (pongWait * 9) / 10

	maxMessageSize = 2048
)

type Client struct {
	username string
	conn     *ws.Conn
	id       int64
	send     chan *Message
	room     *Room
}

func NewClient(conn *ws.Conn, id int64, username string, room *Room) *Client {
	c := &Client{
		username: username,
		id:       id,
		conn:     conn,
		send:     make(chan *Message, 256),
		room:     room,
	}

	go c.Read()
	go c.Write()
	return c
}

func (c *Client) Read() {
	defer func() {
		c.room.unregistry <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	defaultHandle := c.conn.CloseHandler()
	c.conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("receive close message code:%d text:%s", code, text)
		c.room.unregistry <- c

		return defaultHandle(code, text)
	})

	for {
		_, data, err := c.conn.ReadMessage()

		var m = &Message{}
		if err != nil {
			if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseNormalClosure) {
				log.Printf("error:%v", err)
			}

			break
		}
		err = json.Unmarshal(data, &m)
		m.Username = c.username
		m.PlayerID = c.id
		handleMessage(m)
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(ws.CloseMessage, []byte{})
				if err != nil {
					//判断之前是不是已经发过CloseMessage了
					if errors.Is(err, ws.ErrCloseSent) {
						return
					}

					//判断连接是不是已经被关闭了
					if netErr, ok := err.(*net.OpError); ok && netErr.Err.Error() == "use of closed network connection" {
						return
					}

					log.Println(err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			err = c.conn.WriteMessage(ws.BinaryMessage, data)
			if err != nil {
				log.Println(err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(ws.PingMessage, nil)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
