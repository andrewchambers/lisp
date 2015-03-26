package main

import (
	"fmt"
)

func RegisterStandardBuiltins(e *PxiEnv) {
	e.Define("+", PxiBuiltin(PxiBuiltinAdd))
	e.Define("def", PxiBuiltin(PxiBuiltinDef))
}

func PxiBuiltinAdd(p *PxiState, e *PxiEnv, args []PxiVal) (PxiVal, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("+ expects at least two arguments")
	}
	args, err := p.EvalSlice(e, args)
	if err != nil {
		return nil, err
	}
	var v int64
	for idx, a := range args {
		n, ok := a.(PxiNum)
		if !ok {
			return nil, fmt.Errorf("+ arg %d is not a number", idx)
		}
		v += int64(n)
	}
	return PxiNum(v), nil
}

func PxiBuiltinDef(p *PxiState, e *PxiEnv, args []PxiVal) (PxiVal, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("def expects two arguments")
	}
	s, ok := args[0].(PxiSym)
	if !ok {
		return nil, fmt.Errorf("def expects a symbol")
	}
	v, err := p.Eval(e, args[1])
	if err != nil {
		return nil, err
	}
	e.Define(string(s), v)
	return v, nil
}
