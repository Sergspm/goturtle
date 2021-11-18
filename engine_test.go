package main

import "testing"

const ScenarioCaught = `5 5
1,1 1,3 3,3
2 4
0 1 N
R M L M M`

const ScenarioInArea = `5 5
1,1 1,3 3,3
2 4
0 1 N
M R M M M`

const ScenarioSaved = `5 5
1,1 1,3 3,3
2 4
0 1 N
M R M M M M R M M M M R M M M`

func TestCaughtScenario(t *testing.T) {
	behavior := OutOfBoundBehaviorPass
	engine := Engine{
		ConfigRaw:          ScenarioCaught,
		OutOfBoundBehavior: &behavior,
	}

	engine.ParseConfig()
	engine.ProcessCommands()

	if len(engine.Area.Commands) != 1 {
		t.Fatalf("Commands number must be 1")
	}

	cmd := engine.Area.Commands[0]

	if !cmd.IsFinished || !cmd.IsCaught || cmd.IsSaved {
		t.Fatalf("Command must be finished and caught, but not saved")
	}
}

func TestInAreaScenario(t *testing.T) {
	behavior := OutOfBoundBehaviorPass
	engine := Engine{
		ConfigRaw:          ScenarioInArea,
		OutOfBoundBehavior: &behavior,
	}

	engine.ParseConfig()
	engine.ProcessCommands()

	if len(engine.Area.Commands) != 1 {
		t.Fatalf("Commands number must be 1")
	}

	cmd := engine.Area.Commands[0]

	if cmd.IsFinished || cmd.IsCaught || cmd.IsSaved {
		t.Fatalf("Command must be not finished, not caught and not saved")
	}
}

func TestSavedScenario(t *testing.T) {
	behavior := OutOfBoundBehaviorPass
	engine := Engine{
		ConfigRaw:          ScenarioSaved,
		OutOfBoundBehavior: &behavior,
	}

	engine.ParseConfig()
	engine.ProcessCommands()

	if len(engine.Area.Commands) != 1 {
		t.Fatalf("Commands number must be 1")
	}

	cmd := engine.Area.Commands[0]

	if !cmd.IsFinished || cmd.IsCaught || !cmd.IsSaved {
		t.Fatalf("Command must be finished and saved, but not caught")
	}
}
