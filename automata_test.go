package automata

import (
	"reflect"
	"testing"
)

func TestAutomata(t *testing.T) {
	t.Run("state transition for symbol", func(t *testing.T) {
		s1 := State(false)
		s2 := State(false)
		s1.addTransition("a", s2)

		transitions := s1.getTransition("a")

		if len(transitions) != 1 {
			t.Errorf("expected length is %d, got %d", 1, len(transitions))
		}

		isIn := false
		for _, s := range transitions {
			if reflect.DeepEqual(s, s2) {
				isIn = true
				break
			}
		}
		if !isIn {
			t.Errorf("s2 isn't in the transitions of symbol a")
		}
	})

	t.Run("char nfa", func(t *testing.T) {
		nfa := Char("a")
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", false},
			{"a", true},
			{"b", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("epsilon nfa", func(t *testing.T) {
		nfa := Epsilon()

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", true},
			{"a", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
		}
	})

	t.Run("connection pair nfa", func(t *testing.T) {
		nfa := ConcatPair(Char("a"), Char("b"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"ab", true},
			{"ac", false},
			{"aab", false},
			{"a", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("choice pair nfa", func(t *testing.T) {
		nfa := ChoicePair(Char("a"), Char("b"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"a", true},
			{"b", true},
			{"c", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("connection nfa", func(t *testing.T) {
		nfa := Concat(Char("a"), []NFA{Char("b"), Char("c")})
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"abc", true},
			{"aba", false},
			{"ab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("choice nfa", func(t *testing.T) {
		nfa := Choice(Char("a"), []NFA{Char("b"), Char("c")})
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"c", true},
			{"a", true},
			{"b", true},
			{"d", false},
			{"ab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("rep_nfa", func(t *testing.T) {
		nfa := Rep(Char("a"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", true},
			{"a", true},
			{"b", false},
			{"aa", true},
			{"ab", false},
			{"aaa", true},
			{"aab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("repE_nfa", func(t *testing.T) {
		nfa := RepExplicit(Char("a"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", true},
			{"a", true},
			{"b", false},
			{"aa", true},
			{"ab", false},
			{"aaa", true},
			{"aab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})
	t.Run("plus_nfa", func(t *testing.T) {
		nfa := PlusRep(Char("a"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", false},
			{"a", true},
			{"b", false},
			{"aa", true},
			{"ab", false},
			{"aaa", true},
			{"aab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})
	t.Run("plusE_nfa", func(t *testing.T) {
		nfa := PlusRepExplicit(Char("a"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", false},
			{"a", true},
			{"b", false},
			{"aa", true},
			{"ab", false},
			{"aaa", true},
			{"aab", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("ques_nfa", func(t *testing.T) {
		nfa := QuestionExplicit(Char("a"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"", true},
			{"a", true},
			{"b", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

	t.Run("complex composition", func(t *testing.T) {
		// xy*|z
		nfa := ChoicePair(ConcatPair(Char("x"), RepExplicit(Char("y"))), Char("z"))
		dfa := NewDFA(&nfa)

		tests := []struct {
			testStr        string
			expectedResult bool
		}{
			{"x", true},
			{"xy", true},
			{"xyy", true},
			{"xyz", false},
			{"z", true},
			{"a", false},
			{"", false},
		}

		for _, test := range tests {
			got := nfa.Matches(test.testStr)
			if got != test.expectedResult {
				t.Errorf("test %v: got:%v, wanted:%v", test.testStr, got, test.expectedResult)
			}
			gotDfa := dfa.Matches(test.testStr)
			if gotDfa != test.expectedResult {
				t.Errorf("test %v: gotDfa:%v, wanted:%v", test.testStr, gotDfa, test.expectedResult)
			}
		}
	})

}
