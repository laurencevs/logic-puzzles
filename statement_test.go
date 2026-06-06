package puzzles

import (
	"slices"
	"testing"
)

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for a := 2; a*a <= n; a++ {
		if n%a == 0 {
			return false
		}
	}
	return true
}

func TestFilterInPlace(t *testing.T) {
	i100 := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		i100 = append(i100, i)
	}

	isPrimeStatement := Condition[int](isPrime)
	filterInPlace[int](isPrimeStatement, &i100)
	expected := []int{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47,
		53, 59, 61, 67, 71, 73, 79, 83, 89, 97,
	}
	if !slices.Equal(i100, expected) {
		t.Errorf("expected %v, got %v", expected, i100)
	}

	isGt70Statement := Condition[int](func(i int) bool { return i > 70 })
	filterInPlace[int](isGt70Statement, &i100)
	expected = []int{71, 73, 79, 83, 89, 97}
	if !slices.Equal(i100, expected) {
		t.Errorf("expected %v, got %v", expected, i100)
	}
}
