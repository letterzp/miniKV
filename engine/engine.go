package engine

import (
	"testing"

	internal "github.com/letterbeezps/miniKV/internal"
	"github.com/stretchr/testify/assert"
)

type Engine interface {
	Get(key string) ([]byte, bool)

	Set(key string, value []byte)

	Delete(key string)

	Scan(start, end internal.Bound, iter func(key string, value []byte) bool)

	Reverse(start, end internal.Bound, iter func(key string, value []byte) bool)

	Iter(start, end internal.Bound) Iterator

	ReverseIter(start, end internal.Bound) Iterator
}

type Iterator interface {
	Value() []byte

	Key() string

	IsValid() bool

	Next()
}

//////// test ////////////

func EngineTestGetSet(engine Engine, t *testing.T) {
	ret, ok := engine.Get("a")
	assert.False(t, ok)
	assert.Equal(t, []byte{}, ret)

	engine.Set("a", []byte("aaa"))
	ret, ok = engine.Get("a")
	assert.True(t, ok)
	assert.Equal(t, []byte("aaa"), ret)

	engine.Set("b", []byte("bbb"))
	ret, ok = engine.Get("a")
	assert.True(t, ok)
	assert.Equal(t, []byte("aaa"), ret)
	ret, ok = engine.Get("b")
	assert.True(t, ok)
	assert.Equal(t, []byte("bbb"), ret)

	ret, ok = engine.Get("A")
	assert.False(t, ok)
	assert.Equal(t, []byte{}, ret)
	ret, ok = engine.Get("B")
	assert.False(t, ok)
	assert.Equal(t, []byte{}, ret)

	engine.Delete("a")
	ret, ok = engine.Get("a")
	assert.False(t, ok)
	assert.Equal(t, []byte{}, ret)
	ret, ok = engine.Get("b")
	assert.True(t, ok)
	assert.Equal(t, []byte("bbb"), ret)

	engine.Set("b", []byte{})
	ret, ok = engine.Get("b")
	assert.True(t, ok)
	assert.Equal(t, []byte{}, ret)
}

type kvPair struct {
	key string
	val []byte
}

func ScanTest(engine Engine, t *testing.T) {
	datas := []kvPair{
		{"a", []byte("a")},
		{"aa", []byte("aa")},
		{"aaa", []byte("aaa")},
		{"aaaa", []byte("aaaaa")},
		{"b", []byte("b")},
		{"ba", []byte("ba")},
		{"bb", []byte("bb")},
		{"bc", []byte("bc")},
		{"c", []byte("c")},
	}
	for _, data := range datas {
		engine.Set(data.key, data.val)
	}

	///////////////////////// iter ////////////////////

	// start : NoBound
	// end : NoBound
	iterator := engine.Iter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("", internal.NoBound),
	)
	i := 0
	for iterator.IsValid() {
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : NoBound
	// end : Include
	iterator = engine.Iter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("bc", internal.Include),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 7)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	iterator = engine.Iter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("bd", internal.Include),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 7)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : NoBound
	// end : Exclude
	iterator = engine.Iter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("bc", internal.Exclude),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 6)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	iterator = engine.Iter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("bd", internal.Exclude),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 7)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : Include
	// end : NoBound
	iterator = engine.Iter(
		internal.NewBound("a", internal.Include),
		internal.NewBound("", internal.NoBound),
	)
	i = 0
	for iterator.IsValid() {
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	iterator = engine.Iter(
		internal.NewBound("aab", internal.Include),
		internal.NewBound("", internal.NoBound),
	)
	i = 4
	for iterator.IsValid() {
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : Include
	// end : Include
	iterator = engine.Iter(
		internal.NewBound("a", internal.Include),
		internal.NewBound("ba", internal.Include),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 5)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	iterator = engine.Iter(
		internal.NewBound("a", internal.Include),
		internal.NewBound("baa", internal.Include),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 5)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : Include
	// end : Exclude
	iterator = engine.Iter(
		internal.NewBound("a", internal.Include),
		internal.NewBound("ba", internal.Exclude),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 4)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	iterator = engine.Iter(
		internal.NewBound("a", internal.Include),
		internal.NewBound("baa", internal.Exclude),
	)
	i = 0
	for iterator.IsValid() {
		assert.LessOrEqual(t, i, 5)
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}

	// start : exclude
	// end : NoBound
	iterator = engine.Iter(
		internal.NewBound("a", internal.Exclude),
		internal.NewBound("", internal.NoBound),
	)
	i = 1
	for iterator.IsValid() {
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		iterator.Next()
		i++
	}
}

func ReverseScanTest(engine Engine, t *testing.T) {
	datas := []kvPair{
		{"a", []byte("a")},
		{"aa", []byte("aa")},
		{"aaa", []byte("aaa")},
		{"aaaa", []byte("aaaaa")},
		{"b", []byte("b")},
		{"ba", []byte("ba")},
		{"bb", []byte("bb")},
		{"bc", []byte("bc")},
		{"c", []byte("c")},
	}
	for _, data := range datas {
		engine.Set(data.key, data.val)
	}

	////////////////////////// reverse iter//////////////////

	iterator := engine.ReverseIter(
		internal.NewBound("", internal.NoBound),
		internal.NewBound("", internal.NoBound),
	)
	i := len(datas) - 1
	for iterator.IsValid() {
		assert.Equal(t, datas[i].key, iterator.Key())
		assert.Equal(t, datas[i].val, iterator.Value())
		// t.Log("key: ", iterator.Key(), " value: ", iterator.Value())
		iterator.Next()
		i--
	}
}
