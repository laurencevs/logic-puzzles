package puzzles

import (
	"fmt"
	"slices"

	"github.com/laurencevs/logic-puzzles/internal/set"
)

// Knowledge represents the information an Actor is given about the solution
// before any statements are made.
type Knowledge[P any] *Valuation[P]

type Puzzle[P comparable] struct {
	// solutionSpace is the initial set of possible solutions to the Puzzle.
	solutionSpace []P
	// actors are the characters in the Puzzle.
	actors []*Actor[P]
	// possibilitiesByKnowledge represents the set of remaining possible
	// solutions, conditional on the value of a given piece of Knowledge.
	// One can reason about what solutions an Actor considers possible by
	// considering the possibilitiesByKnowledge values for their Knowledge.
	possibilitiesByKnowledge map[Knowledge[P]]map[int][]P
	// externalPossibilities represents the set of remaining possibilities
	// from the perspective of an outside observer who is not privy to any
	// specific Knowledge.
	externalPossibilities []P
}

func NewPuzzle[P comparable](possibilities []P) *Puzzle[P] {
	return &Puzzle[P]{
		solutionSpace:            possibilities,
		externalPossibilities:    slices.Clone(possibilities),
		possibilitiesByKnowledge: make(map[Knowledge[P]]map[int][]P),
	}
}

func (p *Puzzle[P]) SolutionSpace() []P {
	return p.solutionSpace
}

func (p *Puzzle[P]) Actors() []*Actor[P] {
	return p.actors
}

func (p *Puzzle[P]) NewActor() *Actor[P] {
	a := &Actor[P]{
		Id:     len(p.actors),
		puzzle: p,
	}
	p.actors = append(p.actors, a)
	return a
}

func (p *Puzzle[P]) NewKnowledge(v Valuation[P]) Knowledge[P] {
	p.initialiseKnowledge(&v)
	return &v
}

func (p *Puzzle[P]) initialiseKnowledge(k Knowledge[P]) {
	p.possibilitiesByKnowledge[k] = make(map[int][]P)
	for _, poss := range p.solutionSpace {
		val := (*k)(poss)
		p.possibilitiesByKnowledge[k][val] = append(p.possibilitiesByKnowledge[k][val], poss)
	}
}

func (p *Puzzle[P]) NewActorWithKnowledge(v Valuation[P]) *Actor[P] {
	k := p.NewKnowledge(v)
	a := p.NewActor()
	a.knowledge = k
	return a
}

func (p *Puzzle[P]) ExternalPossibilities() []P {
	return p.externalPossibilities
}

// Reset resets the Puzzle to its initial state.
func (p *Puzzle[P]) Reset() {
	p.externalPossibilities = slices.Clone(p.solutionSpace)
	copy(p.externalPossibilities, p.solutionSpace)
	for k := range p.possibilitiesByKnowledge {
		p.initialiseKnowledge(k)
	}
}

type Actor[P comparable] struct {
	Id        int
	puzzle    *Puzzle[P]
	knowledge Knowledge[P]
}

// HasKnowledge sets the Actor's Knowledge without initialising the internal
// Puzzle state for that Knowledge.
// It should only be used with Knowledge values created using
// Puzzle.NewKnowledge.
func (a *Actor[P]) HasKnowledge(k Knowledge[P]) {
	a.knowledge = k
}

func (a *Actor[P]) PossibilitiesByKnowledge() map[int][]P {
	return a.puzzle.possibilitiesByKnowledge[a.knowledge]
}

func NormalisePossibilities[P comparable](ps []P, normalise func(P) P) []P {
	s := set.New[P]()
	for _, p := range ps {
		s.Add(normalise(p))
	}
	return s.Values()
}

func SprintPossibilities[P any](ps []P) string {
	switch len(ps) {
	case 0:
		return "no remaining possibilities"
	case 1:
		return fmt.Sprintf("1 remaining possibility: %v", ps[0])
	default:
		return fmt.Sprintf("%d remaining possibilities", len(ps))
	}
}
