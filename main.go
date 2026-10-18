package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parsePosition(input []string, outBoard *Board) error {
	// gotta be able to read the fenstring to a board,
	// determine whos turn it is
	// for now: no validation, if moves are legal

	var err error

	if len(input) < 4 {
		// all FEN inputs shorter than 4 are invalid
		return errors.New("invalid FEN input, shorter than 4 args")
	}

	position := input[0]
	if position != "position" {
		return errors.New("keyword 'position' triggered, but not found on pos 0 in parsePosition")
	}

	// we get a "fen" or "startpos": fen is more likely studying / debugging, startpos is normal case
	initPos := input[1]
	// startpos scenario
	if initPos == "startpos" {
		// typical input looks like: position startpos moves e2e4 e7e5 g1f3 b8c6 ... - rest is implicit!
		outBoard.EnPassant = "-"
		outBoard.SideToMove = true
		outBoard.Castling = "KQkq"
		outBoard.HalfMove = 0
		outBoard.FullMove = 1
		// moves parsing
		// validate 'moves' keyword
		movesKeyword := input[2]
		if movesKeyword != "moves" {
			return errors.New("unexpected keyword: expected 'moves' on index 2")
		}
		moves := input[3:]

		if len(moves) == 0 {
			return errors.New("no moves provided after 'startpos' keyword")
		}
		// process moves
		for _, move := range moves {
			outBoard.makeMove(move)
		}
		return nil
	}
	// fen scenario
	if initPos == "fen" {
		fmt.Fprintf(os.Stderr, "parsing input 2: %v\n", input[2])
		err := outBoard.parseFen(input[2])
		if err != nil {
			return err
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
		return errors.New("could not find element for en passant indicator in FEN")
	}
	outBoard.EnPassant = input[5]

	// halfmove clock
	if len(input) <= 5 {
		return errors.New("could not parse halfmove clock, input too short")
	}
	val, err := strconv.Atoi(input[6])
	if err != nil {
		return errors.New("could not parse input 6 as int")
	}
	outBoard.HalfMove = val

	// Fullmove number
	if len(input) <= 6 {
		return errors.New("could not parse fullmove clock")
	}
	val, err = strconv.Atoi(input[7])
	if err != nil {
		return errors.New("could not parse input 7 as int to Fullmove")
	}
	outBoard.FullMove = val

	// moves
	if len(input) <= 7 {
		// no moves
		return nil
	}
	movesKeyword := input[8]
	if movesKeyword != "moves" {
		return errors.New("keyword on index 8 was not moves, expected 'moves'")
	}

	// parse moves actually without validation

	return nil
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
	var currentBoard *Board

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
			currentBoard = NewBoard()
		case "position":
			// position [fenstring | startpos] moves
			err = parsePosition(input, currentBoard)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsePos: %v", err)
			}
			// FEN documentation and examples: https://www.chess.com/terms/fen-chess
		case "d":
			// debug, should be used like: position fen .... -> d -> check if position is any good
			fmt.Println(currentBoard)
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
