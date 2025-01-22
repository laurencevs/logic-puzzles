package puzzles

import (
	"fmt"
	"slices"
	"strings"

	"github.com/laurencevs/logic-puzzles/internal/set"
)

type Puzzle[P comparable] struct {
	solutionSpace            []P
	actors                   []*Actor[P]
	possibilitiesByKnowledge map[Knowledge[P]]map[int][]P
	externalPossibilities    []P
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

func (p *Puzzle[P]) NormalisedPossibilities(normalise func(P) P) []P {
	s := set.New[P]()
	for _, p := range p.externalPossibilities {
		s.Add(normalise(p))
	}
	return s.Values()
}

func (p *Puzzle[P]) Reset() {
	p.externalPossibilities = slices.Clone(p.solutionSpace)
	copy(p.externalPossibilities, p.solutionSpace)
	for k := range p.possibilitiesByKnowledge {
		p.initialiseKnowledge(k)
	}
}

type Knowledge[P comparable] *Valuation[P]

type Actor[P comparable] struct {
	Id        int
	puzzle    *Puzzle[P]
	knowledge Knowledge[P]
}

func (a *Actor[P]) HasKnowledge(k Knowledge[P]) {
	a.knowledge = k
}

func (a *Actor[P]) PossibilitiesByKnowledge() map[int][]P {
	return a.puzzle.possibilitiesByKnowledge[a.knowledge]
}

func SprintPossibilities[P comparable](ps []P) string {
	var b strings.Builder
	pss := "possibilities"
	if len(ps) == 1 {
		pss = "possibility"
	}
	b.WriteString(fmt.Sprintf("Puzzle has %d remaining %s", len(ps), pss))
	if len(ps) == 0 {
		return b.String()
	}
	b.WriteString(": ")
	b.WriteString(fmt.Sprint(ps[0]))
	for i, p := range ps[1:] {
		b.WriteString(", ")
		if i == 49 {
			b.WriteString("...")
		}
		b.WriteString(fmt.Sprint(p))
	}
	return b.String()
}
