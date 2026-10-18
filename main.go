package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Board struct {
	WP         uint64
	WR         uint64
	WB         uint64
	WN         uint64 // Knight
	WQ         uint64
	WK         uint64 // King
	BP         uint64
	BR         uint64
	BB         uint64
	BN         uint64
	BQ         uint64
	BK         uint64
	W          uint64 // White
	B          uint64 // black
	Board      uint64 // all pieces on the board
	SideToMove bool   // true is white
	Castling   string // castling rights
	EnPassant  string // En Passant Option Square
	HalfMove   int    // how many moves both players have made since last pawn advance or piece capture (50 move draw rule)
	FullMove   int    // number of completered turns this game. incremented every time black has moved.
}

func convertSquare(square string) int {
	cols := make(map[string]int)
	cols["a"] = 0
	cols["b"] = 1
	cols["c"] = 2
	cols["d"] = 3
	cols["e"] = 4
	cols["f"] = 5
	cols["g"] = 6
	cols["h"] = 7
	file := cols[string(square[0])]
	row, err := strconv.Atoi(string(square[1]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not convert number of square to int")
	}
	row = row - 1
	bitPos := convert0To7(row, file)
	return bitPos
}

func (b *Board) generateBoards() []*uint64 {
	boards := []*uint64{&b.BP, &b.WP, &b.BK, &b.BN, &b.BB, &b.BQ, &b.BR, &b.WB, &b.WK, &b.WN, &b.WQ, &b.WR}
	return boards
}

func (b *Board) updateGlobals() {
	b.W = b.WK | b.WQ | b.WB | b.WR | b.WN | b.WP
	b.B = b.BK | b.BQ | b.BB | b.BR | b.BN | b.BP
	b.Board = b.W | b.B
}

func (b *Board) makeMove(move string) error {
	// compute a move and update the board
	// assumption: GUI sends only legal move, hence: no validation required
	/*
		Normal Case:
		- e2e4
		- length of 4
		- if square already occupied, gotta delete other piece (was taken)
		Edge Cases:
		- promotion to queen: e7e8q
		- castling: e1g1 -> rook also has to teleport
		- en passant: e5d5, switching columns and captured pawn is not on the target square!
	*/
	//split move input into 'origin' and 'target'
	if len(move) == 4 {
		origin := move[:2]
		target := move[2:]
		// convert 'e2' to bitpos
		bitOrigin := convertSquare(origin)
		bitTarget := convertSquare(target)
		var boardOrigin uint64 = 1 << bitOrigin
		var boardTarget uint64 = 1 << bitTarget
		// loop over 'piece boards' to 'AND' check which hits
		boards := b.generateBoards()
		// check for capture before updating board, otherwise we 'find' the newly set move-board in next for loop
		// TODO: En Passant Rule is not implemented here
		for _, tboard := range boards {
			if *tboard&boardTarget != 0 {
				*tboard &^= boardTarget
				break
			}
		}
		for _, board := range boards {
			// we need to use *board to get from the pointer to the actual data stored in the board
			// the if reads as: if board AND origin Board have overlap (!=0; no overlap, means the AND is 0)
			if *board&boardOrigin != 0 {
				// the &^ is the AND with the flipped bit of boardOrigin.
				// if boardOrigin was 0001 0000, we now have 1110 1111.
				// we then perform the AND to clear out exactly that one spot, where the '0' sits in the flipped int
				// *board |= is the same as *board = *board | boardTarget
				*board &^= boardOrigin
				*board |= boardTarget
				break
			}
		}
		b.updateGlobals()
	}
	return nil
}

func convert(row int, col int) int {
	// converts row and col index to bitboard pos
	// expects 1 to 8 numerical
	row = row - 1
	col = col - 1
	return row*8 + col
}

func convert0To7(row int, col int) int {
	// converts row and col index to bitboard pos
	// expects index 0 to 7
	return row*8 + col
}

func printBoard(board Board) {
	// Uppercase = white, lowercase = black
	// 2 loops, first for rows, second for columns.
	// rows has to count backwards from 8 to 1 to go from top to bottom
	for i := 8; i >= 1; i-- {
		// Outer loop going from top to bottom!
		for j := 1; j <= 8; j++ {
			bitIndex := convert(i, j)
			var lookup uint64 = 1 << bitIndex

			if lookup&board.Board == 0 {
				fmt.Print(". ")
			}
			if lookup&board.B != 0 {
				// check black pieces
				if lookup&board.BP != 0 {
					fmt.Print("p ")
				}
				if lookup&board.BN != 0 {
					fmt.Print("n ")
				}
				if lookup&board.BB != 0 {
					fmt.Print("b ")
				}
				if lookup&board.BR != 0 {
					fmt.Print("r ")
				}
				if lookup&board.BQ != 0 {
					fmt.Print("q ")
				}
				if lookup&board.BK != 0 {
					fmt.Print("k ")
				}
			}
			if lookup&board.W != 0 {
				// check white pieces
				if lookup&board.WP != 0 {
					fmt.Print("P ")
				}
				if lookup&board.WN != 0 {
					fmt.Print("N ")
				}
				if lookup&board.WB != 0 {
					fmt.Print("B ")
				}
				if lookup&board.WR != 0 {
					fmt.Print("R ")
				}
				if lookup&board.WQ != 0 {
					fmt.Print("Q ")
				}
				if lookup&board.WK != 0 {
					fmt.Print("K ")
				}
			}
		}
		fmt.Println()
	}
}

func setupBoard() Board {
	var board Board

	board.SideToMove = true // White = true

	// set white pawns
	var wpawns uint64
	for i := 8; i <= 15; i++ {
		var wpawnsNew uint64 = 1 << i
		wpawns = wpawns | wpawnsNew
	}
	board.WP = wpawns

	// set black pawns
	var bpawns uint64
	for i := 48; i <= 55; i++ {
		var bpawnsNew uint64 = 1 << i
		bpawns = bpawns | bpawnsNew
	}
	board.BP = bpawns

	// set white rooks
	var wrook uint64 = 1 << 0
	var wrook2 uint64 = 1 << 7

	wrook = wrook | wrook2
	board.WR = wrook

	// set black rooks
	var brook uint64 = 1 << 56
	var brook2 uint64 = 1 << 63

	brook = brook | brook2
	board.BR = brook

	// set white knights
	var wknight uint64 = 1 << 1
	var wknight2 uint64 = 1 << 6

	wknight = wknight | wknight2
	board.WN = wknight

	// set black knights
	var bknight uint64 = 1 << 62
	var bknight2 uint64 = 1 << 57

	bknight = bknight | bknight2
	board.BN = bknight

	// set white bishop
	var wbish uint64 = 1 << 2
	var wbish2 uint64 = 1 << 5
	wbish = wbish | wbish2
	board.WB = wbish

	// set black bishop
	var bbish uint64 = 1 << 61
	var bbish2 uint64 = 1 << 58
	bbish = bbish | bbish2
	board.BB = bbish

	// set white Queen
	var wqueen uint64 = 1 << 3
	board.WQ = wqueen

	// set black Queen
	var bqueen uint64 = 1 << 59
	board.BQ = bqueen

	// set white King
	var wking uint64 = 1 << 4
	board.WK = wking

	// set black King
	var bking uint64 = 1 << 60
	board.BK = bking

	board.updateGlobals()
	return board
}

func findBitPos(row int, col int) int {
	row = 7 - row
	pos := convert0To7(row, col)
	// fmt.Fprintf(os.Stderr, "findBitPos triggered with values: row %v, col %v, pos %v", row, col, pos)
	return pos
}

func parseFen(fenstring string) (Board, error) {
	fenstringRows := strings.Split(fenstring, "/")
	var board Board
	for i, row := range fenstringRows {
		// iterate over each char to evaluate:
		// row index 0 is Row 8
		columnTracker := 0
		for _, char := range row {
			if _, err := strconv.Atoi(string(char)); err == nil {
				// is number, gotta move the index by number
				num, _ := strconv.Atoi(string(char))
				columnTracker = columnTracker + num
				continue
			}
			// Letter hit, no number found
			// match letter to piece and save on Board
			bitIndex := findBitPos(i, columnTracker)
			switch char {
			case 'p':
				// add black pawn to board
				var bpawn uint64 = 1 << bitIndex
				board.BP = board.BP | bpawn
			case 'P':
				// add white pawn
				var wpawn uint64 = 1 << bitIndex
				board.WP = board.WP | wpawn
			case 'r':
				// add black rooks
				var brook uint64 = 1 << bitIndex
				board.BR = board.BR | brook
			case 'R':
				// add white rooks
				var wrook uint64 = 1 << bitIndex
				board.WR = board.WR | wrook
			case 'n':
				// add black knights
				var bnight uint64 = 1 << bitIndex
				board.BN = board.BN | bnight
			case 'N':
				// add white knights
				var wnight uint64 = 1 << bitIndex
				board.WN = board.WN | wnight
			case 'b':
				// add black bishops
				var bbish uint64 = 1 << bitIndex
				board.BB = board.BB | bbish
			case 'B':
				// add white bishops
				var wbish uint64 = 1 << bitIndex
				board.WB = board.WB | wbish
			case 'q':
				// add black Queen
				var bqueen uint64 = 1 << bitIndex
				board.BQ = board.BQ | bqueen
			case 'Q':
				// add white Queen
				var wqueen uint64 = 1 << bitIndex
				board.WQ = board.WQ | wqueen
			case 'k':
				// add black King
				var bking uint64 = 1 << bitIndex
				board.BK = board.BK | bking
			case 'K':
				// add white King
				var wking uint64 = 1 << bitIndex
				board.WK = board.WK | wking
			}
			// increment columnTracker by 1
			columnTracker++
		}
	}
	// build board, white and black
	board.W = board.WP | board.WN | board.WB | board.WR | board.WQ | board.WK
	board.B = board.BP | board.BN | board.BB | board.BR | board.BQ | board.BK
	board.Board = board.W | board.B
	return board, nil
}

func parsePosition(input []string) (Board, error) {
	// gotta be able to read the fenstring to a board,
	// determine whos turn it is
	// for now: no validation, if moves are legal

	var outBoard Board
	var err error

	if len(input) < 4 {
		// all FEN inputs shorter than 4 are invalid
		return Board{}, errors.New("invalid FEN input, shorter than 4 args")
	}

	position := input[0]
	if position != "position" {
		return Board{}, errors.New("keyword 'position' triggered, but not found on pos 0 in parsePosition")
	}

	// we get a "fen" or "startpos": fen is more likely studying / debugging, startpos is normal case
	initPos := input[1]
	// startpos scenario
	if initPos == "startpos" {
		// typical input looks like: position startpos moves e2e4 e7e5 g1f3 b8c6 ... - rest is implicit!
		outBoard = setupBoard()
		outBoard.EnPassant = "-"
		outBoard.SideToMove = true
		outBoard.Castling = "KQkq"
		outBoard.HalfMove = 0
		outBoard.FullMove = 1
		// moves parsing
		// validate 'moves' keyword
		movesKeyword := input[2]
		if movesKeyword != "moves" {
			return Board{}, errors.New("unexpected keyword: expected 'moves' on index 2")
		}
		moves := input[3:]

		if len(moves) == 0 {
			return Board{}, errors.New("no moves provided after 'startpos' keyword")
		}
		// process moves
		for _, move := range moves {
			outBoard.makeMove(move)
		}
		return outBoard, nil
	}
	// startpos scenario

	if initPos == "fen" {
		fmt.Fprintf(os.Stderr, "parsing input 2: %v\n", input[2])
		outBoard, err = parseFen(input[2])
		if err != nil {
			return Board{}, err
		}
	}

	// side to move next:
	sideMove := input[3]
	if sideMove == "w" {
		// white
		outBoard.SideToMove = true
	} else {
		// black
		outBoard.SideToMove = false
	}

	// castling rights
	// Qk = white can castle Queen side, black can castle King side. "-", no castling options.
	outBoard.Castling = input[4]

	// en passant indicator
	if len(input) <= 4 {
		return Board{}, errors.New("could not find element for en passant indicator in FEN")
	}
	outBoard.EnPassant = input[5]

	// halfmove clock
	if len(input) <= 5 {
		return Board{}, errors.New("could not parse halfmove clock, input too short")
	}
	val, err := strconv.Atoi(input[6])
	if err != nil {
		return Board{}, errors.New("could not parse input 6 as int")
	}
	outBoard.HalfMove = val

	// Fullmove number
	if len(input) <= 6 {
		return Board{}, errors.New("could not parse fullmove clock")
	}
	val, err = strconv.Atoi(input[7])
	if err != nil {
		return Board{}, errors.New("could not parse input 7 as int to Fullmove")
	}
	outBoard.FullMove = val

	// moves
	if len(input) <= 7 {
		// no moves
		return outBoard, nil
	}
	movesKeyword := input[8]
	if movesKeyword != "moves" {
		return Board{}, errors.New("keyword on index 8 was not moves, expected 'moves'")
	}

	// parse moves actually without validation

	return outBoard, nil
}

func main() {
	// CLI based program which speaks UCI
	// Input: Bufio.Scanner reads the os.Stdin
	// Output: fmt.Println to stdout is the default output
	// People usually include a print the board in ascii
	// there is a "perft" command typically included, which counts all legal positions to depth N and compares against known numbers
	// Add custom commands like "-d" to display the board
	// GUI testing: cutechess-cli, arena
	// commands to get started with: uci, isready, ucinewgame, position, go, quit

	/*
		UCI protocoll:
		os.Stdin() and os.stdout() are only ways of input and output.
		- The engine must always be able to process input from stdin, even while thinking
		- all commands end with \n
		- setting up internal parameters of the engine: after "isready" or "setoption" command
		- fail silently

		"-uci":
		tells the engine to use the uci format. First command after program boot.
		should identify itself with the "id" command
		should send "uciok" to acknowledge the uci mode.

		"-isready":
		acts like a heartbeat check
		must always be answered with "readyok"
		can be sent whilst machine is calculating

		"-ucinewgame":
		the next search will be from a different game
		can be used whilst testing -> when tests go through different positions

		"-position [fen <fenstring> | startpos] moves <move1> .... <movei>":
		sets up the position as per the fenstring
		plays the moves from fenstring on internal board

		"-go":
		start calculating the next move
		following commands are allowed, all of which can sent in the same string:
		"-searchmoves"
		"-ponder"
		"-wtime"
		"-btime"
		"-winc"
		"-binc"
		"-movestogo"
		"-depth"
		"-nodes"
		"-mate"
		"-movetime"
		"-infinite"

		"-quit"
		Quit the program as soon as possible

	*/

	// Loop over os.Stdin lines

	scanner := bufio.NewScanner(os.Stdin)
	var currentBoard Board

MainLoop:
	for scanner.Scan() {

		line := scanner.Text()
		input := strings.Fields(line)
		var err error

		if len(input) == 0 {
			// avoid panic when no input
			continue
		}

		firstArg := input[0]

		// keywords matching UCI
		switch firstArg {
		case "uci":
			// first command after program boot
			fmt.Fprintf(os.Stdout, "id name SuperBot\n")
			fmt.Fprintf(os.Stdout, "id author Timothy\n")
			fmt.Fprintf(os.Stdout, "uciok\n")
		case "isready":
			// heartbeat with answer "readyok"
			fmt.Fprintf(os.Stdout, "readyok\n")
		case "ucinewgame":
			// for now just board setup and initialization
			board := setupBoard()
			printBoard(board)
		case "position":
			// position [fenstring | startpos] moves
			currentBoard, err = parsePosition(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsePos: %v", err)
			}
			// FEN documentation and examples: https://www.chess.com/terms/fen-chess
		case "d":
			// debug, should be used like: position fen .... -> d -> check if position is any good
			printBoard(currentBoard)
		case "go":
			// always returns the best move found: bestmove <move>
			fmt.Fprintln(os.Stderr, "not implemented yet")
		case "stop":
			// stop calculating as soon as possible
			// does return best move
			fmt.Fprintln(os.Stderr, "not implemented yet")
		case "quit":
			break MainLoop
		}
	}
}
