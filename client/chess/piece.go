package chess

const (
	PieceTypeNone = iota
	PieceTypeWhitePawn
	PieceTypeWhiteKnight
	PieceTypeWhiteBishop
	PieceTypeWhiteRook
	PieceTypeWhiteQueen
	PieceTypeWhiteKing

	PieceTypeBlackPawn
	PieceTypeBlackKnight
	PieceTypeBlackBishop
	PieceTypeBlackRook
	PieceTypeBlackQueen
	PieceTypeBlackKing
)

func ParsePiece(pieceType int) rune {
	switch pieceType {
	case PieceTypeWhitePawn:
		return '♙'
	case PieceTypeWhiteKnight:
		return '♘'
	case PieceTypeWhiteBishop:
		return '♗'
	case PieceTypeWhiteRook:
		return '♖'
	case PieceTypeWhiteQueen:
		return '♕'
	case PieceTypeWhiteKing:
		return '♔'
	case PieceTypeBlackPawn:
		return '♟'
	case PieceTypeBlackKnight:
		return '♞'
	case PieceTypeBlackBishop:
		return '♝'
	case PieceTypeBlackRook:
		return '♜'
	case PieceTypeBlackQueen:
		return '♛'
	case PieceTypeBlackKing:
		return '♚'
	default:
		return ' '
	}
}
