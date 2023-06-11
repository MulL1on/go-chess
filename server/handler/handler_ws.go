package handler

import (
	"log"
)

var rooms = make(map[string]*Room)

func handleMove(m *Message) {
	room, ok := rooms[m.Move.RoomID]
	playerId := m.PlayerID
	if !ok {
		log.Printf("Room %s not found\n", m.Move.RoomID)
		return
	}
	//判断来自哪个玩家的消息
	var site string
	if playerId == room.WhitePlayer.id {
		site = "white"
	} else {
		site = "black"
	}

	switch site {
	case "white":
		if room.turn == TurnBlack {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToWhite <- msg
			return
		}
	case "black":
		if room.turn == TurnWhite {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToBlack <- msg
			return
		}
	}

	if room.isLegalMove(m.Move.Move) {
		room.makeMove(m.Move.Move)
		msg := &Message{
			Type: MessageTypeUpdateBoard,
			UpdateBoard: updateBoard{
				Board: room.Board,
				Turn:  room.turn,
			},
		}
		//log.Println("Move made", msg)
		room.broadcast <- msg
	} else {
		msg := &Message{
			Type: MessageTypeErr,
			Err: err{
				Error: "Illegal move",
			},
		}

		//send error message to player
		if playerId == room.WhitePlayer.id {
			room.sendToWhite <- msg
		} else {
			room.sendToBlack <- msg
		}
		return
	}
}

func handleChat(m *Message) {
	room, ok := rooms[m.Chat.RoomID]
	if !ok {
		log.Printf("Room %s not found\n", m.Chat.RoomID)
		return
	}
	msg := &Message{
		Type: 1,
		Chat: chat{
			Username: m.Username,
			Content:  m.Chat.Content,
		},
	}
	//log.Print("Chat message", msg)
	room.broadcast <- msg
}

func handleMessage(msg *Message) {

	switch msg.Type {
	case MessageTypeMove:
		handleMove(msg)
	case MessageTypeChat:
		handleChat(msg)
	case MessageTypeCastling:
		handleKingRookShift(msg)
	case MessageTypePromotion:
		handlePawnPromotion(msg)
	}
}

func handleKingRookShift(m *Message) {
	room, ok := rooms[m.Castling.RoomID]
	if !ok {
		log.Printf("Room %s not found\n", m.Move.RoomID)
		return
	}
	playerId := m.PlayerID

	if playerId == room.WhitePlayer.id {
		if room.turn == TurnBlack {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToWhite <- msg
			return
		}
	} else {
		if room.turn == TurnWhite {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToBlack <- msg
			return
		}
	}

	//判断是长移动还是短移动
	if m.Castling.Shift == "O-O" {
		if room.canKingsideCastling() {
			room.makeKingsideCastling()
			msg := &Message{
				Type: MessageTypeUpdateBoard,
				UpdateBoard: updateBoard{
					Board: room.Board,
					Turn:  room.turn,
				},
			}
			room.broadcast <- msg
			return
		}

	} else if m.Castling.Shift == "O-O-O" {
		if room.canQueensideCastling() {
			room.makeQueensideCastling()
			msg := &Message{
				Type: MessageTypeUpdateBoard,
				UpdateBoard: updateBoard{
					Board: room.Board,
					Turn:  room.turn,
				},
			}
			room.broadcast <- msg
			return
		}
	}

	//send error message to player
	msg := &Message{
		Type: MessageTypeErr,
		Err: err{
			Error: "Illegal move",
		},
	}

	if playerId == room.WhitePlayer.id {
		room.sendToWhite <- msg
	} else {
		room.sendToBlack <- msg
	}
}

func handlePawnPromotion(m *Message) {
	room, ok := rooms[m.Move.RoomID]
	if !ok {
		log.Printf("Room %s not found\n", m.Move.RoomID)
		return
	}
	playerId := m.PlayerID

	//判断来自哪个玩家的消息
	if playerId == room.WhitePlayer.id {
		if room.turn == TurnBlack {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToWhite <- msg
			return
		}
	} else {
		if room.turn == TurnWhite {
			msg := &Message{
				Type: MessageTypeErr,
				Err: err{
					Error: "Not your turn",
				},
			}
			room.sendToBlack <- msg
			return
		}
	}

	//判断是否为合法的升变
	if room.isLegalPawnPromotion(m.Move.Move) {
		room.makePawnPromotion(m.Move.Move)
		msg := &Message{
			Type: MessageTypeUpdateBoard,
			UpdateBoard: updateBoard{
				Board: room.Board,
				Turn:  room.turn,
			},
		}
		room.broadcast <- msg
		return
	}

	//send error message to player
	msg := &Message{
		Type: MessageTypeErr,
		Err: err{
			Error: "Illegal move",
		},
	}

	if playerId == room.WhitePlayer.id {
		room.sendToWhite <- msg
	} else {
		room.sendToBlack <- msg
	}
}
