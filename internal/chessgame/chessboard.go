package chessgame

import "fmt"

const (
	aSideCastlingXValue = 3
	hSideCastlingYValue = 5
)

func (cb *chessBoard) initializeBoard() {
	// Initialize whole board with empty squares
	for file := 0; file < 8; file++ {
		cb.board[file] = [8]chessPiece{
			emptyPiece, emptyPiece, emptyPiece, emptyPiece,
			emptyPiece, emptyPiece, emptyPiece, emptyPiece,
		}
	}

	cb.board[0][0] = chessPiece{rook, white}
	cb.board[1][0] = chessPiece{knight, white}
	cb.board[2][0] = chessPiece{bishop, white}
	cb.board[3][0] = chessPiece{queen, white}
	cb.board[4][0] = chessPiece{king, white}
	cb.board[5][0] = chessPiece{bishop, white}
	cb.board[6][0] = chessPiece{knight, white}
	cb.board[7][0] = chessPiece{rook, white}

	cb.board[0][1] = chessPiece{pawn, white}
	cb.board[1][1] = chessPiece{pawn, white}
	cb.board[2][1] = chessPiece{pawn, white}
	cb.board[3][1] = chessPiece{pawn, white}
	cb.board[4][1] = chessPiece{pawn, white}
	cb.board[5][1] = chessPiece{pawn, white}
	cb.board[6][1] = chessPiece{pawn, white}
	cb.board[7][1] = chessPiece{pawn, white}

	cb.board[0][7] = chessPiece{rook, black}
	cb.board[1][7] = chessPiece{knight, black}
	cb.board[2][7] = chessPiece{bishop, black}
	cb.board[3][7] = chessPiece{queen, black}
	cb.board[4][7] = chessPiece{king, black}
	cb.board[5][7] = chessPiece{bishop, black}
	cb.board[6][7] = chessPiece{knight, black}
	cb.board[7][7] = chessPiece{rook, black}

	cb.board[0][6] = chessPiece{pawn, black}
	cb.board[1][6] = chessPiece{pawn, black}
	cb.board[2][6] = chessPiece{pawn, black}
	cb.board[3][6] = chessPiece{pawn, black}
	cb.board[4][6] = chessPiece{pawn, black}
	cb.board[5][6] = chessPiece{pawn, black}
	cb.board[6][6] = chessPiece{pawn, black}
	cb.board[7][6] = chessPiece{pawn, black}

	cb.enpassantSquare = nilSquare

	cb.castlingState.a1RookMoved = false
	cb.castlingState.a8RookMoved = false
	cb.castlingState.h1RookMoved = false
	cb.castlingState.h8RookMoved = false
	cb.castlingState.blackKingMoved = false
	cb.castlingState.whiteKingMoved = false
}

func (cb *chessBoard) deepCopy() chessBoard {
	var newBoard chessBoard

	for x := range 8 {
		for y := range 8 {
			newBoard.board[x][y] = cb.board[x][y]
		}
	}

	newBoard.enpassantSquare = cb.enpassantSquare

	return newBoard
}

func (cb *chessBoard) getPiece(square vector2) chessPiece {
	return cb.board[square.x][square.y]
}

func (cb *chessBoard) getKingSquare(playerColor pieceColor) vector2 {
	for x := range 8 {
		for y := range 8 {
			if cb.getPiece(vector2{x, y}).color == playerColor && cb.getPiece(vector2{x, y}).pieceType == king {
				return vector2{x, y}
			}
		}
	}

	panic(fmt.Sprintf("No %+v king found in board!", playerColor))
}

func (cb *chessBoard) setPiece(square vector2, piece chessPiece) {
	cb.board[square.x][square.y] = piece
}

func (cb *chessBoard) isEnPassantMove(square vector2, targetSquare vector2) (bool, vector2) {
	chessPiece := cb.getPiece(square)
	if chessPiece.pieceType != pawn {
		return false, nilSquare
	}
	if chessPiece.color == white && targetSquare == square.add((vector2{0, 2})) {
		return true, square.add(vector2{0, 1})
	}
	if chessPiece.color == black && targetSquare == square.add((vector2{0, -2})) {
		return true, square.add(vector2{0, -1})
	}
	return false, nilSquare
}

