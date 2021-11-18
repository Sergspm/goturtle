package main

import (
	"errors"
	"sync"
)

const (
	EntityTypeTrap   = "trap"
	EntityTypeRefuge = "refuge"

	TurtleDirectionNorth = "N"
	TurtleDirectionWest  = "W"
	TurtleDirectionSouth = "S"
	TurtleDirectionEast  = "E"

	TurtleCommandMove        = "M"
	TurtleCommandRotateRight = "R"
	TurtleCommandRotateLeft  = "L"

	OutOfBoundBehaviorPass   = "pass"
	OutOfBoundBehaviorStop   = "stop"
	OutOfBoundBehaviorPortal = "portal"
)

var (
	ErrIncorrectTurtleDirection    = errors.New("incorrect turtle direction")
	ErrIncorrectCommand            = errors.New("incorrect command")
	ErrIncorrectOutOfBoundBehavior = errors.New("incorrect out of bound behavior")
	ErrIncorrectPosition           = errors.New("incorrect position")

	RotateMap = map[string]map[string]string{
		TurtleCommandRotateRight: {
			TurtleDirectionNorth: TurtleDirectionEast,
			TurtleDirectionEast:  TurtleDirectionSouth,
			TurtleDirectionSouth: TurtleDirectionWest,
			TurtleDirectionWest:  TurtleDirectionNorth,
		},
		TurtleCommandRotateLeft: {
			TurtleDirectionNorth: TurtleDirectionWest,
			TurtleDirectionWest:  TurtleDirectionSouth,
			TurtleDirectionSouth: TurtleDirectionEast,
			TurtleDirectionEast:  TurtleDirectionNorth,
		},
	}
)

type Area struct {
	Width  int
	Height int

	OutOfBoundBehavior string

	Entities []Entity
	Turtle   Turtle
	Commands []Command
}

func (a *Area) IsOutOfBoundBehaviorStop() bool {
	return a.OutOfBoundBehavior == OutOfBoundBehaviorStop
}

func (a *Area) IsOutOfBoundBehaviorPortal() bool {
	return a.OutOfBoundBehavior == OutOfBoundBehaviorPortal
}

func (a *Area) AddTrap(x int, y int) error {
	return a.addEntity(x, y, EntityTypeTrap)
}

func (a *Area) AddRefuge(x int, y int) error {
	return a.addEntity(x, y, EntityTypeRefuge)
}

func (a *Area) AddTurtle(x int, y int, direction string) error {
	if !isTurtleDirectionCorrect(direction) {
		return ErrIncorrectTurtleDirection
	}
	if x < 0 || x >= a.Width || y < 0 || y >= a.Height {
		return ErrIncorrectPosition
	}
	a.Turtle = Turtle{
		X:         x,
		Y:         y,
		Direction: direction,
	}
	return nil
}

func (a *Area) AddCommand(commands []string) error {
	for _, command := range commands {
		if !isTurtleCommandCorrect(command) {
			return ErrIncorrectCommand
		}
	}
	a.Commands = append(a.Commands, Command{
		Commands: commands,
		Turtle:   a.Turtle,
	})
	return nil
}

func (a *Area) SetOutOfBoundBehavior(behavior string) error {
	if !isOutOfBoundBehaviorCorrect(behavior) {
		return ErrIncorrectOutOfBoundBehavior
	}
	a.OutOfBoundBehavior = behavior
	return nil
}

func (a *Area) addEntity(x int, y int, entityType string) error {
	if x < 0 || x >= a.Width || y < 0 || y >= a.Height {
		return ErrIncorrectPosition
	}

	a.Entities = append(a.Entities, Entity{
		X:    x,
		Y:    y,
		Type: entityType,
	})

	return nil
}

func (a *Area) ProcessCommands() {
	if len(a.Commands) == 0 {
		return
	}

	var wg sync.WaitGroup

	wg.Add(len(a.Commands))

	for index := range a.Commands {
		go func(command *Command) {
			command.Process(*a)
			wg.Done()
		}(&a.Commands[index])
	}

	wg.Wait()
}

type Entity struct {
	X    int
	Y    int
	Type string
}

func (e *Entity) IsTrap() bool {
	return e.Type == EntityTypeTrap
}

func (e *Entity) IsRefuge() bool {
	return e.Type == EntityTypeRefuge
}

type Turtle struct {
	X         int
	Y         int
	Direction string
}

type CommandLog struct {
	Command string
	Action  string
	Index   int
}

type Command struct {
	Commands            []string
	Log                 []CommandLog
	CurrentCommandIndex int
	Turtle              Turtle
	IsFinished          bool
	IsSaved             bool
	IsCaught            bool
	IsOutOfBound        bool
}

