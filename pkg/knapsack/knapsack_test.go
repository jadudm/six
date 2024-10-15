package knapsack

import (
	"fmt"
	"testing"

	gtst "com.jadud.search.six/pkg/types"
	"github.com/stretchr/testify/assert"
)

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}
func TestKnapsack1(t *testing.T) {
	s1 := NewKnapsack(10)
	s1.Add(&gtst.SqliteFile{Name: "Alice", Size: 2})
	s1.Add(&gtst.SqliteFile{Name: "Bob", Size: 3})
	s1.Add(&gtst.SqliteFile{Name: "Clarice", Size: 5})
	s1.Add(&gtst.SqliteFile{Name: "Denzel", Size: 1})
	soln := s1.Solve()
	fmt.Println(soln)
	soln_strings := make([]string, 0)
	for _, s := range soln.Items {
		soln_strings = append(soln_strings, s.Id())
	}

	assert.Equal(t, []string{"Clarice", "Bob", "Alice"}, soln_strings)
}

func TestKnapsack2(t *testing.T) {
	s1 := NewKnapsack(10)
	s1.Add(&gtst.SqliteFile{Name: "Alice", Size: 4})
	s1.Add(&gtst.SqliteFile{Name: "Bob", Size: 3})
	s1.Add(&gtst.SqliteFile{Name: "Clarice", Size: 5})
	s1.Add(&gtst.SqliteFile{Name: "Denzel", Size: 1})
	soln := s1.Solve()
	fmt.Println(soln)
	soln_strings := make([]string, 0)
	for _, s := range soln.Items {
		soln_strings = append(soln_strings, s.Id())
	}

	assert.Equal(t, []string{"Denzel", "Bob", "Alice"}, soln_strings)
}

func TestKnapsack3(t *testing.T) {
	s1 := NewKnapsack(10)
	s1.Add(&gtst.SqliteFile{Name: "Alice", Size: 1})
	s1.Add(&gtst.SqliteFile{Name: "Bob", Size: 1})
	s1.Add(&gtst.SqliteFile{Name: "Clarice", Size: 1})
	s1.Add(&gtst.SqliteFile{Name: "Denzel", Size: 3})
	s1.Add(&gtst.SqliteFile{Name: "Ernie", Size: 9})
	soln := s1.Solve()
	fmt.Println(soln)
	soln_strings := make([]string, 0)
	for _, s := range soln.Items {
		soln_strings = append(soln_strings, s.Id())
	}

	assert.Equal(t, []string{"Denzel", "Clarice", "Bob", "Alice"}, soln_strings)
}

func TestKnapsack4(t *testing.T) {
	s1 := NewKnapsack(10)
	s1.Add(&gtst.SqliteFile{Name: "Denzel", Size: 3})
	s1.Add(&gtst.SqliteFile{Name: "Ernie", Size: 9})
	soln := s1.Solve()
	fmt.Println(soln)
	soln_strings := make([]string, 0)
	for _, s := range soln.Items {
		soln_strings = append(soln_strings, s.Id())
	}

	assert.Equal(t, []string{"Denzel"}, soln_strings)
}
