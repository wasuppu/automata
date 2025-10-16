package automata

type node struct {
	lable    string
	children []node
}

func newNode(ch byte) node {
	return node{string(ch), []node{}}
}

type parser struct {
	pattern string
	pos     int
}

func parse(pattern string) node {
	p := parser{pattern: pattern}
	return p.expr()
}

func (p parser) hasMore() bool {
	return p.pos < len(p.pattern)
}

func (p parser) peek() byte {
	return p.pattern[p.pos]
}
func (p *parser) next() byte {
	ch := p.peek()
	p.match(ch)
	return ch
}

func (p *parser) match(ch byte) {
	if p.peek() != ch {
		panic("Unexpected symbol " + string(ch))
	}

	p.pos++
}

func (p *parser) expr() node {
	term := p.term()
	if p.hasMore() && p.peek() == '|' {
		p.match('|')
		expr := p.expr()
		return node{"Expr", []node{term, newNode('|'), expr}}
	}

	return node{"Expr", []node{term}}
}

func (p *parser) term() node {
	factor := p.factor()

	if p.hasMore() && p.peek() != ')' && p.peek() != '|' {
		term := p.term()
		return node{"Term", []node{factor, term}}
	}

	return node{"Term", []node{factor}}
}

func (p *parser) factor() node {
	atom := p.atom()

	if p.hasMore() && isMetaChar(p.peek()) {
		meta := p.next()
		return node{"Factor", []node{atom, newNode(meta)}}
	}

	return node{"Factor", []node{atom}}
}

func (p *parser) atom() node {
	if p.peek() == '(' {
		p.match('(')
		expr := p.expr()
		p.match(')')
		return node{"Atom", []node{newNode('('), expr, newNode(')')}}
	}

	ch := p.char()
	return node{"Atom", []node{ch}}
}

func (p *parser) char() node {
	if isMetaChar(p.peek()) {
		panic("Unexpected meta char " + string(p.peek()))
	}

	if p.peek() == '\\' {
		p.match('\\')
		return node{"Char", []node{newNode('\\'), newNode(p.next())}}
	}

	return node{"Char", []node{newNode(p.next())}}
}

func isMetaChar(ch byte) bool {
	return ch == '*' || ch == '+' || ch == '?'
}
