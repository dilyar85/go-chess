package game

import "strings"

type PieceCharacter int

const (
    King   PieceCharacter = iota
    Queen  PieceCharacter = iota
    Rook   PieceCharacter = iota
    Bishop PieceCharacter = iota
    Knight PieceCharacter = iota
    Pawn   PieceCharacter = iota
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
    if sign == strings.ToUpper(sign) {
        team = white
    } else {
        team = black
    }

    return Piece{sign, team, row, col}

}

func (piece Piece) String() string {
    switch piece.sign {
    case "K":
        return WhiteKing
    case "k":
        return BlackKing
    case "Q":
        return WhiteQueen
    case "q":
        return BlackQueen
    case "R":
        return WhiteRook
    case "r":
        return BlackRook
    case "B":
        return WhiteBishop
    case "b":
        return BlackBishop
    case "N":
        return WhiteKnight
    case "n":
        return BlackKnight
    case "P":
        return WhitePawn
    case "p":
        return BlackPawn

    default:
        panic("Cannot Print Piece. Unknown Sign: " + piece.sign)
    }
}

func getPieceCharacter(piece Piece) PieceCharacter {
    pureSign := strings.Replace(piece.sign, "+", "", -1)
    switch strings.ToUpper(pureSign) {
    case "K":
        return King
    case "Q":
        return Queen
    case "R":
        return Rook
    case "B":
        return Bishop
    case "N":
        return Knight
    case "P":
        return Pawn
    default:
        panic("Undefined Piece Character: " + piece.sign)
    }
}

// MARK: Piece

type Piece struct {
    sign     string
    team     Team
    row, col int
}
