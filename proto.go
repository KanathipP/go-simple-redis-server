package main

import "fmt"

type Command struct {
	//??
}

func parseCommnand(msg string) (Command, error) {
	t := msg[0]
	fmt.Println(t)
	switch t {
	case '*':
		fmt.Println(msg[1:])
	}
	return Command{}, nil
}
