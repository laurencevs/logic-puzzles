package puzzles

import (
	"github.com/laurencevs/logic-puzzles/internal/set"
)

// Statement represents a statement about the solution to a Puzzle.
// ConsistentWith should return true for any possibility that the Statement
// does not directly contradict.
type Statement[P any] interface {
	ConsistentWith(p P) bool
	Not() Statement[P]
}

type Condition[P any] func(P) bool

func (c Condition[P]) ConsistentWith(p P) bool {
	return c(p)
}

func (c Condition[P]) Not() Statement[P] {
	return Condition[P](func(p P) bool {
		return !c(p)
	})
}

var _ Statement[int] = Condition[int](func(int) bool { return false })

// ValuationStatement is a statement whose truth depends only on the value of a
// particular valuation. A statement by an actor is always a
// ValuationStatement, where the valuation is the same as that pointed to by
// the actor's Knowledge.
type ValuationStatement[P any, V comparable] struct {
	valuation     Valuation[P, V]
	allowedValues set.Set[V]
	invert        bool
}

func (s ValuationStatement[P, V]) ConsistentWith(p P) bool {
	return s.allowedValues.Contains(s.valuation(p)) != s.invert
}

func (s ValuationStatement[P, V]) not() ValuationStatement[P, V] {
	return ValuationStatement[P, V]{
		valuation:     s.valuation,
		allowedValues: s.allowedValues,
		invert:        !s.invert,
	}
}

func (s ValuationStatement[P, V]) Not() Statement[P] {
	return s.not()
}

var _ Statement[int] = (*ValuationStatement[int, string])(nil)

// Evaluate tests whether the given statement holds for all current
// possibilities, from an external observer's perspective.
func (p *Puzzle[P]) Evaluate(s Statement[P]) bool {
	for _, poss := range p.externalPossibilities {
		if !s.ConsistentWith(poss) {
			return false
		}
	}
	return true
}

func ValuationEquals[P, V comparable](v Valuation[P, V], value V) ValuationStatement[P, V] {
	return ValuationStatement[P, V]{
		valuation:     v,
		allowedValues: set.New(value),
	}
}

func ValuationIn[P, V comparable](v Valuation[P, V], values set.Set[V]) ValuationStatement[P, V] {
	return ValuationStatement[P, V]{
		valuation:     v,
		allowedValues: values,
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
// Note that this will cause the puzzle's externalPossibilities and
// possibilitiesByKnowledge to become inconsistent. Narrate should only be used
// to reveal the solution to the audience at the end of the puzzle.
func (p *Puzzle[P]) Narrate(s Statement[P]) {
	filterInPlace(s, &p.externalPossibilities)
}

// Says makes the truth of the given statement s 'common knowledge' within the
// puzzle. It does not account for the information implied by the fact that
// the given actor knows s; for this, it should be combined with a.Knows(s).
func (a *Actor[P]) Says(s Statement[P]) {
	filterInPlace(s, &a.puzzle.externalPossibilities)
	for k := range a.puzzle.knowledgeStates {
		k.Filter(s)
	}
}
