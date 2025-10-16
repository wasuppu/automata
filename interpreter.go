package automata

func Match(line string, pattern string) bool {
	nfa := Interp(pattern)
	return nfa.Matches(line)
}

func MatchDFA(line string, pattern string) bool {
	dfa := NewDFA(Interp(pattern))
	return dfa.Matches(line)
}

func Interp(pattern string) *NFA {
	nfa := expr(parse(pattern))
	return &nfa
}

func expr(root node) NFA {
	if root.lable == "Expr" {
		term := term(root.children[0])
		if len(root.children) == 3 {
			return ChoicePair(term, expr(root.children[2]))
		}
		return term
	} else {
		panic("expr: " + root.lable)
	}
}

func term(root node) NFA {
	if root.lable == "Term" {
		factor := factor(root.children[0])
		if len(root.children) == 2 {
			return ConcatPair(factor, term(root.children[1]))
		}
		return factor
	} else {
		panic("term: " + root.lable)
	}
}

func factor(root node) NFA {
	if root.lable == "Factor" {
		atom := atom(root.children[0])
		if len(root.children) == 2 {
			meta := root.children[1].lable
			if meta == "*" {
				return Rep(atom)
			}
			if meta == "+" {
				return PlusRep(atom)
			}
			if meta == "?" {
				return Question(atom)
			}
		}
		return atom
	} else {
		panic("factor: " + root.lable)
	}
}

func atom(root node) NFA {
	if root.lable == "Atom" {
		if len(root.children) == 3 {
			return expr(root.children[1])
		}
		return char(root.children[0])
	} else {
		panic("atom: " + root.lable)
	}
}

func char(root node) NFA {
	if root.lable == "Char" {
		if len(root.children) == 2 {
			if root.children[1].lable == "d" {
				return Digit()
			}

			if root.children[1].lable == "w" {
				return Word()
			}

			if root.children[1].lable == "s" {
				return Space()
			}
			return Char(root.children[1].lable)
		}
		return Char(root.children[0].lable)
	} else {
		panic("char: " + root.lable)
	}
}