func (cb *chessBoard) updateCastlingState(movedFromSquare vector2) {
	if movedFromSquare.equals(vector2{0, 0}) {
		cb.castlingState.a1RookMoved = true
	}
	if movedFromSquare.equals(vector2{7, 0}) {
		cb.castlingState.h1RookMoved = true
	}
	if movedFromSquare.equals(vector2{0, 7}) {
		cb.castlingState.a8RookMoved = true
	}
	if movedFromSquare.equals(vector2{7, 7}) {
		cb.castlingState.h8RookMoved = true
	}

	if movedFromSquare.equals(vector2{4, 0}) {
		cb.castlingState.whiteKingMoved = true
	}
	if movedFromSquare.equals(vector2{4, 7}) {
		cb.castlingState.blackKingMoved = true
	}
}

func (cb *chessBoard) promotionTriggeredOnSquare(square vector2) bool {
	chessPiece := cb.getPiece(square)
	if chessPiece.pieceType != pawn {
		return false
	}
	return (chessPiece.color == black && square.y == 0) ||
		(chessPiece.color == white && square.y == 7)
}

// Expects given move to be legal
func (cb *chessBoard) movePiece(square vector2, targetSquare vector2) {
	piece := cb.getPiece(square)

	cb.setPiece(targetSquare, piece)
	cb.setPiece(square, emptyPiece)

	// Move Rook if move is a caslting move
	if piece.pieceType == king && abs(square.x-targetSquare.x) > 1 {
		moveARook := square.x > targetSquare.x // False implies that H rook must be moved
		if moveARook {
			rookSquare := vector2{0, targetSquare.y}
			cb.setPiece(rookSquare, emptyPiece)
			cb.setPiece(vector2{aSideCastlingXValue, targetSquare.y}, chessPiece{rook, piece.color})
			cb.updateCastlingState(rookSquare)
		} else {
			rookSquare := vector2{7, targetSquare.y}
			cb.setPiece(rookSquare, emptyPiece)
			cb.setPiece(vector2{hSideCastlingYValue, targetSquare.y}, chessPiece{rook, piece.color})
			cb.updateCastlingState(rookSquare)
		}
	}

	// Updates Enpassant state if applicable
	createsEnpassant, enPassantSquare := cb.isEnPassantMove(square, targetSquare)
	if createsEnpassant {
		cb.enpassantSquare = enPassantSquare
	} else {
		cb.enpassantSquare = nilSquare
	}

	// Update Castling state if applicable
	cb.updateCastlingState(square)
}

func (cb *chessBoard) isEmpty(square vector2) bool {
	return cb.getPiece(square) == emptyPiece
}

func (cb *chessBoard) inBoard(square vector2) bool {
	return square.x >= 0 &&
		square.x < 8 &&
		square.y >= 0 &&
		square.y < 8
}

func (cb *chessBoard) isDarkSquare(square vector2) bool {
	return (square.x+square.y)%2 == 0
}

func (cb *chessBoard) getAllPieceSquares(pieceColor pieceColor) []vector2 {
	pieces := []vector2{}
	for x := range 8 {
		for y := range 8 {
			piece := cb.getPiece(vector2{x, y})
			if piece.color == pieceColor {
				pieces = append(pieces, vector2{x, y})
			}
		}
	}
	return pieces
}

func (cb *chessBoard) moveInducesCheck(pieceSquare vector2, targetSquare vector2, playerColor pieceColor) bool {
	testBoard := cb.deepCopy()

	testBoard.movePiece(pieceSquare, targetSquare)

	return testBoard.playerInCheck(playerColor)
}

func (cb *chessBoard) kingInCheck(kingSquare vector2, kingColor pieceColor) bool {
	opponentSquares := cb.getAllPieceSquares(kingColor.oppositeColor())
	for _, opponentSquare := range opponentSquares {
		attackedSquares := cb.getAttackedSquares(opponentSquare)
		if contains(attackedSquares, kingSquare) {
			return true
		}
	}
	return false
}

func (cb *chessBoard) playerInCheck(playerColor pieceColor) bool {
	kingSquare := cb.getKingSquare(playerColor)

	return cb.kingInCheck(kingSquare, playerColor)
}

func (cb *chessBoard) playerInCheckMate(playerColor pieceColor) bool {

	if !cb.playerInCheck(playerColor) {
		return false
	}

	playerMoves := cb.getAllValidMovesForPlayer(playerColor)
	for _, playerMove := range playerMoves {
		boardCopy := cb.deepCopy()
		boardCopy.movePiece(playerMove[0], playerMove[1])
		if !boardCopy.playerInCheck(playerColor) {
			return false
		}
	}
	return true
}

