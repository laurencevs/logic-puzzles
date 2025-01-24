package puzzles

import (
	"github.com/laurencevs/logic-puzzles/internal/set"
)

// Statement represents a statement about the solution to a Puzzle.
// ConsistentWith should return true for any possibility that the Statement
// does not rule out.
type Statement[P any] interface {
	ConsistentWith(p P) bool
}

type Condition[P any] func(P) bool

func (c Condition[P]) ConsistentWith(p P) bool {
	return c(p)
}

func (c Condition[P]) Not() Condition[P] {
	return Condition[P](func(p P) bool {
		return !c(p)
	})
}

// ValuationStatement is a Statement whose truth depends only on the value of a
// particular Valuation. Statements by Actors are always ValuationStatements.
type ValuationStatement[P any] struct {
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
outer:
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
				continue outer
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

func filterInPlace[P any](s Statement[P], l *[]P) {
	i := 0
	for _, p := range *l {
		if s.ConsistentWith(p) {
			(*l)[i] = p
			i++
		}
	}
	*l = (*l)[:i]
}

// Narrate is the narrator's equivalent of Actor.Says(). It restricts the
// solution space without 'informing' any characters.
//
// Note that this will cause the Puzzle's externalPossibilities and
// possibilitiesByKnowledge to become inconsistent. Narrate should only be used
// to reveal the solution to the audience at the end of the Puzzle.
func (p *Puzzle[P]) Narrate(s Statement[P]) {
	filterInPlace(s, &p.externalPossibilities)
}

// a.Knows(s) is equivalent to saying that s evaluates to true for all
// solutions that a considers possible based on their Knowledge.
func (a *Actor[P]) Knows(s Statement[P]) ValuationStatement[P] {
	allowedValues := set.New[int]()
outer:
	for knowledge, possibilities := range a.puzzle.possibilitiesByKnowledge[a.knowledge] {
		if len(possibilities) == 0 {
			continue
		}
		for _, p := range possibilities {
			if !s.ConsistentWith(p) {
				continue outer
			}
		}
		allowedValues.Add(knowledge)
	}
	return ValuationStatement[P]{
		valuation:     *a.knowledge,
		allowedValues: allowedValues,
	}
}

// a.KnowsWhether(s) is equivalent to the statement that s has the same truth
// value for all solutions that a considers possible based on their Knowledge.
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

// Says makes the truth of Statement s 'common knowledge' within the Puzzle. It
// does not account for the information implied by the fact that a knows s; for
// this, it should be combined with a.Knows(s).
func (a *Actor[P]) Says(s Statement[P]) {
	filterInPlace(s, &a.puzzle.externalPossibilities)
	newPossibilitiesByKnowledge := make(map[Knowledge[P]]map[int][]P, len(a.puzzle.possibilitiesByKnowledge))
	for k, possibilitiesByValue := range a.puzzle.possibilitiesByKnowledge {
		newPossibilitiesByKnowledge[k] = make(map[int][]P, len(possibilitiesByValue))
		for value, possibilities := range possibilitiesByValue {
			filterInPlace(s, &possibilities)
			if len(possibilities) > 0 {
				newPossibilitiesByKnowledge[k][value] = possibilities
			}
		}
	}
	a.puzzle.possibilitiesByKnowledge = newPossibilitiesByKnowledge
}
