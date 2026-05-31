package puzzles

import "github.com/laurencevs/logic-puzzles/internal/set"

type Valuation[P any] func(P) int

// ValuationFromFunc converts a function f: S -> X for some set X to a
// valuation v: S -> int by assigning a unique integer to each value in the
// image of the original function.
func ValuationFromFunc[P any, Q comparable](solutionSpace []P, f func(P) Q) Valuation[P] {
	valuationByFuncValue := make(map[Q]int)
	i := 0
	for _, p := range solutionSpace {
		fp := f(p)
		_, ok := valuationByFuncValue[fp]
		if !ok {
			valuationByFuncValue[fp] = i
			i++
		}
	}
	return func(p P) int {
		return valuationByFuncValue[f(p)]
	}
}

// 'Possibilities by knowledge' are currently stored as a raw map, from each
// value to the possibilities that result in that value. This means all types
// of knowledge must use the same value type. So far I have chosen int, hence
// Valuation[P] is defined as func(P) int. In order to remove the reliance on
// int as a value type, I want to make Valuation generic, but in order to do
// that we need an interface to hide this genericity behind, so that it can be
// used in Puzzle[P] without requiring extra type parameters.

type IValuation[P any] interface {
	KnowsAnswer() Statement[P]
	DoesNotKnowAnswer() Statement[P]

	Knows(s Statement[P]) Statement[P]
	DoesNotKnow(s Statement[P]) Statement[P]

	KnowsWhether(s Statement[P]) Statement[P]
	DoesNotKnowWhether(s Statement[P]) Statement[P]

	KnowsNormalisedAnswer(normalise func(P) P) Statement[P]
	DoesNotKnowAnswerGivenNormalised(normalise func(P) P) Condition[P]
}

var _ IValuation[int] = (*Actor[int])(nil)

type GenericValuation[P any, V comparable] func(P) V

type GenericValuationStatement[P any, V comparable] struct {
	valuation     GenericValuation[P, V]
	allowedValues set.Set[V]
	invert        bool
}

func (s GenericValuationStatement[P, V]) ConsistentWith(p P) bool {
	return s.allowedValues.Contains(s.valuation(p)) != s.invert
}

func (s GenericValuationStatement[P, V]) not() GenericValuationStatement[P, V] {
	return GenericValuationStatement[P, V]{
		valuation:     s.valuation,
		allowedValues: s.allowedValues,
		invert:        !s.invert,
	}
}

func (s GenericValuationStatement[P, V]) Not() Statement[P] {
	return s.not()
}

type PossibilitiesByValuation[P comparable, V comparable] struct {
	f             GenericValuation[P, V]
	possibilities map[V][]P
}

var _ IValuation[int] = PossibilitiesByValuation[int, string]{}

func (pv PossibilitiesByValuation[P, V]) knowsAnswer() GenericValuationStatement[P, V] {
	possibleValues := set.New[V]()
	for knowledgeValue, possiblities := range pv.possibilities {
		if len(possiblities) == 1 {
			possibleValues.Add(knowledgeValue)
		}
	}
	return GenericValuationStatement[P, V]{
		valuation:     pv.f,
		allowedValues: possibleValues,
	}
}

func (pv PossibilitiesByValuation[P, V]) KnowsAnswer() Statement[P] {
	return pv.knowsAnswer()
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowAnswer() Statement[P] {
	return pv.knowsAnswer().not()
}

func (pv PossibilitiesByValuation[P, V]) knows(s Statement[P]) GenericValuationStatement[P, V] {
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
	return GenericValuationStatement[P, V]{
		valuation:     pv.f,
		allowedValues: allowedValues,
	}
}

func (pv PossibilitiesByValuation[P, V]) Knows(s Statement[P]) Statement[P] {
	return pv.knows(s)
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnow(s Statement[P]) Statement[P] {
	return pv.knows(s).not()
}

func (pv PossibilitiesByValuation[P, V]) knowsWhether(s Statement[P]) GenericValuationStatement[P, V] {
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
	return GenericValuationStatement[P, V]{
		valuation:     pv.f,
		allowedValues: allowedValues,
	}
}

func (pv PossibilitiesByValuation[P, V]) KnowsWhether(s Statement[P]) Statement[P] {
	return pv.knowsWhether(s)
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowWhether(s Statement[P]) Statement[P] {
	return pv.knowsWhether(s).not()
}

func (pv PossibilitiesByValuation[P, V]) KnowsNormalisedAnswer(normalise func(P) P) Statement[P] {
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
	return GenericValuationStatement[P, V]{
		valuation:     pv.f,
		allowedValues: possibleValues,
	}
}

type genericPossibilityWithKnowledge[P, V comparable] struct {
	possibility P
	knowledge   V
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowAnswerGivenNormalised(normalise func(P) P) Condition[P] {
	normalCount := make(map[genericPossibilityWithKnowledge[P, V]]int)
	for k, possibilities := range pv.possibilities {
		for _, p := range possibilities {
			normalCount[genericPossibilityWithKnowledge[P, V]{normalise(p), k}]++
		}
	}
	return func(p P) bool {
		return normalCount[genericPossibilityWithKnowledge[P, V]{
			possibility: normalise(p),
			knowledge:   pv.f(p),
		}] > 1
	}
}
