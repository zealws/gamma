package parse

import (
	"bufio"
	"fmt"
	"github.com/zfjagann/gamma/sexpr"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	bufferSize = 1024
)

func Parse(input string) (sexpr.SExpr, error) {
	return NewParser(strings.NewReader(input)).Parse()
}

type Parser struct {
	reader   *bufio.Reader
	offset   int
	lastRead rune
}

func (p *Parser) errorf(f string, items ...interface{}) error {
	return p.error(fmt.Sprintf(f, items...))
}

func (p *Parser) error(msg string) error {
	return fmt.Errorf(msg+" at offset %d", p.offset)
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{bufio.NewReaderSize(reader, bufferSize), 0, '\000'}
}

func (p *Parser) readCh() (rune, error) {
	var err error
	p.lastRead, _, err = p.reader.ReadRune()
	p.offset += len(string(p.lastRead))
	return p.lastRead, err
}

func (p *Parser) unread() {
	p.offset -= len(string(p.lastRead))
	p.reader.UnreadRune()
}

func (p *Parser) Parse() (sexpr.SExpr, error) {
	p.offset = 0
	p.lastRead = '\000'
	expr, _, err := p.readSExpr()
	return expr, err
}

func (p *Parser) readSExpr() (sexpr.SExpr, bool, error) {
	var ch rune
	var err error
	for ch, err = p.readCh(); unicode.IsSpace(ch); ch, err = p.readCh() {
		if err != nil {
			return nil, false, err
		}
	}

	if ch == '\'' {
		literalExpr, eof, err := p.readSExpr()
		if err != nil {
			return nil, false, err
		}
		return sexpr.Quote(literalExpr), eof, nil
	} else if ch == '(' {
		return p.readList()
	} else if ch == '#' {
		return p.readBoolean()
	} else if unicode.IsDigit(ch) {
		p.unread()
		return p.readNumber()
	} else {
		p.unread()
		return p.readSymbol()
	}
}

func (p *Parser) readList() (sexpr.SExpr, bool, error) {
	var ch rune
	var err error
	for ch, err = p.readCh(); unicode.IsSpace(ch); ch, err = p.readCh() {
		if err != nil {
			if err == io.EOF {
				return nil, true, p.error("unexpected EOF in list")
			}
			return nil, false, err
		}
	}

	if ch == ')' {
		return sexpr.Null, false, nil
	} else if ch == '.' {
		e, eof, err := p.readSExpr()
		if err != nil {
			return nil, false, err
		}
		if !eof {
			for ch, err = p.readCh(); unicode.IsSpace(ch); ch, err = p.readCh() {
				if err != nil {
					return nil, false, err
				}
			}
		}
		return e, eof, nil
	} else {
		p.unread()
		head, eof, err := p.readSExpr()
		if err != nil {
			return nil, false, err
		}
		if eof {
			return nil, false, p.error("unexpected EOF in list")
		}
		tail, eof, err := p.readList()
		if err != nil {
			return nil, false, err
		}
		return sexpr.Cons(head, tail), eof, nil
	}
}

func (p *Parser) readBoolean() (sexpr.SExpr, bool, error) {
	ch, err := p.readCh()
	if err != nil {
		if err == io.EOF {
			return nil, false, p.error("unexpected EOF in boolean expression")
		}
		return nil, false, err
	}

	if ch == 't' {
		return sexpr.True, false, nil
	} else {
		return sexpr.False, false, nil
	}
}

func (p *Parser) readSymbol() (sexpr.SExpr, bool, error) {
	name := ""
	var err error
	var ch rune
	for ch, err = p.readCh(); ch != ' ' && ch != '\t' && ch != '\n' && ch != '(' && ch != ')'; ch, err = p.readCh() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, false, err
			}
		}
		name += string(ch)
	}
	if name == "" && err == io.EOF {
		return nil, false, p.error("unexpected EOF in symbol expression")
	} else if name == "" {
		return nil, false, p.errorf("unexpected '%v'. expecting symbol", string(ch))
	}
	p.reader.UnreadRune()
	return sexpr.Symbol(name), err == io.EOF, nil
}
func (p *Parser) readNumber() (sexpr.SExpr, bool, error) {
	numstr := ""
	var err error
	var ch rune
	for ch, err = p.readCh(); ch != ' ' && ch != '\t' && ch != '\n' && ch != '(' && ch != ')'; ch, err = p.readCh() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, false, err
			}
		}
		numstr += string(ch)
	}
	if numstr == "" && err == io.EOF {
		return nil, false, p.error("unexpected EOF in number expression")
	}
	p.reader.UnreadRune()
	if strings.Contains(numstr, ".") {
		f, ferr := strconv.ParseFloat(numstr, 64)
		if ferr != nil {
			return nil, false, p.errorf("invalid float literal %q", numstr)
		}
		return sexpr.Float(f), err == io.EOF, nil
	} else {
		d, derr := strconv.ParseInt(numstr, 0, 64)
		if derr != nil {
			return nil, false, p.errorf("invalid integer literal %q", numstr)
		}
		return sexpr.Integer(d), err == io.EOF, nil
	}
}
