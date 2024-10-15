package knapsack

type Packable interface {
	Weight() int64
	Value() int64
	Id() string
}

type Sack struct {
	Items    []Packable
	Capacity int64
}

func NewKnapsack(limit int64) *Sack {
	s := &Sack{
		Items:    make([]Packable, 0),
		Capacity: limit,
	}
	return s
}

func (k *Sack) Add(i Packable) {
	k.Items = append(k.Items, i)
}

func (k *Sack) Solve() *Sack {
	items := k.Items
	capacity := k.Capacity

	// We store our working solutions in matrices of N+1 x M+1, where N is the number
	// of items and M is the capacity. We add 1 so we can index from 0.
	// `values` stores the sum of a set of items' values.
	values := make([][]int64, len(items)+1)
	for i := range values {
		values[i] = make([]int64, capacity+1)
	}

	// `keep` stores a matrix of bits, 1 meaning we want to keep the item in this
	// combination, 0 means we'll leave it.
	keep := make([][]int, len(items)+1)
	for i := range keep {
		keep[i] = make([]int, capacity+1)
	}

	// Initially, we'll set all combinations in both `values` and `keep` to 0.
	for i := int64(0); i < capacity+1; i++ {
		values[0][i] = 0
		keep[0][i] = 0
	}

	for i := 0; i < len(items)+1; i++ {
		values[i][0] = 0
		keep[i][0] = 0
	}

	// Simply put, for every item in `items` we want to know whether it will
	// fit in our sack for every capacity from 0 to `capacity`.
	// We know that with 0 items or 0 capacity, no outcome is possible, so start
	// from item 1 and capacity of 1.
	for i := 1; i <= len(items); i++ {
		for c := int64(1); c <= capacity; c++ {

			// Does the item fit at this capacity?
			itemFits := (items[i-1].Weight() <= c)
			if !itemFits {
				continue // skip this iteration
			}

			// Is the value of the item, plus the (previously calculated) value of
			// any remaining space after the addition of this item, greater than the
			// value gained from the previous item?
			maxValueAtThisCapacity := items[i-1].Value() + values[i-1][c-items[i-1].Weight()]
			previousValueAtThisCapacity := values[i-1][c]

			// If the max value to be gained by using this item at this level of
			// capacity is greater than the value to be gained from using the previous
			// item at this capacity, then we want to use this item and keep it.
			// Otherwise, we'll just use the previous item's combination.
			if itemFits && (maxValueAtThisCapacity > previousValueAtThisCapacity) {
				values[i][c] = maxValueAtThisCapacity
				keep[i][c] = 1
			} else {
				values[i][c] = previousValueAtThisCapacity
				keep[i][c] = 0
			}
		}
	}

	// We've now calculated the maximum value to be gained from a combination of
	// items. The maximum value will live at `values[len(items)][capacity]`
	// We now want to loop through our `keep` array and return the indices that
	// point to the specific items to pack into our Knapsack.
	n := len(items)
	c := capacity
	var indices []int64

	for n > 0 {
		if keep[n][c] == 1 {
			indices = append(indices, int64(n-1))
			c -= items[n-1].Weight()
		}
		n--
	}

	//return indices
	new_sack := NewKnapsack(capacity)
	for _, i := range indices {
		new_sack.Items = append(new_sack.Items, k.Items[i])
	}
	return new_sack
}
