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

func TestBishopMovement(t *testing.T) {
	tests := []struct {
		name           string
		setupMoves     [][]vector2
		bishopPosition vector2
		expectedMoves  []vector2
		description    string
	}{
		{
			name:           "Bishop in starting position - blocked by pawns",
			setupMoves:     [][]vector2{},
			bishopPosition: vector2{2, 0}, // White bishop at c1
			expectedMoves:  []vector2{},   // No valid moves as it's blocked by pawns
			description:    "Bishop in starting position should have no valid moves",
		},
		{
			name: "Bishop with open diagonals",
			setupMoves: [][]vector2{
				{{3, 1}, {3, 2}}, // Move pawn out of the way
			},
			bishopPosition: vector2{2, 0}, // White bishop at c1
			expectedMoves: []vector2{
				{3, 1}, // d2
				{4, 2}, // e3
				{5, 3}, // f4
				{6, 4}, // g5
				{7, 5}, // h6
			},
			description: "Bishop with open diagonals should have multiple valid moves",
		},
		{
			name: "Bishop with capture opportunity",
			setupMoves: [][]vector2{
				{{2, 1}, {2, 3}}, // Move pawn out of the way
				{{2, 0}, {5, 3}}, // Move bishop to f4
				{{6, 6}, {6, 4}}, // Move black pawn to g5
			},
			bishopPosition: vector2{5, 3}, // White bishop at f4
			expectedMoves: []vector2{
				{2, 0}, // c1
				{3, 1}, // d2
				{4, 2}, // e3
				{6, 4}, // g5 - can capture black pawn
				{7, 5}, // h6
				{4, 4}, // e5
				{3, 5}, // d6
				{2, 6}, // c7
				{6, 2}, // g3
				{7, 1}, // h2
			},
			description: "Bishop should be able to capture opponent's pieces",
		},
		{
			name: "Bishop blocked by friendly pieces",
			setupMoves: [][]vector2{
				{{2, 1}, {2, 3}}, // Move pawn out of the way
				{{2, 0}, {4, 2}}, // Move bishop to e3
				{{3, 1}, {3, 3}}, // Move pawn to e4
			},
			bishopPosition: vector2{4, 2}, // White bishop at e3
			expectedMoves: []vector2{
				{2, 0}, // c1
				{3, 1}, // d2
				{5, 3}, // f4
				{6, 4}, // g5
				{7, 5}, // h6
				{3, 3}, // d4
				{2, 4}, // c5
				{1, 5}, // b6
				{0, 6}, // a7
				{5, 1}, // f2
				{6, 0}, // g1
			},
			description: "Bishop should not be able to move through friendly pieces",
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

			// Get bishop's valid moves
			bishopPiece := board.getPiece(tt.bishopPosition)
			gotMoves := board.getValidMoves(tt.bishopPosition, bishopPiece)

			// Check if the number of moves matches
			if len(gotMoves) != len(tt.expectedMoves) {
				t.Errorf("%s: got %d moves, want %d moves",
					tt.description, len(gotMoves), len(tt.expectedMoves))
			}

			// Check if all expected moves are in the got moves
			for _, expectedMove := range tt.expectedMoves {
				found := false
				for _, gotMove := range gotMoves {
					if gotMove.equals(expectedMove) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: expected move %v not found in valid moves",
						tt.description, expectedMove)
				}
			}

			// Check if there are any unexpected moves
			for _, gotMove := range gotMoves {
				found := false
				for _, expectedMove := range tt.expectedMoves {
					if gotMove.equals(expectedMove) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("%s: got unexpected move %v",
						tt.description, gotMove)
				}
			}
		})
	}
}
