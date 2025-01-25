package main

import (
	"fmt"

	puzzles "github.com/laurencevs/logic-puzzles"
)

/*
This example demonstrates a solution to the 'aces and eights' puzzle from
Exercise 1.1 of the book 'Reasoning About Knowledge' by Fagin, Halpern, Moses,
and Vardi (MIT Press, 1995).

https://mitpress.mit.edu/9780262562003/reasoning-about-knowledge/
*/

func main() {
	solutionSpace := possibleStates()
	puzzle := puzzles.NewPuzzle(solutionSpace)

	p1valuation := puzzles.ValuationFromFunc(solutionSpace, p1Knowledge)
	player1 := puzzle.NewActorWithKnowledge(p1valuation)
	p2valuation := puzzles.ValuationFromFunc(solutionSpace, p2Knowledge)
	player2 := puzzle.NewActorWithKnowledge(p2valuation)
	p3valuation := puzzles.ValuationFromFunc(solutionSpace, p3Knowledge)
	player3 := puzzle.NewActorWithKnowledge(p3valuation)

	// Exercise 1.1(a)
	player1.Says(player1.DoesNotKnowAnswer())
	player2.Says(player2.DoesNotKnowAnswer())
	fmt.Println(puzzles.SprintPossibilities(player3.PossibilitiesByKnowledge()[p3valuation(state{
		p1a: 2,
		p2a: 0,
	})]))
	puzzle.Reset()

	// Exercise 1.1(b)
	player1.Says(player1.DoesNotKnowAnswer())
	player2.Says(player2.DoesNotKnowAnswer())
	player3.Says(player3.DoesNotKnowAnswer())
	fmt.Println(puzzles.SprintPossibilities(player1.PossibilitiesByKnowledge()[p1valuation(state{
		p2a: 0,
		p3a: 1,
	})]))
	puzzle.Reset()

	// Exercise 1.1(c)
	player1.Says(player1.DoesNotKnowAnswer())
	player2.Says(player2.DoesNotKnowAnswer())
	player3.Says(player3.DoesNotKnowAnswer())
	player1.Says(player1.DoesNotKnowAnswer())
	fmt.Println(puzzles.SprintPossibilities(player2.PossibilitiesByKnowledge()[p2valuation(state{
		p1a: 1,
		p3a: 1,
	})]))
}

// The game state is represented by the number of aces held by each character.
type state struct {
	p1a, p2a, p3a int
}

func possibleStates() []state {
	var res []state
	for p1a := 0; p1a <= 2; p1a++ {
		for p2a := 0; p2a <= 2; p2a++ {
			// Player 3 holds between 0 and 2 aces, and together the players
			// must hold between 2 and 4 aces.
			min := 2 - p1a - p2a
			if min < 0 {
				min = 0
			}
			max := 4 - p1a - p2a
			if max > 2 {
				max = 2
			}
			for p3a := min; p3a <= max; p3a++ {
				res = append(res, state{p1a, p2a, p3a})
			}
		}
	}
	return res
}

func (s state) String() string {
	playerStrings := [3]string{
		"88", "A8", "AA",
	}
	return fmt.Sprintf("1:%s 2:%s 3:%s", playerStrings[s.p1a], playerStrings[s.p2a], playerStrings[s.p3a])
}

// Each player knows the number of aces held by the other two characters.
type playerKnowledge struct {
	prevA, nextA int
}

func p1Knowledge(s state) playerKnowledge {
	return playerKnowledge{s.p3a, s.p2a}
}

func p2Knowledge(s state) playerKnowledge {
	return playerKnowledge{s.p1a, s.p3a}
}

func p3Knowledge(s state) playerKnowledge {
	return playerKnowledge{s.p2a, s.p1a}
}