func (cb *chessBoard) getAllValidMovesForPlayer(playerColor pieceColor) [][]vector2 {
	playerSquares := cb.getAllPieceSquares(playerColor)

	var allPlayerMoves [][]vector2

	for _, playerSquare := range playerSquares {
		playerMoveSquares := cb.getValidMoves(playerSquare, cb.getPiece(playerSquare))
		for _, playerMoveSquare := range playerMoveSquares {
			allPlayerMoves = append(allPlayerMoves, []vector2{playerSquare, playerMoveSquare})
		}
	}

	return allPlayerMoves
}

func (cb *chessBoard) getValidMoves(square vector2, chessPiece chessPiece) []vector2 {
	switch chessPiece.pieceType {
	case pawn:
		return cb.getPawnMoves(chessPiece.color, square)
	case bishop:
		return cb.getBishopMoves(chessPiece.color, square)
	case knight:
		return cb.getKnightMoves(chessPiece.color, square)
	case rook:
		return cb.getRookMoves(chessPiece.color, square)
	case queen:
		return cb.getQueenMoves(chessPiece.color, square)
	case king:
		return cb.getKingMoves(chessPiece.color, square)
	default:
		return []vector2{}
	}
}

// Returns squares that a piece can attack.  Differs only from getValidMoves for pawns.  Attacked squares include squares that a piece can take (i.e. squares occupied by a piece of the opposite color)
func (cb *chessBoard) getAttackedSquares(square vector2) []vector2 {
	chessPiece := cb.getPiece(square)
	if chessPiece.pieceType == pawn {
		return cb.getPawnAttackedSquares(chessPiece.color, square)
	} else if chessPiece.pieceType == king {
		return cb.getKingAttackedSquares(chessPiece.color, square)
	} else if chessPiece.pieceType == bishop {
		return cb.getBishopAttackedSquares(chessPiece.color, square)
	} else if chessPiece.pieceType == knight {
		return cb.getKnightAttackedSquares(chessPiece.color, square)
	} else if chessPiece.pieceType == rook {
		return cb.getRookAttackedSquares(chessPiece.color, square)
	} else if chessPiece.pieceType == queen {
		return cb.getQueenAttackedSquares(chessPiece.color, square)
	}
	panic(fmt.Sprintf("%v is a non-existent piece!", chessPiece))
}

func (cb *chessBoard) getAllAttackedSquares(attackingColor pieceColor) []vector2 {

	// Use a map to avoid duplicates
	allAttackedSquaresMap := make(map[vector2]bool)

	attackingSquares := cb.getAllPieceSquares(attackingColor)
	for _, attackingSqaure := range attackingSquares {
		attackedSquares := cb.getAttackedSquares(attackingSqaure)
		for _, attackedSquare := range attackedSquares {
			allAttackedSquaresMap[attackedSquare] = true
		}
	}

	var allAttackedSquares []vector2

	for key := range allAttackedSquaresMap {
		allAttackedSquares = append(allAttackedSquares, key)
	}

	return allAttackedSquares
}

func (cb *chessBoard) isSquareAttacked(square vector2, attackingColor pieceColor) bool {
	for _, attackedSquare := range cb.getAllAttackedSquares(attackingColor) {
		if square.equals(attackedSquare) {
			return true
		}
	}

	return false
}

