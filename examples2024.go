package main

import (
	"fmt"
	"math/rand"
)

func generate2024Puzzles() {
	puzzle, Sophie, Paul, Dave := commonSetup2024(false)

	for {
		fmt.Print(".")
		statements := []string{}

		Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(20)))
		Paul.Says(puzzle.Satisfies(productIsDivisibleBy(24)))
		Sophie.Says(Sophie.Knows(Paul.Knows(Dave.DoesNotKnowAnswer())))
		var lastLen int
		for len(puzzle.externalPossibilities) > 1 {
			lastLen = len(puzzle.externalPossibilities)
			statements = append(statements, randomStatement(puzzle.characters, rand.Intn(4)))
		}
		fmt.Print(lastLen)
		if len(puzzle.externalPossibilities) == 1 || lastLen <= 3 {
			fmt.Println()
			fmt.Println(statements)
			puzzle.PrintPossibilities()
		}
		puzzle.Reset()
	}
}

func shuffleSlice[T any](s []T) {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

func randomStatement[P PuzzlePossibility](characters []*Character[P], n int) string {
	s := []int{0, 1, 2}
	shuffleSlice(s)
	c1 := characters[s[0]]
	c2 := characters[s[1]]
	c3 := characters[s[2]]
	switch n {
	case 0:
		c1.Says(c1.Knows(c2.Knows(c3.DoesNotKnowAnswer())))
	case 1:
		c1.Says(c1.Knows(c2.Knows(c3.KnowsAnswer())))
	case 2:
		c1.Says(c1.Knows(c2.Knows(c3.Knows(c1.KnowsAnswer()))))
	case 3:
		c1.Says(c1.Knows(c2.DoesNotKnowAnswer()))
	}
	return fmt.Sprintf("%d %s %s %s; ", n, c1.name, c2.name, c3.name)
}

func commonSetup2024(allowRepeated bool) (*Puzzle[intPair], *Character[intPair], *Character[intPair], *Character[intPair]) {
	possibilities := UnorderedIntPairs(1, 2024, allowRepeated)
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
	puzzle, Sophie, Paul, Dave := commonSetup2024(false)
	puzzle1(puzzle, Sophie, Paul, Dave)
	puzzle2(puzzle, Sophie, Paul, Dave)
	puzzle3(puzzle, Sophie, Paul, Dave)
	puzzle4(puzzle, Sophie, Paul, Dave)
	puzzle6(puzzle, Sophie, Paul, Dave)
	puzzle7(puzzle, Sophie, Paul, Dave)
	puzzle, Sophie, Paul, Dave = commonSetup2024(true)
	puzzle5(puzzle, Sophie, Paul, Dave)
}

func puzzle1(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.DoesNotKnowAnswer()
	Paul.Says(Paul.Knows(puzzle.Satisfies(hasNumberDivisibleBy(20))))
	Sophie.DoesNotKnowAnswer()
	Sophie.Says(Sophie.Knows(puzzle.Satisfies(hasNumberDivisibleBy(24))))
	Dave.Says(Dave.DoesNotKnowAnswer())

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle2(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(Paul.DoesNotKnowAnswer())
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Sophie.Says(Sophie.DoesNotKnowAnswer())
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Dave.Says(Dave.Knows(Paul.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle3(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Dave.Says(Dave.Knows(Paul.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle4(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Paul.Says(Paul.KnowsAnswer())
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Paul.Says(Paul.Knows(Sophie.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle5(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	/*
		Two numbers from 1 to 2024 are randomly generated. (They could be the
		same.) Tristram is told their product, Walter their sum, and Toby their
		difference.

		Tristram says "The product of the numbers is divisible by 20"
		Toby replies "I can tell that Walter knows that you don't know what the
		  numbers are"
		Walter says "If it helps, the sum of the numbers is divisible by 24"
		Tristram replies, "Ah, now I know that Walter knows that Toby knows what
		  the numbers are"

		What are the numbers?
	*/
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Dave.Says(Dave.Knows(Sophie.Knows(Paul.DoesNotKnowAnswer())))
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Paul.Says(Paul.Knows(Sophie.Knows(Dave.KnowsAnswer())))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle6(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	/*
		The numbers from 1 to 2024 are put into a hat, and two are drawn at
		random. Paul is told their product, Sophie their sum, and Dave their
		difference.
		Paul says "The product of the numbers is divisible by 20."
		Sophie notices that the sum of the numbers is divisible by 24, but
		doesn't tell anyone.
		Dave replies, "Then I know that you know that Sophie doesn't know what
		  the numbers are."
		Sophie interjects, "Well now I do!"
		What are the numbers?
	*/
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Dave.Says(Dave.Knows(Paul.Knows(Sophie.DoesNotKnowAnswer())))
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Sophie.Says(Sophie.KnowsAnswer())

	puzzle.PrintPossibilities()
	puzzle.Reset()
}

func puzzle7(puzzle *Puzzle[intPair], Sophie, Paul, Dave *Character[intPair]) {
	Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
	Sophie.Says(Sophie.Knows(Paul.Knows(Dave.DoesNotKnowAnswer())))
	Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
	Paul.Says(Paul.Knows(Dave.KnowsAnswer()))

	puzzle.PrintPossibilities()
	puzzle.Reset()
}
