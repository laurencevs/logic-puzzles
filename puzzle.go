package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
