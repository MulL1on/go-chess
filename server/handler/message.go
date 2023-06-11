package handler

type Message struct {
	Type        int         `json:"type"`
	PlayerID    int64       `json:"playerId"`
	Username    string      `json:"username"`
	Move        move        `json:"move,omitempty"`
	Chat        chat        `json:"chat,omitempty"`
	UpdateBoard updateBoard `json:"updateBoard,omitempty"`
	GameStart   gameStart   `json:"gameStart,omitempty"`
	Err         err         `json:"err,omitempty"`
	Check       check       `json:"check,omitempty"`
	GameOver    gameOver    `json:"gameOver,omitempty"`
	Castling    castling    `json:"castling,omitempty"`
	Promotion   promotion   `json:"promotion,omitempty"`
}

const (
	MessageTypeMove = iota
	MessageTypeChat
	MessageTypeErr
	MessageTypeGameStart
	MessageTypeUpdateBoard
	MessageTypeCheck
	MessageTypeGameOver
	MessageTypeCastling
	MessageTypePromotion
)

type move struct {
	RoomID string `json:"roomID"`
	Move   string `json:"move"`
}

type chat struct {
	Username string `json:"username"`
	RoomID   string `json:"roomID"`
	Content  string `json:"content"`
}

type updateBoard struct {
	Board [8][8]int `json:"board"`
	Turn  int       `json:"turn"`
}

type gameStart struct {
	Board   [8][8]int `json:"board"`
	Content string    `json:"content"`
	Turn    int       `json:"turn"`
}

type err struct {
	Error string `json:"error"`
}

type check struct {
	Content string `json:"content"`
}

type gameOver struct {
	Content string `json:"content"`
}

type castling struct {
	RoomID string `json:"roomID"`
	Shift  string `json:"shift"`
}

type promotion struct {
	RoomID    string `json:"roomID"`
	Promotion string `json:"promotion"`
}
