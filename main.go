package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func main() {
	run2024Puzzles()
	generate2024Puzzles()
}

func shuffleSlice[T any](s []T) {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

func randomStatement[P PuzzlePossibility](characters []*Character[P]) string {
	s := []int{0, 1, 2}
	shuffleSlice(s)
	c1 := characters[s[0]]
	c2 := characters[s[1]]
	c3 := characters[s[2]]
	n := rand.Intn(4)
	switch n {
	case 0:
		c1.Says(c1.Knows(c2.Knows(c3.DoesNotKnowAnswer())))
	case 1:
		c1.Says(c1.Knows(c2.Knows(c3.KnowsAnswer())))
	case 2:
		c1.Says(c1.Knows(c2.Knows(c3.Knows(c1.KnowsAnswer()))))
	case 3:
		c1.Says(c1.KnowsAnswer())
	}
	return fmt.Sprintf("%d %s %s %s; ", n, c1.name, c2.name, c3.name)
}

type PuzzlePossibility interface {
	comparable
	String() string
}

type Puzzle[P PuzzlePossibility] struct {
	// externalPossibilities represents the set of possible states from the
	// perspective of an external observer, i.e. the user. It takes into
	// account both statements by the characters and constraints imposed in
	// the narrative of the puzzle, e.g. "Paul knows...". It is therefore a
	// subset of internalPossibilities.
	externalPossibilities map[P]struct{}
	// internalPossibilities represents the set of possible states from the
	// perspective of an observer within the puzzle. It only takes into
	// account public statements by the characters.
	internalPossibilities map[P]struct{}
	characters            []*Character[P]
	// originalPossibilities stores the original state of the puzzle so that
	// it can be reset.
	originalPossibilities []P
}

func NewPuzzle[P PuzzlePossibility](possibilities []P) *Puzzle[P] {
	internalPossibilities := make(map[P]struct{},
		len(possibilities))
	externalPossibilities := make(map[P]struct{},
		len(possibilities))
	for _, p := range possibilities {
		internalPossibilities[p] = struct{}{}
		externalPossibilities[p] = struct{}{}
	}
	return &Puzzle[P]{
		internalPossibilities: internalPossibilities,
		externalPossibilities: externalPossibilities,
		originalPossibilities: possibilities,
	}
}

func (p *Puzzle[P]) NewCharacter(name string) *Character[P] {
	privatePossibilities := make(map[P]struct{},
		len(p.internalPossibilities))
	for p := range p.internalPossibilities {
		privatePossibilities[p] = struct{}{}
	}
	character := &Character[P]{
		name:                   name,
		puzzle:                 p,
		possibilities:          privatePossibilities,
		knowledgeByPossibility: make(map[P]string),
	}
	p.characters = append(p.characters, character)
	return character
}

func (p *Puzzle[P]) PrintPossibilities() {
	b := strings.Builder{}
	possibilitiesString := "possibilities"
	if len(p.externalPossibilities) == 1 {
		possibilitiesString = "possibility"
	}
	b.WriteString(fmt.Sprintf("Puzzle has %d remaining %s: ",
		len(p.externalPossibilities), possibilitiesString))
	count := 0
	for poss := range p.externalPossibilities {
		if count != 0 {
			_, err := b.WriteString(", ")
			if err != nil {
				panic(err)
			}
		}
		_, err := b.WriteString(poss.String())
		if err != nil {
			panic(err)
		}
		count += 1
		if count >= 5 {
			b.WriteString("...")
			break
		}
	}
	fmt.Println(b.String())
}

func (p *Puzzle[P]) Reset() {
	for _, poss := range p.originalPossibilities {
		p.internalPossibilities[poss] = struct{}{}
		p.externalPossibilities[poss] = struct{}{}
		for _, c := range p.characters {
			c.possibilities[poss] = struct{}{}
			delete(c.knowledgeByPossibility, poss)
		}
	}
	for _, c := range p.characters {
		knownValues := c.knownValues
		c.knownValues = []Valuation[P]{}
		for _, v := range knownValues {
			c.KnowsValueOf(v)
		}
	}
}

// Statements are represented by the set of internal possibilities they rule out.
type Statement[T comparable] map[T]struct{}

func (p *Puzzle[P]) ValuationEquals(v Valuation[P], value int) Statement[P] {
	return nil
}

type Character[P PuzzlePossibility] struct {
	name                   string
	puzzle                 *Puzzle[P]
	knownValues            []Valuation[P]
	possibilities          map[P]struct{}
	knowledgeByPossibility map[P]string
}

func (c *Character[P]) KnowsValueOf(v Valuation[P]) {
	c.knownValues = append(c.knownValues, v)
	for poss := range c.puzzle.internalPossibilities {
		c.knowledgeByPossibility[poss] += "/" + strconv.Itoa(v(poss))
	}
}

func (c *Character[P]) KnowsAnswer() Statement[P] {
	return c.HasNPossibilities(func(n int) bool { return n == 1 })
}

func (c *Character[P]) DoesNotKnowAnswer() Statement[P] {
	return c.HasNPossibilities(func(n int) bool { return n > 1 })
}

func (c *Character[P]) HasNPossibilities(condition func(n int) bool) Statement[P] {
	numPossibilitiesByKnowledge := make(map[string]int)
	for poss := range c.puzzle.internalPossibilities {
		numPossibilitiesByKnowledge[c.knowledgeByPossibility[poss]] += 1
	}
	possibilitiesToDelete := make(map[P]struct{},
		len(c.puzzle.internalPossibilities))
	for poss := range c.puzzle.internalPossibilities {
		if !condition(numPossibilitiesByKnowledge[c.knowledgeByPossibility[poss]]) {
			possibilitiesToDelete[poss] = struct{}{}
		}
	}
	for poss := range possibilitiesToDelete {
		delete(c.puzzle.externalPossibilities, poss)
		delete(c.possibilities, poss)
	}
	return possibilitiesToDelete
}

func (p *Puzzle[P]) AnswerSatisfies(condition Condition[P]) Statement[P] {
	possibilitiesToDelete := make(map[P]struct{},
		len(p.internalPossibilities))
	for poss := range p.internalPossibilities {
		if !condition(poss) {
			possibilitiesToDelete[poss] = struct{}{}
		}
	}
	for poss := range possibilitiesToDelete {
		delete(p.externalPossibilities, poss)
	}
	return possibilitiesToDelete
}

func (c *Character[P]) Knows(f Statement[P]) Statement[P] {
	// delete possibilities corresponding to knowledge for which any possibilities are to be deleted
	excludedKnowledge := make(map[string]struct{})
	for poss := range f {
		excludedKnowledge[c.knowledgeByPossibility[poss]] = struct{}{}
	}

	possibilitiesToDelete := make(Statement[P])
	for poss, knowledge := range c.knowledgeByPossibility {
		if _, ok := excludedKnowledge[knowledge]; ok {
			possibilitiesToDelete[poss] = struct{}{}
		}
	}
	for poss := range possibilitiesToDelete {
		delete(c.puzzle.externalPossibilities, poss)
		delete(c.possibilities, poss)
	}

	return possibilitiesToDelete
}

func (c *Character[P]) Says(f Statement[P]) {
	for poss := range f {
		delete(c.puzzle.internalPossibilities, poss)
		delete(c.puzzle.externalPossibilities, poss)
		for _, c := range c.puzzle.characters {
			delete(c.possibilities, poss)
			delete(c.knowledgeByPossibility, poss)
		}
	}
}

type Valuation[T any] func(T) int
type Condition[T any] func(T) bool

type intPair struct {
	a, b int
}

func (p intPair) String() string {
	return fmt.Sprintf("(%d, %d)", p.a, p.b)
}

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

func hasNumberDivisibleBy(n int) Condition[intPair] {
	return func(p intPair) bool {
		return p.a%n == 0 || p.b%n == 0
	}
}

func productIsDivisibleBy(n int) Condition[intPair] {
	return func(p intPair) bool {
		return (p.a*p.b)%n == 0
	}
}

func sumIsDivisibleBy(n int) Condition[intPair] {
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
