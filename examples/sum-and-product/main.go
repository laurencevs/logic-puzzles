package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/types/intpair"
)

// Source: https://en.wikipedia.org/wiki/Sum_and_Product_Puzzle

func main() {
	solutionSpace := intpair.IntPairs(2, 100, false, false)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	S := puzzles.NewActorWithKnowledge(puzzle, intpair.Sum)
	P := puzzles.NewActorWithKnowledge(puzzle, intpair.Product)

	S.Says(S.Knows(P.DoesNotKnowAnswer()))
	P.Says(P.KnowsAnswer())
	S.Says(S.KnowsAnswer())

	fmt.Println(puzzles.SprintPossibilities(puzzle.ExternalPossibilities()))
}
