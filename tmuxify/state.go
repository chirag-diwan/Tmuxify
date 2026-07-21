package tmuxify

import (
	"sync"
	"time"
)

type ProgramState struct {
	Paths []string
	Display []int
	
	CurrentQuery string
	Lps []int

	Cursor int

	Ticker *time.Ticker
	
	Done chan bool
	PathChan chan string
	InputChan chan string
	KeyChan chan int

	PathRWmutex sync.RWMutex
	QueryMutex sync.Mutex
};

func NewProgramState() ProgramState{
	return ProgramState{
		Paths: []string{},
		Display: []int{},
		Cursor: 0,
		Ticker: time.NewTicker(16 * time.Millisecond),
		Done : make(chan bool),
		PathChan: make(chan string),
		InputChan: make(chan string , 2),
		KeyChan: make(chan int),
		PathRWmutex: sync.RWMutex{},
		QueryMutex: sync.Mutex{},
	}
}

