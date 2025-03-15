package chessgame

import "github.com/hajimehoshi/ebiten/v2"

type piece int

type pieceColor int

var emptyPiece = chessPiece{empty, nocolor}
var nilSquare = vector2{-1, -1}

const (
	nocolor pieceColor = iota
	white
	black
)

const (
	empty piece = iota
	pawn
	bishop
	knight
	rook
	queen
	king
)

type chessPiece struct {
	piece piece
	color pieceColor
}

type castlingState struct {
	a1RookMoved    bool
	a8RookMoved    bool
	h1RookMoved    bool
	h8RookMoved    bool
	blackKingMoved bool
	whiteKingMoved bool
}

type chessBoard struct {
	board           [8][8]chessPiece
	enpassantSquare vector2
	castlingState
}

type vector2 struct {
	x int
	y int
}

type mouseLifeCycle struct {
	selectedPiece               chessPiece
	selectedSquare              vector2
	mouseClickedOnInvalidSquare bool
	possibleMoveSquares         []vector2 //representing possible moves from currently selected piece
}

type mouseState struct {
	mousePressed   bool
	lastMouseState bool
}

type graphicsState struct {
	windowWidth    int
	windowHeight   int
	pieceEbitenMap map[chessPiece]*ebiten.Image
}

type ChessGame struct {
	chessBoard
	mouseState
	mouseLifeCycle
	graphicsState
	whitesTurn bool
}
