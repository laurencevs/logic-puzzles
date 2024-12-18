## Logic puzzle solver

This is a tool for solving 'perfect logician' puzzles.

### Example

> The integers from 1 to 2024 are written on pieces of paper, which are then put into a (big) hat, and two are drawn at random. Stifado is told their sum, Pastitsio their product, and Dolmadakia the difference between them.
>
> Pastitsio says "The product of the numbers is divisible by 20."
>
> Stifado replies "Then I know that you know that Dolmadakia doesn't know what the numbers are."
>
> "If it's any help," she continues, "the sum of the numbers is divisible by 24."
>
> Pastitsio exclaims "Ah! Now I know that Dolmadakia does know what the numbers are!"
>
> What are the numbers?

With this module, the puzzle above is easy to solve (`examples/readme-example/main.go`):

```go
solutionSpace := intpair.IntPairs(1, 2024, false, false)
puzzle := puzzles.NewPuzzle(solutionSpace)

Stifado := puzzle.NewActorWithKnowledge(intpair.Sum)
Pastitsio := puzzle.NewActorWithKnowledge(intpair.Product)
Dolmadakia := puzzle.NewActorWithKnowledge(intpair.AbsDifference)

Pastitsio.Says(Pastitsio.KnowsHolds(intpair.ProductIsDivisibleBy(20)))
Stifado.Says(Stifado.Knows(Pastitsio.Knows(Dolmadakia.DoesNotKnowAnswer())))
Stifado.Says(Stifado.KnowsHolds(intpair.SumIsDivisibleBy(24)))
Pastitsio.Says(Pastitsio.Knows(Dolmadakia.KnowsAnswer()))

fmt.Println(puzzles.SprintPossibilities(puzzle.ExternalPossibilities()))
```

```
$ go run ./examples/readme-example
Puzzle has 1 remaining possibility: (10, 1982)
```

### Puzzle specification

This module can be used to solve problems matching the following specification.

* An element $s$ is chosen from a finite set $S$, the solution space.
* The puzzle contains $n$ actors (or agents), $A_1, \ldots, A_n$.
* Each actor $A_i$ is told the value of $f_i(s)$ for some function $f_i: S \rightarrow X_i$. Initially, this is the only information each actor has on $s$ (apart from the fact that it was chosen from $S$).
* The functions $f_i$, the fact that actor $A_i$ initially knows the value of $f_i(s)$ and nothing else about $s$, and the fact that $s \in S$ are all [*common knowledge*](https://en.wikipedia.org/wiki/Common_knowledge_(logic)).
* The actors make a series of true statements based on their own knowledge (including their knowledge of others' knowledge, and so on), with the truth of each statement becoming common knowledge.
* Each actor is a 'perfect logician', who will instantly and correctly deduce any information about $s$ implied by each statement. This fact is also common knowledge.

### Acknowledgement

This project was inspired by the last question on [this example sheet](https://www.dpmms.cam.ac.uk/study/IA/Numbers%2BSets/2023-2024/numset1_2023.pdf) from the Cambridge mathematics tripos, which seems to have been based on the [sum and product puzzle](https://en.wikipedia.org/wiki/Sum_and_Product_Puzzle).
