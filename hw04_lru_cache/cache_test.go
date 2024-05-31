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

func TestFixedCacheSize(t *testing.T) {
	vals := []int{1, 2, 3, 4}
	c := NewCache(len(vals) - 1)
	for _, v := range vals {
		c.Set(Key(strconv.Itoa(v)), v)
	}

	for _, v := range vals[1:] {
		_, ok := c.Get(Key(strconv.Itoa(v)))
		require.True(t, ok, "Expect true for %d", v)
	}
	_, ok := c.Get(Key(strconv.Itoa(vals[0])))
	require.False(t, ok, "Expect false for %d", vals[0])
}

func TestLRU(t *testing.T) {
	vals := []int{1, 2, 3}
	additionalValue := 4
	droppedValue := 2
	access := []int{1, 2, 3, 1, 1, 3, 3, 1}
	c := NewCache(len(vals))
	for _, v := range vals {
		c.Set(Key(strconv.Itoa(v)), v)
	}

	for _, v := range access {
		c.Get(Key(strconv.Itoa(v)))
	}

	c.Set(Key(strconv.Itoa(additionalValue)), additionalValue)
	vals = append(vals, additionalValue)

	for _, v := range vals {
		if v == droppedValue {
			_, ok := c.Get(Key(strconv.Itoa(v)))
			require.False(t, ok, "Expect false for %d", v)
			continue
		}
		_, ok := c.Get(Key(strconv.Itoa(v)))
		require.True(t, ok, "Expect true for %d", v)
	}
}
