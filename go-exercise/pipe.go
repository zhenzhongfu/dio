package main

import "fmt"

type S struct {
	v1 int
	v2 int
	v3 int
}

type fn func(*S)

func (s *S) fn1(a int) *S {
	s.v1 = a
	return s
}

func (s *S) fn2(a int) *S {
	s.v2 = a
	return s
}

/*
func fn3(a int) fn{
	return func(s *S) {
		s.v3 = a
	}
}
*/

func main() {
	ps := &S{}
	ps.fn1(1).fn2(2)
	fmt.Println(ps)
}
