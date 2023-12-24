package main

import (
	"fmt"
	"math/rand"
)

func generate2024Puzzles() {
	p := commonSetup2024(false)

	for {
		fmt.Print(".")
		statements := []string{}

		p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(20)))
		p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(24)))
		p.Sophie.Says(p.Sophie.Knows(p.Paul.Knows(p.Dave.DoesNotKnowAnswer())))
		var lastLen int
		for len(p.puzzle.externalPossibilities) > 1 {
			lastLen = len(p.puzzle.externalPossibilities)
			statements = append(statements, randomStatement(p.puzzle.characters, rand.Intn(4)))
		}
		fmt.Print(lastLen)
		if len(p.puzzle.externalPossibilities) == 1 || lastLen <= 3 {
			fmt.Println()
			fmt.Println(statements)
			p.puzzle.PrintPossibilities()
		}
		p.puzzle.Reset()
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

func commonSetup2024(allowRepeated bool) *puzzle2024[intPair] {
	possibilities := UnorderedIntPairs(1, 2024, allowRepeated)
	puzzle := NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter("S")
	Sophie.KnowsValueOf(sum)
	Paul := puzzle.NewCharacter("P")
	Paul.KnowsValueOf(product)
	Dave := puzzle.NewCharacter("D")
	Dave.KnowsValueOf(absDifference)

	return &puzzle2024[intPair]{
		puzzle: puzzle,
		Sophie: Sophie,
		Paul:   Paul,
		Dave:   Dave,
	}
}

func run2024Puzzles() {
	puzzle := commonSetup2024(false)
	example2024_1(puzzle)
	example2024_2(puzzle)
	example2024_3(puzzle)
	example2024_4(puzzle)
	example2024_6(puzzle)
	example2024_7(puzzle)
	puzzle = commonSetup2024(true)
	example2024_5(puzzle)
}

type puzzle2024[P PuzzlePossibility] struct {
	puzzle *Puzzle[P]
	Sophie *Character[P]
	Paul   *Character[P]
	Dave   *Character[P]
}

func example2024_1(p *puzzle2024[intPair]) {
	p.Paul.DoesNotKnowAnswer()
	p.Paul.Says(p.Paul.Knows(p.puzzle.Satisfies(hasNumberDivisibleBy(20))))
	p.Sophie.DoesNotKnowAnswer()
	p.Sophie.Says(p.Sophie.Knows(p.puzzle.Satisfies(hasNumberDivisibleBy(24))))
	p.Dave.Says(p.Dave.DoesNotKnowAnswer())

	p.puzzle.PrintPossibilities() // (256, 480)
	p.puzzle.Reset()
}

func example2024_2(p *puzzle2024[intPair]) {
	p.Paul.Says(p.Paul.DoesNotKnowAnswer())
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.DoesNotKnowAnswer())
	p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}

func example2024_3(p *puzzle2024[intPair]) {
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}

func example2024_4(p *puzzle2024[intPair]) {
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.Paul.Says(p.Paul.KnowsAnswer())
	p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (1046, 1090)
	p.puzzle.Reset()
}

func example2024_5(p *puzzle2024[intPair]) {
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
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.Dave.Says(p.Dave.Knows(p.Sophie.Knows(p.Paul.DoesNotKnowAnswer())))
	p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.Knows(p.Dave.KnowsAnswer())))

	p.puzzle.PrintPossibilities() // (20, 2020)
	p.puzzle.Reset()
}

func example2024_6(p *puzzle2024[intPair]) {
	/*
		The numbers from 1 to 2024 are put into a hat, and two are drawn at
		random. Paul is told their product, Sophie their sum, and Dave their
		difference.
		Paul says "The product of the numbers is divisible by 20."
		Sophie notices that the sum of the numbers is divisible by 24, but
		doesn't say anything.
		Dave replies to Paul, "Then I know that you know that Sophie doesn't
		  know what the numbers are."
		Sophie interjects, "Well now I do!"
		What are the numbers?
	*/
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.puzzle.Satisfies(sumIsDivisibleBy(24))
	p.Dave.Says(p.Dave.Knows(p.Paul.Knows(p.Sophie.DoesNotKnowAnswer())))
	p.Sophie.Says(p.Sophie.KnowsAnswer())

	p.puzzle.PrintPossibilities() // (2010, 2022)
	p.puzzle.Reset()
}

func example2024_7(p *puzzle2024[intPair]) {
	p.Paul.Says(p.puzzle.Satisfies(productIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.Knows(p.Paul.Knows(p.Dave.DoesNotKnowAnswer())))
	p.Sophie.Says(p.puzzle.Satisfies(sumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Dave.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}
