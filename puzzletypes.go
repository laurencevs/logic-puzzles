package puzzles

import "fmt"

type IntPair struct {
	a, b int
}

func (p IntPair) String() string {
	return fmt.Sprintf("(%d, %d)", p.a, p.b)
}

// Some intPair Valuations

func Sum(p IntPair) int {
	return p.a + p.b
}

func Product(p IntPair) int {
	return p.a * p.b
}

func AbsDifference(p IntPair) int {
	if p.a >= p.b {
		return p.a - p.b
	}
	return p.b - p.a
}

// Some intPair Conditions

func HasNumberDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return p.a%n == 0 || p.b%n == 0
	}
}

func HasNumberIn(s map[int]struct{}) func(IntPair) bool {
	return func(p IntPair) bool {
		_, ok1 := s[p.a]
		_, ok2 := s[p.b]
		return ok1 || ok2
	}
}

func HasOneNumberIn(s map[int]struct{}) func(IntPair) bool {
	return func(p IntPair) bool {
		_, ok1 := s[p.a]
		_, ok2 := s[p.b]
		return ok1 != ok2
	}
}

func ProductIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.a*p.b)%n == 0
	}
}

func ProductIsNotDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.a*p.b)%n != 0
	}
}

func SumIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.a+p.b)%n == 0
	}
}

func AbsDifferenceIsDivisibleBy(n int) func(IntPair) bool {
	return func(p IntPair) bool {
		return (p.a-p.b)%n == 0
	}
}

func UnorderedIntPairs(min, max int, withRepetition bool) []IntPair {
	if max < min {
		return nil
	}
	size := (max - min + 1) * (max - min) / 2
	if withRepetition {
		size += max - min + 1
	}
	pairs := make([]IntPair, 0, size)
	for i := min; i <= max; i++ {
		if withRepetition {
			pairs = append(pairs, IntPair{i, i})
		}
		for j := i + 1; j <= max; j++ {
			pairs = append(pairs, IntPair{i, j})
		}
	}
	return pairs
}
