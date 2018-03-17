package main

import (
    "github.com/dilyar85/chess/game"
    "os"
)

func main() {
    chessGame := game.New()

    arg := os.Args
    isInteractiveMode := len(arg) <= 1

    if isInteractiveMode {
        chessGame.StartInteractiveMode()
    } else {
        chessGame.StartFileMode(arg[1])
    }

}