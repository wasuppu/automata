package automata

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	EPSILON         = "ε"
	EPSILON_CLOSURE = "ε*"
)

type state struct {
	IsAccepted     bool
	Transitions    map[string][]*state
	EpsilonClosure map[*state]bool
	Number         int
}

func (s state) addTransition(symbol string, state *state) {
	s.Transitions[symbol] = append(s.Transitions[symbol], state)
}

func (s state) getTransition(symbol string) []*state {
	return s.Transitions[symbol]
}

func (s *state) Matches(str string, visited map[*state]bool) bool {
	if visited[s] {
		return false
	}
	visited[s] = true

	if len(str) == 0 {
		if s.IsAccepted {
			return true
		}

		for _, nextState := range s.getTransition(EPSILON) {
			if nextState.Matches("", visited) {
				return true
			}
		}
		return false
	}

	symbol := string([]rune(str)[0])

	rest := str[1:]

	for _, nextState := range s.getTransition(symbol) {
		visited = make(map[*state]bool)
		if nextState.Matches(rest, visited) {
			return true
		}
	}

	for _, nextState := range s.getTransition(EPSILON) {
		if nextState.Matches(str, visited) {
			return true
		}
	}

	return false
}

func (s *state) getEpsilonClosure() map[*state]bool {
	if s.EpsilonClosure == nil {
		epsilonTrans := s.getTransition(EPSILON)
		s.EpsilonClosure = make(map[*state]bool)
		s.EpsilonClosure[s] = true
		for _, nextState := range epsilonTrans {
			if !s.EpsilonClosure[nextState] {
				s.EpsilonClosure[nextState] = true
				nextClosure := nextState.getEpsilonClosure()
				for state := range nextClosure {
					s.EpsilonClosure[state] = true
				}
			}
		}
	}
	return s.EpsilonClosure
}

func State(isAccepted bool) *state {
	s := state{IsAccepted: isAccepted}
	s.Transitions = make(map[string][]*state)
	return &s
}

type NFA struct {
	in                 *state
	out                *state
	transTable         map[int]map[string][]int
	acceptingStates    map[*state]bool
	acceptingStateNums map[int]bool

	alphabet map[string]bool
}

func (nfa *NFA) SetLabel() {
	visited := make(map[*state]bool)
	var visitState func(st *state)
	visitState = func(st *state) {
		if visited[st] {
			return
		}
		visited[st] = true
		st.Number = len(visited)
		transitions := st.Transitions

		for _, symTransitions := range transitions {
			for _, nextState := range symTransitions {
				visitState(nextState)
			}
		}
	}

	visitState(nfa.in)
}

func (nfa *NFA) GetTransitionTable() map[int]map[string][]int {
	if nfa.transTable == nil {
		nfa.transTable = make(map[int]map[string][]int)
		nfa.acceptingStates = make(map[*state]bool)
		visited := make(map[*state]bool)
		symbols := make(map[string]bool)

		var visitState func(st *state)
		visitState = func(st *state) {
			if visited[st] {
				return
			}
			visited[st] = true
			st.Number = len(visited)
			nfa.transTable[st.Number] = make(map[string][]int)
			if st.IsAccepted {
				nfa.acceptingStates[st] = true
			}

			transitions := st.Transitions

			for sym, symTransitions := range transitions {
				var combineState []int
				symbols[sym] = true
				for _, nextState := range symTransitions {
					visitState(nextState)
					combineState = append(combineState, nextState.Number)
				}
				nfa.transTable[st.Number][sym] = combineState
			}
		}

		visitState(nfa.in)

		for state := range visited {
			delete(nfa.transTable[state.Number], EPSILON)
			for closureState := range state.getEpsilonClosure() {
				nfa.transTable[state.Number][EPSILON_CLOSURE] = append(nfa.transTable[state.Number][EPSILON_CLOSURE], closureState.Number)
			}
		}
	}
	return nfa.transTable
}

func (nfa NFA) Matches(str string) bool {
	visited := make(map[*state]bool)
	return nfa.in.Matches(str, visited)
}

func (nfa *NFA) GetAlphabet() map[string]bool {
	if nfa.alphabet == nil {
		nfa.alphabet = make(map[string]bool)
		table := nfa.GetTransitionTable()
		for stateNum := range table {
			transitions := table[stateNum]
			for symbol := range transitions {
				if symbol != EPSILON_CLOSURE {
					nfa.alphabet[symbol] = true
				}
			}
		}
	}
	return nfa.alphabet
}

func (nfa *NFA) getAcceptingStates() map[*state]bool {
	if nfa.acceptingStates == nil {
		nfa.GetTransitionTable()
	}
	return nfa.acceptingStates
}

func (nfa NFA) getAcceptingStateNums() map[int]bool {
	if nfa.acceptingStateNums == nil {
		nfa.acceptingStateNums = make(map[int]bool)
		for acceptingState := range nfa.getAcceptingStates() {
			nfa.acceptingStateNums[acceptingState.Number] = true
		}
	}
	return nfa.acceptingStateNums
}

