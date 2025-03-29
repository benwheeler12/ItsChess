package itschess

type piece int

type pieceColor int

var emptyPiece = chessPiece{empty, nocolor}
var nilSquare = vector2{-1, -1}

const rotationFactor int = 2

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
	pieceType piece
	color     pieceColor
}

type castlingState struct {
	a1RookMoved    bool
	a8RookMoved    bool
	h1RookMoved    bool
	h8RookMoved    bool
	blackKingMoved bool
	whiteKingMoved bool
}

type vector2 struct {
	x int
	y int
}

type point struct {
	x float64
	y float64
}

type mouseLifeCycle struct {
	selectedPiece               chessPiece
	selectedSquare              vector2
	mouseClickedOnInvalidSquare bool
	possibleMoveSquares         []vector2 //representing possible moves from currently selected piece
}

type promotionLifeCycle struct {
	promotionSquare     vector2
	promotionInProgress bool
}

type mouseState struct {
	mousePressed   bool
	lastMouseState bool
}

// Initializers:
func (ms *mouseLifeCycle) resetMouseState() {
	ms.selectedPiece = emptyPiece
	ms.mouseClickedOnInvalidSquare = false
	ms.selectedSquare = nilSquare
	ms.possibleMoveSquares = nil
}

func (pfc *promotionLifeCycle) resetPromotionLifeCycle() {
	pfc.promotionInProgress = false
	pfc.promotionSquare = nilSquare
}
