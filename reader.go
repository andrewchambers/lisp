package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strconv"
)

type PxiReader struct {
	br *bufio.Reader
}

func NewReader(r io.Reader) *PxiReader {
	br := bufio.NewReader(r)
	return &PxiReader{
		br: br,
	}
}

func (rdr *PxiReader) skipWS() error {
	for {
		r, _, err := rdr.br.ReadRune()
		if err != nil {
			return err
		}
		if !isWhiteSpace(r) {
			break
		}
	}
	rdr.br.UnreadRune()
	return nil
}

func (rdr *PxiReader) Read() (PxiVal, error) {
	err := rdr.skipWS()
	if err != nil {
		return nil, err
	}
	r, _, err := rdr.br.ReadRune()
	if err != nil {
		return nil, err
	}
	err = rdr.br.UnreadRune()
	if err != nil {
		return nil, err
	}
	switch {
	case r == ';':
		for {
			r, _, err := rdr.br.ReadRune()
			if err != nil {
				return nil, err
			}
			if r == '\n' {
				return rdr.Read()
			}
		}
	case r == '(':
		return rdr.readList()
	case r == '-':
		return rdr.readNegNumberOrSymbol()
	case isNumeric(r):
		return rdr.readNumber()
	default:
		return rdr.readSymbol()
	}
}

func (rdr *PxiReader) readList() (PxiVal, error) {
	r, _, err := rdr.br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '(' {
		return nil, errors.New("list should start with '('")
	}
	var lst *PxiList = nil
	for {
		err = rdr.skipWS()
		if err != nil {
			return nil, err
		}
		r, _, err = rdr.br.ReadRune()
		if err != nil {
			return nil, err
		}
		if r == ')' {
			break
		}
		rdr.br.UnreadRune()
		v, err := rdr.Read()
		if err != nil {
			return nil, err
		}
		lst = lst.Cons(v)
	}
	return lst.Reverse(), nil
}

func (rdr *PxiReader) readNumber() (PxiVal, error) {
	var buff bytes.Buffer
	r, _, err := rdr.br.ReadRune()
	if err != nil {
		return nil, err
	}
	if !isNumeric(r) {
		return nil, errors.New("bad number")
	}
	buff.WriteRune(r)
	for {
		r, _, err = rdr.br.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if !isNumeric(r) {
			break
		}
		buff.WriteRune(r)
	}
	rdr.br.UnreadRune()
	n, err := strconv.ParseInt(buff.String(), 0, 64)
	if err != nil {
		return nil, err
	}
	return PxiNum(n), nil
}

func (rdr *PxiReader) readNegNumberOrSymbol() (PxiVal, error) {
	r, _, err := rdr.br.ReadRune()
	if err != nil {
		return nil, err
	}
	if r != '-' {
		return nil, errors.New("internal error")
	}
	r, _, err = rdr.br.ReadRune()
	if err != nil {
		return nil, err
	}
	rdr.br.UnreadRune()
	if isNumeric(r) {
		return rdr.readNumber()
	}
	return rdr.readSymbol()
}

func (rdr *PxiReader) readSymbol() (PxiVal, error) {
	var buff bytes.Buffer
	for {
		r, _, err := rdr.br.ReadRune()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if endOfSymbol(r) {
			break
		}
		buff.WriteRune(r)
	}
	err := rdr.br.UnreadRune()
	if err != nil {
		return nil, err
	}
	return PxiSym(buff.String()), nil
}

func endOfSymbol(r rune) bool {
	if isWhiteSpace(r) {
		return true
	}
	switch r {
	case '(', ')', '{', '}', '[', ']':
		return true
	}
	return false
}

func isWhiteSpace(r rune) bool {
	return r == ' ' || r == '\r' || r == '\n' || r == '\t'
}

func isNumeric(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	return false
}
