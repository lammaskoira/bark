package util_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/lammaskoira/bark/pkg/util"
)

func TestShouldReturnNilIfEmpty(t *testing.T) {
	t.Parallel()

	s := util.NewStack[*int]()
	require.Nil(t, s.Pop(), "Should return nil if empty")
}

func TestShouldAcceptValuesNormally(t *testing.T) {
	t.Parallel()

	s := util.NewStack[int]()
	s.Push(1)
	s.Push(2)
	s.Push(3)

	require.False(t, s.IsEmpty(), "Should return false")
	require.Equal(t, 3, s.Pop(), "Should return the last value")
	require.Equal(t, 2, s.Pop(), "Should return the last value")
	require.Equal(t, 1, s.Pop(), "Should return the last value")
	require.True(t, s.IsEmpty(), "Should return true")
}

func TestParallelAccess(t *testing.T) {
	t.Parallel()

	s := util.NewStack[int]()

	wg := sync.WaitGroup{}
	iterations := 800
	wg.Add(iterations)

	for i := 0; i < iterations; i++ {
		go func(i int) {
			s.Push(i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	got := 0

	for !s.IsEmpty() {
		got++
		s.Pop()
	}

	require.Equal(t, iterations, got, "Should return the same number of values as pushed")
}
