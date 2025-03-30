package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/internal/set"
	"github.com/laurencevs/logic-puzzles/types/intpair"
	"github.com/laurencevs/logic-puzzles/types/inttriple"
)

func main() {
	puzzle1()
	puzzle2()
}

/*
Two numbers are drawn from a giant bingo machine containing 2025 balls
labelled from 1 to 2025. Stifado is told their sum, Pastitsio their product,
and Dolmadakia the difference between them.

Dolmadakia says "I can tell that Stifado doesn't know what the numbers are."

Stifado thinks for a moment, before replying "And I can tell that Pastitsio
doesn't know, either."

Pastitsio nods sadly. "You're quite right... But I can say with certainty that
at least one of the numbers is a side length in a Pythagorean triangle in
which the square of one of the sides is 2025."

"Oh, now I know what they are!" exclaims Stifado.

Pastitsio and Dolmadakia sit in silence. They won't admit it, but neither of
them can figure out what the numbers are.

What are the two numbers?
*/
func puzzle1() {
	solutionSpace := intpair.IntPairs(1, 2025, false, false)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	Stifado := puzzle.NewActorWithKnowledge(intpair.Sum)
	Pastitsio := puzzle.NewActorWithKnowledge(intpair.Product)
	Dolmadakia := puzzle.NewActorWithKnowledge(intpair.AbsDifference)

	Dolmadakia.Says(Dolmadakia.Knows(Stifado.DoesNotKnowAnswer()))
	Stifado.Says(Stifado.Knows(Pastitsio.DoesNotKnowAnswer()))
	Pastitsio.Says(Pastitsio.Knows(intpair.HasNumberIn(sideLengths2025))) // HasOneNumberIn also works
	Stifado.Says(Stifado.KnowsAnswer())

	puzzle.Narrate(Pastitsio.DoesNotKnowAnswer())
	puzzle.Narrate(Dolmadakia.DoesNotKnowAnswer())

	fmt.Println(puzzles.SprintPossibilities(puzzle.ExternalPossibilities())) // (59, 108)
}

/*
Three numbers are drawn from a giant bingo machine containing 2025 balls
labelled from 1 to 2025. Stifado, Pastitsio, and Dolmadakia are each told the
product of a distinct sub-pair of the three numbers. It is announced that the
three numbers sum to 2025.

Stifado: "I can tell that Pastitsio doesn't know what the numbers are."

Pastitsio: "And I can say the same about Dolmadakia."

Dolmadakia: "You're quite right... But I can say with certainty that at least
one of the numbers is a side length in a Pythagorean triangle in which the
square of one of the sides is 2025."

Stifado: "That narrows it down a bit, but I'm still not sure what the numbers
are..."

Pastitsio: "If I told you, you still wouldn't know which of me and Dolmadakia
was told which product!"

What are the numbers?
*/
func puzzle2() {
	solutionSpace := inttriple.IntTriplesWithSumWithoutRepetition(2025, true)
	puzzle := puzzles.NewPuzzle(solutionSpace)

	Stifado := puzzle.NewActorWithKnowledge(inttriple.Pair1Product)
	Pastitsio := puzzle.NewActorWithKnowledge(inttriple.Pair2Product)
	Dolmadakia := puzzle.NewActorWithKnowledge(inttriple.Pair3Product)

	Stifado.Says(Stifado.Knows(Pastitsio.KnowsNormalisedAnswer(inttriple.Normalise).Not()))
	Pastitsio.Says(Pastitsio.Knows(Dolmadakia.KnowsNormalisedAnswer(inttriple.Normalise).Not()))
	Dolmadakia.Says(Dolmadakia.Knows(inttriple.HasNumberIn(sideLengths2025)))
	Stifado.Says(Stifado.KnowsNormalisedAnswer(inttriple.Normalise).Not())
	Pastitsio.Says(Pastitsio.Knows(Stifado.DoesNotKnowAnswerGivenNormalised(inttriple.Normalise)))

	fmt.Println(puzzles.SprintPossibilities(puzzles.NormalisePossibilities(puzzle.ExternalPossibilities(), inttriple.Normalise)))
}

var sideLengths2025 = set.New(
	24, 27, 28, 36, 45, 51, 53, 60, 75, 108, 117, 200, 205, 336, 339, 1012, 1013,
)
