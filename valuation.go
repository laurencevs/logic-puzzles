package puzzles

import "github.com/laurencevs/logic-puzzles/internal/set"

type Valuation[P any, V comparable] func(P) V

// ValuationStatement is a statement whose truth depends only on the value of a
// particular Valuation. A statement by an Actor is always a
// ValuationStatement, where the Valuation is the same as that pointed to by
// the actor's KnowledgeState.
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

// valuationKnowledgeState represents the state of knowledge of a character who
// has been told the value of some function (Valuation) of the solution.
type valuationKnowledgeState[P comparable, V comparable] struct {
	valuation     Valuation[P, V]
	possibilities map[V][]P
}

func (pv *valuationKnowledgeState[P, V]) Initialise(solutionSpace []P) {
	pv.possibilities = make(map[V][]P)
	for _, poss := range solutionSpace {
		val := pv.valuation(poss)
		pv.possibilities[val] = append(pv.possibilities[val], poss)
	}
}

func (pv *valuationKnowledgeState[P, V]) Update(s Statement[P]) {
	newPossibilities := make(map[V][]P)
	for value, possibilities := range pv.possibilities {
		filterInPlace(s, &possibilities)
		if len(possibilities) > 0 {
			newPossibilities[value] = possibilities
		}
	}
	pv.possibilities = newPossibilities
}

func (pv *valuationKnowledgeState[P, V]) PossibilitiesInCase(s P) []P {
	return pv.possibilities[pv.valuation(s)]
}

func (pv *valuationKnowledgeState[P, V]) Possibilities() []P {
	var pp []P
	for _, p := range pv.possibilities {
		pp = append(pp, p...)
	}
	return pp
}

func (pv *valuationKnowledgeState[P, V]) knowsAnswer() ValuationStatement[P, V] {
	possibleValues := set.New[V]()
	for knowledgeValue, possiblities := range pv.possibilities {
		if len(possiblities) == 1 {
			possibleValues.Add(knowledgeValue)
		}
	}
	return ValuationIn(pv.valuation, possibleValues)
}

// KnowsAnswer is the statement that the solution has a unique value under the
// given valuation, among all remaining possibilities.
func (pv *valuationKnowledgeState[P, V]) KnowsAnswer() Statement[P] {
	return pv.knowsAnswer()
}

// DoesNotKnowAnswer is the statement that the solution does not have a unique
// value under the given valuation, among all remaining possibilities.
func (pv *valuationKnowledgeState[P, V]) DoesNotKnowAnswer() Statement[P] {
	return pv.knowsAnswer().not()
}

func (pv *valuationKnowledgeState[P, V]) knows(s Statement[P]) ValuationStatement[P, V] {
	allowedValues := set.New[V]()
outer:
	for knowledge, possibilities := range pv.possibilities {
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
	return ValuationIn(pv.valuation, allowedValues)
}

// Knows is the statement that the given statement evaluates to true for all
// solutions that the given actor considers possible based on their knowledge.
func (pv *valuationKnowledgeState[P, V]) Knows(s Statement[P]) Statement[P] {
	return pv.knows(s)
}

// DoesNotKnow is the statement that the given statement does not evaluate to
// true for all solutions that the given actor considers possible based on
// their knowledge.
//
// Note that k.Knows(s).Not() is not the same as k.Knows(s.Not()).
func (pv *valuationKnowledgeState[P, V]) DoesNotKnow(s Statement[P]) Statement[P] {
	return pv.knows(s).not()
}

func (pv *valuationKnowledgeState[P, V]) knowsWhether(s Statement[P]) ValuationStatement[P, V] {
	allowedValues := set.New[V]()
knowledgeLoop:
	for knowledge, possibilities := range pv.possibilities {
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
	return ValuationIn(pv.valuation, allowedValues)
}

// KnowsWhether is the statement that the given statement has the same truth
// value for all solutions that the given actor considers possible based on
// their knowledge.
func (pv *valuationKnowledgeState[P, V]) KnowsWhether(s Statement[P]) Statement[P] {
	return pv.knowsWhether(s)
}

// DoesNotKnowWhether is the statement that the given statement does not have
// the same truth value for all solutions that the given actor considers
// possible based on their knowledge.
func (pv *valuationKnowledgeState[P, V]) DoesNotKnowWhether(s Statement[P]) Statement[P] {
	return pv.knowsWhether(s).not()
}

func (pv *valuationKnowledgeState[P, V]) KnowsNormalisedAnswer(normalise func(P) P) Statement[P] {
	possibleValues := set.New[V]()
outer:
	for knowledge, possibilities := range pv.possibilities {
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
	return ValuationIn(pv.valuation, possibleValues)
}

type possibilityWithKnowledge[P, V comparable] struct {
	possibility P
	knowledge   V
}

// KnowsAnswerGivenNormalised is the statement that if the given actor were
// told the normalised value of the solution, they would know the solution.
// That is, the solution has a unique normalised value among possibilities
// consistent with the actor's known valuation.
func (pv *valuationKnowledgeState[P, V]) KnowsAnswerGivenNormalised(normalise func(P) P) Statement[P] {
	normalCount := make(map[possibilityWithKnowledge[P, V]]int)
	for k, possibilities := range pv.possibilities {
		for _, p := range possibilities {
			normalCount[possibilityWithKnowledge[P, V]{normalise(p), k}]++
		}
	}
	return Condition[P](func(p P) bool {
		return normalCount[possibilityWithKnowledge[P, V]{
			possibility: normalise(p),
			knowledge:   pv.valuation(p),
		}] <= 1
	})
}