func (cb *chessBoard) getCastlingMoves(kingColor pieceColor) []vector2 {
	// King has already moved, no legal caslting squares
	if (kingColor == white && cb.castlingState.whiteKingMoved) ||
		(kingColor == black && cb.castlingState.blackKingMoved) {
		return nil
	}

	// Initialize color generic state
	var aRookMoved, hRookMoved bool
	var rank int
	if kingColor == white {
		aRookMoved, hRookMoved = cb.castlingState.a1RookMoved, cb.castlingState.h1RookMoved
		rank = 0
	} else {
		aRookMoved, hRookMoved = cb.castlingState.a8RookMoved, cb.castlingState.h8RookMoved
		rank = 7
	}

	// Squares that need to be empty for castling to be legal
	aSideSquares := []vector2{{1, rank}, {2, rank}, {3, rank}}
	hSideSquares := []vector2{{6, rank}, {5, rank}}

	// Squares that cannot be attacked for castling to be legal
	aSideUnattacked := []vector2{{2, rank}, {3, rank}}
	hSideUnattacked := hSideSquares

	// Check legality of A side castling
	aSideCastleLegal := true
	if aRookMoved {
		aSideCastleLegal = false
	}
	for _, aSideSquare := range aSideSquares {
		if cb.getPiece(aSideSquare) != emptyPiece {
			aSideCastleLegal = false
			break
		}
	}
	for _, aSideUnattackedSquare := range aSideUnattacked {
		if cb.isSquareAttacked(aSideUnattackedSquare, kingColor.oppositeColor()) {
			aSideCastleLegal = false
			break
		}
	}

	// Check legality of H side castling
	hSideCastleLegal := true
	if hRookMoved {
		hSideCastleLegal = false
	}
	for _, hSideSquare := range hSideSquares {
		if cb.getPiece(hSideSquare) != emptyPiece {
			hSideCastleLegal = false
			break
		}
	}
	for _, hSideUnattackedSquare := range hSideUnattacked {
		if cb.isSquareAttacked(hSideUnattackedSquare, kingColor.oppositeColor()) {
			hSideCastleLegal = false
			break
		}
	}

	var castlingMoves []vector2
	if aSideCastleLegal {
		castlingMoves = append(castlingMoves, vector2{2, rank})
	}
	if hSideCastleLegal {
		castlingMoves = append(castlingMoves, vector2{6, rank})
	}

	return castlingMoves

}

func (cb *chessBoard) getKingMoves(kingColor pieceColor, kingSquare vector2) []vector2 {
	kingMoveSquares := []vector2{
		kingSquare.add(vector2{1, 1}),
		kingSquare.add(vector2{1, 0}),
		kingSquare.add(vector2{1, -1}),
		kingSquare.add(vector2{0, 1}),
		kingSquare.add(vector2{0, -1}),
		kingSquare.add(vector2{-1, 1}),
		kingSquare.add(vector2{-1, 0}),
		kingSquare.add(vector2{-1, -1}),
	}

	kingMoves := []vector2{}

	for _, kingMoveSquare := range kingMoveSquares {
		if !cb.inBoard(kingMoveSquare) {
			continue
		}
		if cb.getPiece(kingMoveSquare).color == kingColor {
			continue
		}
		if cb.moveInducesCheck(kingSquare, kingMoveSquare, kingColor) {
			continue
		}
		kingMoves = append(kingMoves, kingMoveSquare)
	}

	kingMoves = append(kingMoves, cb.getCastlingMoves(kingColor)...)

	return kingMoves
}

func (cb *chessBoard) getKingAttackedSquares(kingColor pieceColor, kingSquare vector2) []vector2 {
	kingMoveSquares := []vector2{
		kingSquare.add(vector2{1, 1}),
		kingSquare.add(vector2{1, 0}),
		kingSquare.add(vector2{1, -1}),
		kingSquare.add(vector2{0, 1}),
		kingSquare.add(vector2{0, -1}),
		kingSquare.add(vector2{-1, 1}),
		kingSquare.add(vector2{-1, 0}),
		kingSquare.add(vector2{-1, -1}),
	}

	kingMoves := []vector2{}

	for _, kingMoveSquare := range kingMoveSquares {
		if !cb.inBoard(kingMoveSquare) {
			continue
		}
		if cb.getPiece(kingMoveSquare).color == kingColor {
			continue
		}
		kingMoves = append(kingMoves, kingMoveSquare)

	}

	return kingMoves
}

