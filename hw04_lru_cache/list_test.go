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

func TestList_Remove(t *testing.T) {
	t.Run("when item is nil, then do nothing", func(t *testing.T) {
		l := NewList()
		l.Remove(nil)

		require.Equal(t, l.Len(), 0)
	})

	t.Run("when item is front, then front must be nil", func(t *testing.T) {
		l := NewList()

		item := l.PushFront(1)

		require.Equal(t, item, l.Front())

		l.Remove(item)

		require.Nil(t, l.Front())
	})
}

func TestList_PushBack(t *testing.T) {
	t.Run("when list is empty, then item must be front and back", func(t *testing.T) {
		l := NewList()

		item := l.PushBack(1)

		require.Equal(t, l.Front(), item)
		require.Equal(t, l.Back(), item)
	})
}

func TestList_MoveToFront(t *testing.T) {
	t.Run("when item is nil, then do nothing", func(t *testing.T) {
		l := NewList()

		l.MoveToFront(nil)
		require.Equal(t, l.Len(), 0)
	})

	t.Run("move second item", func(t *testing.T) {
		l := NewList()

		first := l.PushBack(1)
		second := l.PushBack(2)
		third := l.PushBack(3)

		l.MoveToFront(second)

		require.Equal(t, l.Front(), second)
		require.Equal(t, l.Back(), third)
		require.Equal(t, l.Front().Next, first)
	})
}
