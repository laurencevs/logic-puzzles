package main

import "fmt"

type intPair struct {
	a, b int
}

func (p intPair) String() string {
	return fmt.Sprintf("(%d, %d)", p.a, p.b)
}

// Some intPair Valuations

func sum(p intPair) int {
	return p.a + p.b
}

func product(p intPair) int {
	return p.a * p.b
}

func absDifference(p intPair) int {
	if p.a >= p.b {
		return p.a - p.b
	}
	return p.b - p.a
}

// Some intPair Conditions

func hasNumberDivisibleBy(n int) func(intPair) bool {
	return func(p intPair) bool {
		return p.a%n == 0 || p.b%n == 0
	}
}

func productIsDivisibleBy(n int) func(intPair) bool {
	return func(p intPair) bool {
		return (p.a*p.b)%n == 0
	}
}

func sumIsDivisibleBy(n int) func(intPair) bool {
	return func(p intPair) bool {
		return (p.a+p.b)%n == 0
	}
}

func UnorderedIntPairs(min, max int, withRepetition bool) []intPair {
	if max < min {
		return []intPair{}
	}
	size := (max - min + 1) * (max - min) / 2
	if withRepetition {
		size += max - min + 1
	}
	pairs := make([]intPair, 0, size)
	for i := min; i <= max; i++ {
		if withRepetition {
			pairs = append(pairs, intPair{i, i})
		}
		for j := i + 1; j <= max; j++ {
			pairs = append(pairs, intPair{i, j})
		}
	}
	return pairs
}
