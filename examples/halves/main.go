package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/types/intpair"
)

// Source: https://www.reddit.com/r/puzzles/comments/17vzw6f/can_you_figure_out_as_number_in_your_head_from/

func main() {
	solutionSpace := intpair.IntPairs(1, 30, true, true)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	A := puzzle.NewActorWithKnowledge(intpair.First)
	B := puzzle.NewActorWithKnowledge(intpair.Second)

	B.Says(B.KnowsWhether(SecondIsDoubleFirst).Not())
	A.Says(A.KnowsWhether(FirstIsDoubleSecond).Not())
	B.Says(B.KnowsWhether(FirstIsDoubleSecond).Not())
	A.Says(A.KnowsWhether(SecondIsDoubleFirst).Not())

	fmt.Println("B knows A's number:", puzzle.Evaluate(B.KnowsAnswer()))
	for _, poss := range B.PossibilitiesByKnowledge() {
		if len(poss) > 0 {
			fmt.Println("A's number:", poss[0].A)
			break
		}
	}
}

var (
	FirstIsDoubleSecond puzzles.Condition[intpair.IntPair] = func(p intpair.IntPair) bool {
		return p.A == p.B*2
	}
	SecondIsDoubleFirst puzzles.Condition[intpair.IntPair] = func(p intpair.IntPair) bool {
		return p.B == p.A*2
	}
)
