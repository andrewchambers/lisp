package main

import (
	"bytes"
	"fmt"
)

type PxiVal interface {
	fmt.Stringer
}

type PxiBuiltin func(*PxiState, *PxiEnv, []PxiVal) (PxiVal, error)

func (PxiBuiltin) String() string { return "builtin" }

type PxiSym string

func (s PxiSym) String() string { return string(s) }

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
