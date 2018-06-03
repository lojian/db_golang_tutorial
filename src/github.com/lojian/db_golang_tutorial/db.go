package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type MetaCommandResult int
type PrepareResult int
type StatementType int

const (
	ExitSuccess = 0
	ExitFailure = -1

	META_COMMAND_SUCCESS              MetaCommandResult = 0
	META_COMMAND_UNRECOGNIZED_COMMAND MetaCommandResult = 1

	PREPARE_SUCCESS                PrepareResult = 0
	PREPARE_UNRECOGNIZED_STATEMENT PrepareResult = 1
	PREPARE_SYNTAX_ERROR           PrepareResult = 2

	STATEMENT_INSERT StatementType = 0
	STATEMENT_SELECT StatementType = 1

	COLUMN_USERNAME_SIZE uint32 = 32
	COLUMN_EMAIL_SIZE    uint32 = 255

	ID_SIZE         uint32 = 4
	USERNAME_SIZE   uint32 = COLUMN_USERNAME_SIZE
	EMAIL_SIZE      uint32 = COLUMN_EMAIL_SIZE
	ID_OFFSET              = 0
	USERNAME_OFFSET        = ID_OFFSET + ID_SIZE
	EMAIL_OFFSET           = USERNAME_OFFSET + USERNAME_SIZE
	ROW_SIZE               = ID_SIZE + USERNAME_SIZE + EMAIL_SIZE
)

type Row struct {
	id       uint32
	username [COLUMN_USERNAME_SIZE]byte
	email    [COLUMN_EMAIL_SIZE]byte
}

type Statement struct {
	stmtType      StatementType
	row_to_insert Row
}

var (
	scanner     *bufio.Scanner
	inputBuffer string
)

func init() {
	scanner = bufio.NewScanner(os.Stdin)
}

func printPrompt() {
	fmt.Print("db> ")
}

func doMetaCommand(commandText *string) MetaCommandResult {
	if *commandText == ".exit" {
		os.Exit(ExitSuccess)
	} else {
		return META_COMMAND_UNRECOGNIZED_COMMAND
	}
	return META_COMMAND_SUCCESS
}

func prepareStatement(commandText *string, statement *Statement) PrepareResult {
	if strings.HasPrefix(*commandText, "insert") {
		statement.stmtType = STATEMENT_INSERT
		argsAssigned, _ := fmt.Sscanf("insert %d %s %s", *commandText)
		if argsAssigned < 3 {
			return PREPARE_SYNTAX_ERROR
		}

		return PREPARE_SUCCESS
	}
	if strings.HasPrefix(*commandText, "select") {
		statement.stmtType = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}
	return PREPARE_UNRECOGNIZED_STATEMENT
}

func executeStatement(statement *Statement) {
	switch statement.stmtType {
	case STATEMENT_INSERT:
		fmt.Printf("This is where we would do an insert.\n")
		break
	case STATEMENT_SELECT:
		fmt.Printf("This is where we would do a select.\n")
		break
	}
}

func serializeRow(source *Row, destination [ROW_SIZE]byte) {
	binary.LittleEndian.PutUint32(destination[0:3], source.id)
	//copy(destination[USERNAME_OFFSET:USERNAME_OFFSET+32], source.username)
}
func main() {

	for {
		printPrompt()
		fmt.Println(ID_SIZE)
		fmt.Println(USERNAME_SIZE)
		fmt.Println(EMAIL_SIZE)
		fmt.Println(ROW_SIZE)
		if scanner.Scan() {
			inputBuffer = scanner.Text()
			if inputBuffer[0] == '.' {
				switch doMetaCommand(&inputBuffer) {
				case META_COMMAND_SUCCESS:
					continue
				case META_COMMAND_UNRECOGNIZED_COMMAND:
					fmt.Printf("Unrecognized command '%s'\n", inputBuffer)
					continue
				}
			}

			statement := new(Statement)
			switch prepareStatement(&inputBuffer, statement) {
			case PREPARE_SUCCESS:
				break
			case PREPARE_SYNTAX_ERROR:
				fmt.Printf("Syntax error. Could not parse statement.\n")
				continue
			case PREPARE_UNRECOGNIZED_STATEMENT:
				fmt.Printf("Unrecognized keyword at start of '%s'.\n", inputBuffer)
				continue
			}

			executeStatement(statement)
			fmt.Println("Executed.")

		} else {
			fmt.Println("Error reading input")
			os.Exit(ExitFailure)
		}
	}

}
