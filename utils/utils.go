package utils

import (
    "bytes"
    "strconv"
    "os"
    "bufio"
    "strings"
)

// MARK: Static methods of Utils class
func StringifyBoard(board [][]string) string {

    row := len(board)
    col := len(board[0])

    var buffer bytes.Buffer

    //Print top column letters
    for i := 0; i < col; i++ {
        if i == 0 {
            buffer.WriteString("   ") //padding for col letters
        }
        colLetter := (string)('a' + i)
        buffer.WriteString(" " + colLetter + "  ")
    }
    buffer.WriteString("\n")

    //Print each row
    for i := row - 1; i >= 0; i-- {
        //Print left row numbers following with " |"
        buffer.WriteString(strconv.Itoa(i + 1)) //has to convert string for WriteString() method
        buffer.WriteString(" |")

        //Print each square in same row
        for j := 0; j < col; j++ {
            buffer.WriteString(stringifySquare(board[-i+len(board)-1][j]))
        }

        //Print right row numbers
        buffer.WriteString(" " + strconv.Itoa(i + 1)) //has to convert string for WriteString() method
        buffer.WriteString("\n")
        if i != 0 {
            buffer.WriteString("\n")
        }

    }

    //Print bottom column letters
    for i := 0; i < col; i++ {
        if i == 0 {
            buffer.WriteString("   ") //padding for col letters
        }
        colLetter := (string)('a' + i)
        buffer.WriteString(" " + colLetter + "  ")
    }

    buffer.WriteString("\n")

    return buffer.String()
}

func stringifySquare(sq string) string {
    if len(sq) == 0 {
        return " _ |"
    } else {
        return " " + sq + " |"
    }
}

func ParseTestCase(path string) TestCase {
    file, err := os.Open(path)

    defer file.Close()

    if err != nil {
        panic(err)
    }

    var line string
    reader := bufio.NewReader(file)
    line, err = reader.ReadString('\n')
    line = strings.TrimSpace(line)

    var initialPositions []InitialPosition

    for line != "" {
        lineParts := strings.Split(line, " ")
        initialPositions = append(initialPositions, InitialPosition{lineParts[0], lineParts[1]})
        line, _ = reader.ReadString('\n')
        line = strings.TrimSpace(line)
    }


    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
    whiteCaptures := strings.Split(line[1: len(line)-1], " ")

    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)
    blackCaptures := strings.Split(line[1: len(line)-1], " ")

    line, _ = reader.ReadString('\n')
    line = strings.TrimSpace(line)

    var moves []string
    for line != "" {
        line = strings.TrimSpace(line)
        moves = append(moves, line)
        line, _ = reader.ReadString('\n')
    }

    return TestCase{initialPositions,whiteCaptures,blackCaptures, moves}
}


// MARK: InitialPosition class
type InitialPosition struct {
    Sign     string
    Position string
}

func (ip InitialPosition) String() string {
    return ip.Sign + " " + ip.Position
}

// MARK: TestCase class
type TestCase struct {
    InitialPositions []InitialPosition
    WhiteCaptures, BlackCaptures []string
    Moves            []string

}

func (tc TestCase) String() string {
    var buffer bytes.Buffer

    buffer.WriteString("initialPieces: [\n")
    for _, piece := range tc.InitialPositions {
        buffer.WriteString(piece.String())
        buffer.WriteString("\n")
    }
    buffer.WriteString("]\n")

    buffer.WriteString("moves: [\n")
    for _, move := range tc.Moves {
        buffer.WriteString(move)
        buffer.WriteString("\n")
    }
    buffer.WriteString("]")

    return buffer.String()
}
