package chess

func CreatInitialBoard() [8][8]int {
	board := [8][8]int{
		{10, 8, 9, 11, 12, 9, 8, 10},
		{7, 7, 7, 7, 7, 7, 7, 7},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0},
		{1, 1, 1, 1, 1, 1, 1, 1},
		{4, 2, 3, 5, 6, 3, 2, 4},
	}
	return board
}

func IsLegalMove(board [8][8]int, move string) bool {
	fromFile := int(move[0] - 'a')
	fromRank := int('8' - move[1])
	toFile := int(move[2] - 'a')
	toRank := int('8' - move[3])

	//判断是否出界
	if fromFile < 0 || fromFile > 7 || fromRank < 0 || fromRank > 7 || toFile < 0 || toFile > 7 || toRank < 0 || toRank > 7 {
		return false
	}

	// 判断起始位置和目标位置是否是同一方的棋子
	if board[fromRank][fromFile] > PieceTypeWhiteKing && board[toRank][toFile] > PieceTypeWhiteKing || board[fromRank][fromFile] < PieceTypeWhitePawn && board[toRank][toFile] < PieceTypeBlackPawn {
		return false
	}

	switch board[fromRank][fromFile] {
	case PieceTypeWhitePawn, PieceTypeBlackPawn:
		return validatePawnMove(board, fromRank, fromFile, toRank, toFile)
	case PieceTypeWhiteRook, PieceTypeBlackRook:
		return validateRookMove(board, fromRank, fromFile, toRank, toFile)
	case PieceTypeWhiteKnight, PieceTypeBlackKnight:
		return validateKnightMove(board, fromRank, fromFile, toRank, toFile)
	case PieceTypeWhiteBishop, PieceTypeBlackBishop:
		return validateBishopMove(board, fromRank, fromFile, toRank, toFile)
	case PieceTypeWhiteQueen, PieceTypeBlackQueen:
		return validateQueenMove(board, fromRank, fromFile, toRank, toFile)
	case PieceTypeWhiteKing, PieceTypeBlackKing:
		return validateKingMove(board, fromRank, fromFile, toRank, toFile)
	}
	return false
}

// Rank 行 File 列

func validatePawnMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {
	// 获取起始位置和目标位置的棋子值
	fromPiece := board[fromRank][fromFile]
	toPiece := board[toRank][toFile]

	// 计算移动的行数和列数
	deltaRank := toRank - fromRank
	deltaFile := toFile - fromFile

	// 根据兵的颜色确定移动的方向（正方向或反方向）
	direction := -1                     //白棋
	if fromPiece > PieceTypeWhiteKing { //黑棋
		direction = 1
	}

	//TODO:吃过路兵

	// 判断兵的移动规则
	if deltaRank == direction {
		// 兵向前移动一格的情况
		if deltaFile == 0 && toPiece == PieceTypeNone {
			return true
		}

		// 兵斜向前方攻击对手的棋子
		if abs(deltaFile) == 1 && toPiece > PieceTypeNone {
			return true
		}
	}

	// 兵在初始位置可以选择向前移动两格的情况
	if deltaRank == 2*direction && deltaFile == 0 && (fromRank == 1 || fromRank == 6) && toPiece == PieceTypeNone {
		return true
	}

	return false
}

func validateRookMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {

	// 判断起始位置是不是在同一列
	if fromFile == toFile {
		// 判断起始位置和目标位置是否有棋子阻挡
		for i := fromRank + 1; i < toRank; i++ {
			if board[i][fromFile] != PieceTypeNone {
				return false
			}
		}
		for i := fromRank - 1; i > toRank; i-- {
			if board[i][fromFile] != PieceTypeNone {
				return false
			}
		}
		return true
	}

	// 判断起始位置是不是在同一行
	if fromRank == toRank {
		// 判断起始位置和目标位置是否有棋子阻挡
		for i := fromFile + 1; i < toFile; i++ {
			if board[fromRank][i] != 0 {
				return false
			}
		}
		for i := fromFile - 1; i > toFile; i-- {
			if board[fromRank][i] != 0 {
				return false
			}
		}
		return true
	}

	return false
}

func validateKnightMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {

	deltaRank := abs(toRank - fromRank)
	deltaFile := abs(toFile - fromFile)

	if deltaRank == 2 && deltaFile == 1 || deltaRank == 1 && deltaFile == 2 {
		return true
	}
	return false
}

func validateBishopMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {

	deltaRank := toRank - fromRank
	deltaFile := toFile - fromFile

	//判断是不是在一对角线上
	if abs(deltaRank) != abs(deltaFile) {
		return false
	}

	//判断是否有棋子阻挡
	if deltaRank > 0 && deltaFile > 0 {
		for i := 1; i < deltaRank; i++ {
			if board[fromRank+i][fromFile+i] != PieceTypeNone {
				return false
			}
		}
		return true
	}

	if deltaRank > 0 && deltaFile < 0 {
		for i := 1; i < deltaRank; i++ {
			if board[fromRank+i][fromFile-i] != PieceTypeNone {
				return false
			}
		}
		return true
	}

	if deltaRank < 0 && deltaFile > 0 {
		for i := 1; i < deltaRank; i++ {
			if board[fromRank-i][fromFile+i] != PieceTypeNone {
				return false
			}
		}
		return true
	}

	if deltaRank < 0 && deltaFile < 0 {
		for i := 1; i < deltaRank; i++ {
			if board[fromRank-i][fromFile-i] != PieceTypeNone {
				return false
			}
		}
		return true
	}

	return false
}

func validateQueenMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {

	deltaRank := toRank - fromRank
	deltaFile := toFile - fromFile

	// 判断是不是在一条直线上
	if deltaRank == 0 || deltaFile == 0 {
		// 判断起始位置和目标位置是否有棋子阻挡
		if deltaRank == 0 {
			for i := fromFile + 1; i < toFile; i++ {
				if board[fromRank][i] != PieceTypeNone {
					return false
				}
			}
			for i := fromFile - 1; i > toFile; i-- {
				if board[fromRank][i] != PieceTypeNone {
					return false
				}
			}
			return true
		}

		if deltaFile == 0 {
			for i := fromRank + 1; i < toRank; i++ {
				if board[i][fromFile] != PieceTypeNone {
					return false
				}
			}
			for i := fromRank - 1; i > toRank; i-- {
				if board[i][fromFile] != PieceTypeNone {
					return false
				}
			}
			return true
		}
	}

	// 判断是不是在一对角线上
	if abs(deltaRank) == abs(deltaFile) {
		// 判断起始位置和目标位置是否有棋子阻挡
		if deltaRank > 0 && deltaFile > 0 {
			for i := 1; i < deltaRank; i++ {
				if board[fromRank+i][fromFile+i] != PieceTypeNone {
					return false
				}
			}
			return true
		}

		if deltaRank > 0 && deltaFile < 0 {
			for i := 1; i < deltaRank; i++ {
				if board[fromRank+i][fromFile-i] != PieceTypeNone {
					return false
				}
			}
			return true
		}

		if deltaRank < 0 && deltaFile > 0 {
			for i := 1; i < deltaRank; i++ {
				if board[fromRank-i][fromFile+i] != PieceTypeNone {
					return false
				}
			}
			return true
		}

		if deltaRank < 0 && deltaFile < 0 {
			for i := 1; i < deltaRank; i++ {
				if board[fromRank-i][fromFile-i] != PieceTypeNone {
					return false
				}
			}
			return true
		}
	}

	return false
}

