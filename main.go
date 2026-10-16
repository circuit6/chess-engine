package main

import (
	"bufio"
	"fmt"
	"os"
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
}

var board Board

func printBoard(board Board) {
	// Uppercase = white, lowercase = black
	// 2 loops, first for rows, second for columns.
	// rows has to count backwards from 8 to 1 to go from top to bottom
	convert := func(row int, col int) int {
		row = row - 1
		col = col - 1
		return row*8 + col
	}

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

	white := wking | wqueen | wbish | wrook | wknight | wpawns
	fmt.Printf("white: %064b\n", white)
	black := bking | bqueen | bbish | brook | bknight | bpawns
	fmt.Printf("black: %064b\n", black)

	board.W = white
	board.B = black

	allboard := white | black
	board.Board = allboard

	return board
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

MainLoop:
	for scanner.Scan() {

		line := scanner.Text()
		input := strings.Fields(line)

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
			board = setupBoard()
			printBoard(board)
		case "position":
			// position [fenstring | startpos] moves
			fmt.Fprintln(os.Stderr, "not implemented yet")
			// FEN documentation and examples: https://www.chess.com/terms/fen-chess
		case "d":
			// debug, should be used like: position fen .... -> d -> check if position is any good
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
