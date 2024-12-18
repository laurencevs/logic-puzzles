package puzzles

import (
	"github.com/laurencevs/logic-puzzles/internal/set"
)

type Valuation[T any] func(T) int
type Condition[T any] func(T) bool

type Statement[P comparable] struct {
	valuation     Valuation[P]
	allowedValues set.Set[int]
	invert        bool
}

func (s Statement[P]) Not() Statement[P] {
	return Statement[P]{
		valuation:     s.valuation,
		allowedValues: s.allowedValues,
		invert:        !s.invert,
	}
}

func (s Statement[P]) evaluate(p P) bool {
	return s.allowedValues.Contains(s.valuation(p)) != s.invert
}

func (s Statement[P]) filterInPlace(l *[]P) {
	i := 0
	for _, p := range *l {
		if s.evaluate(p) {
			(*l)[i] = p
			i++
		}
	}
	*l = (*l)[:i]
}

func (c Condition[P]) Not() Condition[P] {
	return func(p P) bool {
		return !c(p)
	}
}

func (p *Puzzle[P]) Evaluate(s Statement[P]) bool {
	for _, poss := range p.externalPossibilities {
		if !s.evaluate(poss) {
			return false
		}
	}
	return true
}

func (p *Puzzle[P]) ValuationEquals(v Valuation[P], value int) Statement[P] {
	return Statement[P]{
		valuation:     v,
		allowedValues: set.New(value),
	}
}

func (a *Actor[P]) KnowsAnswer() Statement[P] {
	possibleValues := set.New[int]()
	for knowledge, possiblities := range a.puzzle.possibilitiesByKnowledge[a.knowledge.id] {
		if len(possiblities) == 1 {
			possibleValues.Add(knowledge)
		}
	}
	return Statement[P]{
		valuation:     a.puzzle.knowledge[a.knowledge.id],
		allowedValues: possibleValues,
	}
}

func (a *Actor[P]) DoesNotKnowAnswer() Statement[P] {
	return a.KnowsAnswer().Not()
}

// Narrate is the narrator's equivalent of Character.Says(). It restricts the
// solution space without informing any characters. Note that this will cause
// the puzzle's internalPossibilities and externalPossibilities to differ.
func (p *Puzzle[P]) Narrate(s Statement[P]) {
	s.filterInPlace(&p.externalPossibilities)
}

func (a *Actor[P]) knows(eval func(P) bool) Statement[P] {
	allowedValues := set.New[int]()
	v := a.puzzle.knowledge[a.knowledge.id]
knowledgeLoop:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge.id] {
		if len(possibilities) == 0 {
			continue
		}
		for _, p := range possibilities {
			if !eval(p) {
				continue knowledgeLoop
			}
		}
		allowedValues.Add(knowledge)
	}
	return Statement[P]{
		valuation:     v,
		allowedValues: allowedValues,
	}
}

func (a *Actor[P]) Knows(s Statement[P]) Statement[P] {
	return a.knows(s.evaluate)
}

func (a *Actor[P]) KnowsHolds(c Condition[P]) Statement[P] {
	return a.knows(c)
}

func (a *Actor[P]) knowsWhether(eval func(P) bool) Statement[P] {
	allowedValues := set.New[int]()
	v := a.puzzle.knowledge[a.knowledge.id]
knowledgeLoop:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge.id] {
		if len(possibilities) == 0 {
			continue
		}
		if len(possibilities) == 1 {
			allowedValues.Add(knowledge)
			continue
		}
		truthValue := eval(possibilities[0])
		for _, p := range possibilities[1:] {
			if eval(p) != truthValue {
				continue knowledgeLoop
			}
		}
		allowedValues.Add(knowledge)
	}
	return Statement[P]{
		valuation:     v,
		allowedValues: allowedValues,
	}
}

func (a *Actor[P]) KnowsWhether(s Statement[P]) Statement[P] {
	return a.knowsWhether(s.evaluate)
}

func (a *Actor[P]) KnowsWhetherHolds(c Condition[P]) Statement[P] {
	return a.knowsWhether(c)
}

func (a *Actor[P]) Says(s Statement[P]) {
	s.filterInPlace(&a.puzzle.externalPossibilities)
	newPossibilitiesByKnowledge := make([]map[int][]P, len(a.puzzle.possibilitiesByKnowledge))
	for id, possibilitiesByValue := range a.puzzle.possibilitiesByKnowledge {
		newPossibilitiesByKnowledge[id] = make(map[int][]P, len(possibilitiesByValue))
		if id == a.knowledge.id {
			for value, possibilities := range possibilitiesByValue {
				if s.allowedValues.Contains(value) != s.invert {
					newPossibilitiesByKnowledge[id][value] = possibilities
				}
			}
			continue
		}
		for value, possibilities := range possibilitiesByValue {
			s.filterInPlace(&possibilities)
			newPossibilitiesByKnowledge[id][value] = possibilities
		}
	}
	a.puzzle.possibilitiesByKnowledge = newPossibilitiesByKnowledge
}
