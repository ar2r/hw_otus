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
		require.False(t, ok, "Ключ 'aaa' должен отсутствовать в пустом кэше")

		_, ok = c.Get("aaa")
		require.False(t, ok, "Вызов Get не должен приводить к записи ключа в кэш")

		// Ничего другого не придумал, как проверить на панику)))
		require.NotPanics(t, func() {
			c.Clear()
		}, "Очистка пустого кэша не должна вызывать паники")
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

	t.Run("Очистка кэша", func(t *testing.T) {
		c := NewCache(3)
		c.Set("meKey", 100)
		c.Clear()
		_, ok := c.Get("myKey")
		require.False(t, ok, "Кэш должен быть пустым после вызова Clear")
	})

	t.Run("Превышение емкости", func(t *testing.T) {
		c := NewCache(2)
		c.Set("a", 100)
		c.Set("b", 200)
		c.Set("c", 300)

		_, ok := c.Get("a")
		require.False(t, ok, "Ключ 'a' должен быть удален из кэша")
		_, ok = c.Get("b")
		require.True(t, ok, "Ключ 'b' должен остаться в кэше")
		_, ok = c.Get("c")
		require.True(t, ok, "Ключ 'c' должен остаться в кэше")
	})

	t.Run("Выталкивание неиспользуемых", func(t *testing.T) {
		c := NewCache(3)
		c.Set("Сбер", 100)
		c.Set("Т-банк", 200)
		c.Set("PayPal", "редко используется")

		// Читаем две первых и добавляем новое значение
		c.Get("Сбер")
		c.Set("Т-Банк", 201)
		c.Set("Альфа", 300)

		_, ok := c.Get("PayPal")
		require.False(t, ok, "Ключ 'PayPal' должен быть удален из кэша")
		_, ok = c.Get("Сбер")
		require.True(t, ok, "Ключ 'Сбер' должен остаться в кэше")
		_, ok = c.Get("Т-Банк")
		require.True(t, ok, "Ключ 'Т-Банк' должен остаться в кэше")
		_, ok = c.Get("Альфа")
		require.True(t, ok, "Ключ 'Альфа' должен остаться в кэше")
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
