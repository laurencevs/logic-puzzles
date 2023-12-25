from collections import defaultdict

def verify_example_2024_6():

    # possibilities := UnorderedIntPairs(1, 2024, false)

    possibilities_1 = [(i, j) for i in range(1, 2025)
                              for j in range(i+1, 2025)]
    assert(len(possibilities_1) == 2047276)

    # Paul.Says(puzzle.Satisfies(productIsDivisibleBy(20)))

    possibilities_2 = [x for x in possibilities_1 if (x[0] * x[1]) % 20 == 0]
    assert(len(possibilities_2) == 367943)

    # Dave.Says(Dave.Knows(Paul.Knows(Sophie.DoesNotKnowAnswer())))

    possibilities_by_sum = defaultdict(list)
    possibilities_by_product = defaultdict(list)
    possibilities_by_difference = defaultdict(list)
    for x in possibilities_2:
        possibilities_by_sum[sum(x)].append(x)
        possibilities_by_product[x[0]*x[1]].append(x)
        possibilities_by_difference[x[1]-x[0]].append(x)
    # Sophie.DoesNotKnowAnswer() => sum does not correspond to only one
    # possibility
    sums_ruled_out = set(s for s in possibilities_by_sum
                         if len(possibilities_by_sum[s]) == 1)
    # Paul.Knows(Sophie.DoesNotKnowAnswer()) => product does not correspond to
    # any sum that corresponds to only one possibility
    products_ruled_out = set(p for p in possibilities_by_product
                             if any(sum(x) in sums_ruled_out
                                    for x in possibilities_by_product[p]))
    # Dave.Knows(Paul.Knows(Sophie.DoesNotKnowAnswer())) => difference does not
    # correspond to any product that corresponds to any sum that corresponds to
    # only one possibility
    differences_ruled_out = set(d for d in possibilities_by_difference
                             if any(x[0]*x[1] in products_ruled_out
                                    for x in possibilities_by_difference[d]))
    possibilities_3 = [x for x in possibilities_2
                       if x[1]-x[0] not in differences_ruled_out]
    assert(len(possibilities_3) == 355312)

    # Sophie.Says(Sophie.KnowsAnswer())

    new_possibilities_by_sum = defaultdict(list)
    for x in possibilities_3:
        new_possibilities_by_sum[sum(x)].append(x)
    possibilities_4 = [x for x in possibilities_3
                       if len(new_possibilities_by_sum[sum(x)]) == 1]
    assert(len(possibilities_4) == 35)

    # puzzle.Satisfies(sumIsDivisibleBy(24))
    
    final_possibilities = [x for x in possibilities_4 if sum(x)%24 == 0]
    assert(final_possibilities == [(2010, 2022)])

if __name__ == "__main__":
    verify_example_2024_6()
