package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	ascii = `
   ______      __          __    __             ____  __________  __ 
  / ____/___  / /   ____  / /_  / /_  __  __   / __ \/ ____/ __ \/ / 
 / / __/ __ \/ /   / __ \/ __ \/ __ \/ / / /  / /_/ / __/ / /_/ / /  
/ /_/ / /_/ / /___/ /_/ / /_/ / /_/ / /_/ /  / _, _/ /___/ ____/ /___
\____/\____/_____/\____/_.___/_.___/\__, /  /_/ |_/_____/_/   /_____/
                                   /____/                            
`
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	session, err := newSession(wd)
	if err != nil {
		panic(err)
	}
	fmt.Println(ascii)
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">> ")
		code, err := r.ReadString(';')
		if err != nil {
			panic(err)
		}
		session.add(code)

		err = session.writeToFile()
		if err != nil {
			session.displayError(err)
			continue
		}
		err = session.run(os.Stdout, os.Stdout)
		if err != nil {
			session.displayError(err)
		}
	}
}
