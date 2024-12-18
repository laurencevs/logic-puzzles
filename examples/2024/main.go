package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
	"github.com/laurencevs/logic-puzzles/types/intpair"
)

func main() {
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

func commonSetup2024(allowRepeated bool) *puzzle2024[intpair.IntPair] {
	possibilities := intpair.IntPairs(1, 2024, false, allowRepeated)
	puzzle := puzzles.NewPuzzle(possibilities)

	Sophie := puzzle.NewActorWithKnowledge(intpair.Sum)
	Paul := puzzle.NewActorWithKnowledge(intpair.Product)
	Dave := puzzle.NewActorWithKnowledge(intpair.AbsDifference)

	return &puzzle2024[intpair.IntPair]{
		puzzle: puzzle,
		Sophie: Sophie,
		Paul:   Paul,
		Dave:   Dave,
	}
}

type puzzle2024[P comparable] struct {
	puzzle *puzzles.Puzzle[P]
	Sophie *puzzles.Actor[P]
	Paul   *puzzles.Actor[P]
	Dave   *puzzles.Actor[P]
}

func example2024_1(p *puzzle2024[intpair.IntPair]) {
	p.puzzle.Narrate(p.Paul.DoesNotKnowAnswer())
	p.Paul.Says(p.Paul.KnowsHolds(intpair.HasNumberDivisibleBy(20)))
	p.puzzle.Narrate(p.Sophie.DoesNotKnowAnswer())
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.HasNumberDivisibleBy(24)))
	p.Dave.Says(p.Dave.DoesNotKnowAnswer())

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (256, 480)
	p.puzzle.Reset()
}

func example2024_2(p *puzzle2024[intpair.IntPair]) {
	p.Paul.Says(p.Paul.DoesNotKnowAnswer())
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.DoesNotKnowAnswer())
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (10, 1982)
	p.puzzle.Reset()
}

func example2024_3(p *puzzle2024[intpair.IntPair]) {
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	p.Dave.Says(p.Dave.Knows(p.Paul.KnowsAnswer()))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (10, 1982)
	p.puzzle.Reset()
}

func example2024_4(p *puzzle2024[intpair.IntPair]) {
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Paul.Says(p.Paul.KnowsAnswer())
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.KnowsAnswer()))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (1046, 1090)
	p.puzzle.Reset()
}

func example2024_5(p *puzzle2024[intpair.IntPair]) {
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
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Dave.Says(p.Dave.Knows(p.Sophie.Knows(p.Paul.DoesNotKnowAnswer())))
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Sophie.Knows(p.Dave.KnowsAnswer())))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (20, 2020)
	p.puzzle.Reset()
}

func example2024_6(p *puzzle2024[intpair.IntPair]) {
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
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Dave.Says(p.Dave.Knows(p.Paul.Knows(p.Sophie.DoesNotKnowAnswer())))
	p.Sophie.Says(p.Sophie.KnowsAnswer())
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (2010, 2022)
	p.puzzle.Reset()
}

func example2024_7(p *puzzle2024[intpair.IntPair]) {
	p.Paul.Says(p.Paul.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
	p.Sophie.Says(p.Sophie.Knows(p.Paul.Knows(p.Dave.DoesNotKnowAnswer())))
	p.Sophie.Says(p.Sophie.KnowsHolds(intpair.SumIsDivisibleBy(24)))
	p.Paul.Says(p.Paul.Knows(p.Dave.KnowsAnswer()))

	fmt.Println(puzzles.SprintPossibilities(p.puzzle.ExternalPossibilities())) // (10, 1982)
	p.puzzle.Reset()
}
