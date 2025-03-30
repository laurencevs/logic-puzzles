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
