package game

import (
    "fmt"
    "bufio"
    "os"
    "github.com/dilyar85/chess/utils"
)

func New() ChessGame {
    game := ChessGame{NewBoard(), 0, undecided}
    return game
}

const (
    MovesLimitCount      = 400
    InitialBoardFileName = "./playbook/initialBoard.txt"
)

type ChessGame struct {
    board      *Board
    movesCount int
    curTeam    Team
}

func (game ChessGame) StartInteractiveMode() {

    game.init(InitialBoardFileName)
    game.printGameStatus()

    reader := bufio.NewReader(os.Stdin)

    for {
        game.nextTurn()

        fmt.Print(game.getCuPlayerName(), "> ")

        command, _ := reader.ReadString('\n')

        gameEnd := game.execute(command)
        if gameEnd {
            return
        }

        game.printAction(command)
        game.printGameStatus()
    }

}

func (game ChessGame) StartFileMode(path string) {
    fmt.Println("Entered file path:", path)
}

func (game ChessGame) init(path string) utils.TestCase {

    defer func() {
        if err := recover(); err != nil {
            fmt.Println("Error: Unable to init the game from file: " + path)
        }
    }()

    testCase := utils.ParseTestCase(path)
    game.board.setupBoard(testCase)
    return testCase
}

func (game ChessGame) printGameStatus() {
    //Print current board
    fmt.Println(game.board.String())

    //TODO: Add captured stuff
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

func (game ChessGame) execute(command string) (gameEnd bool) {

    defer func() {
        if err := recover(); err != nil {
            game.endGameWithWinner(game.getOtherPlayerName(), err, command)
            gameEnd = true
        }
    }()

    //checkmate := game.board.Execute(command, game.curTeam) //TODO
    checkmate := false

    if checkmate {
        game.endGameWithWinner(game.getCuPlayerName(), "Checkmate", command)
        gameEnd = true
    }

    if game.isTie() {
        game.endGameByTie(command)
        gameEnd = true
    }

    return gameEnd
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

func (game ChessGame) printAction(action string) {
    fmt.Print(game.getCuPlayerName(), " player action: ", action)
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

func (game ChessGame) isTie() bool {
    return game.movesCount >= MovesLimitCount
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