type DFA struct {
	nfa                        *NFA
	acceptingStateNums         map[string]bool
	originalTransitonTable     map[string]map[string]string
	originalAcceptingStateNums map[string]bool

	originalStartState string
	startState         string

	transTable map[string]map[string]string
}

func NewDFA(nfa *NFA) *DFA {
	dfa := DFA{nfa: nfa}
	dfa.GetTransitionTable()
	return &dfa
}

func (dfa *DFA) GetAlphabet() map[string]bool {
	return dfa.nfa.GetAlphabet()
}

func (dfa *DFA) GetAcceptingStateNums() map[string]bool {
	if dfa.acceptingStateNums == nil {
		dfa.GetTransitionTable()
	}
	return dfa.acceptingStateNums
}

func (dfa *DFA) GetTransitionTable() map[string]map[string]string {
	if dfa.transTable != nil {
		return dfa.transTable
	}

	nfaTable := dfa.nfa.GetTransitionTable()
	dfa.acceptingStateNums = make(map[string]bool)

	startState := nfaTable[dfa.nfa.in.Number][EPSILON_CLOSURE]
	dfa.originalStartState = intlistToString(startState, ",")

	worklist := [][]int{startState}

	alphabet := dfa.GetAlphabet()
	nfaAcceptingStates := dfa.nfa.getAcceptingStateNums()

	dfaTable := make(map[string]map[string]string)

	var updateAcceptingStates func(stateNums []int) = func(stateNums []int) {
		for nfaAcceptingState := range nfaAcceptingStates {
			isContain := slices.Contains(stateNums, nfaAcceptingState)
			if isContain {
				dfa.acceptingStateNums[intlistToString(stateNums, ",")] = true
				break
			}
		}
	}

	for len(worklist) > 0 {
		stateNums := worklist[0]
		worklist = worklist[1:]
		dfaStateLabel := intlistToString(stateNums, ",")
		dfaTable[dfaStateLabel] = make(map[string]string)

		for symbol := range alphabet {
			onSymbol := []int{}
			updateAcceptingStates(stateNums)

			for _, stateNum := range stateNums {
				nfaStateNumsOnSymbol := nfaTable[stateNum][symbol]
				if len(nfaStateNumsOnSymbol) <= 0 {
					continue
				}

				for _, nfaStateNumOnSymbol := range nfaStateNumsOnSymbol {
					if nfaTable[nfaStateNumOnSymbol] == nil {
						continue
					}
					onSymbol = append(onSymbol, nfaTable[nfaStateNumOnSymbol][EPSILON_CLOSURE]...)
				}
			}

			dfaStateNumsOnSymbolSet := make(map[int]bool)
			for _, symbol := range onSymbol {
				dfaStateNumsOnSymbolSet[symbol] = true
			}
			dfaStateNumsOnSymbol := []int{}
			for symbol := range dfaStateNumsOnSymbolSet {
				dfaStateNumsOnSymbol = append(dfaStateNumsOnSymbol, symbol)
			}

			if len(dfaStateNumsOnSymbol) > 0 {
				dfaOnSymbolStr := intlistToString(dfaStateNumsOnSymbol, ",")

				dfaTable[dfaStateLabel][symbol] = dfaOnSymbolStr

				if dfaTable[dfaOnSymbolStr] == nil {
					worklist = append([][]int{dfaStateNumsOnSymbol}, worklist...)
				}
			}

		}

	}
	dfa.transTable = dfa.remapStateNumbers(dfaTable)
	return dfa.transTable
}

func (dfa *DFA) remapStateNumbers(calculatedDFATable map[string]map[string]string) map[string]map[string]string {
	newStatesMap := make(map[string]string)
	dfa.originalTransitonTable = calculatedDFATable
	transitionTable := make(map[string]map[string]string)

	n := 1
	for origianlNumber := range calculatedDFATable {
		newStatesMap[origianlNumber] = strconv.Itoa(n)
		if origianlNumber == dfa.originalStartState {
			dfa.startState = strconv.Itoa(n)
		}
		n++
	}
	for origianlNumber := range calculatedDFATable {
		originalRow := calculatedDFATable[origianlNumber]
		row := make(map[string]string)
		for symbol := range originalRow {
			row[symbol] = newStatesMap[originalRow[symbol]]
		}
		transitionTable[newStatesMap[origianlNumber]] = row
	}

	dfa.originalAcceptingStateNums = dfa.acceptingStateNums
	dfa.acceptingStateNums = make(map[string]bool)

	for originalNumber := range dfa.originalAcceptingStateNums {
		dfa.acceptingStateNums[newStatesMap[originalNumber]] = true
	}

	return transitionTable
}