func (c *Command) Process(area Area) {
	turtle := &c.Turtle
	outOfBoundBehaviorPortal := area.IsOutOfBoundBehaviorPortal()
	outOfBoundBehaviorStop := area.IsOutOfBoundBehaviorStop()

	for index, cmd := range c.Commands {
		if !c.checkCommandConditions(index, cmd, area) {
			return
		}

		c.CurrentCommandIndex = index

		if cmd == TurtleCommandRotateRight {
			turtle.Direction = RotateMap[TurtleCommandRotateRight][turtle.Direction]
			c.pushToLog(index, cmd, "changed direction to "+turtle.Direction)
		} else if cmd == TurtleCommandRotateLeft {
			turtle.Direction = RotateMap[TurtleCommandRotateLeft][turtle.Direction]
			c.pushToLog(index, cmd, "changed direction to "+turtle.Direction)
		} else if cmd == TurtleCommandMove {
			if turtle.Direction == TurtleDirectionNorth {
				if turtle.Y <= 0 {
					if outOfBoundBehaviorPortal {
						turtle.Y = area.Height - 1
						c.pushToLog(index, cmd, "portal by y step to bottom")
					} else if outOfBoundBehaviorStop {
						c.pushToLog(index, cmd, "stop by y step to bottom")
						return
					} else {
						c.pushToLog(index, cmd, "pass decrement y step")
					}
				} else {
					turtle.Y--
					c.pushToLog(index, cmd, "decrement y step")
				}
			} else if turtle.Direction == TurtleDirectionSouth {
				if turtle.Y >= area.Height-1 {
					if outOfBoundBehaviorPortal {
						turtle.Y = 0
						c.pushToLog(index, cmd, "portal by y step to top")
					} else if outOfBoundBehaviorStop {
						c.pushToLog(index, cmd, "stop by y step to top")
						return
					} else {
						c.pushToLog(index, cmd, "pass increment y step")
					}
				} else {
					turtle.Y++
					c.pushToLog(index, cmd, "increment y step")
				}
			} else if turtle.Direction == TurtleDirectionEast {
				if turtle.X >= area.Width-1 {
					if outOfBoundBehaviorPortal {
						turtle.X = 0
						c.pushToLog(index, cmd, "portal by x step to top")
					} else if outOfBoundBehaviorStop {
						c.pushToLog(index, cmd, "stop by x step to top")
						return
					} else {
						c.pushToLog(index, cmd, "pass increment x step")
					}
				} else {
					turtle.X++
					c.pushToLog(index, cmd, "increment x step")
				}
			} else if turtle.Direction == TurtleDirectionWest {
				if turtle.X <= 0 {
					if outOfBoundBehaviorPortal {
						turtle.X = area.Width - 1
						c.pushToLog(index, cmd, "portal by x step to bottom")
					} else if outOfBoundBehaviorStop {
						c.pushToLog(index, cmd, "stop by x step to bottom")
						return
					} else {
						c.pushToLog(index, cmd, "pass decrement x step")
					}
				} else {
					turtle.X--
					c.pushToLog(index, cmd, "decrement x step")
				}
			}
		}

		if !c.checkCommandConditions(index, cmd, area) {
			return
		}
	}
}

func (c *Command) checkCommandConditions(index int, cmd string, area Area) bool {
	if (c.Turtle.X < 0 || c.Turtle.X >= area.Width || c.Turtle.Y < 0 || c.Turtle.Y >= area.Height) &&
		area.IsOutOfBoundBehaviorStop() {
		c.pushToLog(index, cmd, "out of bound")
		c.IsOutOfBound = true
		return false
	}

	for _, entity := range area.Entities {
		if entity.X == c.Turtle.X && entity.Y == c.Turtle.Y {
			c.IsFinished = true
			if entity.IsTrap() {
				c.IsCaught = true
			} else if entity.IsRefuge() {
				c.IsSaved = true
			}
			c.pushToLog(index, cmd, "finished")
			return false
		}
	}

	return true
}

func (c *Command) pushToLog(index int, command string, action string) {
	c.Log = append(c.Log, CommandLog{
		Command: command,
		Action:  action,
		Index:   index,
	})
}

// TODO тесты

func isTurtleDirectionCorrect(direction string) bool {
	return direction == TurtleDirectionNorth ||
		direction == TurtleDirectionWest ||
		direction == TurtleDirectionSouth ||
		direction == TurtleDirectionEast
}

func isTurtleCommandCorrect(command string) bool {
	return command == TurtleCommandMove || command == TurtleCommandRotateRight || command == TurtleCommandRotateLeft
}

func isOutOfBoundBehaviorCorrect(behavior string) bool {
	return behavior == OutOfBoundBehaviorPass ||
		behavior == OutOfBoundBehaviorStop ||
		behavior == OutOfBoundBehaviorPortal
}
