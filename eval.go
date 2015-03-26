package main

import (
	"fmt"
	"sync"
)

type PxiState struct {
	genv *PxiEnv
}

func NewPxiState() *PxiState {
	genv := NewEnv(nil)
	RegisterStandardBuiltins(genv)
	return &PxiState{
		genv: genv,
	}
}

type PxiEnv struct {
	parent *PxiEnv
	mutex  sync.RWMutex
	kv     map[string]PxiVal
}

func (e *PxiEnv) Define(s string, v PxiVal) {
	e.mutex.Lock()
	e.kv[s] = v
	e.mutex.Unlock()
}

func (e *PxiEnv) Lookup(s string) (PxiVal, error) {
	e.mutex.RLock()
	v, ok := e.kv[s]
	e.mutex.RUnlock()
	if !ok {
		if e.parent == nil {
			return nil, fmt.Errorf("%s is undefined", s)
		}
		return e.parent.Lookup(s)
	}
	return v, nil
}

func NewEnv(parent *PxiEnv) *PxiEnv {
	return &PxiEnv{
		parent: parent,
		kv:     make(map[string]PxiVal),
	}
}

func (ps *PxiState) Eval(env *PxiEnv, v PxiVal) (PxiVal, error) {
	switch v := v.(type) {
	case *PxiList:
		if v == nil {
			// invoking the empty list
			// is the empty list.
			return v, nil
		}
		fn, err := ps.Eval(env, v.val)
		if err != nil {
			return nil, err
		}
		v = v.tail
		var args []PxiVal
		for v != nil {
			arg := v.val
			args = append(args, arg)
			v = v.tail
		}
		return ps.Apply(env, fn, args)
	case PxiSym:
		return env.Lookup(string(v))
	case PxiNum:
		return v, nil
	}
	panic("internal error")
}

func (ps *PxiState) EvalSlice(env *PxiEnv, vals []PxiVal) ([]PxiVal, error) {
	var ret []PxiVal
	for _, v := range vals {
		r, err := ps.Eval(env, v)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
	}
	return ret, nil
}

func (ps *PxiState) Apply(e *PxiEnv, fn PxiVal, args []PxiVal) (PxiVal, error) {
	switch fn := fn.(type) {
	case PxiBuiltin:
		return fn(ps, e, args)
	}
	return nil, fmt.Errorf("Cannot call %T", fn)
}
