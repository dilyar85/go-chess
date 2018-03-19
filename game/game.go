package game

import (
    "fmt"
    "bufio"
    "os"
    "github.com/dilyar85/chess/utils"
    "strings"
)

const (
    MovesLimitCount      = 400
    InitialBoardFileName = "./playbook/initialBoard.txt"
    illegalMoveMessage   = "Illegal move! Please enter again."
    causingSelfInCheckMessage = "This move will cause yourself in check! Please enter again."

)

func New() ChessGame {
    game := ChessGame{NewBoard(), 0, undecided, *bufio.NewReader(os.Stdin)}
    return game
}


type ChessGame struct {
    board      *Board
    movesCount int
    curTeam    Team
    inputReader bufio.Reader

}

func (game *ChessGame) StartInteractiveMode() {

    game.setupBoard(InitialBoardFileName)
    game.printGameStatus()

    for {
        game.changeTurn(true)
        game.printAvailableMovesInCheck()
        input := game.promptInput(game.inputReader)
        gameEnd := game.execute(input)
        if gameEnd {
            return
        }
    }

}

func (game *ChessGame) StartFileMode(path string) {
    fmt.Println("Entered file path:", path)
    //TODO: Implement File Mode
}


func (game *ChessGame) setupBoard(path string) utils.TestCase {

    defer func() {
        if err := recover(); err != nil {
            fmt.Println("Error: Unable to setup board from file: " + path)
        }
    }()

    testCase := utils.ParseTestCase(path)
    game.board.setup(testCase)
    return testCase
}


//Execute the command entered by user and return if game should end
func (game *ChessGame) execute(command string)  bool {

    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
            game.changeTurn(false)
        }
    }()

    checkmate := game.board.execute(command, game.curTeam)

    if checkmate {
        game.endGameWithWinner(getTeamName(game.curTeam), "Checkmate", command)
        return true
    }

    if game.isTie() {
        game.endGameByTie(command)
        return true
    }

    game.printAction(command)
    game.printGameStatus()
    return false
}


func (game ChessGame) promptInput(reader bufio.Reader) string {
    fmt.Print(getTeamName(game.curTeam), "> ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimRight(input, "\n") //remove "\n" from input

    return input
}


func (game *ChessGame) changeTurn(goesNext bool) {

    if goesNext {
        game.movesCount++
    } else {
        game.movesCount--
    }

    //Switch team
    switch game.curTeam {
    case undecided:
        game.curTeam = white
    case black:
        game.curTeam = white
    case white:
        game.curTeam = black
    }
}

func (game ChessGame) printAvailableMovesInCheck() {
    curTeam := game.curTeam

    if !game.board.inCheck(curTeam) {
        return
    }

    fmt.Println(getTeamName(game.curTeam) + " is in check!")
    fmt.Println("Available moves:")
    availableMoves := game.board.getAvailableMovesInCheck(curTeam)
    for _, move := range availableMoves {
        fmt.Println(move)
    }
    fmt.Println()
}


func (game ChessGame) isTie() bool {
    return game.movesCount >= MovesLimitCount
}


func (game ChessGame) endGameWithWinner(winnerPlayer string,  reason interface{}, lastCommand string) {
    game.printAction(lastCommand)
    game.printGameStatus()
    fmt.Println()
    fmt.Println(winnerPlayer, "player wins. ", reason)
}

func (game ChessGame) endGameByTie(lastCommand string) {
    game.printAction(lastCommand)
    game.printGameStatus()
    fmt.Println("Tie game.  Too many moves.")
}


func (game ChessGame) printGameStatus() {
    fmt.Println(game.board.String())
}

func (game ChessGame) printAction(action string) {
    fmt.Println(getTeamName(game.curTeam), " player action: ", action)
}


// Player.Team Enum
type Team int
const (
    undecided Team = iota
    white     Team = iota
    black     Team = iota
)

func getTeamName(team Team) string {
    switch team {
    case white:
        return "WHITE Player"
    case black:
        return "BLACK Player"

    default:
        return "Unknown Player"
    }
}

func getOpponentTeamName(team Team) string {
    var opponent Team
    if team == white {
        opponent = black
    } else {
        opponent = white
    }
    return getTeamName(opponent)
}


func getOpponentTeam(curTeam Team) Team {

    if curTeam == black {
        return white
    } else {
        return black
    }

}


