## Logic puzzle generator and solver

This is a tool for generating and solving puzzles like the following.

> The numbers from 1 to 2024 are put into a hat, and two are drawn at random. Paul is told their product, Sophie their sum, and Dave their difference.
>
> Paul says "The product of the numbers is divisible by 20."
>
> Sophie replies "Then I know that you know that Dave doesn't know what the numbers are."
>
> "If it's any help," she continues, "the sum of the numbers is divisible by 24."
>
> Paul exclaims "Ah! Now I know that Dave does know what the numbers are!"
>
> What are the numbers?

In this module, the puzzle above is easy to solve:

```go
possibilities := UnorderedIntPairs(1, 2024, false)
puzzle := NewPuzzle(possibilities)

Sophie := puzzle.NewCharacter("S")
Sophie.KnowsValueOf(sum)
Paul := puzzle.NewCharacter("P")
Paul.KnowsValueOf(product)
Dave := puzzle.NewCharacter("D")
Dave.KnowsValueOf(absDifference)

Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))
Sophie.Says(Sophie.Knows(Paul.Knows(Dave.DoesNotKnowAnswer())))
Sophie.Says(puzzle.Satisfies(sumIsDivisibleBy(24)))
Paul.Says(Paul.Knows(Dave.KnowsAnswer()))

puzzle.PrintPossibilities()
```

Inspired by the last question on [this example sheet](https://www.dpmms.cam.ac.uk/study/IA/Numbers%2BSets/2023-2024/numset1_2023.pdf) from the Cambridge mathematics tripos.
