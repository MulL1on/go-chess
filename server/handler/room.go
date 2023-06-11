package handler

import (
	"github.com/google/uuid"
	"go-chess/server/chess"
)

const (
	TurnNone = iota
	TurnWhite
	TurnBlack
)

type Room struct {
	ID                      string
	WhitePlayer             *Client
	BlackPlayer             *Client
	Board                   [8][8]int
	turn                    int
	registry                chan *Client
	unregistry              chan *Client
	broadcast               chan *Message
	sendToWhite             chan *Message
	sendToBlack             chan *Message
	canWhiteKingsideCastle  bool
	canWhiteQueensideCastle bool
	canBlackKingsideCastle  bool
	canBlackQueensideCastle bool
}

func NewRoom() *Room {
	return &Room{
		ID:                      newUUID(),
		WhitePlayer:             nil,
		BlackPlayer:             nil,
		Board:                   chess.CreatInitialBoard(),
		turn:                    1,
		broadcast:               make(chan *Message),
		registry:                make(chan *Client),
		sendToWhite:             make(chan *Message),
		sendToBlack:             make(chan *Message),
		unregistry:              make(chan *Client),
		canBlackKingsideCastle:  true,
		canBlackQueensideCastle: true,
		canWhiteKingsideCastle:  true,
		canWhiteQueensideCastle: true,
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.registry:
			if r.WhitePlayer == nil {
				r.WhitePlayer = client
			} else if r.BlackPlayer == nil {
				r.BlackPlayer = client
			}
		case client := <-r.unregistry:
			if r.WhitePlayer == client {
				r.WhitePlayer = nil
			} else if r.BlackPlayer == client {
				r.BlackPlayer = nil
			}
		case message := <-r.broadcast:
			if r.WhitePlayer != nil {
				r.WhitePlayer.send <- message
			}
			if r.BlackPlayer != nil {
				r.BlackPlayer.send <- message
			}
		case message := <-r.sendToWhite:
			if r.WhitePlayer != nil {
				r.WhitePlayer.send <- message
			}
		case message := <-r.sendToBlack:
			if r.BlackPlayer != nil {
				r.BlackPlayer.send <- message
			}
		}
	}
}

func newUUID() string {
	return uuid.New().String()
}

func (r *Room) makeMove(move string) {
	fromFile := int(move[0] - 'a')
	fromRank := int('8' - move[1])
	toFile := int(move[2] - 'a')
	toRank := int('8' - move[3])
	r.Board[toRank][toFile] = r.Board[fromRank][fromFile]
	r.Board[fromRank][fromFile] = 0

	//判断是否移动了王
	if r.Board[fromRank][fromFile] == chess.PieceTypeBlackKing {
		r.canBlackKingsideCastle = false
		r.canBlackQueensideCastle = false
	} else if r.Board[fromRank][fromFile] == chess.PieceTypeWhiteKing {
		r.canWhiteKingsideCastle = false
		r.canWhiteQueensideCastle = false
	}

	//判断是否移动了车
	if fromRank == 0 && fromFile == 0 {
		r.canBlackQueensideCastle = false
	} else if fromRank == 0 && fromFile == 7 {
		r.canBlackKingsideCastle = false
	} else if fromRank == 7 && fromFile == 0 {
		r.canWhiteQueensideCastle = false
	} else if fromRank == 7 && fromFile == 7 {
		r.canWhiteKingsideCastle = false
	}

	//判断是否将军
	if r.turn == TurnBlack {
		//获取白方将军的位置
		//TODO: 优化
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if chess.IsKingInCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "白方将军",
				},
			}
			r.broadcast <- msg
		}

	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if chess.IsKingInCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "黑方将军",
				},
			}
			r.broadcast <- msg
		}
	}

	//判断游戏是否结束
	if r.turn == TurnBlack {
		//获取白方将军的位置
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "黑方胜利",
				},
			}
			r.broadcast <- msg
			return
		}

	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "白方胜利",
				},
			}
			r.broadcast <- msg
			return
		}
	}

	//回合交换
	r.turn = r.turn%2 + 1
}

