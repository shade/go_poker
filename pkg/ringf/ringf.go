package ringf

/**

Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

*/

// A RingF is an element of a circular list, or ringf.
// RingFs do not have a beginning or end; a pointer to any ringf element
// serves as reference to the entire ringf. Empty ringfs are represented
// as nil RingF pointers. The zero value for a RingF is a one-element
// ringf with a nil Value.
//
type RingF struct {
	next, prev *RingF
	Value      interface{} // for use by client; untouched by this library
}

func (r *RingF) init() *RingF {
	r.next = r
	r.prev = r
	return r
}

// Next returns the next ringf element. r must not be empty.
func (r *RingF) Next() *RingF {
	if r.next == nil {
		return r.init()
	}
	return r.next
}

// Prev returns the previous ringf element. r must not be empty.
func (r *RingF) Prev() *RingF {
	if r.next == nil {
		return r.init()
	}
	return r.prev
}

// Move moves n % r.Len() elements backward (n < 0) or forward (n >= 0)
// in the ringf and returns that ringf element. r must not be empty.
//
func (r *RingF) Move(n int) *RingF {
	if r.next == nil {
		return r.init()
	}
	switch {
	case n < 0:
		for ; n < 0; n++ {
			r = r.prev
		}
	case n > 0:
		for ; n > 0; n-- {
			r = r.next
		}
	}
	return r
}

// New creates a ringf of n elements.
func New(n int) *RingF {
	if n <= 0 {
		return nil
	}
	r := new(RingF)
	p := r
	for i := 1; i < n; i++ {
		p.next = &RingF{prev: p}
		p = p.next
	}
	p.next = r
	r.prev = p
	return r
}

// Link connects ringf r with ringf s such that r.Next()
// becomes s and returns the original value for r.Next().
// r must not be empty.
//
// If r and s point to the same ringf, linking
// them removes the elements between r and s from the ringf.
// The removed elements form a subringf and the result is a
// reference to that subringf (if no elements were removed,
// the result is still the original value for r.Next(),
// and not nil).
//
// If r and s point to different ringfs, linking
// them creates a single ringf with the elements of s inserted
// after r. The result points to the element following the
// last element of s after insertion.
//
func (r *RingF) Link(s *RingF) *RingF {
	n := r.Next()
	if s != nil {
		p := s.Prev()
		// Note: Cannot use multiple assignment because
		// evaluation order of LHS is not specified.
		r.next = s
		s.prev = r
		n.prev = p
		p.next = n
	}
	return n
}

// Unlink removes n % r.Len() elements from the ringf r, starting
// at r.Next(). If n % r.Len() == 0, r remains unchanged.
// The result is the removed subringf. r must not be empty.
//
func (r *RingF) Unlink(n int) *RingF {
	if n <= 0 {
		return nil
	}
	return r.Link(r.Move(n + 1))
}

// Len computes the number of elements in ringf r.
// It executes in time proportional to the number of elements.
//
func (r *RingF) Len() int {
	n := 0
	if r != nil {
		n = 1
		for p := r.Next(); p != r; p = p.next {
			n++
		}
	}
	return n
}

// Do calls function f on each element of the ringf, in forward order.
// The behavior of Do is undefined if f changes *r.
func (r *RingF) Do(f func(interface{})) {
	if r != nil {
		f(r.Value)
		for p := r.Next(); p != r; p = p.next {
			f(p.Value)
		}
	}
}

func (f *RingF) Filter(cb func(interface{}) bool) *RingF {
	if f == nil {
		return nil
	}

	p := f.Next()

	for p != f {
		if cb(p.Value) {
			p = p.Next()
		} else {
			p = p.Next()
			p.Prev().Prev().Unlink(1)
		}
	}

	if !cb(p.Value) {
		if p.Len() == 1 {
			return nil
		} else {
			p = f.Next()
			p.Prev().Prev().Unlink(1)
			return p
		}
	}

	return p
}

func (f *RingF) Any(cb func(interface{}) bool) bool {
	if f != nil {
		if cb(f.Value) {
			return true
		}
		for p := f.Next(); p != f; p = p.Next() {
			if cb(p.Value) {
				return true
			}
		}
	}

	return false
}

func (f *RingF) All(cb func(interface{}) bool) bool {
	if f != nil {
		if !cb(f.Value) {
			return false
		}
		for p := f.Next(); p != f; p = p.Next() {
			if !cb(p.Value) {
				return false
			}
		}
	}

	return false
}
