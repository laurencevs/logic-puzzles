package main

import "fmt"

func generate2024Puzzles() {
	puzzle, Sophie, Paul, _ := commonSetup2024()

	for {
		fmt.Print(".")
		statements := []string{}

		Paul.Says(puzzle.AnswerSatisfies(productIsDivisibleBy(20)))
		statements = append(statements, randomStatement(puzzle.characters))
		Sophie.Says(puzzle.AnswerSatisfies(sumIsDivisibleBy(24)))
		var lastLen int
		for len(puzzle.externalPossibilities) > 1 {
			lastLen = len(puzzle.externalPossibilities)
			statements = append(statements, randomStatement(puzzle.characters))
		}
		fmt.Print(lastLen)
		if len(puzzle.externalPossibilities) == 1 {
			fmt.Println()
			fmt.Println(statements)
			puzzle.PrintPossibilities()
		}
		puzzle.Reset()
	}
}

func commonSetup2024() (*Puzzle[intPair], *Character[intPair], *Character[intPair], *Character[intPair]) {
	possibilities := UnorderedIntPairs(1, 2024, true)
	puzzle := NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter("S")
	Sophie.KnowsValueOf(sum)
	Paul := puzzle.NewCharacter("P")
	Paul.KnowsValueOf(product)
	Dave := puzzle.NewCharacter("D")
	Dave.KnowsValueOf(absDifference)

	return puzzle, Sophie, Paul, Dave
}

func run2024Puzzles() {
	puzzle, Sophie, Paul, Dave := commonSetup2024()
	puzzle1(puzzle, Sophie, Paul, Dave)
	puzzle2(puzzle, Sophie, Paul, Dave)
	puzzle3(puzzle, Sophie, Paul, Dave)
	puzzle4(puzzle, Sophie, Paul, Dave)
	puzzle5(puzzle, Sophie, Paul, Dave)
}

func puzzle1(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.DoesNotKnowAnswer()
	Paul.Says(Paul.Knows(puzzle.AnswerSatisfies(hasNumberDivisibleBy(20))))
	Sophie.DoesNotKnowAnswer()
	Sophie.Says(Sophie.Knows(puzzle.AnswerSatisfies(hasNumberDivisibleBy(24))))
	Dave.Says(Dave.DoesNotKnowAnswer())

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle2(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(Paul.DoesNotKnowAnswer())
	Paul.Says(puzzle.AnswerSatisfies(productIsDivisibleBy(20)))
	Sophie.Says(Sophie.DoesNotKnowAnswer())
	Sophie.Says(puzzle.AnswerSatisfies(sumIsDivisibleBy(24)))
	Dave.Says(Dave.Knows(Paul.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle3(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.AnswerSatisfies(productIsDivisibleBy(20)))
	Sophie.Says(puzzle.AnswerSatisfies(sumIsDivisibleBy(24)))
	Dave.Says(Dave.Knows(Paul.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle4(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.AnswerSatisfies(productIsDivisibleBy(20)))
	Paul.Says(Paul.KnowsAnswer())
	Sophie.Says(puzzle.AnswerSatisfies(sumIsDivisibleBy(24)))
	Paul.Says(Paul.Knows(Sophie.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle5(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.AnswerSatisfies(productIsDivisibleBy(20)))
	Dave.Says(Dave.Knows(Sophie.Knows(Paul.DoesNotKnowAnswer())))
	Sophie.Says(puzzle.AnswerSatisfies(sumIsDivisibleBy(24)))
	Paul.Says(Paul.Knows(Sophie.Knows(Dave.KnowsAnswer())))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}
