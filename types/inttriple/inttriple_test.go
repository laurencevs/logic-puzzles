package inttriple

import (
	"reflect"
	"testing"
)

func Test_UnorderedWithSum(t *testing.T) {
	s6 := IntTriplesWithSumWithoutRepetition(6, false)
	exp6 := []IntTriple{{1, 2, 3}}
	if !reflect.DeepEqual(s6, exp6) {
		t.Logf("Got %v, expected %v\n", s6, exp6)
		t.Fail()
	}

	s7 := IntTriplesWithSumWithoutRepetition(7, false)
	exp7 := []IntTriple{{1, 2, 4}}
	if !reflect.DeepEqual(s7, exp7) {
		t.Logf("Got %v, expected %v\n", s7, exp7)
		t.Fail()
	}

	s8 := IntTriplesWithSumWithoutRepetition(8, false)
	exp8 := []IntTriple{{1, 2, 5}, {1, 3, 4}}
	if !reflect.DeepEqual(s8, exp8) {
		t.Logf("Got %v, expected %v\n", s8, exp8)
		t.Fail()
	}

	s20 := IntTriplesWithSumWithoutRepetition(20, false)
	exp20 := []IntTriple{{1, 2, 17}, {1, 3, 16}, {1, 4, 15}, {1, 5, 14}, {1, 6, 13}, {1, 7, 12}, {1, 8, 11}, {1, 9, 10}, {2, 3, 15}, {2, 4, 14}, {2, 5, 13}, {2, 6, 12}, {2, 7, 11}, {2, 8, 10}, {3, 4, 13}, {3, 5, 12}, {3, 6, 11}, {3, 7, 10}, {3, 8, 9}, {4, 5, 11}, {4, 6, 10}, {4, 7, 9}, {5, 6, 9}, {5, 7, 8}}
	if !reflect.DeepEqual(s20, exp20) {
		t.Logf("Got %v, expected %v\n", s20, exp20)
		t.Fail()
	}

	s2025 := IntTriplesWithSumWithoutRepetition(2025, false)
	exp2025 := 340707
	if len(s2025) != exp2025 {
		t.Logf("Got slice of length %d, expected %d\n", len(s2025), exp2025)
		t.Fail()
	}
}

func Test_OrderedWithSum(t *testing.T) {
	s6 := IntTriplesWithSumWithoutRepetition(6, true)
	exp6 := []IntTriple{{1, 2, 3}, {1, 3, 2}, {2, 1, 3}, {2, 3, 1}, {3, 1, 2}, {3, 2, 1}}
	if !reflect.DeepEqual(s6, exp6) {
		t.Logf("Got %v, expected %v\n", s6, exp6)
		t.Fail()
	}

	s7 := IntTriplesWithSumWithoutRepetition(7, true)
	exp7 := []IntTriple{{1, 2, 4}, {1, 4, 2}, {2, 1, 4}, {2, 4, 1}, {4, 1, 2}, {4, 2, 1}}
	if !reflect.DeepEqual(s7, exp7) {
		t.Logf("Got %v, expected %v\n", s7, exp7)
		t.Fail()
	}

	s8 := IntTriplesWithSumWithoutRepetition(8, true)
	exp8 := []IntTriple{{1, 2, 5}, {1, 3, 4}, {1, 4, 3}, {1, 5, 2}, {2, 1, 5}, {2, 5, 1}, {3, 1, 4}, {3, 4, 1}, {4, 1, 3}, {4, 3, 1}, {5, 1, 2}, {5, 2, 1}}
	if !reflect.DeepEqual(s8, exp8) {
		t.Logf("Got %v, expected %v\n", s8, exp8)
		t.Fail()
	}

	s2025 := IntTriplesWithSumWithoutRepetition(2025, true)
	exp2025 := 340707 * 6
	if len(s2025) != exp2025 {
		t.Logf("Got slice of length %d, expected %d\n", len(s2025), exp2025)
		t.Fail()
	}
}

func Test_CompareOrderedUnorderedWithSum(t *testing.T) {
	for i := 6; i < 500; i++ {
		ord := IntTriplesWithSumWithoutRepetition(i, true)
		unord := IntTriplesWithSumWithoutRepetition(i, false)
		if len(ord) != len(unord)*6 {
			t.Logf("Number of ordered triples without repetition (%d) is not 6 times number of unordered triples (%d) for n=%d\n", len(ord), len(unord), i)
			t.Fail()
		}
		if cap(ord) != i*i/2 {
			t.Logf("Ordered triples without repetition exceed original capacity for n=%d", i)
			t.Fail()
		}
		if cap(unord) != i*i/12 {
			t.Logf("Unordered triples without repetition exceed original capacity for n=%d", i)
			t.Fail()
		}
	}
}
