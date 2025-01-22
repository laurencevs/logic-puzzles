package inttriple

import (
	"fmt"
	"sort"

	puzzles "github.com/laurencevs/logic-puzzles"
)

type IntTriple struct {
	A, B, C int
}

func (t IntTriple) String() string {
	return fmt.Sprintf("(%d, %d, %d)", t.A, t.B, t.C)
}

/*
IntTriplesWithSumWithoutRepetition returns a list of all positive integer
triples with the given sum n. If ordered is false, the triples are normalised
to satisfy A < B < C.
*/
func IntTriplesWithSumWithoutRepetition(n int, ordered bool) []IntTriple {
	if n < 6 {
		return nil
	}
	size := n * n / 2
	if !ordered {
		size /= 6
	}
	res := make([]IntTriple, 0, size)
	for a := 1; a < n-2; a++ {
		b := 1
		if !ordered {
			if a == n/3 {
				break
			}
			b = a + 1
		}
		for ; b < n-a; b++ {
			c := n - a - b
			if !ordered && c <= b {
				break
			}
			if a == b || a == c || b == c {
				continue
			}
			res = append(res, IntTriple{a, b, c})
		}
	}
	return res
}

func Normalise(t IntTriple) IntTriple {
	vals := []int{t.A, t.B, t.C}
	sort.Ints(vals)
	return IntTriple{
		A: vals[0],
		B: vals[1],
		C: vals[2],
	}
}

// Some IntTriple Valuations

func Pair1Product(t IntTriple) int {
	return t.A * t.B
}

func Pair2Product(t IntTriple) int {
	return t.A * t.C
}

func Pair3Product(t IntTriple) int {
	return t.B * t.C
}

func SumOfPairwiseProducts(t IntTriple) int {
	return t.A*(t.B+t.C) + t.B*t.C
}

// Some IntTriple Conditions

func HasNumberIn(s map[int]struct{}) puzzles.Condition[IntTriple] {
	return func(t IntTriple) bool {
		_, ok1 := s[t.A]
		_, ok2 := s[t.B]
		_, ok3 := s[t.C]
		return ok1 || ok2 || ok3
	}
}

func SumIsDivisibleBy(n int) puzzles.Condition[IntTriple] {
	return func(t IntTriple) bool {
		return (t.A+t.B+t.C)%n == 0
	}
}
