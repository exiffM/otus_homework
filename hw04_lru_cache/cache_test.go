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
		c := NewCache(3)

		wasIn := c.Set("first", 10) // ["first" - 10]
		require.False(t, wasIn)

		wasIn = c.Set("second", 20) // ["second" - 20, "first" - 10]
		require.False(t, wasIn)

		wasIn = c.Set("third", 30) // ["third" - 30, "second" - 20, "first" - 10]
		require.False(t, wasIn)

		wasIn = c.Set("fourth", 40) // ["fourth" - 40, "third" - 30, "second" - 20]
		require.False(t, wasIn)

		val, ok := c.Get("fourth") // ["fourth" - 40, "third" - 30, "second" - 20]
		require.True(t, ok)
		require.Equal(t, 40, val)

		val, ok = c.Get("first")
		require.False(t, ok)
		require.Nil(t, val)

		c.Get("third")             // ["third" - 30, "fourth" - 40, "second" - 20]
		c.Get("second")            // ["second" - 20, "third" - 30, "fourth" - 40]
		wasIn = c.Set("first", 10) // ["first" - 10, "second" - 20, "third" - 30]
		require.False(t, wasIn)

		_, ok = c.Get("fourth")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(3)

		wasIn := c.Set("first", 10) // ["first" - 10]
		require.False(t, wasIn)

		wasIn = c.Set("second", 20) // ["second" - 20, "first" - 10]
		require.False(t, wasIn)

		wasIn = c.Set("third", 30) // ["third" - 30, "second" - 20, "first" - 10]
		require.False(t, wasIn)

		c.Clear()
		require.Equal(t, &lruCache{}, c)
	})
}

func TestCacheMultithreading(_ *testing.T) {
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
