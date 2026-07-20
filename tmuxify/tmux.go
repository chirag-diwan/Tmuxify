package tmuxify

import (
	"context"
	"fmt"
	"github.com/owenthereal/tmux"
)

func SetupTmux(projectContext context.Context , config Config , start_directory string) {
	tmuxHandler , err := tmux.Default()
	if err != nil {
		fmt.Printf("tmux.Default() failed with %s\n" , err.Error())
		return
	}

	ses , err := tmuxHandler.NewSession(projectContext , &tmux.SessionOptions{
		Name: config.Session.Name,
		StartDirectory:start_directory,
	})

	if err != nil {
		fmt.Printf("tmux.NewSession() failed with %s\n" , err.Error())
		return
	}

	for i := range len(config.Window) - 1{
		_ , err := ses.NewWindow(projectContext , &tmux.NewWindowOptions{
			WindowName: config.Window[i].Name,
			StartDirectory: start_directory,
		})

		if err != nil {
			fmt.Printf("ses.NewWindow() failed with %s\n" , err.Error())
			continue
		}
	}

	windows , err := ses.ListWindows(projectContext)

	if err != nil {
		fmt.Printf("ses.ListWindows() failed with %s\n" , err.Error())
		return
	}

	for i , win := range windows{
		if i == config.Session.Main {
			win.Select(projectContext)
		}

		panes , err := win.ListPanes(projectContext)
		if err != nil || len(panes) == 0 {
			fmt.Printf("win.ListPanes() failed with %s\n" , err.Error())
			continue;
		}

		pane := panes[0]

		for _ , cmd := range config.Window[i].Cmds{
			err = pane.SendLine(projectContext , cmd)
			if err != nil {
				fmt.Printf("pane.SendKeys() failed with %s\n" , err.Error())
				continue;
			}
		}

		i++
	}

	if err = ses.Attach(projectContext) ; err != nil {
		err = tmuxHandler.SwitchClient(projectContext , &tmux.SwitchClientOptions{
			TargetSession: config.Session.Name,
		})

		if err != nil {
			fmt.Printf("tmuxHandler.SwitchClient() failed with %s\n" , err.Error())
		}
	}
}
