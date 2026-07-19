package main

import (
	"fmt"

	"github.com/chirag-diwan/tmuxify/tmuxify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main(){
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	tview.Styles.ContrastBackgroundColor  = tcell.ColorDefault
	app := tview.NewApplication()
	root := tview.NewFlex().SetDirection(tview.FlexRow)
	list := tview.NewList().ShowSecondaryText(false)

	list.SetMainTextColor(tcell.ColorLightGray)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorWhite)
	list.SetSelectedBackgroundColor(tcell.ColorGray)
	list.SetSelectedBackgroundColor(tcell.ColorDarkGray)

	program_state := tmuxify.NewProgramState();

	input_chan := make(chan string)
	paths_chan := tmuxify.ReadDir([]string{} , []string{})

	go func(){
		for path := range paths_chan{
			tmuxify.AddPath(&program_state.RadixNodeRoot , path)
			program_state.Paths = append(program_state.Paths, path)
		}
	}()

	go func(){
		for {
			current_input := <-input_chan
			if len(current_input) > 0{
				program_state.Display_paths = tmuxify.GetMatch(&program_state.RadixNodeRoot , current_input)
			}else{
				program_state.Display_paths = program_state.Paths
			}

			list.Clear()
			for _ , path := range program_state.Display_paths {
				list.AddItem(path , "" , 0 , nil)
			}
		}
	}()

	input := tview.NewInputField(). SetLabel(" > "). SetPlaceholder("Enter the damm path")
	
	var selecte_dir string

	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune, tcell.KeyBackspace , tcell.KeyBackspace2:
			text := input.GetText();
			input_chan <- text

		case tcell.KeyUp:
			current := list.GetCurrentItem()
			if current > 0 {
				list.SetCurrentItem(current - 1)
			}
			return nil // consume the event

		case tcell.KeyDown:
			current := list.GetCurrentItem()
			if current < list.GetItemCount()-1 {
				list.SetCurrentItem(current + 1)
			}
			return nil
		case tcell.KeyEnter:
			app.Stop()
		}
		return event
	})

	root.AddItem(list , 0 , 100 , false)
	root.AddItem(input , 0 , 1 , true)

	if err := app.SetRoot(root , true).SetFocus(input).Run(); err != nil {
		fmt.Print(err)
	}

	close(input_chan)

}
