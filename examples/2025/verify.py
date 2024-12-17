import collections
import itertools

"""
Two numbers are drawn from a giant bingo machine containing 2025 balls labelled
from 1 to 2025. Stifado is told their sum, Pastitsio their product, and
Dolmadakia the difference between them.

Dolmadakia says "I can tell that Stifado doesn't know what the numbers are."

Stifado thinks for a moment, before replying "And I can tell that Pastitsio
doesn't know, either."

Pastitsio nods. "Quite right... But I can tell that one of the numbers is a
side length in a Pythagorean triangle in which the square of one of the sides
is 2025."

Stifado chimes in "Oh, now I know what they are!"

Pastitsio and Dolmadakia sit in silence. Neither of them can figure out what
the numbers are.

What are the two numbers?

-------------------------

Answer: 59, 108
"""

def apply_filters(filters, possibilities):
    for f in filters:
        possibilities = filter(f, possibilities)
    return possibilities

def diff(p):
    return p[0]-p[1]

def prod(p):
    return p[0]*p[1]

pythagorean_sides_2025 = set([
    24, 27, 28, 36, 45, 51, 53, 60, 75, 108, 117, 200, 205, 336, 339, 1012, 1013
])

def solve_2025():
    original_possibilities = [(i+1, j+1) for i in range(2025) for j in range(i)]

    possibilities_0_by_sum = collections.defaultdict(list)
    possibilities_0_by_prod = collections.defaultdict(list)
    possibilities_0_by_diff = collections.defaultdict(list)
    for poss in original_possibilities:
        possibilities_0_by_sum[sum(poss)].append(poss)
        possibilities_0_by_prod[prod(poss)].append(poss)
        possibilities_0_by_diff[diff(poss)].append(poss)

    assert(sum(len(x) for x in possibilities_0_by_diff.values()) ==
           sum(len(x) for x in possibilities_0_by_prod.values()) ==
           sum(len(x) for x in possibilities_0_by_sum.values()) == 2049300)
    
    # Dolmadakia says "I can tell that Stifado doesn't know what the numbers are."
    # sum(p) corresponds to multiple possibilities
    # diff(p) only corresponds to possibilities for which the sum corresponds to multiple possibilities
    # => delete all possibilities p for which diff(p) corresponds to any possibility for which the sum corresponds to only one possibility

    allowed_diff_1 = set()
    possibilities_1_by_diff = dict()
    for d, poss in possibilities_0_by_diff.items():
        if all(len(possibilities_0_by_sum[sum(p)]) > 1 for p in poss):
            allowed_diff_1.add(d)
            possibilities_1_by_diff[d] = poss

    # sum(p) corresponds to some possibility for which the diff only corresponds to possibilities for which the sum corresponds to multiple possibilities
    # This already implies sum(p) corresponds to multiple possibilities
    possibilities_1_by_sum = dict()
    for s, poss in possibilities_0_by_sum.items():
        possibilities_1_by_sum[s] = list(filter(lambda p: diff(p) in allowed_diff_1, poss))

    # prod(p) corresponds to some possibility for which the diff only corresponds to possibilities for which the sum corresponds to multiple possibilities
    possibilities_1_by_prod = dict()
    for pr, poss in possibilities_0_by_prod.items():
        possibilities_1_by_prod[pr] = list(filter(lambda p: diff(p) in allowed_diff_1, poss))

    assert(sum(len(x) for x in possibilities_1_by_diff.values()) ==
           sum(len(x) for x in possibilities_1_by_prod.values()) ==
           sum(len(x) for x in possibilities_1_by_sum.values()) == 2045253)

    # Stifado thinks for a moment, before replying "And I can tell that Pastitsio doesn't know, either."
    # prod(p) corresponds to multiple possibilities
    # sum(p) only corresponds to possibilities for which the prod corresponds to multiple possibilities
    # => delete all possibilities p for which sum(p) corresponds to any possibility for which the prod corresponds to only one possibility

    allowed_sum_2 = set()
    possibilities_2_by_sum = dict()
    for s, poss in possibilities_1_by_sum.items():
        if all(len(possibilities_1_by_prod[prod(p)]) > 1 for p in poss):
            allowed_sum_2.add(s)
            possibilities_2_by_sum[s] = poss
    
    possibilities_2_by_prod = dict()
    for pr, poss in possibilities_1_by_prod.items():
        possibilities_2_by_prod[pr] = list(filter(lambda p: sum(p) in allowed_sum_2, poss))
    
    possibilities_2_by_diff = dict()
    for d, poss in possibilities_1_by_diff.items():
        possibilities_2_by_diff[d] = list(filter(lambda p: sum(p) in allowed_sum_2, poss))

    assert(sum(len(x) for x in possibilities_2_by_diff.values()) ==
           sum(len(x) for x in possibilities_2_by_prod.values()) ==
           sum(len(x) for x in possibilities_2_by_sum.values()) == 98686)

    # Pastitsio nods. "Quite right... But I can tell that one of the numbers is a side length in a Pythagorean triangle in which the square of one of the sides is 2025."
    # prod(p) only corresponds to possibilities for which a number is in pythagorean_sides_2025
    # => delete all possibilities p for which prod(p) corresponds to any possibility without a number in pythagorean_sides_2025

    allowed_prod_3 = set()
    possibilities_3_by_prod = dict()
    for pr, poss in possibilities_2_by_prod.items():
        if all(p[0] in pythagorean_sides_2025 or p[1] in pythagorean_sides_2025 for p in poss):
            allowed_prod_3.add(pr)
            possibilities_3_by_prod[pr] = poss
    
    possibilities_3_by_sum = dict()
    for s, poss in possibilities_2_by_sum.items():
        possibilities_3_by_sum[s] = list(filter(lambda p: prod(p) in allowed_prod_3, poss))
    
    possibilities_3_by_diff = dict()
    for d, poss in possibilities_2_by_diff.items():
        possibilities_3_by_diff[d] = list(filter(lambda p: prod(p) in allowed_prod_3, poss))

    assert(sum(len(x) for x in possibilities_3_by_diff.values()) ==
           sum(len(x) for x in possibilities_3_by_prod.values()) ==
           sum(len(x) for x in possibilities_3_by_sum.values()) == 897)

    # Stifado chimes in "Oh, now I know what they are!"
    # sum(p) is unique
    # => delete all possibilities p for which sum(p) is non-unique

    allowed_sum_4 = set()
    possibilities_4_by_sum = dict()
    for s, poss in possibilities_3_by_sum.items():
        if len(poss) == 1:
            possibilities_4_by_sum[s] = poss
            allowed_sum_4.add(s)
    
    possibilities_4_by_diff = dict()
    for d, poss in possibilities_3_by_diff.items():
        possibilities_4_by_diff[d] = list(filter(lambda p: sum(p) in allowed_sum_4, poss))
    
    possibilities_4_by_prod = dict()
    for pr, poss in possibilities_3_by_prod.items():
        possibilities_4_by_prod[pr] = list(filter(lambda p: sum(p) in allowed_sum_4, poss))

    assert(sum(len(x) for x in possibilities_4_by_diff.values()) ==
           sum(len(x) for x in possibilities_4_by_prod.values()) ==
           sum(len(x) for x in possibilities_4_by_sum.values()) == 69)

    # Pastitsio and Dolmadakia sit in silence. Neither of them can figure out what the numbers are.

    candidates = itertools.chain(*possibilities_4_by_sum.values())
    solutions = list(filter(lambda p: len(possibilities_4_by_prod[prod(p)]) > 1 and len(possibilities_4_by_diff[diff(p)]) > 1, candidates))
    assert(len(solutions) == 1)
    assert(solutions[0] == (108, 59))
    print("Got correct solution")

