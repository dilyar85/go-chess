package game

import (
    "github.com/dilyar85/chess/utils"
    "strings"
    "bytes"
)

const boardSize = 8

type Board struct {
    squares                      [][]Square
    whiteCaptures, blackCaptures []string
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
    var buffer bytes.Buffer

    boardStr := make([][]string, boardSize)
    for i := 0; i < boardSize; i++ {
        boardStr[i] = make([]string, boardSize)
        for j := 0; j < boardSize; j++ {
            boardStr[i][j] = board.squares[i][j].String()
        }
    }
    buffer.WriteString(utils.StringifyBoard(boardStr))

    //White Captures: [symbols...]
    buffer.WriteString("White Captures: [")
    for _, capturedSign := range board.whiteCaptures {
        if capturedSign != "" {
            buffer.WriteString(getPieceSymbol(capturedSign) + " ") //cannot do getPieceSymbol("")
        }
    }
    buffer.WriteString("] \n")

    //Black Captures: [symbols...]
    buffer.WriteString("Black Captures: [")
    for _, capturedSign := range board.blackCaptures {
        if capturedSign != "" {
            buffer.WriteString(getPieceSymbol(capturedSign) + " ") //cannot do getPieceSymbol("")
        }
    }
    buffer.WriteString("] \n")

    return buffer.String()
}

func (board *Board) setup(testCase utils.TestCase) {
    for _, ip := range testCase.InitialPositions {
        board.initPiece(ip.Position, ip.Sign)
    }

    board.whiteCaptures = testCase.WhiteCaptures
    board.blackCaptures = testCase.BlackCaptures
}

func (board *Board) initPiece(position string, sign string) {
    square := board.getSquare(position)
    if square == nil || square.hasPiece() {
        panic("initPiece() failed on the position: " + position)
        return
    }

    piece := createPiece(sign, square.row, square.col)
    square.setPiece(&piece)
}

// Execute the command passed, return true if it's inCheckmate
func (board *Board) execute(command string, team Team) bool {

    //Handle "in check situation" first
    if board.inCheck(team) {
        validMoves := board.getAvailableMovesInCheck(team)
        if !containsMove(validMoves, command) {
            panic(illegalMoveMessage)
        }
    }

    tokens := strings.Split(command, " ")
    if len(tokens) != 2 {
        panic(illegalMoveMessage)
    }

    //Check the movePiece first, panic if it's illegal movePiece on board
    move := board.checkMove(tokens[0], tokens[1], team)

    //Move Piece
    board.movePiece(move.piece, move.squareFrom, move.squareTo)

    //Promote
    //TODO: Handle promotion for Pawn

    //return if the opponent team is in checkmate
    return board.inCheckmate(getOpponentTeam(team))

}

func (board *Board) movePiece(piece *Piece, squareFrom, squareTo *Square) {

    capturedPiece := squareTo.piece
    if capturedPiece != nil {
        board.captured(*capturedPiece)
    }

    //Update squares on board
    squareFrom.setPiece(nil)
    squareTo.setPiece(piece)

    //Update Piece row and col
    piece.row = squareTo.row
    piece.col = squareTo.col

}

func (board *Board) captured(capturedPiece Piece) {
    team := getOpponentTeam(capturedPiece.team)
    sign := capturedPiece.sign //board will print the symbols from piece.sign
    switch team {
    case white:
        board.whiteCaptures = append(board.whiteCaptures, sign)
    case black:
        board.blackCaptures = append(board.blackCaptures, sign)
    }
}

func (board Board) inCheck(curTeam Team) bool {
    kingPosition := board.getKingPosition(curTeam)
    opponentPositions := board.getReachablePositions(getOpponentTeam(curTeam))
    for _, position := range opponentPositions {
        if kingPosition == position {
            return true
        }
    }
    return false
}

func (board Board) inCheckmate(curTeam Team) bool {
    return board.inCheck(curTeam) && len(board.getAvailableMovesInCheck(curTeam)) == 0
}


//Get all available moves to escape check for the given team
func (board Board) getAvailableMovesInCheck(team Team) []string {
    var availableMoves []string

    //Escape check by moving King
    availableMoves = append(availableMoves, board.getAvailableMovesByMovingKing(team)...)

    //Escape check by moving other pieces than King
    availableMoves = append(availableMoves, board.getAvailableMovesByMovingOtherPieces(team)...)

    return availableMoves
}

func (board Board) getAvailableMovesByMovingKing(team Team) []string {
    //No need to check threatenKingPiece size cause we are moving the King to safe place directly
    var moves []string
    kingPosition := board.getKingPosition(team)
    kingPiece := board.getSquare(kingPosition).piece

    opponentMoves := board.getReachablePositions(getOpponentTeam(team))

    for _, kingMove := range getMoves(board, *kingPiece) {
        if !containsMove(opponentMoves, kingMove) {
            moves = append(moves, kingPosition+" "+kingMove)
        }
    }

    return moves
}

