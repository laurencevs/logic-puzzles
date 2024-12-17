package main

import (
	"fmt"
	"math/rand"

	puzzles "github.com/laurencevs/logic-puzzles"
)

func generate2024Puzzles() {
	p := commonSetup2024(false)

	for {
		fmt.Print(".")
		statements := []string{}

		p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(20)))
		p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(24)))
		p.Sophie.Says(p.Sophie.Knows(p.Paul.Knows(p.Dave.DoesNotKnowAnswer())))
		var lastLen int
		for len(p.puzzle.ExternalPossibilities()) > 1 {
			lastLen = len(p.puzzle.ExternalPossibilities())
			statements = append(statements, randomStatement(p.puzzle.Characters(), rand.Intn(4)))
		}
		fmt.Print(lastLen)
		if len(p.puzzle.ExternalPossibilities()) == 1 || lastLen <= 3 {
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

func randomStatement[P puzzles.PuzzlePossibility](characters []*puzzles.Character[P], n int) string {
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
	return fmt.Sprintf("%d %s %s %s; ", n, c1.Name(), c2.Name(), c3.Name())
}

func commonSetup2024(allowRepeated bool) *puzzle2024[puzzles.IntPair] {
	possibilities := puzzles.UnorderedIntPairs(1, 2024, allowRepeated)
	puzzle := puzzles.NewPuzzle(possibilities)

	Sophie := puzzle.NewCharacter().KnowsValueOf(puzzles.Sum)
	Paul := puzzle.NewCharacter().KnowsValueOf(puzzles.Product)
	Dave := puzzle.NewCharacter().KnowsValueOf(puzzles.AbsDifference)

	return &puzzle2024[puzzles.IntPair]{
		puzzle: puzzle,
		Sophie: Sophie,
		Paul:   Paul,
		Dave:   Dave,
	}
}

func run2024Puzzles() {
	puzzle := commonSetup2024(false)
	// example2024_1(puzzle)
	// example2024_2(puzzle)
	// example2024_3(puzzle)
	// example2024_4(puzzle)
	example2024_6(puzzle)
	// example2024_7(puzzle)
	// puzzle = commonSetup2024(true)
	// example2024_5(puzzle)
}

type puzzle2024[P puzzles.PuzzlePossibility] struct {
	puzzle *puzzles.Puzzle[P]
	Sophie *puzzles.Character[P]
	Paul   *puzzles.Character[P]
	Dave   *puzzles.Character[P]
}

func example2024_1(p *puzzle2024[puzzles.IntPair]) {
	p.Paul.DoesNotKnowAnswer()
	p.Paul.Says(p.Paul.Knows(p.puzzle.Satisfies(puzzles.HasNumberDivisibleBy(20))))
	p.Sophie.DoesNotKnowAnswer()
	p.Sophie.Says(p.Sophie.Knows(p.puzzle.Satisfies(puzzles.HasNumberDivisibleBy(24))))
	p.Dave.Says(p.Dave.DoesNotKnowAnswer())

	p.puzzle.PrintPossibilities() // (256, 480)
	p.puzzle.Reset()
}

func example2024_2(p *puzzle2024[puzzles.IntPair]) {
	p.Paul.Says(p.Paul.DoesNotKnowAnswer())
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.DoesNotKnowAnswer())
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}

func example2024_3(p *puzzle2024[puzzles.IntPair]) {
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}

func example2024_4(p *puzzle2024[puzzles.IntPair]) {
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Paul.Says(p.Paul.KnowsAnswer())
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (1046, 1090)
	p.puzzle.Reset()
}

func example2024_5(p *puzzle2024[puzzles.IntPair]) {
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
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Dave.Says(p.Dave.Knows(p.Sophie.Knows(p.Paul.DoesNotKnowAnswer())))
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.Knows(p.Dave.KnowsAnswer())))

	p.puzzle.PrintPossibilities() // (20, 2020)
	p.puzzle.Reset()
}

func example2024_6(p *puzzle2024[puzzles.IntPair]) {
	/*
		The integers from 1 to 2024 are written on pieces of paper and put into
		a hat, and two are drawn at random. Paul is told their product, Sophie
		their sum, and Dave the difference between them.
		Paul says "The product of the two numbers is divisible by 20."
		Dave replies, "In that case, I know that you know that Sophie doesn't
		  know what the numbers are."
		Sophie interjects, "Well now I do! ... If it helps, the sum of the
		  numbers is divisible by 24"
		What are the two numbers?
	*/
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Dave.Says(p.Dave.Knows(p.Paul.Knows(p.Sophie.DoesNotKnowAnswer())))
	p.Sophie.Says(p.Sophie.KnowsAnswer())
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))

	p.puzzle.PrintPossibilities() // (2010, 2022)
	p.puzzle.Reset()
}

func example2024_7(p *puzzle2024[puzzles.IntPair]) {
	p.Paul.Says(p.puzzle.Satisfies(puzzles.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.Knows(p.Paul.Knows(p.Dave.DoesNotKnowAnswer())))
	p.Sophie.Says(p.puzzle.Satisfies(puzzles.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Dave.KnowsAnswer()))

	p.puzzle.PrintPossibilities() // (10, 1982)
	p.puzzle.Reset()
}
