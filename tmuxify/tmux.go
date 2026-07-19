package tmuxify

import (
	"context"
	"fmt"
	"os"

	"github.com/owenthereal/tmux"
)

func SetupTmux(config Config , start_directory string){
	context := context.Background();

	tmux_ , err := tmux.Default()
	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}

	ses , err := tmux_.NewSession(context , &tmux.SessionOptions{
		Name: config.Session.Name,
		StartDirectory:start_directory,
	})

	if err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
	
	if err := ses.Attach(context) ; err != nil {
		fmt.Print(err)
		os.Exit(-1)
	}
}