func (r *Room) isLegalMove(move string) bool {
	if len(move) != 4 {
		return false
	}
	fromFile := int(move[0] - 'a')
	fromRank := int('8' - move[1])
	toFile := int(move[2] - 'a')
	toRank := int('8' - move[3])
	if r.Board[fromRank][fromFile] == 0 {
		return false
	}

	//判断是否是自己的棋子
	if r.Board[fromRank][fromFile] < chess.PieceTypeBlackPawn && r.turn == TurnBlack {
		return false
	}
	if r.Board[fromRank][fromFile] > chess.PieceTypeWhiteKing && r.turn == TurnWhite {
		return false
	}

	//判断是否会导致将军
	if r.turn == TurnBlack {
		board := r.Board
		board[toRank][toFile] = board[fromRank][fromFile]
		board[fromRank][fromFile] = 0
		if chess.IsKingInCheck(r.Board, toRank, toFile) {
			return false
		}
	} else if r.turn == TurnWhite {
		board := r.Board
		board[toRank][toFile] = board[fromRank][fromFile]
		board[fromRank][fromFile] = 0
		if chess.IsKingInCheck(r.Board, toRank, toFile) {
			return false
		}
	}

	//进一步判断
	return chess.IsLegalMove(r.Board, move)
}

func (r *Room) canKingsideCastling() bool {
	if r.turn == TurnBlack {
		if !r.canBlackKingsideCastle {
			return false
		}

		//判断是否会导致将军
		board := r.Board
		board[0][4] = chess.PieceTypeNone
		board[0][7] = chess.PieceTypeNone
		board[0][6] = chess.PieceTypeBlackKing
		board[0][5] = chess.PieceTypeBlackRook
		if chess.IsKingInCheck(board, 0, 6) {
			return false
		}

		//进一步判断
		if chess.CanBlackKingsideCastling(r.Board) {
			return true
		}
	}

	if r.turn == TurnWhite {
		if !r.canWhiteKingsideCastle {
			return false
		}

		//判断是否会导致将军
		board := r.Board
		board[7][4] = chess.PieceTypeNone
		board[7][7] = chess.PieceTypeNone
		board[7][6] = chess.PieceTypeWhiteKing
		board[7][5] = chess.PieceTypeWhiteRook
		if chess.IsKingInCheck(board, 7, 6) {
			return false
		}

		//进一步判断
		if chess.CanWhiteKingsideCastling(r.Board) {
			return true
		}
	}

	return false
}

func (r *Room) canQueensideCastling() bool {
	if r.turn == TurnBlack {
		if !r.canBlackQueensideCastle {
			return false
		}
		board := r.Board
		board[0][4] = chess.PieceTypeNone
		board[0][0] = chess.PieceTypeNone
		board[0][2] = chess.PieceTypeBlackKing
		board[0][3] = chess.PieceTypeBlackRook
		if chess.IsKingInCheck(board, 0, 2) {
			return false
		}

		if chess.CanBlackQueensideCastling(r.Board) {
			return true
		}
	}

	if r.turn == TurnWhite {
		if !r.canWhiteQueensideCastle {
			return false
		}
		board := r.Board
		board[7][4] = chess.PieceTypeNone
		board[7][0] = chess.PieceTypeNone
		board[7][2] = chess.PieceTypeWhiteKing
		board[7][3] = chess.PieceTypeWhiteRook
		if chess.IsKingInCheck(board, 7, 2) {
			return false
		}

		if chess.CanWhiteQueensideCastling(r.Board) {
			return true
		}
	}

	return false
}

func (r *Room) makeKingsideCastling() {
	if r.turn == TurnBlack {
		r.Board[0][4] = chess.PieceTypeNone
		r.Board[0][7] = chess.PieceTypeNone
		r.Board[0][6] = chess.PieceTypeBlackKing
		r.Board[0][5] = chess.PieceTypeBlackRook
		r.canBlackKingsideCastle = false
		r.canBlackQueensideCastle = false

	} else if r.turn == TurnWhite {
		r.Board[7][4] = chess.PieceTypeNone
		r.Board[7][7] = chess.PieceTypeNone
		r.Board[7][6] = chess.PieceTypeWhiteKing
		r.Board[7][5] = chess.PieceTypeWhiteRook
		r.canWhiteKingsideCastle = false
		r.canWhiteQueensideCastle = false
	}

	//回合改变
	r.turn = r.turn%2 + 1

	//判断是否将军
	if r.turn == TurnBlack {
		//获取白方将军的位置
		//TODO: 优化
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if chess.IsKingInCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "白方将军",
				},
			}
			r.broadcast <- msg
		}
	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if chess.IsKingInCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "黑方将军",
				},
			}
			r.broadcast <- msg
		}
	}

	//判断游戏是否结束
	if r.turn == TurnBlack {
		//获取白方将军的位置
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "黑方胜利",
				},
			}
			r.broadcast <- msg
			return
		}

	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "白方胜利",
				},
			}
			r.broadcast <- msg
			return
		}
	}
}

