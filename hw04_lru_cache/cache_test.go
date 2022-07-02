package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		const count = 10

		list := NewList()
		cache := newLruCache(count, list)
		for i := 0; i < count; i++ {
			ok := cache.Set(Key(strconv.Itoa(i)), i)
			require.False(t, ok)
		}

		cache.Clear()

		for i := 0; i < count; i++ {
			v, ok := cache.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
			require.Nil(t, v)
		}

	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

var makeKey = func(i int) Key {
	return Key(strconv.Itoa(i))
}

func TestLruCache_Set(t *testing.T) {
	const count = 10

	list := NewList()
	cache := newLruCache(10, list)

	for i := 0; i < count; i++ {
		ok := cache.Set(makeKey(i), i)
		require.False(t, ok)
	}

	require.Equal(t, count, list.Len())

	curr := list.Front()
	k := count - 1
	for i := 0; i < count; i++ {
		require.Equal(t, k, curr.Value.(int))

		if curr.Next == nil {
			break
		}

		curr = curr.Next
		k--
	}
	require.Equal(t, 0, k)

	// если элемент присутствует в словаре,
	//то обновить его значение и переместить элемент в начало очереди;
	const multiplier = 100
	for i := 0; i < count; i++ {
		ok := cache.Set(makeKey(i), i*multiplier)
		require.True(t, ok)
		require.Equal(t, i*multiplier, list.Front().Value.(int))
		require.Equal(t, count, list.Len())
	}

	//если элемента нет в словаре, то добавить в словарь и в начало очереди
	//(при этом, если размер очереди больше ёмкости кэша,
	//то необходимо удалить последний элемент из очереди и его значение из словаря);
	for i := 1; i < count+1; i++ {
		// проверяю, что удалится (прямой вызов, чтобы не сдвигать очередь)
		checkKey := makeKey(i - 1)
		_, ok := cache.items[checkKey]
		require.True(t, ok)
		ok = cache.Set(makeKey(i+10), i)
		require.False(t, ok)
		require.Equal(t, 10, len(cache.items))
		require.Equal(t, 10, cache.queue.Len())
		require.Equal(t, i, list.Front().Value.(int))
		//
		_, ok = cache.items[checkKey]
		require.False(t, ok, strconv.Itoa(i))
		//
		head := list.Front()
		require.Equal(t, i, head.Value.(int))
	}

	// проверяем все значения
	curr = list.Front()
	k = 10
	for curr != nil {
		require.Equal(t, k, curr.Value.(int))

		if curr.Next == nil {
			break
		}

		curr = curr.Next
		k--
	}
	require.Equal(t, 1, k)
}

func TestLruCache_Get(t *testing.T) {
	const count = 10

	list := NewList()
	cache := newLruCache(count, list)
	for i := 0; i < count; i++ {
		cache.Set(makeKey(i), i)
	}

	// проверяем порядок
	k := count - 1
	curr := list.Front()
	for curr != nil {
		require.Equal(t, k, curr.Value.(int))

		if curr.Next == nil {
			break
		}

		curr = curr.Next
		k--
	}

	require.Equal(t, 0, k)

	// если элемент присутствует в словаре,
	//то переместить элемент в начало очереди и вернуть его значение и true;
	for i := 0; i < count; i++ {
		v, ok := cache.Get(makeKey(i))
		require.True(t, ok)
		require.Equal(t, i, v.(int))

		head := list.Front()
		require.Equal(t, i, head.Value.(int))
	}

	//если элемента нет в словаре,
	//то вернуть nil и false
	v, ok := cache.Get(makeKey(count))
	require.False(t, ok)
	require.Nil(t, v)
}
