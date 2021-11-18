package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var whiteSpacesRegExp = regexp.MustCompile("\\s+")

type Engine struct {
	ConfigFilePath        *string
	ConfigRaw             string
	OutOfBoundBehavior    *string
	PrintCommandsAtTheEnd *bool

	Area Area
}

func (e *Engine) InitCmdParams() {
	e.ConfigFilePath = flag.String("config", "", "Path to config file")
	e.OutOfBoundBehavior = flag.String("out-of-bound-behavior", "", "Path to config file")
	e.PrintCommandsAtTheEnd = flag.Bool("print-commands-at-the-end", false, "Print array of processed commands")

	flag.Parse()
}

func (e *Engine) LoadConfig() {
	if *e.ConfigFilePath == "" {
		panic("Config file path is not set")
	}

	content, err := ioutil.ReadFile(*e.ConfigFilePath)

	if err != nil {
		log.Fatal("Fail reading config: ", err)
	}

	e.ConfigRaw = string(content)
}

func (e *Engine) ParseConfig() {
	if e.ConfigRaw == "" {
		panic("Config data is empty")
	}

	lines := strings.Split(e.ConfigRaw, "\n")

	if len(lines) <= 4 {
		panic("Config have no enough lines")
	}

	if *e.OutOfBoundBehavior != "" {
		if err := e.Area.SetOutOfBoundBehavior(*e.OutOfBoundBehavior); err != nil {
			log.Fatal("Incorrect out of bound behavior: ", *e.OutOfBoundBehavior)
		}
	} else {
		_ = e.Area.SetOutOfBoundBehavior(OutOfBoundBehaviorPass)
	}

	e.parseAreaDimensions(lines[0])
	e.parseTraps(lines[1])
	e.parseRefuge(lines[2])
	e.parseTurtle(lines[3])
	e.parseCommands(lines[4:])
}

func (e *Engine) parseAreaDimensions(row string) {
	chunks := preProcessRow(row)

	width, height := getNumbersPair(chunks)

	e.Area.Width = width
	e.Area.Height = height
}

func (e *Engine) parseTraps(row string) {
	if row == "" {
		return
	}

	chunks := preProcessRow(row)

	for _, pair := range chunks {
		pairChunks := strings.Split(whiteSpacesRegExp.ReplaceAllString(strings.TrimSpace(pair), ""), ",")

		x, y := getNumbersPair(pairChunks)

		if err := e.Area.AddTrap(x, y); err != nil {
			log.Fatal("Incorrect trap position: ", pair)
		}
	}
}

func (e *Engine) parseRefuge(row string) {
	chunks := preProcessRow(row)

	x, y := getNumbersPair(chunks)

	if err := e.Area.AddRefuge(x, y); err != nil {
		log.Fatal("Incorrect refuge position: ", row)
	}
}

func (e *Engine) parseTurtle(row string) {
	chunks := preProcessRow(row)

	if len(chunks) != 3 {
		log.Fatal("Incorrect turtle row: ", row)
	}

	x, y := getNumbersPair(chunks)

	if err := e.Area.AddTurtle(x, y, chunks[2]); err != nil {
		log.Fatal("Incorrect turtle direction or position: ", row)
	}
}

func (e *Engine) parseCommands(commands []string) {
	for _, command := range commands {
		chunks := preProcessRow(command)

		if err := e.Area.AddCommand(chunks); err != nil {
			log.Fatal("Incorrect turtle commands: ", chunks)
		}
	}
}

func (e *Engine) ProcessCommands() {
	e.Area.ProcessCommands()
}

func (e *Engine) PrintResult() {
	for index, command := range e.Area.Commands {
		turtleNo := index + 1
		commandsLeft := len(command.Commands) - command.CurrentCommandIndex - 1
		commandsLeftText := ""

		if commandsLeft > 0 {
			commandsLeftText = fmt.Sprintf(", commands left: %d", commandsLeft)
		}

		if command.IsFinished {
			if command.IsSaved {
				fmt.Printf("Turtle %d finish travel and saved%s\n", turtleNo, commandsLeftText)
			} else if command.IsCaught {
				fmt.Printf("Turtle %d finish travel and caught%s\n", turtleNo, commandsLeftText)
			}
		} else {
			fmt.Printf("Turtle %d not finish travel and still in area%s\n", turtleNo, commandsLeftText)
		}
	}

	if *e.PrintCommandsAtTheEnd {
		fmt.Printf("%+v\n", e.Area.Commands)
	}
}

func preProcessRow(row string) []string {
	return strings.Split(whiteSpacesRegExp.ReplaceAllString(strings.TrimSpace(row), " "), " ")
}

func getNumbersPair(pair []string) (int, int) {
	if len(pair) < 2 {
		log.Fatal("Incorrect pair: ", pair)
	}

	first, err := strconv.Atoi(pair[0])

	if err != nil {
		log.Fatal("Incorrect first number: ", pair)
	}

	second, err := strconv.Atoi(pair[1])

	if err != nil {
		log.Fatal("Incorrect second number: ", pair)
	}

	return first, second
}
