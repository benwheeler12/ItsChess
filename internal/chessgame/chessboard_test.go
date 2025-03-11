package chessgame

import "testing"

func TestPlayerInCheckMate(t *testing.T) {
	tests := []struct {
		name        string
		setupMoves  [][]vector2
		playerColor pieceColor
		expectMate  bool
		description string
	}{
		{
			name: "Fool's Mate",
			setupMoves: [][]vector2{
				{{5, 1}, {5, 2}}, // White f2-f3
				{{4, 6}, {4, 4}}, // Black e7-e5
				{{6, 1}, {6, 3}}, // White g2-g4
				{{3, 7}, {7, 3}}, // Black Qd8-h4#
			},
			playerColor: white,
			expectMate:  true,
			description: "Testing fool's mate checkmate position",
		},
		{
			name:        "Not In Checkmate - Starting Position",
			setupMoves:  [][]vector2{},
			playerColor: white,
			expectMate:  false,
			description: "Testing initial board position - should not be checkmate",
		},
		{
			name: "Not In Checkmate - King In Check But Can Move",
			setupMoves: [][]vector2{
				{{4, 1}, {4, 2}}, // White e2-e3
				{{3, 7}, {7, 3}}, // Black Qd8-h4+
			},
			playerColor: white,
			expectMate:  false,
			description: "King is in check but can escape",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new board
			var board chessBoard
			board.initializeBoard()

			// Execute setup moves
			for _, move := range tt.setupMoves {
				from, to := move[0], move[1]
				board.movePiece(from, to)
			}

			// Check if player is in checkmate
			gotCheckmate := board.playerInCheckMate(tt.playerColor)

			if gotCheckmate != tt.expectMate {
				t.Errorf("%s: playerInCheckMate() = %v, want %v",
					tt.description, gotCheckmate, tt.expectMate)
			}
		})
	}
}
