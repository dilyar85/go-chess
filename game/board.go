package game

import "github.com/dilyar85/chess/utils"

const boardSize = 8

type Board struct {
    squares [][]Square
}

func NewBoard() *Board {
    board := new(Board)
    //Init Squares
    squares := make([][]Square, boardSize)
    for i := 0; i < boardSize; i++ {
        squares[i] = make([]Square, boardSize)
        for j := 0; j < boardSize; j++ {
            square := Square{i, j, nil}
            squares[i][j] = square
        }
    }
    board.squares = squares
    return board
}

func (board Board) String() string {
    boardStr := make([][]string, boardSize)
    for i := 0; i < boardSize; i++ {
        boardStr[i] = make([]string, boardSize)
        for j := 0; j < boardSize; j++ {
            boardStr[i][j] = board.squares[i][j].String()
        }
    }
    return utils.StringifyBoard(boardStr)
}



func (board Board) setupBoard(testCase utils.TestCase) {
    for _, ip := range testCase.InitialPositions {
        board.initPiece(ip.Position, ip.Sign)
    }
}

func (board Board) initPiece(position string, sign string) {
    square := board.getSquare(position)
    if square == nil || square.hasPiece() {
        panic("initPiece() failed on the position: " + position)
        return
    }

    piece := createPiece(sign, square.row, square.col)
    square.setPiece(&piece)
}

func (board Board) getSquare(position string) *Square {
    col := int(position[0] - 'a')
    if col < 0 || col >= boardSize {
        return nil
    }
    row := ((int)(position[1] - '0')) * -1 + boardSize
    if row < 0 || row >= boardSize {
        return nil
    }
    return &board.squares[row][col]
}

// MARK: Square struct
type Square struct {
    row, col int
    piece    *Piece
}

func (square Square) String() string {
    if square.piece == nil {
        return ""
    }

    return square.piece.String()
}

func (square *Square) setPiece(piece *Piece) {
    square.piece = piece
}

func (square Square) hasPiece() bool {
    return square.piece != nil
}

func (square Square) getPiece() *Piece {
    return square.piece
}

