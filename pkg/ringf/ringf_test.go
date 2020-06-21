package ringf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ring2Slice(r *RingF) []interface{} {
	var vals []interface{}

	for i := 0; i < r.Len(); i++ {
		vals = append(vals, r.Value)
		r = r.Next()
	}
	return vals
}

func slice2Ring(vals []interface{}) *RingF {
	r := New(len(vals))

	for i := 0; i < len(vals); i++ {
		r.Value = vals[i]
		r = r.Next()
	}

	return r
}

func int2Interface(ints []int) []interface{} {
	a := make([]interface{}, len(ints))
	for i := range ints {
		a[i] = ints[i]
	}
	return a
}

func TestFilter(t *testing.T) {
	got := slice2Ring(int2Interface([]int{0, 1, 2, 3, 4}))
	got = got.Filter(func(a interface{}) bool {
		return a.(int) > 2
	})

	assert.EqualValues(t, int2Interface([]int{3, 4}), ring2Slice(got))

	got = slice2Ring(int2Interface([]int{12}))
	got = got.Filter(func(_ interface{}) bool {
		return false
	})
	assert.EqualValues(t, New(0), got)

	got = slice2Ring(int2Interface([]int{}))
	got = got.Filter(func(_ interface{}) bool {
		return false
	})
	assert.EqualValues(t, New(0), got)

	got = slice2Ring(int2Interface([]int{1, 3}))
	got = got.Filter(func(a interface{}) bool {
		return a.(int) != 1
	})
	assert.EqualValues(t, int2Interface([]int{3}), ring2Slice(got))
}