func (dfa *DFA) Matches(str string) bool {
	state := dfa.startState
	table := dfa.GetTransitionTable()

	for _, symbol := range str {
		row, ok := table[state]
		if !ok {
			return false
		}
		state, ok = row[string(symbol)]
		if !ok {
			return false
		}
	}

	acceptingStateNums := dfa.GetAcceptingStateNums()
	_, ok := acceptingStateNums[state]
	return ok
}

func intlistToString(il []int, sep string) string {
	stringArray := make([]string, len(il))
	for i, v := range il {
		stringArray[i] = strconv.Itoa(v)
	}
	return strings.Join(stringArray, sep)
}

// in -s-> out
func Char(symbol string) NFA {
	instate := State(false)
	outstate := State(true)
	instate.addTransition(symbol, outstate)
	return NFA{in: instate, out: outstate}
}

func Epsilon() NFA {
	return Char(EPSILON)
}

func ConcatPair(first NFA, second NFA) NFA {
	first.out.IsAccepted = false
	second.out.IsAccepted = true

	first.out.addTransition(EPSILON, second.in)

	return NFA{in: first.in, out: second.out}
}

func Concat(first NFA, rest []NFA) NFA {
	for _, fragment := range rest {
		first = ConcatPair(first, fragment)
	}
	return first
}

func ChoicePair(first NFA, second NFA) NFA {
	instate := State(false)
	outstate := State(true)

	instate.addTransition(EPSILON, first.in)
	instate.addTransition(EPSILON, second.in)

	first.out.IsAccepted = false
	second.out.IsAccepted = false

	first.out.addTransition(EPSILON, outstate)
	second.out.addTransition(EPSILON, outstate)

	return NFA{in: instate, out: outstate}
}

func Choice(first NFA, rest []NFA) NFA {
	for _, fragment := range rest {
		first = ChoicePair(first, fragment)
	}
	return first
}

func RepExplicit(fragment NFA) NFA {
	instate := State(false)
	outstate := State(true)

	instate.addTransition(EPSILON, fragment.in)
	instate.addTransition(EPSILON, outstate)

	fragment.out.IsAccepted = false

	fragment.out.addTransition(EPSILON, outstate)
	outstate.addTransition(EPSILON, fragment.in)

	return NFA{in: instate, out: outstate}
}

func Rep(fragment NFA) NFA {
	fragment.in.addTransition(EPSILON, fragment.out)
	fragment.out.addTransition(EPSILON, fragment.in)
	return fragment
}

func PlusRepExplicit(fragment NFA) NFA {
	fragmentRep := RepExplicit(fragment)
	return ConcatPair(fragment, fragmentRep)
}

func PlusRep(fragment NFA) NFA {
	fragment.out.addTransition(EPSILON, fragment.in)
	return fragment
}

func QuestionExplicit(fragment NFA) NFA {
	return ChoicePair(fragment, Epsilon())
}

func Question(fragment NFA) NFA {
	fragment.in.addTransition(EPSILON, fragment.out)
	return fragment
}

func (s state) String() string {
	str := ""
	if !s.IsAccepted {
		str = fmt.Sprintf("(%v)", s.Number)

		str += "-"
	} else {
		str = fmt.Sprintf("((%v))", s.Number)
		if len(s.Transitions) > 0 {
			str += "-"
		}
	}
	i := 0

	for k := range s.Transitions {
		if i == 0 {
			str += string(k) + "->"
		} else {
			str += "|\n"
			str += "-" + string(k) + "->"
		}
	}
	return str
}

func (nfa NFA) String() string {
	visited := make(map[*state]bool)
	var transStr func(state *state, dep int) (string, *state)
	transStr = func(state *state, dep int) (string, *state) {
		if visited[state] {
			return "", state
		}
		visited[state] = true

		str := ""
		for _, ss := range state.Transitions {
			i := 0
			var length int
			for _, s := range ss {
				if i == 0 {
					str += s.String()
					length = len(s.String())
				} else {
					spaces := strings.Repeat(" ", length)
					str += "\n\t" + strings.Repeat(" ", length-3) + strings.Repeat(spaces, dep-1) + "|_ "
					str += s.String()
				}
				i++
				if len(s.Transitions) > 0 {
					msg, _ := transStr(s, dep+1)
					str += msg
				}
			}
		}
		return str, nil
	}

	str := "{ NFA:\n"
	str += "\t" + nfa.in.String()
	msg, _ := transStr(nfa.in, 1)
	str += msg
	str += "\n}"
	return str
}

func Digit() NFA {
	digits := "123456789"
	first := Char("0")
	for _, d := range digits {
		first = ChoicePair(first, Char(string(d)))
	}
	return first
}

func Word() NFA {
	words := "123456789abcdefghiklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	first := Char("0")
	for _, d := range words {
		first = ChoicePair(first, Char(string(d)))
	}
	return first
}

func Space() NFA {
	spaces := "\t\n\r\f"
	first := Char(" ")
	for _, d := range spaces {
		first = ChoicePair(first, Char(string(d)))
	}
	return first
}
