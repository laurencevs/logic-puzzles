package puzzles

// KnowledgeState represents a character's understanding of the puzzle in all
// possible cases. That is, for every possible solution s, KnowledgeState
// should be able to tell us what a character would consider possible in a
// world where s is the solution (this information is exposed by
// PossibilitiesInCase(s P) []P).
type KnowledgeState[P any] interface {
	// Methods for updating a character's knowledge

	// Initialise the KnowledgeState with a set of possible worlds (solutions)
	Initialise(solutionSpace []P)

	// Update the KnowledgeState by providing a statement s which is to be
	// taken to be true
	Update(s Statement[P])

	// Methods for inspecting a character's knowledge as an external observer

	// The solutions that a character with this state of knowledge would
	// consider possible in a world where s is the solution. If s is not in
	// PossibilitiesInCase(s) then such a character knows that s is not the
	// solution because of statements made so far.
	PossibilitiesInCase(s P) []P

	// All the solutions that a character with this state of knowledge might
	// consider possible depending on what the actual solution is (that is, in
	// some possible world)
	Possibilities() []P

	// Statements made by a character with this state of knowledge

	KnowsAnswer() Statement[P]
	DoesNotKnowAnswer() Statement[P]

	Knows(s Statement[P]) Statement[P]
	DoesNotKnow(s Statement[P]) Statement[P]

	KnowsWhether(s Statement[P]) Statement[P]
	DoesNotKnowWhether(s Statement[P]) Statement[P]

	KnowsNormalisedAnswer(normalise func(P) P) Statement[P]
	KnowsAnswerGivenNormalised(normalise func(P) P) Statement[P]
}

// Statement represents a statement about the solution to a Puzzle.
// ConsistentWith should return true for any possibility that the Statement
// does not directly contradict.
type Statement[P any] interface {
	ConsistentWith(p P) bool
	Not() Statement[P]
}

type Condition[P any] func(P) bool

func (c Condition[P]) ConsistentWith(p P) bool {
	return c(p)
}

func (c Condition[P]) Not() Statement[P] {
	return Condition[P](func(p P) bool {
		return !c(p)
	})
}

func filterInPlace[P any](s Statement[P], l *[]P) {
	i := 0
	for _, p := range *l {
		if s.ConsistentWith(p) {
			(*l)[i] = p
			i++
		}
	}
	*l = (*l)[:i]
}
