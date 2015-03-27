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
		return ps.Apply(env, fn, v.tail)
	case PxiSym:
		return env.Lookup(string(v))
	}
	return v, nil
}

func (ps *PxiState) EvalList(env *PxiEnv, vals *PxiList) ([]PxiVal, error) {
	var ret []PxiVal
	for vals != nil {
		r, err := ps.Eval(env, vals.val)
		if err != nil {
			return nil, err
		}
		ret = append(ret, r)
		vals = vals.tail
	}
	return ret, nil
}

func (ps *PxiState) Apply(env *PxiEnv, fn PxiVal, args *PxiList) (PxiVal, error) {
	var err error
	switch fn := fn.(type) {
	case *PxiFn:
		argnames := fn.args
		newenv := NewEnv(fn.env)
		if len(argnames) != args.Len() {
			return nil, fmt.Errorf("incorrect number of args")
		}
		for idx := range argnames {
			name := argnames[idx]
			argval, err := ps.Eval(env, args.val)
			if err != nil {
				return nil, err
			}
			newenv.Define(string(name), argval)
			args = args.tail
		}
		var ret PxiVal
		for _, expr := range fn.body {
			ret, err = ps.Eval(newenv, expr)
			if err != nil {
				return nil, err
			}
		}
		return ret, nil
	case PxiBuiltin:
		return fn(ps, env, args)
	}
	return nil, fmt.Errorf("Cannot call %T", fn)
}
