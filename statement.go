package puzzles

import (
	"github.com/laurencevs/logic-puzzles/internal/set"
)

type Valuation[T any] func(T) int
type Condition[T any] func(T) bool

type Statement[P comparable] interface {
	ConsistentWith(p P) bool
}

type ValuationStatement[P comparable] struct {
	valuation     Valuation[P]
	allowedValues set.Set[int]
	invert        bool
}

func (s ValuationStatement[P]) ConsistentWith(p P) bool {
	return s.allowedValues.Contains(s.valuation(p)) != s.invert
}

func (s ValuationStatement[P]) Not() ValuationStatement[P] {
	return ValuationStatement[P]{
		valuation:     s.valuation,
		allowedValues: s.allowedValues,
		invert:        !s.invert,
	}
}

func (c Condition[P]) ConsistentWith(p P) bool {
	return c(p)
}

func (c Condition[P]) Not() Condition[P] {
	return Condition[P](func(p P) bool {
		return !c(p)
	})
}

func filterInPlace[P comparable](s Statement[P], l *[]P) {
	i := 0
	for _, p := range *l {
		if s.ConsistentWith(p) {
			(*l)[i] = p
			i++
		}
	}
	*l = (*l)[:i]
}

func (p *Puzzle[P]) Evaluate(s Statement[P]) bool {
	for _, poss := range p.externalPossibilities {
		if !s.ConsistentWith(poss) {
			return false
		}
	}
	return true
}

func (p *Puzzle[P]) ValuationEquals(v Valuation[P], value int) ValuationStatement[P] {
	return ValuationStatement[P]{
		valuation:     v,
		allowedValues: set.New(value),
	}
}

func (a *Actor[P]) KnowsAnswer() ValuationStatement[P] {
	possibleValues := set.New[int]()
	for knowledgeValue, possiblities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		if len(possiblities) == 1 {
			possibleValues.Add(knowledgeValue)
		}
	}
	return ValuationStatement[P]{
		valuation:     *a.knowledge,
		allowedValues: possibleValues,
	}
}

func (a *Actor[P]) DoesNotKnowAnswer() ValuationStatement[P] {
	return a.KnowsAnswer().Not()
}

func (a *Actor[P]) KnowsNormalisedAnswer(normalise func(P) P) ValuationStatement[P] {
	possibleValues := set.New[int]()
knowledgeLoop:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		if len(possibilities) == 0 {
			continue
		}
		if len(possibilities) == 1 {
			possibleValues.Add(knowledge)
			continue
		}
		first := normalise(possibilities[0])
		for _, p := range possibilities[1:] {
			if normalise(p) != first {
				continue knowledgeLoop
			}
		}
		possibleValues.Add(knowledge)
	}
	return ValuationStatement[P]{
		valuation:     *a.knowledge,
		allowedValues: possibleValues,
	}
}

type possibilityWithKnowledge[P comparable] struct {
	possibility P
	knowledge   int
}

// Knowing normalise(s) is insufficient to determine s for all possibilities
func (a *Actor[P]) IsInsufficient(normalise func(P) P) Condition[P] {
	normalCount := make(map[possibilityWithKnowledge[P]]int)
	for k, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		for _, p := range possibilities {
			normalCount[possibilityWithKnowledge[P]{normalise(p), k}]++
		}
	}
	return func(p P) bool {
		return normalCount[possibilityWithKnowledge[P]{
			possibility: normalise(p),
			knowledge:   (*a.knowledge)(p),
		}] > 1
	}
}

// Narrate is the narrator's equivalent of Character.Says(). It restricts the
// solution space without informing any characters. Note that this will cause
// the puzzle's internalPossibilities and externalPossibilities to differ.
func (p *Puzzle[P]) Narrate(s Statement[P]) {
	filterInPlace(s, &p.externalPossibilities)
}

func (a *Actor[P]) Knows(s Statement[P]) ValuationStatement[P] {
	allowedValues := set.New[int]()
knowledgeLoop:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		if len(possibilities) == 0 {
			continue
		}
		for _, p := range possibilities {
			if !s.ConsistentWith(p) {
				continue knowledgeLoop
			}
		}
		allowedValues.Add(knowledge)
	}
	return ValuationStatement[P]{
		valuation:     *a.knowledge,
		allowedValues: allowedValues,
	}
}

func (a *Actor[P]) KnowsWhether(s Statement[P]) ValuationStatement[P] {
	allowedValues := set.New[int]()
knowledgeLoop:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		if len(possibilities) == 0 {
			continue
		}
		if len(possibilities) == 1 {
			allowedValues.Add(knowledge)
			continue
		}
		truthValue := s.ConsistentWith(possibilities[0])
		for _, p := range possibilities[1:] {
			if s.ConsistentWith(p) != truthValue {
				continue knowledgeLoop
			}
		}
		allowedValues.Add(knowledge)
	}
	return ValuationStatement[P]{
		valuation:     *a.knowledge,
		allowedValues: allowedValues,
	}
}

func (a *Actor[P]) Says(s ValuationStatement[P]) {
	filterInPlace(Statement[P](s), &a.puzzle.externalPossibilities)
	newPossibilitiesByKnowledge := make(map[Knowledge[P]]map[int][]P, len(a.puzzle.possibilitiesByKnowledge))
	for k, possibilitiesByValue := range a.puzzle.possibilitiesByKnowledge {
		newPossibilitiesByKnowledge[k] = make(map[int][]P, len(possibilitiesByValue))
		if k == a.knowledge {
			for value, possibilities := range possibilitiesByValue {
				if s.ConsistentWith(possibilities[0]) {
					newPossibilitiesByKnowledge[k][value] = possibilities
				}
			}
			continue
		}
		for value, possibilities := range possibilitiesByValue {
			filterInPlace(Statement[P](s), &possibilities)
			if len(possibilities) > 0 {
				newPossibilitiesByKnowledge[k][value] = possibilities
			}
		}
	}
	a.puzzle.possibilitiesByKnowledge = newPossibilitiesByKnowledge
}
