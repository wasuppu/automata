# automata

A RegExp Machine based on the Finite Automata written in Go

## Example

```go
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: automata -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]
	nfa := automata.Interp(pattern)

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, "> ")
		line, _, err := inputReader.ReadLine()

		if errors.Is(err, io.EOF) {
			break
		}

		if !nfa.Matches(string(line)) {
			fmt.Printf("%s\t\t=> Not matched.\n", string(line))
		} else {
			fmt.Printf("%s\t\t=> Matched.\n", string(line))
		}
	}
}
```

```
$ go run example/main.go -E "(a|b)*c"
> abc
abc             => Matched.
> ab
ab              => Not matched.
> ac
ac              => Matched.
> ababbbac
ababbbac              => Matched.
> c
c               => Matched.
```

## Reference

Automata part follows [automata theory building a regexp machine](https://www.udemy.com/course/automata-theory-building-a-regexp-machine/).
The DFA minimization covered in the course remains unimplemented in the code

The parsing part is modified from the article [implementing a regular expression engine](https://deniskyashif.com/2019/02/17/implementing-a-regular-expression-engine/)
