package main

// Statements are represented by the set of internal possibilities they rule out.
type Statement[T comparable] map[T]struct{}

type Valuation[T any] func(T) int
type Condition[T any] func(T) bool

func (p *Puzzle[P]) ValuationEquals(v Valuation[P], value int) Statement[P] {
	return nil
}

func (c *Character[P]) KnowsAnswer() Statement[P] {
	return c.HasNPossibilities(func(n int) bool { return n == 1 })
}

func (c *Character[P]) DoesNotKnowAnswer() Statement[P] {
	return c.HasNPossibilities(func(n int) bool { return n > 1 })
}

func (c *Character[P]) HasNPossibilities(condition Condition[int]) Statement[P] {
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

func (p *Puzzle[P]) Satisfies(condition Condition[P]) Statement[P] {
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
