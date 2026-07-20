package main

import (
	"context"
	"fmt"

	"github.com/chirag-diwan/tmuxify/tmuxify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	appConfig := tmuxify.GetAppConfig();

	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor = tcell.ColorDefault
	tview.Styles.MoreContrastBackgroundColor = tcell.ColorDefault

	app := tview.NewApplication()
	root := tview.NewFlex().SetDirection(tview.FlexRow)

	list := tview.NewList().
		ShowSecondaryText(false).
		SetHighlightFullLine(true)

	list.SetMainTextColor(tcell.ColorWhite)
	list.SetSelectedTextColor(tcell.ColorWhite)
	list.SetSelectedBackgroundColor(tcell.ColorDarkGray)
	
	list.SetShortcutColor(tcell.ColorDefault)

	program_state := tmuxify.NewProgramState()
	input_chan := make(chan string , 2)

	paths_chan := tmuxify.ReadDir(appConfig.Roots, appConfig.Ignore, appConfig.MaxDepth)

	programContext := context.Background()
	ctx, cancel := context.WithCancel(programContext)

	go func() {
		for path := range paths_chan {
			tmuxify.AddPath(&program_state.RadixNodeRoot, path)
			program_state.Paths = append(program_state.Paths, path)
			list.AddItem(path, "", 0, nil)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-input_chan:
				current_input := <-input_chan
				if len(current_input) > 0 {
					program_state.Display_paths = tmuxify.GetMatch(&program_state.RadixNodeRoot, current_input)
				} else {
					program_state.Display_paths = program_state.Paths
				}

				app.QueueUpdateDraw(func() {
					list.Clear()
					for _, path := range program_state.Display_paths {
						list.AddItem(path, "", 0, nil)
					}
				})
			}
		}
	}()

	input := tview.NewInputField().
	SetLabel("> ").
	SetLabelColor(tcell.ColorAqua).
	SetFieldBackgroundColor(tcell.ColorDefault)

	var selected_dir string
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune, tcell.KeyBackspace, tcell.KeyBackspace2:
			text := input.GetText()
			input_chan <- text
		case tcell.KeyUp:
			current := list.GetCurrentItem()
			if current > 0 {
				list.SetCurrentItem(current - 1)
			}
			return nil
		case tcell.KeyDown:
			current := list.GetCurrentItem()
			if current < list.GetItemCount()-1 {
				list.SetCurrentItem(current + 1)
			}
			return nil
		case tcell.KeyEnter:
			cancel()
			app.Stop()
			if list.GetItemCount() > 0 {
				selected_dir = program_state.Display_paths[list.GetCurrentItem()]
				config := tmuxify.GetConfig(selected_dir)
				tmuxify.SetupTmux(programContext, config, selected_dir)
			}
		}
		return event
	})

	root.AddItem(list, 0, 1, false)
	root.AddItem(input, 1, 0, true) 

	if err := app.SetRoot(root, true).SetFocus(input).Run(); err != nil {
		fmt.Print(err)
	}
}