func validateKingMove(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {

	deltaRank := abs(toRank - fromRank)
	deltaFile := abs(toFile - fromFile)
	if deltaRank > 1 || deltaFile > 1 {
		return false
	}
	if IsKingInCheck(board, toRank, toFile) {
		return false
	}

	return true
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func IsKingInCheck(board [8][8]int, rank int, file int) bool {

	//检查横向纵向是否敌方有车或后
	for i := 0; i < 8; i++ {
		if isRankFileInBounds(rank, file+i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank][file+i] == PieceTypeBlackRook || board[rank][file+i] == PieceTypeBlackQueen {
					for j := file - 1; j >= 0; j-- {
						if isRankFileInBounds(rank, j) {
							if board[rank][j] != PieceTypeNone && board[rank][j] != PieceTypeWhiteKing {
								break
							} else if board[rank][j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank][file+i] == PieceTypeWhiteRook || board[rank][file+i] == PieceTypeWhiteQueen {
					for j := file - 1; j >= 0; j-- {
						if isRankFileInBounds(rank, j) {
							if board[rank][j] != PieceTypeNone && board[rank][j] != PieceTypeBlackKing {
								break
							} else if board[rank][j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank, file-i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank][file-i] == PieceTypeBlackRook || board[rank][file-i] == PieceTypeBlackQueen {
					//加回去看看能不能吃掉
					for j := file + 1; j < 8; j++ {
						if isRankFileInBounds(rank, j) {
							if board[rank][j] != PieceTypeNone && board[rank][j] != PieceTypeWhiteKing {
								break
							} else if board[rank][j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank][file-i] == PieceTypeWhiteRook || board[rank][file-i] == PieceTypeWhiteQueen {
					for j := file + 1; j < 8; j++ {
						if isRankFileInBounds(rank, j) {
							if board[rank][j] != PieceTypeNone && board[rank][j] != PieceTypeBlackKing {
								break
							} else if board[rank][j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank+i, file) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+i][file] == PieceTypeBlackRook || board[rank+i][file] == PieceTypeBlackQueen {
					for j := rank - 1; j >= 0; j-- {
						if isRankFileInBounds(j, file) {
							if board[j][file] != PieceTypeNone && board[j][file] != PieceTypeWhiteKing {
								break
							} else if board[j][file] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+i][file] == PieceTypeWhiteRook || board[rank+i][file] == PieceTypeWhiteQueen {
					for j := rank - 1; j >= 0; j-- {
						if isRankFileInBounds(j, file) {
							if board[j][file] != PieceTypeNone && board[j][file] != PieceTypeBlackKing {
								break
							} else if board[j][file] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank-i, file) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank-i][file] == PieceTypeBlackRook || board[rank-i][file] == PieceTypeBlackQueen {
					for j := rank + 1; j < 8; j++ {
						if isRankFileInBounds(j, file) {
							if board[j][file] != PieceTypeNone && board[j][file] != PieceTypeWhiteKing {
								break
							} else if board[j][file] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank-i][file] == PieceTypeWhiteRook || board[rank-i][file] == PieceTypeWhiteQueen {
					for j := rank + 1; j < 8; j++ {
						if isRankFileInBounds(j, file) {
							if board[j][file] != PieceTypeNone && board[j][file] != PieceTypeBlackKing {
								break
							} else if board[j][file] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
	}

	//检查斜向是否有象或后

	for i := 0; i < 8; i++ {
		if isRankFileInBounds(rank+i, file+i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+i][file+i] == PieceTypeBlackBishop || board[rank+i][file+i] == PieceTypeBlackQueen {
					// 检查后是否能够攻击到王
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank+i-j, file+i-j) {
							if board[rank+i-j][file+i-j] != PieceTypeNone && board[rank+i-j][file+i-j] != PieceTypeWhiteKing {
								break
							} else if board[rank+i-j][file+i-j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+i][file+i] == PieceTypeWhiteBishop || board[rank+i][file+i] == PieceTypeWhiteQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank+i-j, file+i-j) {
							if board[rank+i-j][file+i-j] != PieceTypeNone && board[rank+i-j][file+i-j] != PieceTypeBlackKing {
								break
							} else if board[rank+i-j][file+i-j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank+i, file-i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+i][file-i] == PieceTypeBlackBishop || board[rank+i][file-i] == PieceTypeBlackQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank+i-j, file-i+j) {
							if board[rank+i-j][file-i+j] != PieceTypeNone && board[rank+i-j][file-i+j] != PieceTypeWhiteKing {
								break
							} else if board[rank+i-j][file-i+j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+i][file-i] == PieceTypeWhiteBishop || board[rank+i][file-i] == PieceTypeWhiteQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank+i-j, file-i+j) {
							if board[rank+i-j][file-i+j] != PieceTypeNone && board[rank+i-j][file-i+j] != PieceTypeBlackKing {
								break
							} else if board[rank+i-j][file-i+j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank-i, file+i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank-i][file+i] == PieceTypeBlackBishop || board[rank-i][file+i] == PieceTypeBlackQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank-i+j, file+i-j) {
							if board[rank-i+j][file+i-j] != PieceTypeNone && board[rank-i+j][file+i-j] != PieceTypeWhiteKing {
								break
							} else if board[rank-i+j][file+i-j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank-i][file+i] == PieceTypeWhiteBishop || board[rank-i][file+i] == PieceTypeWhiteQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank-i+j, file+i-j) {
							if board[rank-i+j][file+i-j] != PieceTypeNone && board[rank-i+j][file+i-j] != PieceTypeBlackKing {
								break
							} else if board[rank-i+j][file+i-j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
		if isRankFileInBounds(rank-i, file-i) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank-i][file-i] == PieceTypeBlackBishop || board[rank-i][file-i] == PieceTypeBlackQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank-i+j, file-i+j) {
							if board[rank-i+j][file-i+j] != PieceTypeNone && board[rank-i+j][file-i+j] != PieceTypeWhiteKing {
								break
							} else if board[rank-i+j][file-i+j] == PieceTypeWhiteKing {
								return true
							}
						}
					}
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank-i][file-i] == PieceTypeWhiteBishop || board[rank-i][file-i] == PieceTypeWhiteQueen {
					for j := 1; j < 8; j++ {
						if isRankFileInBounds(rank-i+j, file-i+j) {
							if board[rank-i+j][file-i+j] != PieceTypeNone && board[rank-i+j][file-i+j] != PieceTypeBlackKing {
								break
							} else if board[rank-i+j][file-i+j] == PieceTypeBlackKing {
								return true
							}
						}
					}
				}
			}
		}
	}

	//检查马的位置
	direction := [][]int{{1, 2}, {2, 1}, {-1, 2}, {-2, 1}, {1, -2}, {2, -1}, {-1, -2}, {-2, -1}}
	for _, d := range direction {
		if isRankFileInBounds(rank+d[0], file+d[1]) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeBlackKnight {
					return true
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeWhiteKnight {
					return true
				}
			}
		}
	}

	//检查兵的位置
	direction = [][]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, d := range direction {
		if isRankFileInBounds(rank+d[0], file+d[1]) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeBlackPawn {
					return true
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeWhitePawn {
					return true
				}
			}
		}
	}

	//检查王的位置
	direction = [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, d := range direction {
		if isRankFileInBounds(rank+d[0], file+d[1]) {
			if board[rank][file] == PieceTypeWhiteKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeBlackKing {
					return true
				}
			} else if board[rank][file] == PieceTypeBlackKing {
				if board[rank+d[0]][file+d[1]] == PieceTypeWhiteKing {
					return true
				}
			}
		}
	}
	return false
}

func CanKingEscapeCheck(board [8][8]int, rank int, file int) bool {
	direction := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	for _, d := range direction {
		if isRankFileInBounds(rank+d[0], file+d[1]) {
			if !IsKingInCheck(board, rank+d[0], file+d[1]) {
				return true
			}
		}
	}
	return false
}

func GetPosition(PieceType int, board [8][8]int) (int, int) {
	for i, rank := range board {
		for j, piece := range rank {
			if piece == PieceType {
				return i, j
			}
		}
	}
	return -1, -1
}

func isRankFileInBounds(rank int, file int) bool {
	return rank >= 0 && rank < 8 && file >= 0 && file < 8
}

func CanWhiteKingsideCastling(board [8][8]int) bool {
	//判断车和王之间是否有棋子
	if board[7][5] == PieceTypeNone && board[7][6] == PieceTypeNone {
		//判断王是否在被将军
		if !IsKingInCheck(board, 7, 4) {
			//判断王是否经过被将军的位置
			if !IsKingInCheck(board, 7, 5) && !IsKingInCheck(board, 7, 6) {
				return true
			}
		}
	}
	return false
}

func CanBlackKingsideCastling(board [8][8]int) bool {
	//判断车和王之间是否有棋子
	if board[0][5] == PieceTypeNone && board[0][6] == PieceTypeNone {
		//判断王是否在被将军
		if !IsKingInCheck(board, 0, 4) {
			//判断王是否经过被将军的位置
			if !IsKingInCheck(board, 0, 5) && !IsKingInCheck(board, 0, 6) {
				return true
			}
		}
	}
	return false
}

func CanWhiteQueensideCastling(board [8][8]int) bool {
	//判断车和王之间是否有棋子
	if board[7][1] == PieceTypeNone && board[7][2] == PieceTypeNone && board[7][3] == PieceTypeNone {
		//判断王是否在被将军
		if !IsKingInCheck(board, 7, 4) {
			//判断王是否经过被将军的位置
			if !IsKingInCheck(board, 7, 3) && !IsKingInCheck(board, 7, 2) {
				return true
			}
		}
	}
	return false
}

func CanBlackQueensideCastling(board [8][8]int) bool {
	//判断车和王之间是否有棋子
	if board[0][1] == PieceTypeNone && board[0][2] == PieceTypeNone && board[0][3] == PieceTypeNone {
		//判断王是否在被将军
		if !IsKingInCheck(board, 0, 4) {
			//判断王是否经过被将军的位置
			if !IsKingInCheck(board, 0, 3) && !IsKingInCheck(board, 0, 2) {
				return true
			}
		}
	}
	return false
}

func IsLegalPawnPromotion(board [8][8]int, fromRank int, fromFile int, toRank int, toFile int) bool {
	deltaFile := abs(toFile - fromFile)

	if deltaFile > 1 {
		return false
	}

	if board[fromRank][fromFile] == PieceTypeWhitePawn {
		if fromRank == 6 && toRank == 7 {
			return true
		}
	} else if board[fromRank][fromFile] == PieceTypeBlackPawn {
		if fromRank == 1 && toRank == 0 {
			return true
		}
	}
	return false
}
