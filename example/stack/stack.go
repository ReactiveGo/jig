// Code generated by jig; DO NOT EDIT.

//go:generate jig

package stack

//jig:name Stack

type Stack []interface{}

var zero interface{}

//jig:name StackPush

func (s *Stack) Push(v interface{}) {
	*s = append(*s, v)
}

//jig:name StackPop

func (s *Stack) Pop() (interface{}, bool) {
	if len(*s) == 0 {
		return zero, false
	}
	i := len(*s) - 1
	v := (*s)[i]
	*s = (*s)[:i]
	return v, true
}

//jig:name StackTop

func (s *Stack) Top() (interface{}, bool) {
	if len(*s) == 0 {
		return zero, false
	}
	return (*s)[len(*s)-1], true
}
