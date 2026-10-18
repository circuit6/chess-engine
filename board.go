package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Color int

const (
	None Color = iota
	White
	Black
)

func (c Color) String() string {
	switch c {
	case None:
		return "none"
	case White:
		return "white"
	case Black:
		return "black"
	default:
		return "unknown"
	}
}

type PieceType int

const (
	Empty PieceType = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

func (p PieceType) String() string {
	switch p {
	case Empty:
		return "empty"
	case Pawn:
		return "pawn"
	case Rook:
		return "rook"
	case Knight:
		return "knight"
	case Bishop:
		return "bishop"
	case Queen:
		return "queen"
	case King:
		return "king"
	default:
		return "unkown"
	}
}

type Piece struct {
	Type  PieceType
	Color Color
}

func (p Piece) isEmpty() bool {
	return p.Type == Empty
}

func (p Piece) String() string {
	if p.isEmpty() {
		return "."
	}

	symbol := ""
	switch p.Type {
	case Pawn:
		symbol = "P"
	case Rook:
		symbol = "R"
	case Knight:
		symbol = "N"
	case Bishop:
		symbol = "B"
	case Queen:
		symbol = "Q"
	case King:
		symbol = "K"
	}
	if p.Color == Black {
		symbol = strings.ToLower(symbol)
	}
	return symbol
}

type Square int

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// methods: String() -> returns the algebraic notation
// Square from String -> "e4" -> number matching

func SquareFromString(s string) (Square, error) {
	file := s[0] - 'a' // s[0] is a byte (uint8 representation) and the other value is an untyped constant. It adopts to the other value in that line - hence we are subtracting byte from byte to derive the outputs.
	rank := s[1] - '1' // same applies here, also single quotes '' are runes -> runes would be uint32, however, as explained above, adopted top mechanism and treated as byte

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return 0, fmt.Errorf("invalid algebraic notation")
	}
	return Square(int(rank)*8 + int(file)), nil
}

func SquareFromIndex(i int, j int) (Square, error) {
	// find square from indizes 0 to 7
	if i < 0 || i > 7 || j < 0 || j > 7 {
		return 0, fmt.Errorf("invalid index input, outside of 0 to 7")
	}
	file := i * 8
	squareIndex := file + j
	sq := Square(squareIndex)
	return sq, nil
}

func (s Square) String() string {
	if s < 0 || s > 63 {
		return "invalid"
	}
	file := s % 8
	row := s / 8
	return fmt.Sprintf("%c%c", 'a'+file, '1'+row)
}

type Board struct {
	squares    [64]Piece
	SideToMove bool
	Castling   string
	EnPassant  string
	HalfMove   int
	FullMove   int
}

func (b *Board) makeMove(move string) error {
	if len(move) == 4 {
		origin := move[:2]
		target := move[2:]
		sqOrigin, err := SquareFromString(origin)
		if err != nil {
			return fmt.Errorf("SquareFromString failed on origin: %w", err)
		}
		pieceOrigin := b.squares[sqOrigin]
		sqTarget, err := SquareFromString(target)
		if err != nil {
			return fmt.Errorf("SquareFromString failed on target: %w", err)
		}
		// move to new field (even when capturing)
		b.squares[sqTarget].Type = pieceOrigin.Type
		b.squares[sqTarget].Color = pieceOrigin.Color
		// reset square of origin
		b.squares[sqOrigin].Type = Empty
		b.squares[sqOrigin].Color = None
		return nil
	}
	return nil
}