def verify_2025():
    nums = (108, 59)

    """
    Two numbers are drawn from a giant bingo machine containing 2025 balls labelled
    from 1 to 2025. Stifado is told their sum, Pastitsio their product, and
    Dolmadakia the difference between them.
    """
    original_possibilities = [(i+1, j+1) for i in range(2025) for j in range(i)]
    stifado_possibilities = filter(lambda p: p[0]+p[1] == nums[0]+nums[1], public_possibilities)
    pastitsio_possibilities = filter(lambda p: p[0]*p[1] == nums[0]*nums[1], public_possibilities)
    dolmadakia_possibilities = filter(lambda p: p[0]-p[1] == nums[0]-nums[1], public_possibilities)

    public_possibilities = [(i+1, j+1) for i in range(2025) for j in range(i)]
    public_filters = []

    """
    Dolmadakia says "I can tell that Stifado doesn't know what the numbers are."
    """
    assert(len(list(stifado_possibilities)) > 1)
    # All of Dolmadakia's possibilities have non-unique sums.
    # Unique sums are: 3, 4, 2028, 2029
    sum_counts = collections.Counter(map(lambda p: p[0]+p[1], public_possibilities))
    assert(all(sum_counts[sum(p)] > 1 for p in dolmadakia_possibilities))
    # Others now know this.
    differences_corresponding_to_non_unique_sums = set(map(diff, filter(lambda p: sum_counts[sum(p)] == 1, public_possibilities)))

    public_filter_1 = lambda p: diff(p) not in differences_corresponding_to_non_unique_sums
    stifado_possibilities = filter(public_filter_1, stifado_possibilities)
    pastitsio_possibilities = filter(public_filter_1, pastitsio_possibilities)
    public_possibilities = filter(public_filter_1, public_possibilities)
    public_filters.append(public_filter_1)

    """
    Stifado thinks for a moment, before replying "And I can tell that Pastitsio
    doesn't know, either."
    """
    assert(len(list(pastitsio_possibilities)) > 1)
    stifado_possible_products = set(map(prod, stifado_possibilities))
    pastitsio_possibilities_acc_to_stifado = apply_filters(public_filters + [lambda p: prod(p) in stifado_possible_products], original_possibilities)
    product_counts = collections.Counter(map(prod, pastitsio_possibilities_acc_to_stifado))
    assert(all(x > 1 for x in product_counts.values()))
    sums_corresponding_to_non_unique_products = set(map(sum, filter(lambda p: prod)))
    public_filter_2 = lambda p: sum(p) 

if __name__ == "__main__":
    solve_2025()
    verify_2025()
