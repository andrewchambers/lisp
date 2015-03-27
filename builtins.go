package main

import (
	"fmt"
)

func RegisterStandardBuiltins(e *PxiEnv) {
	e.Define("=", PxiBuiltin(PxiBuiltinEq))
	e.Define("+", PxiBuiltin(PxiBuiltinAdd))
	e.Define("cond", PxiBuiltin(PxiBuiltinCond))
	e.Define("def", PxiBuiltin(PxiBuiltinDef))
	e.Define("quote", PxiBuiltin(PxiBuiltinQuote))
	e.Define("fn", PxiBuiltin(PxiBuiltinFn))
	e.Define("print", PxiBuiltin(PxiBuiltinPrint))
}

func PxiBuiltinEq(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	if arglist.Len() < 2 {
		return nil, fmt.Errorf("= expects two arguments")
	}
	args, err := p.EvalList(e, arglist)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(args)-1; i++ {
		n1, ok := args[i].(PxiNum)
		if !ok {
			return nil, fmt.Errorf("= arg %d is not a number", i)
		}
		n2, ok := args[i+1].(PxiNum)
		if !ok {
			return nil, fmt.Errorf("= arg %d is not a number", i+1)
		}
		if int64(n1) != int64(n2) {
			return PxiBool(false), nil
		}
	}
	return PxiBool(true), nil
}

func PxiBuiltinAdd(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	if arglist.Len() < 2 {
		return nil, fmt.Errorf("+ expects at least two arguments")
	}
	args, err := p.EvalList(e, arglist)
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

func PxiBuiltinCond(p *PxiState, env *PxiEnv, arglist *PxiList) (PxiVal, error) {
	for arglist != nil {
		if arglist.tail == nil {
			v, err := p.Eval(env, arglist.val)
			if err != nil {
				return nil, err
			}
			return v, err
		}
		v, err := p.Eval(env, arglist.val)
		if err != nil {
			return nil, err
		}
		condition, ok := v.(PxiBool)
		if ok && bool(condition) {
			v, err = p.Eval(env, arglist.tail.val)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
		arglist = arglist.tail.tail
	}
	return PxiBool(false), nil
}

func PxiBuiltinDef(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	args := arglist.ToSlice()
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
	p.genv.Define(string(s), v)
	return v, nil
}

func PxiBuiltinQuote(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	if arglist.Len() != 1 {
		return nil, fmt.Errorf("quote expects one argument")
	}
	return arglist.val, nil
}

func PxiBuiltinFn(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	var fnargs []PxiSym
	var fnbody []PxiVal
	if arglist.Len() < 2 {
		return nil, fmt.Errorf("fn expects at least two arguments")
	}
	argnames, ok := arglist.val.(*PxiList)
	if !ok {
		return nil, fmt.Errorf("fn needs a parameter list")
	}
	body := arglist.tail
	for argnames != nil {
		name, ok := argnames.val.(PxiSym)
		if !ok {
			return nil, fmt.Errorf("fn param not a symbol")
		}
		fnargs = append(fnargs, name)
		argnames = argnames.tail
	}
	for body != nil {
		fnbody = append(fnbody, body.val)
		body = body.tail
	}
	return &PxiFn{
		env:  e,
		args: fnargs,
		body: fnbody,
	}, nil
}

func PxiBuiltinPrint(p *PxiState, e *PxiEnv, arglist *PxiList) (PxiVal, error) {
	for arglist != nil {
		v, err := p.Eval(e, arglist.val)
		if err != nil {
			return nil, err
		}
		fmt.Println(v.String())
		arglist = arglist.tail
	}
	return EmptyList, nil
}
