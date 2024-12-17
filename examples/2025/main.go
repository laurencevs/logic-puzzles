package main

import (
	"fmt"
	"math/rand"

	puzzles "github.com/laurencevs/logic-puzzles"
)

func main() {

	option0()
	return
	option4()
	return

	possibilities := puzzles.UnorderedIntPairs(1, 2025, false)
	puzzle := puzzles.NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter()
	Sophie.KnowsValueOf(puzzles.Sum)
	Paul := puzzle.NewCharacter()
	Paul.KnowsValueOf(puzzles.Product)
	Dave := puzzle.NewCharacter()
	Dave.KnowsValueOf(puzzles.AbsDifference)

	people := []*puzzles.Character[puzzles.IntPair]{Sophie, Paul, Dave}

	for {
		switch rand.Intn(3) {
		case 0:
			fmt.Println("20|Product")
			Paul.Says(puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
		case 1:
			fmt.Println("20|Sum")
			Sophie.Says(puzzle.Satisfies(puzzles.SumIsDivisibleBy(20)))
		case 2:
			fmt.Println("20|Difference")
			Dave.Says(puzzle.Satisfies(puzzles.AbsDifferenceIsDivisibleBy(20)))
		}
		fmt.Println(len(puzzle.ExternalPossibilities()))

		rand.Shuffle(len(people), func(i, j int) { people[i], people[j] = people[j], people[i] })
		fmt.Println(people[0].Name(), people[1].Name(), people[2].Name())
		people[0].Says(people[0].Knows(people[1].Knows(people[2].DoesNotKnowAnswer())))
		fmt.Println(len(puzzle.ExternalPossibilities()))
		people[2].Says(people[2].KnowsAnswer())
		fmt.Println(len(puzzle.ExternalPossibilities()))
		fmt.Println()

		puzzle.Reset()
	}
}

type testPuzzle struct {
	Puzzle  *puzzles.Puzzle[puzzles.IntPair]
	P, S, D *puzzles.Character[puzzles.IntPair]
}

func commonSetup() testPuzzle {
	possibilities := puzzles.UnorderedIntPairs(1, 2025, false)
	puzzle := puzzles.NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter().KnowsValueOf(puzzles.Sum)
	Paul := puzzle.NewCharacter().KnowsValueOf(puzzles.Product)
	Dave := puzzle.NewCharacter().KnowsValueOf(puzzles.AbsDifference)

	return testPuzzle{
		puzzle, Paul, Sophie, Dave,
	}
}

func test206() {
	p := commonSetup()
	p.P.Says(p.P.Knows(p.D.DoesNotKnowAnswer()))
	p.D.Says(p.D.Knows(p.S.DoesNotKnowAnswer()))
	p.S.Says(p.S.Knows(p.Puzzle.Satisfies(puzzles.HasNumberIn(sideLengths2025))))
	p.D.Says(p.D.KnowsAnswer())
	p.Puzzle.PrintPossibilities()
}

func test212() {
	p := commonSetup()
	p.P.Says(p.P.Knows(p.S.DoesNotKnowAnswer()))
	p.S.Says(p.S.Knows(p.D.DoesNotKnowAnswer()))
	p.D.Says(p.D.Knows(p.Puzzle.Satisfies(puzzles.HasNumberIn(sideLengths2025))))
	p.S.Says(p.S.KnowsAnswer())
	p.Puzzle.PrintPossibilities()
}

var sideLengths2025 = map[int]struct{}{
	24:   {},
	27:   {},
	28:   {},
	36:   {},
	45:   {},
	51:   {},
	53:   {},
	60:   {},
	75:   {},
	108:  {},
	117:  {},
	200:  {},
	205:  {},
	336:  {},
	339:  {},
	1012: {},
	1013: {},
}

func test1() {
	p := commonSetup()

	p.D.Says(p.D.Knows(p.S.DoesNotKnowAnswer()))
	p.S.Says(p.S.Knows(p.P.DoesNotKnowAnswer()))
	p.P.Says(p.P.Knows(p.Puzzle.Satisfies(puzzles.HasNumberIn(sideLengths2025)))) // HasOneNumberIn also works
	p.S.Says(p.S.KnowsAnswer())

	p.P.DoesNotKnowAnswer()
	p.D.DoesNotKnowAnswer()

	p.Puzzle.PrintPossibilities()
}

/*
Two numbers are drawn from a giant bingo machine containing 2025 balls labelled from 1 to 2025. Stifado is told their sum, Pastitsio their product, and Dolmadakia the difference between them.

Dolmadakia says "I can tell that Stifado doesn't know what the numbers are."

Stifado thinks for a moment, before replying "And I can tell that Pastitsio doesn't know, either."

Pastitsio nods sadly. "You're quite right... But I can say with certainty that at least one of the numbers is a side length in a Pythagorean triangle in which the square of one of the sides is 2025."

"Oh, now I know what they are!" exclaims Stifado.

Pastitsio and Dolmadakia sit in silence, giving nothing away. In truth, neither of them can figure out what the numbers are.

What are the two numbers?
*/
func option0() {
	possibilities := puzzles.UnorderedIntPairs(1, 2025, false)
	puzzle := puzzles.NewPuzzle(possibilities)

	Stifado := puzzle.NewCharacter().KnowsValueOf(puzzles.Sum)
	Pastitsio := puzzle.NewCharacter().KnowsValueOf(puzzles.Product)
	Dolmadakia := puzzle.NewCharacter().KnowsValueOf(puzzles.AbsDifference)

	fmt.Println(len(puzzle.ExternalPossibilities()))
	Dolmadakia.Says(Dolmadakia.Knows(Stifado.DoesNotKnowAnswer()))
	fmt.Println(len(puzzle.ExternalPossibilities()))
	Stifado.Says(Stifado.Knows(Pastitsio.DoesNotKnowAnswer()))
	fmt.Println(len(puzzle.ExternalPossibilities()))
	Pastitsio.Says(Pastitsio.Knows(puzzle.Satisfies(puzzles.HasNumberIn(sideLengths2025)))) // HasOneNumberIn also works
	fmt.Println(len(puzzle.ExternalPossibilities()))
	Stifado.Says(Stifado.KnowsAnswer())
	fmt.Println(len(puzzle.ExternalPossibilities()))

	Pastitsio.DoesNotKnowAnswer()
	Dolmadakia.DoesNotKnowAnswer()

	puzzle.PrintPossibilities() // 59, 108
}

func option1() {
	p := commonSetup()
	p.P.Says(p.Puzzle.Satisfies(puzzles.AbsDifferenceIsDivisibleBy(20)))
	p.P.Says(p.P.Knows(p.D.Knows(p.S.DoesNotKnowAnswer())))
	p.S.Says(p.S.KnowsAnswer())
	p.Puzzle.PrintPossibilities()
}

func option2() {
	p := commonSetup()
	p.S.Says(p.Puzzle.Satisfies(puzzles.SumIsDivisibleBy(20)))
	p.P.Says(p.P.Knows(p.S.Knows(p.D.DoesNotKnowAnswer())))
	p.D.Says(p.D.KnowsAnswer())
	p.Puzzle.PrintPossibilities()
}

func option3() {
	p := commonSetup()
	p.S.Says(p.Puzzle.Satisfies(puzzles.SumIsDivisibleBy(20)))
	p.S.Says(p.S.Knows(p.P.Knows(p.D.DoesNotKnowAnswer())))
	p.D.Says(p.D.KnowsAnswer())
	p.Puzzle.PrintPossibilities()
}

func option4() {
	p := commonSetup()
	p.P.Says(p.P.Knows(p.Puzzle.Satisfies(puzzles.SumIsDivisibleBy(20))))
	p.S.Says(p.S.Knows(p.P.Knows(p.D.DoesNotKnowAnswer())))
	p.D.Says(p.D.KnowsAnswer())
	p.D.Says(p.Puzzle.Satisfies(puzzles.SumIsDivisibleBy(25)))
	p.Puzzle.PrintPossibilities() // 101, 1999
}

func option4NonTrivial() {
	p := commonSetup()
	p.P.Says(p.P.Knows(p.Puzzle.Satisfies(puzzles.SumIsDivisibleBy(100))))
	p.S.Says(p.S.Knows(p.P.Knows(p.D.DoesNotKnowAnswer())))
	p.D.Says(p.D.KnowsAnswer())
	fmt.Println(len(p.Puzzle.ExternalPossibilities()))
}
