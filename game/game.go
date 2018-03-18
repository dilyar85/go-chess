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

func (game ChessGame) StartInteractiveMode() {

    game.setupBoard(InitialBoardFileName)
    game.printGameStatus()

    for {
        game.nextTurn()
        input := game.promptInput(game.inputReader)

        gameEnd := game.execute(input)
        if gameEnd {
            return
        }

        game.printAction(input)
        game.printGameStatus()
    }

}

func (game ChessGame) StartFileMode(path string) {
    fmt.Println("Entered file path:", path)
}


func (game ChessGame) setupBoard(path string) utils.TestCase {

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
func (game ChessGame) execute(command string)  bool {

    defer func() {
        if err := recover(); err != nil {
            fmt.Println(err)
            game.promptInput(game.inputReader)
        }
    }()

    checkmate := game.board.execute(command, game.curTeam)

    if checkmate {
        game.endGameWithWinner(game.getCuPlayerName(), "Checkmate", command)
        return true
    }

    if game.isTie() {
        game.endGameByTie(command)
        return true
    }

    return false
}


func (game ChessGame) promptInput(reader bufio.Reader) string {
    fmt.Print(game.getCuPlayerName(), "> ")
    input, _ := reader.ReadString('\n')
    input = strings.TrimRight(input, "\n") //remove "\n" from input

    return input
}


func (game *ChessGame) nextTurn() {
    game.movesCount++

    //Set lower player if undecided
    if game.curTeam == undecided {
        game.curTeam = white
        return
    }

    //Switch team
    switch game.curTeam {
    case black:
        game.curTeam = white
    case white:
        game.curTeam = black
    }
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
    fmt.Println(game.getCuPlayerName(), " player action: ", action)
}


func (game ChessGame) getCuPlayerName() string {

    if game.curTeam == black {
        return "BLACK"
    } else {
        return "WHITE"
    }

}

func (game ChessGame) getOtherPlayerName() string {
    if game.curTeam == black {
        return "white"
    } else {
        return "lower"
    }
}


// Player.Team Enum
type Team int
const (
    undecided Team = iota
    white     Team = iota
    black     Team = iota
)


func getOpponentTeam(curTeam Team) Team {
    if curTeam == black {
        return white
    } else {
        return black
    }
}