func (b *Board) parseFen(fenstring string) error {
	fenstringRows := strings.Split(fenstring, "/")
	for i, rank := range fenstringRows {
		// iterate over each char to evaluate:
		// rank index 0 is row 8
		file := 0
		for j, char := range rank {
			if _, err := strconv.Atoi(string(char)); err == nil {
				// is number, gotta move the index by number
				num, _ := strconv.Atoi(string(char))
				file = file + num
			}
			// letter hit, no number found
			// match letter to piece and save on board
			i = 7 - i
			sq, err := SquareFromIndex(i, j)
			if err != nil {
				return fmt.Errorf("converting index to square failed: %w", err)
			}
			piece := b.squares[sq]
			switch char {
			case 'p':
				piece.Type = Pawn
				piece.Color = Black
			case 'r':
				piece.Type = Rook
				piece.Color = Black
			case 'n':
				piece.Type = Knight
				piece.Color = Black
			case 'b':
				piece.Type = Bishop
				piece.Color = Black
			case 'q':
				piece.Type = Queen
				piece.Color = Black
			case 'k':
				piece.Type = King
				piece.Color = Black
			case 'P':
				piece.Type = Pawn
				piece.Color = White
			case 'R':
				piece.Type = Rook
				piece.Color = White
			case 'N':
				piece.Type = Knight
				piece.Color = White
			case 'B':
				piece.Type = Bishop
				piece.Color = White
			case 'Q':
				piece.Type = Queen
				piece.Color = White
			case 'K':
				piece.Type = King
				piece.Color = White
			}
			// increase file tracker
			file++
		}
	}
	return nil
}

func (b *Board) SetupStartingPos() {
	// clean up Board
	for i := range b.squares {
		b.squares[i] = Piece{Empty, None}
	}
	// Place white colors
	b.squares[A1] = Piece{Type: Rook, Color: White}
	b.squares[B1] = Piece{Type: Knight, Color: White}
	b.squares[C1] = Piece{Type: Bishop, Color: White}
	b.squares[D1] = Piece{Type: Queen, Color: White}
	b.squares[E1] = Piece{Type: King, Color: White}
	b.squares[F1] = Piece{Type: Bishop, Color: White}
	b.squares[G1] = Piece{Type: Knight, Color: White}
	b.squares[H1] = Piece{Type: Rook, Color: White}
	// Place black colors
	b.squares[A8] = Piece{Type: Rook, Color: Black}
	b.squares[B8] = Piece{Type: Knight, Color: Black}
	b.squares[C8] = Piece{Type: Bishop, Color: Black}
	b.squares[D8] = Piece{Type: Queen, Color: Black}
	b.squares[E8] = Piece{Type: King, Color: Black}
	b.squares[F8] = Piece{Type: Bishop, Color: Black}
	b.squares[G8] = Piece{Type: Knight, Color: Black}
	b.squares[H8] = Piece{Type: Rook, Color: Black}
	// place white pawns
	for i := A2; i <= H2; i++ {
		b.squares[i] = Piece{Type: Pawn, Color: White}
	}
	for j := A7; j <= H7; j++ {
		b.squares[j] = Piece{Type: Pawn, Color: Black}
	}
}

func NewBoard() *Board {
	b := &Board{}
	b.SetupStartingPos()
	return b
}

func (b *Board) GetPiece(sq Square) Piece {
	if sq < 0 || sq > 63 {
		return Piece{Type: Empty, Color: None}
	}
	return b.squares[sq]
}

func (b *Board) SetPiece(sq Square, piece Piece) {
	if sq < 0 || sq > 63 {
		return
	}
	b.squares[sq] = piece
}

func (b *Board) String() string {
	// iterate over ranks from top to bottom
	boardString := ""
	rowString := ""
	for i := 7; i >= 0; i-- {
		// iterate over files from left to right
		var bunchPieces []Piece
		rowString = ""
		for j := 0; j <= 7; j++ {
			// get bitpos
			bitpos := i*8 + j
			piece := b.squares[bitpos]
			bunchPieces = append(bunchPieces, piece)
		}
		for _, v := range bunchPieces {
			ps := fmt.Sprintf("%v", v)
			rowString = rowString + ps
		}
		boardString = boardString + rowString + "\n"
	}
	return boardString
}
