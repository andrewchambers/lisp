package main

import (
	"bytes"
	"fmt"
)

type PxiVal interface {
	fmt.Stringer
}

type PxiFn struct {
	env  *PxiEnv
	args []PxiSym
	body []PxiVal
}

func (PxiFn) String() string { return "function" }

type PxiBuiltin func(*PxiState, *PxiEnv, *PxiList) (PxiVal, error)

func (PxiBuiltin) String() string { return "builtin" }

type PxiSym string

func (s PxiSym) String() string { return string(s) }

type PxiBool bool

func (b PxiBool) String() string {
	if !b {
		return "false"
	}
	return "true"
}

type PxiNum int64

func (n PxiNum) String() string { return fmt.Sprint(int64(n)) }

type PxiList struct {
	val  PxiVal
	tail *PxiList
}

var EmptyList *PxiList = nil

func (l *PxiList) String() string {
	var buff bytes.Buffer
	buff.WriteString("(")
	first := true
	for l != nil {
		if !first {
			buff.WriteString(" ")
		}
		buff.WriteString(l.val.String())
		first = false
		l = l.tail
	}
	buff.WriteString(")")
	return buff.String()
}

func (l *PxiList) Len() int {
	r := 0
	for l != nil {
		r += 1
		l = l.tail
	}
	return r
}

func (l *PxiList) ToSlice() []PxiVal {
	var r []PxiVal
	for l != nil {
		r = append(r, l.val)
		l = l.tail
	}
	return r
}

func (l *PxiList) Reverse() *PxiList {
	var lst *PxiList = nil
	for l != nil {
		lst = lst.Cons(l.val)
		l = l.tail
	}
	return lst
}

func (l *PxiList) Cons(v PxiVal) *PxiList {
	return &PxiList{
		val:  v,
		tail: l,
	}
}

func (l *PxiList) Head() PxiVal {
	if l == nil {
		return nil
	}
	return l.val
}

func (l *PxiList) Rest() *PxiList {
	if l == nil {
		return nil
	}
	return l.tail
}
