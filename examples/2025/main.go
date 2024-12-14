package main

import (
	"fmt"
	"math/rand"

	puzzles "github.com/laurencevs/logic-puzzles"
)

func main() {

	option0()
	option4()
	return

	possibilities := puzzles.UnorderedIntPairs(1, 2025, false)
	puzzle := puzzles.NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter("S")
	Sophie.KnowsValueOf(puzzles.Sum)
	Paul := puzzle.NewCharacter("P")
	Paul.KnowsValueOf(puzzles.Product)
	Dave := puzzle.NewCharacter("D")
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

	Sophie := puzzle.NewCharacter("S")
	Sophie.KnowsValueOf(puzzles.Sum)
	Paul := puzzle.NewCharacter("P")
	Paul.KnowsValueOf(puzzles.Product)
	Dave := puzzle.NewCharacter("D")
	Dave.KnowsValueOf(puzzles.AbsDifference)

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

func option0() {
	p := commonSetup()

	p.D.Says(p.D.Knows(p.S.DoesNotKnowAnswer()))
	p.S.Says(p.S.Knows(p.P.DoesNotKnowAnswer()))
	p.P.Says(p.P.Knows(p.Puzzle.Satisfies(puzzles.HasNumberIn(sideLengths2025)))) // HasOneNumberIn also works
	p.S.Says(p.S.KnowsAnswer())

	p.P.DoesNotKnowAnswer()
	p.D.DoesNotKnowAnswer()

	p.Puzzle.PrintPossibilities() // 59, 108
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
