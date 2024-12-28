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
	knowledge                []Valuation[P]
	possibilitiesByKnowledge []map[int][]P
	externalPossibilities    []P
}

func NewPuzzle[P comparable](possibilities []P) *Puzzle[P] {
	return &Puzzle[P]{
		solutionSpace:         possibilities,
		externalPossibilities: slices.Clone(possibilities),
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

func (p *Puzzle[P]) NewKnowledge(v Valuation[P]) Knowledge {
	k := Knowledge{
		id: len(p.knowledge),
	}
	p.knowledge = append(p.knowledge, v)
	p.possibilitiesByKnowledge = append(p.possibilitiesByKnowledge, nil)
	p.initialiseKnowledge(k.id)
	return k
}

func (p *Puzzle[P]) initialiseKnowledge(id int) {
	p.possibilitiesByKnowledge[id] = make(map[int][]P)
	for _, poss := range p.solutionSpace {
		val := p.knowledge[id](poss)
		p.possibilitiesByKnowledge[id][val] = append(p.possibilitiesByKnowledge[id][val], poss)
	}
}

func (p *Puzzle[P]) NewActorWithKnowledge(v Valuation[P]) *Actor[P] {
	k := p.NewKnowledge(v)
	a := &Actor[P]{
		Id:        len(p.actors),
		knowledge: k,
		puzzle:    p,
	}
	p.actors = append(p.actors, a)
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
	for i := range p.knowledge {
		p.initialiseKnowledge(i)
	}
}

// Knowledge is a reference to a valuation in a specific puzzle's knowledge and possibilitiesByKnowledge lists
type Knowledge struct {
	id int
}

type Actor[P comparable] struct {
	Id        int
	puzzle    *Puzzle[P]
	knowledge Knowledge
}

func (a *Actor[P]) HasKnowledge(k Knowledge) {
	a.knowledge = k
}

func (a *Actor[P]) PossibilitiesByKnowledge() map[int][]P {
	return a.puzzle.possibilitiesByKnowledge[a.knowledge.id]
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
