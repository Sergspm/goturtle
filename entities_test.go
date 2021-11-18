package main

import "testing"

func TestTurtleDirection(t *testing.T) {
	if !isTurtleDirectionCorrect(TurtleDirectionNorth) {
		t.Fatalf("North direction must be correct")
	}
	if !isTurtleDirectionCorrect(TurtleDirectionWest) {
		t.Fatalf("West direction must be correct")
	}
	if !isTurtleDirectionCorrect(TurtleDirectionSouth) {
		t.Fatalf("South direction must be correct")
	}
	if !isTurtleDirectionCorrect(TurtleDirectionEast) {
		t.Fatalf("East direction must be correct")
	}
	if isTurtleDirectionCorrect("non-existed direction") {
		t.Fatalf("Not existing direction must not be correct")
	}
}

func TestTurtleCommand(t *testing.T) {
	if !isTurtleCommandCorrect(TurtleCommandMove) {
		t.Fatalf("Move command must be correct")
	}
	if !isTurtleCommandCorrect(TurtleCommandRotateRight) {
		t.Fatalf("Rotate right command must be correct")
	}
	if !isTurtleCommandCorrect(TurtleCommandRotateLeft) {
		t.Fatalf("Rotate left command must be correct")
	}
	if isTurtleCommandCorrect("non-existed command") {
		t.Fatalf("Not existed command must not be correct")
	}
}

func TestOutOfBoundBehavior(t *testing.T) {
	if !isOutOfBoundBehaviorCorrect(OutOfBoundBehaviorPass) {
		t.Fatalf("Pass out of bound behavior must be correct")
	}
	if !isOutOfBoundBehaviorCorrect(OutOfBoundBehaviorStop) {
		t.Fatalf("Stop out of bound behavior must be correct")
	}
	if !isOutOfBoundBehaviorCorrect(OutOfBoundBehaviorPortal) {
		t.Fatalf("Portal out of bound behavior must be correct")
	}
	if isOutOfBoundBehaviorCorrect("non-existed out of bound behavior") {
		t.Fatalf("Not existed out of bound behavior must not be correct")
	}
}

func TestEntityTypeCheck(t *testing.T) {
	entity := Entity{Type: EntityTypeTrap}
	if !entity.IsTrap() {
		t.Fatalf("Entity type must be detected as trap")
	}
	entity.Type = EntityTypeRefuge
	if !entity.IsRefuge() {
		t.Fatalf("Entity type must be detected as refuge")
	}
	entity.Type = "test"
	if entity.IsTrap() || entity.IsRefuge() {
		t.Fatalf("Entity type must be not detected as trap or refuge")
	}
}