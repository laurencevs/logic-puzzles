package puzzles

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

type PossibilitiesByValuation[P any, V comparable] struct {
	f             GenericValuation[P, V]
	possibilities map[V][]P
}

var _ IValuation[int] = PossibilitiesByValuation[int, string]{}

func (pv PossibilitiesByValuation[P, V]) KnowsAnswer() Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowAnswer() Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) Knows(s Statement[P]) Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnow(s Statement[P]) Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) KnowsWhether(s Statement[P]) Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowWhether(s Statement[P]) Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) KnowsNormalisedAnswer(normalise func(P) P) Statement[P] {
	// TODO
	return nil
}

func (pv PossibilitiesByValuation[P, V]) DoesNotKnowAnswerGivenNormalised(normalise func(P) P) Condition[P] {
	// TODO
	return nil
}