func (cb *chessBoard) getPawnMoves(pawnColor pieceColor, pawnSquare vector2) []vector2 {
	var pawnMoves []vector2
	direction := 1
	if pawnColor == black {
		direction = -1
	}

	// Advance 1 space case
	advanceSquare := vector2{pawnSquare.x, pawnSquare.y + direction}
	if cb.inBoard(advanceSquare) && cb.isEmpty(advanceSquare) {
		pawnMoves = append(pawnMoves, advanceSquare)
	}

	// Move 2 spaces case
	if (pawnColor == white && pawnSquare.y == 1) ||
		(pawnColor == black && pawnSquare.y == 6) {
		doubleAdvanceSquare := vector2{pawnSquare.x, pawnSquare.y + (2 * direction)}
		if cb.isEmpty(advanceSquare) && cb.isEmpty(doubleAdvanceSquare) {
			pawnMoves = append(pawnMoves, doubleAdvanceSquare)
		}
	}

	// Take cases
	takeSquares := []vector2{{pawnSquare.x - 1, pawnSquare.y + direction}, {pawnSquare.x + 1, pawnSquare.y + direction}}
	for _, square := range takeSquares {
		if !cb.inBoard(square) {
			continue
		}
		if cb.getPiece(square).color == pawnColor.oppositeColor() {
			pawnMoves = append(pawnMoves, square)
		}
		if cb.enpassantSquare == square {
			pawnMoves = append(pawnMoves, square)
		}
	}

	return filter(pawnMoves, func(targetSquare vector2) bool { return !cb.moveInducesCheck(pawnSquare, targetSquare, pawnColor) })
}

func (cb *chessBoard) getPawnAttackedSquares(color pieceColor, pawnSquare vector2) []vector2 {
	var pawnAttackedSquares []vector2
	direction := 1
	if color == black {
		direction = -1
	}

	takeSquares := []vector2{{pawnSquare.x - 1, pawnSquare.y + direction}, {pawnSquare.x + 1, pawnSquare.y + direction}}
	for _, square := range takeSquares {
		if !cb.inBoard(square) {
			continue
		}
		if cb.getPiece(square).color == color {
			continue
		}
		pawnAttackedSquares = append(pawnAttackedSquares, square)
	}

	return pawnAttackedSquares
}

func (cb *chessBoard) getBishopAttackedSquares(bishopColor pieceColor, bishopSquare vector2) []vector2 {
	var bishopAttackedSquares []vector2

	bishopMoveVectors := [4]vector2{
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	for _, moveVector := range bishopMoveVectors {
		moveSquare := bishopSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == bishopColor {
				break
			}
			bishopAttackedSquares = append(bishopAttackedSquares, moveSquare)
			if cb.getPiece(moveSquare).color == bishopColor.oppositeColor() {
				break
			}
			moveSquare = moveSquare.add(moveVector)
		}
	}

	return bishopAttackedSquares
}

func (cb *chessBoard) getKnightAttackedSquares(knightColor pieceColor, knightSquare vector2) []vector2 {
	var knightAttackedSquares []vector2

	knightMoveVectors := []vector2{
		{2, 1}, {2, -1},
		{-2, 1}, {-2, -1},
		{1, 2}, {1, -2},
		{-1, 2}, {-1, -2},
	}

	for _, moveVector := range knightMoveVectors {
		moveSquare := knightSquare.add(moveVector)
		if !cb.inBoard(moveSquare) {
			continue
		}
		if cb.getPiece(moveSquare).color == knightColor {
			continue
		}
		knightAttackedSquares = append(knightAttackedSquares, moveSquare)
	}

	return knightAttackedSquares
}

func (cb *chessBoard) getQueenAttackedSquares(queenColor pieceColor, queenSquare vector2) []vector2 {
	var queenAttackedSquares []vector2

	queenMoveVectors := [8]vector2{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, // Rook-like moves
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Bishop-like moves
	}

	for _, moveVector := range queenMoveVectors {
		moveSquare := queenSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == queenColor {
				break
			}
			queenAttackedSquares = append(queenAttackedSquares, moveSquare)
			if cb.getPiece(moveSquare).color == queenColor.oppositeColor() {
				break
			}
			moveSquare = moveSquare.add(moveVector)
		}
	}

	return queenAttackedSquares
}

func (cb *chessBoard) getRookAttackedSquares(rookColor pieceColor, rookSquare vector2) []vector2 {
	var rookAttackedSquares []vector2

	rookMoveVectors := [4]vector2{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}

	for _, moveVector := range rookMoveVectors {
		moveSquare := rookSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == rookColor {
				break
			}
			rookAttackedSquares = append(rookAttackedSquares, moveSquare)
			if cb.getPiece(moveSquare).color == rookColor.oppositeColor() {
				break
			}
			moveSquare = moveSquare.add(moveVector)
		}
	}

	return rookAttackedSquares
}

