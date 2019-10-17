package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/golobby/gshell/interpreter"

	"github.com/c-bata/go-prompt"
)

var (
	currentInterpreter *interpreter.Interpreter
	DEBUG              bool
)

const (
	version = "0.1.1a"
)

func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func handler(input string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic: %v\n%s", err, debug.Stack())
		}
	}()
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	err := currentInterpreter.Add(input)
	if err != nil {
		fmt.Print(err.Error())

		return
	}

	fmt.Print(currentInterpreter.Eval())
	if DEBUG {
		fmt.Printf(":::::: D => %v\n", time.Since(start))
	}
}

func main() {

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	debug := flag.Bool("timestamp", false, "turns timestamp mode on")
	flag.Parse()
	DEBUG = *debug

	currentInterpreter, err = interpreter.NewSession(wd)
	if err != nil {
		panic(err)
	}
	err = currentInterpreter.Add(":e 1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("gshell v%s\n", version)

	p := prompt.New(handler, completer, prompt.OptionPrefix("gshell> "))
	p.Run()
}
