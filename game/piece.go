package game

import (
    "strings"
)


//type PieceSymbol string
const (
    WhiteKing = "\u2654"
    BlackKing = "\u265A"

    WhiteQueen = "\u2655"
    BlackQueen = "\u265B"

    WhiteRook = "\u2656"
    BlackRook = "\u265C"

    WhiteBishop = "\u2657"
    BlackBishop = "\u265D"

    WhiteKnight = "\u2658"
    BlackKnight = "\u265E"

    WhitePawn = "\u2659"
    BlackPawn = "\u265F"
)

func createPiece(sign string, row, col int) Piece {

    var team Team
    //Upper case represents Black Player and lower case represents White Player
    if sign == strings.ToUpper(sign) {
        team = black
    } else {
        team = white
    }

    return Piece{sign, team, row, col}

}

func getPieceSymbol(sign string) string {
    switch sign {
    case "k":
        return WhiteKing
    case "K":
        return BlackKing
    case "q":
        return WhiteQueen
    case "Q":
        return BlackQueen
    case "r":
        return WhiteRook
    case "R":
        return BlackRook
    case "b":
        return WhiteBishop
    case "B":
        return BlackBishop
    case "n":
        return WhiteKnight
    case "N":
        return BlackKnight
    case "p":
        return WhitePawn
    case "P":
        return BlackPawn

    default:
        panic("Cannot Print Piece. Unknown Sign:" + sign)
    }
}



// MARK: Functions about Piece's movement
func getMoves(board Board, piece Piece) []string {


    row := piece.row
    col := piece.col
    team := piece.team

    switch piece.String() {

    case WhiteKing, BlackKing:
        return getKingMoves(row, col, board, team)

    case WhiteQueen, BlackQueen:
        return getQueenMoves(row, col, board, team)

    case WhiteRook, BlackRook:
        return getRookMoves(row, col, board, team)

    case WhiteBishop, BlackBishop:
        return getBishopMovesAt(row, col, board, team)

    case WhiteKnight, BlackKnight:
        return getKnightMoves(row, col, board, team)

    case WhitePawn, BlackPawn:
        return getPawnMoves(row, col, board, team)


    default:
        panic("piece.String() hasn't been defined: " + piece.String())
    }
}

func getKingMoves(row, col int, board Board, team Team) []string {

    var moves []string

    for i := -1; i <= 1; i++ {
        for j := -1; j <= 1; j++ {
            if i == 0 && j == 0 {
                continue //original position
            }
            position := getCoordinatePosition(row+i, col+j)
            if board.canMoveTo(position, team) {
                moves = append(moves, position)
            }
        }
    }
    return moves
}

func getQueenMoves(row, col int, board Board, team Team) []string {
    var moves []string
    moves = append(moves, getRookMoves(row, col, board, team)...)
    moves = append(moves, getBishopMovesAt(row, col, board, team)...)
    return moves
}

func getRookMoves(row, col int, board Board, team Team) []string {
    var moves []string

    //Moving left
    for j := col - 1; j >= 0; j-- {
        position := getCoordinatePosition(row, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team){
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Moving Right
    for j := col + 1; j < boardSize; j++ {
        position := getCoordinatePosition(row, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team){
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Moving Up
    for i := row - 1; i >= 0; i-- {
        position := getCoordinatePosition(i, col)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Moving Down
    for i := row + 1; i < boardSize; i++ {
        position := getCoordinatePosition(i, col)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    return moves
}


func getBishopMovesAt(row, col int, board Board, team Team) []string {

    var moves []string

    //Move Top Left
    for i, j := row - 1, col - 1; i >= 0 && j >= 0; i, j = i - 1, j - 1 {
        position := getCoordinatePosition(i, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Move Top Right
    for i, j := row - 1, col + 1; i >= 0 && j < boardSize; i, j = i - 1, j + 1 {
        position := getCoordinatePosition(i, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Move Bottom Left
    for i, j := row + 1, col - 1; i < boardSize && j >= 0; i, j = i + 1, j - 1 {
        position := getCoordinatePosition(i, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    //Move Bottom Right
    for i, j := row + 1, col + 1; i < boardSize && j < boardSize; i, j = i + 1, j + 1 {
        position := getCoordinatePosition(i, j)
        if board.isEmptyAt(position) {
            moves = append(moves, position)
        } else if board.canMoveTo(position, team) {
            moves = append(moves, position)
            break
        } else {
            break
        }
    }

    return moves
}


func getKnightMoves(row, col int, board Board, team Team) []string {
    var moves []string

    for i := -2; i <= 2; i++ {
        if i == 0 {
            continue
        }
        var j int
        //up
        if i == -2 || i == 2 {
            j = 1
        } else {
            j = 2
        }
        leftward := getCoordinatePosition(row + i, col - j)
        if board.canMoveTo(leftward, team) {
            moves = append(moves, leftward)
        }
        rightward := getCoordinatePosition(row + i, col + j)
        if board.canMoveTo(rightward, team) {
            moves = append(moves, rightward)
        }
    }

    return moves
}

func getPawnMoves(row, col int, board Board, team Team) []string {

    var moves []string

    var firstMove bool

    oneStepRow := row
    twoStepsRow := row


    switch team {
    case white:
        oneStepRow--
        twoStepsRow -= 2
        firstMove = row == boardSize - 2
    case black:
        oneStepRow++
        twoStepsRow += 2
        firstMove = row == 1
    }

    //Get one step forwards positions
    position := getCoordinatePosition(oneStepRow, col)
    if board.isEmptyAt(position) {
        moves = append(moves, position)
    }

    //Get two step forwards positions if it's first move
    position = getCoordinatePosition(twoStepsRow, col)
    if firstMove && board.isEmptyAt(position) {
        moves = append(moves, position)
    }

    //Get Two Killing positions if there is enemy nearby to kill
    position = getCoordinatePosition(oneStepRow, col - 1)
    if board.canMoveTo(position, team) {
        moves = append(moves, position)
    }
    position = getCoordinatePosition(oneStepRow, col + 1)
    if board.canMoveTo(position, team) {
        moves = append(moves, position)
    }

    return moves
}


// MARK: Piece

type Piece struct {
    sign     string
    team     Team
    row, col int
}

func (piece Piece) String() string {
    return getPieceSymbol(piece.sign)
}

