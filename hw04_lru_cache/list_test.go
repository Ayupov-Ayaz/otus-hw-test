package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := newList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := newList()

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

func TestList_Push(t *testing.T) {
	const count = 10

	push := func(add func(v interface{}) *ListItem) {
		for i := 0; i < count; i++ {
			add(i)
		}
	}

	t.Helper()

	check := func(t *testing.T, list List, exps []int) {
		t.Helper()

		curr := list.Front()
		require.Nil(t, curr.Prev)
		require.Equal(t, count, list.Len())
		require.Equal(t, count, len(exps))

		for _, i := range exps {
			item := curr.Value.(int)
			require.Equal(t, i, item)
			if curr.Next == nil {
				break
			}
			curr = curr.Next
		}
	}

	t.Run("PushFront", func(t *testing.T) {
		list := newList()
		push(list.PushFront)

		check(t, list, []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0})
	})

	t.Run("PushBack", func(t *testing.T) {
		list := newList()
		push(list.PushBack)
		require.Equal(t, list.Len(), count)
		check(t, list, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	})
}
