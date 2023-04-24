package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
	// Test for full coverage
	t.Run("first pushing back", func(t *testing.T) {
		lst := NewList()

		lst.PushBack(100)              // [100]
		lst.PushBack(50)               // [100, 50]
		lst.PushFront(15)              // [15, 100, 50]
		prelastItem := lst.Back().Prev // 100

		lst.PushFront(90) // [90, 15, 100, 50]
		lst.PushBack(27)  // [90, 15, 100, 50, 15]
		require.Equal(t, 5, lst.Len())
		require.Equal(t, 100, prelastItem.Value)

		lst.MoveToFront(prelastItem) // [100, 90, 15, 50, 15]
		require.Equal(t, prelastItem.Value, lst.Front().Value)

		lst.Remove(lst.Front()) // [90, 15, 50, 15]
		require.Equal(t, 90, lst.Front().Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