func (board Board) getAvailableMovesByMovingOtherPieces(team Team) []string {
    var moves []string
    threatenKingPieces := board.getThreateningKingPieces(team)
    if len(threatenKingPieces) != 1 {
        return moves //Too many threatening pieces to capture
    }

    threatenPiece := threatenKingPieces[0]
    threatenPosition := getCoordinatePosition(threatenPiece.row, threatenPiece.col)

    for _, ownPiece := range board.getAllPieces(team) {
        if isKing(ownPiece) {
            continue //king's move was already considered in method getAvailableMovesByMovingKing()
        }
        positionFrom := getCoordinatePosition(ownPiece.row, ownPiece.col)

        for _, positionTo := range getMoves(board, ownPiece) {
            if !board.moveWillCauseSelfCheck(positionFrom, positionTo, team) {
                moves = append(moves, positionFrom + " " + threatenPosition)
            }
        }
    }

    return moves
}


func (board Board) getThreateningKingPieces(team Team) []Piece {
    var threatenPieces []Piece
    kingPosition := board.getKingPosition(team)
    for _, opponentPiece := range board.getAllPieces(getOpponentTeam(team)) {
        opponentMoves := getMoves(board, opponentPiece)
        if containsMove(opponentMoves, kingPosition) {
            threatenPieces = append(threatenPieces, opponentPiece)
        }
    }
    return threatenPieces
}

func (board Board) moveWillCauseSelfCheck(positionFrom, positionTo string, team Team) bool {

    squareFrom := board.getSquare(positionFrom)
    squareTo := board.getSquare(positionTo)

    pieceFrom := squareFrom.piece
    pieceTo := squareTo.piece

    //Move piece and check
    board.movePiece(pieceFrom, squareFrom, squareTo)
    selfInCheck := board.inCheck(team)

    //Move pieces back
    board.movePiece(pieceFrom, squareTo, squareFrom)
    squareTo.setPiece(pieceTo)

    return selfInCheck
}

func (board Board) canMoveTo(position string, team Team) bool {

    square := board.getSquare(position)
    if square == nil {
        return false
    }
    piece := square.getPiece()
    return piece == nil || piece.team != team
}

func (board Board) isEmptyAt(position string) bool {
    square := board.getSquare(position)
    return square != nil && !square.hasPiece()
}

func (board Board) getKingPosition(team Team) string {
    var kingSymbol string
    if team == white {
        kingSymbol = WhiteKing
    } else {
        kingSymbol = BlackKing
    }

    for i := 0; i < boardSize; i++ {
        for j := 0; j < boardSize; j++ {
            square := board.squares[i][j]
            piece := square.getPiece()
            if piece != nil && piece.String() == kingSymbol {
                return getSquarePosition(square)
            }
        }
    }

    panic("Error: Cannot find King from the board")
}

func (board Board) getReachablePositions(team Team) []string {
    var moves []string
    pieces := board.getAllPieces(team)
    for _, piece := range pieces {
        moves = append(moves, getMoves(board, piece)...)
    }
    return moves
}

func (board Board) getAllPieces(team Team) []Piece {
    var pieces []Piece
    for i := 0; i < boardSize; i++ {
        for j := 0; j < boardSize; j++ {
            piece := board.squares[i][j].getPiece()
            if piece != nil && piece.team == team {
                pieces = append(pieces, *piece)
            }

        }
    }
    return pieces
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

// MARK: Move & Drop (helper struct for executing "movePiece" and "drop" commands )
type Move struct {
    piece                *Piece
    squareFrom, squareTo *Square
}

func (board Board) checkMove(origin, destination string, team Team) Move {

    //Check input positions
    squareFrom := board.getSquare(origin)
    if squareFrom == nil {
        panic(illegalMoveMessage)
    }
    squareTo := board.getSquare(destination)
    if squareTo == nil {
        panic(illegalMoveMessage)
    }
    piece := squareFrom.getPiece()
    if piece == nil || piece.team != team {
        panic(illegalMoveMessage)
    }

    //Check if the piece's movement is valid
    moves := getMoves(board, *piece)
    if !containsMove(moves, destination) {
        panic(illegalMoveMessage)
    }

    //Check if causing self in check (considered as invalid movePiece in current rule)
    if board.moveWillCauseSelfCheck(origin, destination, team) {
        panic(causingSelfInCheckMessage)
    }

    return Move{piece, squareFrom, squareTo}

}

// MARK: Helper package functions
func containsMove(moves []string, move string) bool {
    for _, element := range moves {
        if move == element {
            return true
        }
    }
    return false
}

func getCoordinatePosition(row, col int) string {
    return (string)('a'+col) + (string)('0'-row+boardSize)
}

func getSquarePosition(square Square) string {
    return getCoordinatePosition(square.row, square.col)
}
