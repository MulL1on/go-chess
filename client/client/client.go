package client

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	ws "github.com/gorilla/websocket"
	"go-chess/client/chess"
	"log"
	"net"
)

const (
	TurnNone = iota
	TurnWhite
	TurnBlack
)

type Client struct {
	roomID string
	conn   *ws.Conn
	send   chan *Message
	board  [8][8]int
	Exit   chan struct{}
	Turn   int

	screen   tcell.Screen
	CmdChan  chan string
	ChatChan chan string
	Input    string

	chatHistory []string
}

func NewClient(roomID string, conn *ws.Conn, screen tcell.Screen) *Client {

	c := &Client{
		roomID: roomID,
		board:  chess.CreatInitialBoard(),
		screen: screen,
		conn:   conn,
		Exit:   make(chan struct{}),
		send:   make(chan *Message, 256),

		CmdChan:  make(chan string, 256),
		ChatChan: make(chan string, 256),

		chatHistory: make([]string, 0),
	}

	//设置CustomCloseHandler
	defaultHandler := conn.CloseHandler()
	customHandler := func(code int, text string) error {
		log.Printf("receive close message code:%d text:%s", code, text)
		c.Exit <- struct{}{}
		conn.Close()
		return defaultHandler(code, text)
	}
	c.conn.SetCloseHandler(customHandler)

	go c.Run()
	go c.ReadPump()
	go c.WritePump()
	return c
}

func (c *Client) Run() {
	for {
		select {
		case input := <-c.CmdChan:
			log.Println("get cmd message: ", input)
			c.handleCmdChan(input)
		case input := <-c.ChatChan:
			c.handleChatChan(input)
		}
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Exit <- struct{}{}
		c.conn.Close()
	}()
	for {
		msgType, data, e := c.conn.ReadMessage()
		if e != nil {
			//判断是不是连接不正常关闭
			if ws.IsUnexpectedCloseError(e, ws.CloseNormalClosure, ws.CloseGoingAway) {
				log.Println(e)
			}

			//判断是不是与服务端断开连接
			if e, ok := e.(net.Error); ok {
				log.Printf("network error : %v", e)
			}
		}
		switch msgType {
		case ws.BinaryMessage:
			//log.Println("get message from server: ", data)
			c.route(data)
		}
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Exit <- struct{}{}
		c.conn.Close()
	}()
	for {
		select {
		case <-c.Exit:
			return
		case message := <-c.send:

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
			}
			err = c.conn.WriteMessage(ws.BinaryMessage, data)
			if err != nil {
				log.Println(err)
			}
			//	err = json.Unmarshal(data, &message)
			//	if err != nil {
			//	log.Println(err)
			//}
			//log.Println("send message to server: ", message)
		}
	}
}

func (c *Client) route(data []byte) {
	var m = &Message{}
	err := json.Unmarshal(data, m)
	if err != nil {
		log.Println(err)
		return

	}

	switch m.Type {
	case MessageTypeChat:
		c.handleChat(m)
	case MessageTypeGameStart:
		c.handleGameStart(m)
	case MessageTypeErr:
		c.handleErr(m)
	case MessageTypeUpdateBoard:
		c.handleUpdateBoard(m)
	case MessageTypeCheck:
		c.handleCheck(m)
	case MessageTypeGameOver:
		c.handleGameOver(m)
	}
}

func (c *Client) handleCmdChan(cmd string) {
	var m *Message

	//用正则匹配 cmd 是移动还是王车易位还是升变
	if RegexTypeMove.MatchString(cmd) {
		m = &Message{
			Type: MessageTypeMove,
			Move: move{
				RoomID: c.roomID,
				Move:   cmd,
			},
		}
	} else if RegexTypeCastle.MatchString(cmd) {
		m = &Message{
			Type: MessageTypeCastling,
			Castling: castling{
				RoomID: c.roomID,
				Shift:  cmd,
			},
		}
	} else if RegexTypePromotion.MatchString(cmd) {
		m = &Message{
			Type: MessageTypePromotion,
			Promotion: promotion{
				RoomID:    c.roomID,
				Promotion: cmd,
			},
		}
	}

	c.send <- m
}

func (c *Client) handleChatChan(content string) {
	m := Message{
		Type: MessageTypeChat,
		Chat: chat{
			RoomID:  c.roomID,
			Content: content,
		},
	}

	c.send <- &m
}

func (c *Client) handleUpdateBoard(m *Message) {
	log.Println("get update board message: ", m.UpdateBoard.Board)
	c.board = m.UpdateBoard.Board
	c.Turn = m.UpdateBoard.Turn
}

func (c *Client) handleChat(m *Message) {
	//log.Println("get chat message: ", m.Chat)
	//append to chat history
	if len(c.chatHistory) > 10 {
		c.chatHistory = c.chatHistory[1:]
	}
	body := fmt.Sprintf("%s:%s", m.Chat.Username, m.Chat.Content)
	c.chatHistory = append(c.chatHistory, body)
}

func (c *Client) Show() {
	//draw number & letter
	for i := 0; i < 8; i++ {
		c.screen.SetContent(10, i+1, '8'-rune(i), nil, tcell.StyleDefault)
		c.screen.SetContent(i+1, 9, 'a'+rune(i), nil, tcell.StyleDefault)
	}
	//draw piece
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			c.screen.SetContent(j+1, i+1, chess.ParsePiece(c.board[i][j]), nil, tcell.StyleDefault)
		}
	}
	//show chat history
	for i := 0; i < len(c.chatHistory); i++ {
		for j, v := range c.chatHistory[i] {
			c.screen.SetContent(j+12, 0+i, v, nil, tcell.StyleDefault)
		}
	}
}

func (c *Client) handleGameStart(m *Message) {
	if len(c.chatHistory) > 10 {
		c.chatHistory = c.chatHistory[1:]
	}
	log.Println("get game start message: ", m.GameStart)
	c.chatHistory = append(c.chatHistory, "Game Start "+m.GameStart.Content)
	c.Turn = m.GameStart.Turn
}

func (c *Client) handleCheck(m *Message) {
	if len(c.chatHistory) > 10 {
		c.chatHistory = c.chatHistory[1:]
	}
	c.chatHistory = append(c.chatHistory, m.Check.Content)
}

func (c *Client) handleErr(m *Message) {
	log.Println("get err message: ", m.Err.Error)
	if len(c.chatHistory) > 10 {
		c.chatHistory = c.chatHistory[1:]
	}
	c.chatHistory = append(c.chatHistory, m.Err.Error)
}

func (c *Client) handleGameOver(m *Message) {
	if len(c.chatHistory) > 10 {
		c.chatHistory = c.chatHistory[1:]
	}
	c.chatHistory = append(c.chatHistory, m.GameOver.Content)
}
