package puzzles_test

import (
	"cmp"
	"slices"
	"testing"

	puzzles "github.com/laurencevs/logic-puzzles"
)

func assertActorPossibilities[P cmp.Ordered](t *testing.T, puzzle *puzzles.Puzzle[P], expected []P) {
	for i, char := range puzzle.Actors() {
		charPossibilities := char.Possibilities()
		slices.Sort(charPossibilities)
		if !slices.Equal(charPossibilities, expected) {
			t.Errorf("char%d's possibilities %v not equal to expected %v", i, charPossibilities, expected)
		}
	}
	if !slices.Equal(puzzle.ExternalPossibilities(), expected) {
		t.Errorf("external possibilities %v not equal to expected %v", puzzle.ExternalPossibilities(), expected)
	}
}

func TestActorPossibilities(t *testing.T) {
	solutionSpace := []string{
		"11", "12", "13",
		"21", "22", "23",
		"31" /* */, "33",
	}
	puzzle := puzzles.NewPuzzle(solutionSpace)

	char0 := puzzles.NewActorWithKnowledge(puzzle, func(s string) byte { return s[0] })
	char1 := puzzles.NewActorWithKnowledge(puzzle, func(s string) int {
		total := 0
		for _, c := range s {
			total += int(c - '0')
		}
		return total
	})
	char2 := puzzles.NewActorWithKnowledge(puzzle, func(s string) string { return s })

	expectedActors := []*puzzles.Actor[string]{char0, char1, char2}
	if !slices.Equal(puzzle.Actors(), expectedActors) {
		t.Errorf("puzzle.Actors() %v not equal to expected %v", puzzle.Actors(), expectedActors)
	}

	assertActorPossibilities(t, puzzle, puzzle.SolutionSpace())

	char2.Says(char2.Knows(char1.DoesNotKnowAnswer())) // rules out 11, 23, 33
	expectedPossibilities := []string{
		/* */ "12", "13",
		"21", "22", /* */
		"31", /*       */
	}
	assertActorPossibilities(t, puzzle, expectedPossibilities)

	char0.Says(char0.KnowsAnswer()) // solution must be 31
	expectedPossibilities = []string{"31"}
	assertActorPossibilities(t, puzzle, expectedPossibilities)

	for i, char := range puzzle.Actors() {
		if !puzzle.Evaluate(char.KnowsAnswer()) {
			t.Errorf("char%d should know the answer, but doesn't", i)
		}
	}
}

func assertSamePossibilitiesInAllCases[P cmp.Ordered](t *testing.T, a, b *puzzles.Actor[P], solutionSpace []P) {
	for _, s := range solutionSpace {
		aPossibilities := a.PossibilitiesInCase(s)
		bPossibilities := b.PossibilitiesInCase(s)
		slices.Sort(aPossibilities)
		slices.Sort(bPossibilities)
		if !slices.Equal(aPossibilities, bPossibilities) {
			t.Errorf("different possibilities for actors %d, %d in case %v", a.Id, b.Id, s)
		}
	}
}

func TestReusedKnowledge(t *testing.T) {
	solutionSpace := []string{
		"11", "12", "13",
		"21", "22", "23",
		"31" /* */, "33",
	}
	puzzle := puzzles.NewPuzzle(solutionSpace)

	char0 := puzzles.NewActorWithKnowledge(puzzle, func(s string) byte { return s[0] })
	char1 := puzzle.NewActor()
	char1.HasKnowledge(char0.KnowledgeState)

	char2Knowledge := puzzles.NewKnowledge(puzzle, func(s string) int {
		total := 0
		for _, c := range s {
			total += int(c - '0')
		}
		return total
	})
	char2 := puzzle.NewActor()
	char2.KnowledgeState = char2Knowledge
	char3 := puzzle.NewActor()
	char3.KnowledgeState = char2Knowledge

	expectedActors := []*puzzles.Actor[string]{char0, char1, char2, char3}
	if !slices.Equal(puzzle.Actors(), expectedActors) {
		t.Errorf("puzzle.Actors() %v not equal to expected %v", puzzle.Actors(), expectedActors)
	}
	assertActorPossibilities(t, puzzle, puzzle.SolutionSpace())

	char0.Says(char0.Knows(char1.DoesNotKnowAnswer())) // already true in all cases
	assertActorPossibilities(t, puzzle, puzzle.SolutionSpace())
	assertSamePossibilitiesInAllCases(t, char0, char1, solutionSpace)
	assertSamePossibilitiesInAllCases(t, char2, char3, solutionSpace)

	char3.Says(char3.DoesNotKnowAnswer()) // rules out 11, 23, 33
	expectedPossibilities := []string{
		/* */ "12", "13",
		"21", "22", /* */
		"31", /*       */
	}
	assertActorPossibilities(t, puzzle, expectedPossibilities)
	assertSamePossibilitiesInAllCases(t, char0, char1, solutionSpace)
	assertSamePossibilitiesInAllCases(t, char2, char3, solutionSpace)
}
