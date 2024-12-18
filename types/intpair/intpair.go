package intpair

import "fmt"

type IntPair struct {
	A, B int
}

func (p IntPair) String() string {
	return fmt.Sprintf("(%d, %d)", p.A, p.B)
}

// IntPairs returns all (un)ordered integer pairs with values in the closed
// interval [min, max], with(out) repetition.
func IntPairs(min, max int, ordered, withRepetition bool) []IntPair {
	if max < min || (max == min && !withRepetition) {
		return nil
	}
	size := (max - min + 1) * (max - min)
	if !ordered {
		size /= 2
	}
	if withRepetition {
		size += max - min + 1
	}
	pairs := make([]IntPair, 0, size)
	for i := min; i <= max; i++ {
		if withRepetition {
			pairs = append(pairs, IntPair{i, i})
		}
		if ordered {
			for j := min; j < i; j++ {
				pairs = append(pairs, IntPair{i, j})
			}
		}
		for j := i + 1; j <= max; j++ {
			pairs = append(pairs, IntPair{i, j})
		}
	}
	return pairs
}

// Some IntPair Valuations

func First(p IntPair) int {
	return p.A
}

func Second(p IntPair) int {
	return p.B
}

func Sum(p IntPair) int {
	return p.A + p.B
}

func Product(p IntPair) int {
	return p.A * p.B
}

func AbsDifference(p IntPair) int {
	if p.A >= p.B {
		return p.A - p.B
	}
	return p.B - p.A
}

// Some intPair Conditions

func HasNumberDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return p.A%n == 0 || p.B%n == 0
	}
}

func HasNumberIn(s map[int]struct{}) func(IntPair) bool {
	return func(p IntPair) bool {
		_, ok1 := s[p.A]
		_, ok2 := s[p.B]
		return ok1 || ok2
	}
}

func HasOneNumberIn(s map[int]struct{}) func(IntPair) bool {
	return func(p IntPair) bool {
		_, ok1 := s[p.A]
		_, ok2 := s[p.B]
		return ok1 != ok2
	}
}

func ProductIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.A*p.B)%n == 0
	}
}

func ProductIsNotDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.A*p.B)%n != 0
	}
}

func SumIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.A+p.B)%n == 0
	}
}

func AbsDifferenceIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.A-p.B)%n == 0
	}
}
