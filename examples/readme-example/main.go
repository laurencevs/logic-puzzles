package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/types/intpair"
)

func main() {
	solutionSpace := intpair.IntPairs(1, 2024, false, false)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	Stifado := puzzle.NewActorWithKnowledge(intpair.Sum)
	Pastitsio := puzzle.NewActorWithKnowledge(intpair.Product)
	Dolmadakia := puzzle.NewActorWithKnowledge(intpair.AbsDifference)

	Pastitsio.Says(Pastitsio.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	Stifado.Says(Stifado.Knows(Pastitsio.Knows(Dolmadakia.DoesNotKnowAnswer())))
	Stifado.Says(Stifado.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	Pastitsio.Says(Pastitsio.Knows(Dolmadakia.KnowsAnswer()))

	fmt.Println(puzzles.SprintPossibilities(puzzle.ExternalPossibilities()))
}
