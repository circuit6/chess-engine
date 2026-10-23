package main

import (
	"fmt"
	"strings"
	"testing"
)

/*
Castling Rights:
- Normal King Move -> both rights gone
- Rook move from H1 to H2 -> one sided rights gone
- Actual castles -> all rights gone
- Rook getting captured -> one sided rights gone
EnPassant:
- Normal Pawn capture does not remove adjacent pawn
- A double move actually results in an EnPassant Field
*/

type boardFlags struct {
	Castling   string
	SideToMove bool
	HalfMove   int
	FullMove   int
	EnPassant  Square
}

func TestFlagParser(t *testing.T) {
	flagtests := []struct {
		name string
		in   string
		want boardFlags
	}{
		{name: "EnPassantTest", in: "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", want: boardFlags{Castling: "KQkq", SideToMove: false, HalfMove: 0, FullMove: 1, EnPassant: E3}},
		// Side To Move testing
		{name: "SideToMove", in: "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", want: boardFlags{Castling: "KQkq", SideToMove: false, HalfMove: 0, FullMove: 1, EnPassant: E3}},
		{name: "SideToMove+1", in: "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1 moves e7e6", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 0, FullMove: 2, EnPassant: 0}},
		{name: "EnPassantNormal", in: "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1 moves e7e5", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 0, FullMove: 2, EnPassant: E6}},
		// Full Move Half Move testing
		{name: "NoFullMove+1AfterWhite", in: "position startpos moves g1f3", want: boardFlags{Castling: "KQkq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "FullMoveUpAfterBlack", in: "position startpos moves g1f3 g8f6", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 2, FullMove: 2, EnPassant: 0}},
		{name: "PawnMoveResetsHalfMove", in: "position startpos moves g1f3 g8f6 e2e4", want: boardFlags{Castling: "KQkq", SideToMove: false, HalfMove: 0, FullMove: 2, EnPassant: E3}},
		{name: "CaptureResetsHalfMove", in: "position startpos moves g1f3 d7d5 f3e5 b8c6 e5c6", want: boardFlags{Castling: "KQkq", SideToMove: false, HalfMove: 0, FullMove: 3, EnPassant: 0}},
		// Long Game from Morphy
		{name: "MorphyMove1", in: "position startpos moves e2e4 e7e5", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 0, FullMove: 2, EnPassant: E6}},
		{name: "MorphyMove2", in: "position startpos moves e2e4 e7e5 g1f3 d7d6", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 0, FullMove: 3, EnPassant: 0}},
		{name: "MorphyMove3", in: "position startpos moves e2e4 e7e5 g1f3 d7d6 d2d4 c8g4 d4e5 g4f3 d1f3 d6e5 f1c4 g8f6 f3b3 d8e7 b1c3 c7c6 c1g5 b7b5", want: boardFlags{Castling: "KQkq", SideToMove: true, HalfMove: 0, FullMove: 10, EnPassant: B6}},
		{name: "MorphyMove4Castling", in: "position startpos moves e2e4 e7e5 g1f3 d7d6 d2d4 c8g4 d4e5 g4f3 d1f3 d6e5 f1c4 g8f6 f3b3 d8e7 b1c3 c7c6 c1g5 b7b5 c3b5 c6b5 c4b5 b8d7 e1c1 a8d8", want: boardFlags{Castling: "k", SideToMove: true, HalfMove: 3, FullMove: 13, EnPassant: 0}},
		{name: "MorphyMove5", in: "position startpos moves e2e4 e7e5 g1f3 d7d6 d2d4 c8g4 d4e5 g4f3 d1f3 d6e5 f1c4 g8f6 f3b3 d8e7 b1c3 c7c6 c1g5 b7b5 c3b5 c6b5 c4b5 b8d7 e1c1 a8d8 d1d7 d8d7 h1d1 e7e6 b5d7 f6d7 b3b8 d7b8 d1d8", want: boardFlags{Castling: "k", SideToMove: false, HalfMove: 1, FullMove: 17, EnPassant: 0}},
		// Testing Castling Rights Flags
		{name: "w_rook_h1_moves", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves h1h2", want: boardFlags{Castling: "Qkq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_rook_a1_moves", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves a1a2", want: boardFlags{Castling: "Kkq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_king_steps", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1e2", want: boardFlags{Castling: "kq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_king_steps", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1e2", want: boardFlags{Castling: "kq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_king_side_step", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1f1", want: boardFlags{Castling: "kq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_castles_short", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1g1", want: boardFlags{Castling: "kq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_castles_long", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves e1c1", want: boardFlags{Castling: "kq", SideToMove: false, HalfMove: 1, FullMove: 1, EnPassant: 0}},
		{name: "w_rook_takes_rook_a8", in: "position fen r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1 moves a1a8", want: boardFlags{Castling: "Kk", SideToMove: false, HalfMove: 0, FullMove: 1, EnPassant: 0}},
		{name: "w_bishop_takes_rook_h8", in: "position fen r3k2r/8/8/8/8/8/1B6/R3K2R w KQkq - 0 1 moves b2h8", want: boardFlags{Castling: "KQq", SideToMove: false, HalfMove: 0, FullMove: 1, EnPassant: 0}},
	}

	for _, tc := range flagtests {
		t.Run(tc.name, func(t *testing.T) {
			board := NewBoard()
			input := strings.Fields(tc.in)
			err := parsePosition(input, board)
			if err != nil {
				t.Errorf("failed to parse position to outBoard")
			}
			got := boardFlags{
				Castling:   board.Castling,
				SideToMove: board.SideToMove,
				HalfMove:   board.HalfMove,
				FullMove:   board.FullMove,
				EnPassant:  board.EnPassant,
			}
			out := tc.want
			if out != got {
				t.Errorf("got: %v, wanted: %v", got, out)
			}
		})
	}
}

func TestParsePositionBasic(t *testing.T) {
	board := NewBoard()
	input := strings.Fields("position startpos moves e2e4 e7e5 g1f3 b8c6")
	err := parsePosition(input, board)
	if err != nil {
		t.Errorf("failed parsePosition: %v", err)
	}
	check1 := board.squares[E4] == Piece{Type: Pawn, Color: White}
	check2 := board.squares[E2] == Piece{Type: Empty, Color: None}
	check3 := board.squares[E5] == Piece{Type: Pawn, Color: Black}
	check4 := board.squares[F3] == Piece{Type: Knight, Color: White}
	check5 := board.squares[C6] == Piece{Type: Knight, Color: Black}
	if !check1 || !check2 || !check3 || !check4 || !check5 {
		t.Errorf("failed basic position checks")
		fmt.Print(board)
	}
}

func TestParsePositionEnPassant(t *testing.T) {
	board := NewBoard()
	input := strings.Fields("position startpos moves e2e4 a7a6 e4e5 d7d5 e5d6")
	err := parsePosition(input, board)
	if err != nil {
		t.Errorf("failed parsePosition: %v", err)
	}
	check1 := board.squares[D6] == Piece{Type: Pawn, Color: White}
	check2 := board.squares[D5] == Piece{Type: Empty, Color: None}
	if !check1 || !check2 {
		t.Errorf("failed en passant rules check after d7d5 and e5d6")
		fmt.Print(board)
	}
}

func TestParsePositionCastling(t *testing.T) {
	board := NewBoard()
	input := strings.Fields("position startpos moves e2e4 e7e5 g1f3 b8c6 f1c4 g8f6 e1g1")
	err := parsePosition(input, board)
	if err != nil {
		t.Errorf("failed parsePosition: %v", err)
	}
	check1 := board.squares[H1] == Piece{Type: Empty, Color: None}
	check2 := board.squares[F1] == Piece{Type: Rook, Color: White}
	check3 := board.squares[G1] == Piece{Type: King, Color: White}
	if !check1 || !check2 || !check3 {
		t.Errorf("failed castling white short e1g1")
		fmt.Print(board)
	}
}

func TestParsePositionPromotion(t *testing.T) {
	board := NewBoard()
	input := strings.Fields("position startpos moves a2a4 b7b5 a4b5 a7a6 b5a6 g8f6 a6a7 e7e6 a7b8q")
	err := parsePosition(input, board)
	if err != nil {
		t.Errorf("failed parsePosition: %v", err)
	}
	check1 := board.squares[B8] == Piece{Type: Queen, Color: White}
	check2 := board.squares[F6] == Piece{Type: Knight, Color: Black}
	if !check1 || !check2 {
		t.Errorf("failed promotion to B8")
		fmt.Print(board)
	}
}