func (r *Room) makeQueensideCastling() {
	if r.turn == TurnBlack {
		r.Board[0][4] = chess.PieceTypeNone
		r.Board[0][0] = chess.PieceTypeNone
		r.Board[0][2] = chess.PieceTypeBlackKing
		r.Board[0][3] = chess.PieceTypeBlackRook
		r.canBlackKingsideCastle = false
		r.canBlackQueensideCastle = false
	} else if r.turn == TurnWhite {
		r.Board[7][4] = chess.PieceTypeNone
		r.Board[7][0] = chess.PieceTypeNone
		r.Board[7][2] = chess.PieceTypeWhiteKing
		r.Board[7][3] = chess.PieceTypeWhiteRook
		r.canWhiteKingsideCastle = false
		r.canWhiteQueensideCastle = false
	}

	// 回合改变
	r.turn = r.turn%2 + 1

	//判断是否将军
	if r.turn == TurnBlack {
		//获取白方将军的位置
		//TODO: 优化
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if chess.IsKingInCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "白方将军",
				},
			}
			r.broadcast <- msg
		}
	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if chess.IsKingInCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "黑方将军",
				},
			}
			r.broadcast <- msg
		}
	}

	//判断游戏是否结束
	if r.turn == TurnBlack {
		//获取白方将军的位置
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "黑方胜利",
				},
			}
			r.broadcast <- msg
			return
		}

	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "白方胜利",
				},
			}
			r.broadcast <- msg
			return
		}
	}
}

func (r *Room) isLegalPawnPromotion(promotion string) bool {
	if len(promotion) != 6 {
		return false
	}
	fromRank := int('8' - promotion[1])
	fromFile := int(promotion[0] - 'a')
	toRank := int('8' - promotion[3])
	toFile := int(promotion[2] - 'a')

	//判断是否会导致将军
	if r.turn == TurnBlack {
		board := r.Board
		board[toRank][toFile] = board[fromRank][fromFile]
		board[fromRank][fromFile] = 0
		if chess.IsKingInCheck(r.Board, toRank, toFile) {
			return false
		}
	} else if r.turn == TurnWhite {
		board := r.Board
		board[toRank][toFile] = board[fromRank][fromFile]
		board[fromRank][fromFile] = 0
		if chess.IsKingInCheck(r.Board, toRank, toFile) {
			return false
		}
	}

	//进一步判断是否合法
	if chess.IsLegalPawnPromotion(r.Board, fromRank, fromFile, toRank, toFile) {
		return true
	}

	return false
}

func (r *Room) makePawnPromotion(promotion string) {
	toRank := int('8' - promotion[3])
	toFile := int(promotion[2] - 'a')
	var toPiece int
	if r.turn == TurnBlack {
		switch promotion[5] {
		case 'Q':
			toPiece = chess.PieceTypeBlackQueen
		case 'R':
			toPiece = chess.PieceTypeBlackRook
		case 'B':
			toPiece = chess.PieceTypeBlackBishop
		case 'N':
			toPiece = chess.PieceTypeBlackKnight
		}
	} else {
		switch promotion[5] {
		case 'Q':
			toPiece = chess.PieceTypeWhiteQueen
		case 'R':
			toPiece = chess.PieceTypeWhiteRook
		case 'B':
			toPiece = chess.PieceTypeWhiteBishop
		case 'N':
			toPiece = chess.PieceTypeWhiteKnight
		}
	}
	r.Board[toRank][toFile] = toPiece

	//判断是否将军
	if r.turn == TurnBlack {
		//获取白方将军的位置
		//TODO: 优化
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if chess.IsKingInCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "白方将军",
				},
			}
			r.broadcast <- msg
		}
	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if chess.IsKingInCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeCheck,
				Check: check{
					Content: "黑方将军",
				},
			}
			r.broadcast <- msg
		}
	}

	//判断游戏是否结束
	if r.turn == TurnBlack {
		//获取白方将军的位置
		whiteRank, whiteFile := chess.GetPosition(chess.PieceTypeWhiteKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, whiteRank, whiteFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "黑方胜利",
				},
			}
			r.broadcast <- msg
			return
		}

	} else if r.turn == TurnWhite {
		//获取黑方将军的位置
		blackRank, blackFile := chess.GetPosition(chess.PieceTypeBlackKing, r.Board)
		if !chess.CanKingEscapeCheck(r.Board, blackRank, blackFile) {
			msg := &Message{
				Type: MessageTypeGameOver,
				GameOver: gameOver{
					Content: "白方胜利",
				},
			}
			r.broadcast <- msg
			return
		}
	}

	r.turn = r.turn%2 + 1
}
