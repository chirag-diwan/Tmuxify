package tmuxify

import (
)

type ProgramState struct {
	Paths []string
	Display_paths []string
	RadixNodeRoot RadixNode
};

func NewProgramState() ProgramState{
	return ProgramState{
		Paths: []string{},
		Display_paths: []string{},
		RadixNodeRoot: RadixNode{},
	}
}

