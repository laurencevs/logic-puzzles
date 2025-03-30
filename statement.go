package puzzles

import (
	"github.com/laurencevs/logic-puzzles/internal/set"
)

// Statement represents a statement about the solution to a puzzle.
// ConsistentWith should return true for any possibility that the Statement
// does not directly rule out.
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

// ValuationStatement is a statement whose truth depends only on the value of a
// particular valuation. A statement by an actor is always a
// ValuationStatement, where the valuation is the same as that pointed to by
// the actor's knowledge.
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

func (p *Puzzle[P]) ValuationEquals(v Valuation[P], value int) ValuationStatement[P] {
	return ValuationStatement[P]{
		valuation:     v,
		allowedValues: set.New(value),
	}
}

// KnowsAnswer is the statement that the solution has a unique value under the
// given actor's known valuation, among all remaining possibilities.
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

// DoesNotKnowAnswerGivenNormalised is the statement that even if the given
// actor were told the normalised value of the solution, they would not know
// the solution. That is, the solution has a non-unique normalised value among
// possibilities consistent with the actor's known valuation.
func (a *Actor[P]) DoesNotKnowAnswerGivenNormalised(normalise func(P) P) Condition[P] {
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
// Note that this will cause the puzzle's externalPossibilities and
// possibilitiesByKnowledge to become inconsistent. Narrate should only be used
// to reveal the solution to the audience at the end of the puzzle.
func (p *Puzzle[P]) Narrate(s Statement[P]) {
	filterInPlace(s, &p.externalPossibilities)
}

// Knows is the statement that the given statement evaluates to true for all
// solutions that the given actor considers possible based on their knowledge.
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

// KnowsWhether is the statement that the given statement has the same truth
// value for all solutions that the given actor considers possible based on
// their knowledge.
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

// Says makes the truth of the given statement s 'common knowledge' within the
// puzzle. It does not account for the information implied by the fact that
// the given actor knows s; for this, it should be combined with a.Knows(s).
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
