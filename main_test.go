package main

import (
	"fmt"
	"strings"
	"testing"
)

var flagtests = []struct {
	name string
	in   string
	out  Board
}{
	{name: "EnPassantTest", in: "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1", out: Board{EnPassant: E3}},
}

func TestFlagParser(t *testing.T) {
	for _, tc := range flagtests {
		t.Run(tc.name, func(t *testing.T) {
			board := NewBoard()
			input := strings.Fields(tc.in)
			err := parsePosition(input, board)
			if err != nil {
				t.Errorf("failed to parse position to outBoard")
			}
			out := tc.out
			if out.EnPassant != board.EnPassant {
				t.Errorf("got: %v, wanted: %v", board.EnPassant, out.EnPassant)
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
