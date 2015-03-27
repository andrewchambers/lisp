package main

/*

type Opc byte

const (
	OP_LOAD Opc = iota
	OP_RET
	OP_LOOKUP
	OP_CALL
	OP_TAILCALL
	OP_POP
)

type ByteCode struct {
	op   byte
	arg1 PxiVal
	next *ByteCode
}

// Compile a pxi function
type Function struct {
	finalized bool
	ent       *ByteCode
	cur       *ByteCode
	bcs       []*Bytecode
}

func (c *Function) emit(bc *Bytecode) {
	if c.finalized {
		panic("internal error")
	}
	if c.cur != nil {
		c.cur.next = bc
	} else {
		c.ent = bc
	}
	c.cur = bc
	c.bcs = append(c.bcs, bc)
}

func NewFunction(fnargs []PxiSym) *Function {
	c := &Function{}
	return c
}

func (c *Function) finalize() {
	for _, bc := range c.bcs {
		if bc.next == nil {
			if bc.op == OP_CALL {
				bc.op = OP_TAILCALL
			} else {
				bc.next = &Bytecode{
					op: OP_RET,
				}
			}
		}
	}
	c.bcs = nil
	c.cur = nil
	c.finalized = true
}

func (c *Function) compileForm(v PxiVal) error {
	switch v := v.(type) {
	case PxiNum:
		c.emit(&Bytecode{
			op:   OP_LOAD,
			arg1: v,
		})
	case PxiBool:
		c.emit(&Bytecode{
			op:   OP_LOAD,
			arg1: v,
		})
	case *PxiList:
		c.compileCall(v)
	}
}

func (c *Function) compileCall(l *PxiList) error {
	if l.Len() == 0 {
		c.emit(&Bytecode{
			op:   OP_LOAD,
			arg1: v,
		})
	}
	f, args := l.Head(), l.Rest()
	sym, ok := f.(PxiSym)
	if !ok {
		return fmt.Errorf("cannot call %s", sym.String())
	}
	switch string(sym) {
	case "fn":
		//XXX args
		newf := NewFunction(nil)
		newf.finalize()
		c.emit(&Bytecode{
			op:   OP_LOAD,
			arg1: v,
		})
	case "cond":
	default:
		nargs := 0
		for args != nil {
			nargs += 1
			c.compileForm(args.Head())
			args = args.Rest()
		}
		c.emit(&Bytecode{
			op:   OP_LOOKUP,
			arg1: v,
		})
		c.emit(&Bytecode{
			op:   OP_CALL,
			arg1: PxiNum(nargs),
		})
	}
}
*/