func (cb *chessBoard) getBishopMoves(bishopColor pieceColor, bishopSquare vector2) []vector2 {
	bishopMoveVectors := [4]vector2{
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	var bishopMoves []vector2

	for _, moveVector := range bishopMoveVectors {
		moveSquare := bishopSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == bishopColor {
				break
			}
			if cb.getPiece(moveSquare).color == bishopColor.oppositeColor() {
				bishopMoves = append(bishopMoves, moveSquare)
				break
			}
			if cb.isEmpty(moveSquare) {
				bishopMoves = append(bishopMoves, moveSquare)
				moveSquare = moveSquare.add(moveVector)
				continue
			}
			panic("Unhandled Case In Bishop Movement Logic!")
		}
	}

	return filter(bishopMoves, func(targetSquare vector2) bool { return !cb.moveInducesCheck(bishopSquare, targetSquare, bishopColor) })

}

func (cb *chessBoard) getKnightMoves(knightColor pieceColor, knightSquare vector2) []vector2 {
	knightMoveVectors := []vector2{
		{2, 1}, {2, -1},
		{-2, 1}, {-2, -1},
		{1, 2}, {1, -2},
		{-1, 2}, {-1, -2},
	}

	var knightMoves []vector2
	for _, moveVector := range knightMoveVectors {
		moveSquare := knightSquare.add(moveVector)
		if !cb.inBoard(moveSquare) {
			continue
		}
		if cb.isEmpty(moveSquare) || cb.getPiece(moveSquare).color == knightColor.oppositeColor() {
			knightMoves = append(knightMoves, moveSquare)
		}
	}
	return filter(knightMoves, func(targetSquare vector2) bool { return !cb.moveInducesCheck(knightSquare, targetSquare, knightColor) })
}

func (cb *chessBoard) getRookMoves(rookColor pieceColor, rookSquare vector2) []vector2 {
	rookMoveVectors := []vector2{
		{0, 1}, {0, -1},
		{1, 0}, {-1, 0},
	}

	var rookMoves []vector2
	for _, moveVector := range rookMoveVectors {
		moveSquare := rookSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == rookColor {
				break
			}
			if cb.getPiece(moveSquare).color == rookColor.oppositeColor() {
				rookMoves = append(rookMoves, moveSquare)
				break
			}
			if cb.isEmpty(moveSquare) {
				rookMoves = append(rookMoves, moveSquare)
				moveSquare = moveSquare.add(moveVector)
				continue
			}
			break
		}
	}
	return filter(rookMoves, func(targetSquare vector2) bool { return !cb.moveInducesCheck(rookSquare, targetSquare, rookColor) })
}

func (cb *chessBoard) getQueenMoves(queenColor pieceColor, queenSquare vector2) []vector2 {
	var queenMoves []vector2

	queenMoveVectors := []vector2{
		{1, 0}, {-1, 0}, {0, 1}, {0, -1}, // Rook-like moves
		{1, 1}, {1, -1}, {-1, 1}, {-1, -1}, // Bishop-like moves
	}

	for _, moveVector := range queenMoveVectors {
		moveSquare := queenSquare.add(moveVector)
		for {
			if !cb.inBoard(moveSquare) {
				break
			}
			if cb.getPiece(moveSquare).color == queenColor {
				break
			}
			if cb.getPiece(moveSquare).color == queenColor.oppositeColor() {
				queenMoves = append(queenMoves, moveSquare)
				break
			}
			if cb.isEmpty(moveSquare) {
				queenMoves = append(queenMoves, moveSquare)
				moveSquare = moveSquare.add(moveVector)
				continue
			}
			break
		}
	}

	return filter(queenMoves, func(targetSquare vector2) bool { return !cb.moveInducesCheck(queenSquare, targetSquare, queenColor) })
}

func (cb *chessBoard) isValidMove(piece chessPiece, piecePosition vector2, newPosition vector2) bool {
	switch piece.pieceType {
	case pawn:
		return contains(cb.getPawnMoves(piece.color, piecePosition), newPosition)
	case bishop:
		return contains(cb.getBishopMoves(piece.color, piecePosition), newPosition)
	case rook:
		return contains(cb.getRookMoves(piece.color, piecePosition), newPosition)
	case knight:
		return contains(cb.getKnightMoves(piece.color, piecePosition), newPosition)
	case queen:
		return contains(cb.getQueenMoves(piece.color, piecePosition), newPosition)
	case king:
		return contains(cb.getKingMoves(piece.color, piecePosition), newPosition)
	default:
		return true
	}
}
