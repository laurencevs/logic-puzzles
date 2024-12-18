package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/types/intpair"
)

// Source: https://math.stackexchange.com/questions/2534847/a-deduction-question-for-perfect-logicians

func main() {
	solutionSpace := intpair.IntPairs(1, 9, false, true)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	A := puzzle.NewActorWithKnowledge(intpair.Product)
	B := puzzle.NewActorWithKnowledge(intpair.Sum)

	A.Says(A.DoesNotKnowAnswer())
	B.Says(B.DoesNotKnowAnswer())
	A.Says(A.DoesNotKnowAnswer())
	B.Says(B.DoesNotKnowAnswer())
	A.Says(A.DoesNotKnowAnswer())
	B.Says(B.DoesNotKnowAnswer())
	A.Says(A.DoesNotKnowAnswer())
	B.Says(B.DoesNotKnowAnswer())
	A.Says(A.KnowsAnswer())

	fmt.Println(puzzles.SprintPossibilities(puzzle.ExternalPossibilities()))
}
