package puzzles

import (
	"fmt"
	"slices"

	"github.com/laurencevs/logic-puzzles/internal/set"
)

type Puzzle[P comparable] struct {
	// solutionSpace is the initial set of possible solutions to the Puzzle.
	solutionSpace []P
	// actors are the characters in the Puzzle.
	actors []*Actor[P]
	// TODO: docs
	knowledgeStates set.Set[KnowledgeState[P]]
	// externalPossibilities represents the set of remaining possibilities from
	// the perspective of an outside observer who is not privy to any specific
	// knowledge.
	externalPossibilities []P
}

func NewPuzzle[P comparable](possibilities []P) *Puzzle[P] {
	return &Puzzle[P]{
		solutionSpace:         possibilities,
		externalPossibilities: slices.Clone(possibilities),
		knowledgeStates:       set.New[KnowledgeState[P]](),
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

func (p *Puzzle[P]) initialiseKnowledge(k KnowledgeState[P]) {
	k.Initialise(p.solutionSpace)
}

// NewKnowledge should be a method on Puzzle, but this requires generic
// methods (coming soon to Go).
func NewKnowledge[P, V comparable](p *Puzzle[P], f func(P) V) KnowledgeState[P] {
	pv := &valuationKnowledgeState[P, V]{valuation: f}
	p.initialiseKnowledge(pv)
	p.knowledgeStates.Add(pv)
	return pv
}

// NewActorWithKnowledge should be a method on Puzzle, but this requires
// generic methods (coming soon to Go).
func NewActorWithKnowledge[P, V comparable](p *Puzzle[P], f func(P) V) *Actor[P] {
	k := NewKnowledge(p, f)
	a := p.NewActor()
	a.KnowledgeState = k
	return a
}

func (p *Puzzle[P]) ExternalPossibilities() []P {
	return p.externalPossibilities
}

// Reset resets the puzzle to its initial state.
func (p *Puzzle[P]) Reset() {
	p.externalPossibilities = slices.Clone(p.solutionSpace)
	copy(p.externalPossibilities, p.solutionSpace)
	for k := range p.knowledgeStates {
		p.initialiseKnowledge(k)
	}
}

type Actor[P comparable] struct {
	Id     int
	puzzle *Puzzle[P]
	KnowledgeState[P]
}

// HasKnowledge sets the actor's knowledge without initialising the internal
// puzzle state for that knowledge. It should only be used with knowledge
// values created using Puzzle.NewKnowledge.
func (a *Actor[P]) HasKnowledge(k KnowledgeState[P]) {
	a.KnowledgeState = k
}

func (a *Actor[P]) PossibilitiesByKnowledge() KnowledgeState[P] {
	return a.KnowledgeState
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
