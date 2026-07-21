package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"github.com/chirag-diwan/tmuxify/tmuxify"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	appConfig := tmuxify.GetAppConfig();
	program_state := tmuxify.NewProgramState()
	tmuxify.ReadDir(&program_state.PathChan, appConfig.Roots, appConfig.Ignore, appConfig.MaxDepth)
	programContext := context.Background()

	tview.Styles.PrimitiveBackgroundColor = tcell.NewHexColor(0x1E1E2E)
	tview.Styles.ContrastBackgroundColor = tcell.NewHexColor(0x313244)
	tview.Styles.MoreContrastBackgroundColor = tcell.NewHexColor(0x313244)
	tview.Styles.PrimaryTextColor = tcell.NewHexColor(0xCDD6F4)

	app := tview.NewApplication()
	root := tview.NewFlex().SetDirection(tview.FlexRow)

	list := tview.NewList().
	ShowSecondaryText(false).
	SetHighlightFullLine(true)

	list.SetMainTextColor(tcell.NewHexColor(0xA6ADC8))
	list.SetSelectedTextColor(tcell.NewHexColor(0x1E1E2E))
	list.SetSelectedBackgroundColor(tcell.NewHexColor(0xCBA6F7))

	list.SetShortcutColor(tcell.ColorDefault)

	input_element := tview.NewInputField().
	SetLabel(" > ").
	SetLabelColor(tcell.NewHexColor(0x89B4FA)).
	SetFieldBackgroundColor(tcell.NewHexColor(0x1E1E2E)).
	SetFieldTextColor(tcell.NewHexColor(0xCDD6F4))

	go func() {
		i := 0
		for{
			select {
			case <-program_state.Done:
				return

			case key := <-program_state.KeyChan:
				current := list.GetCurrentItem()
				switch key {
				case int(tcell.KeyUp):
					if current > 0 {
						program_state.Cursor = current - 1
					}

				case int(tcell.KeyDown):
					if current < list.GetItemCount()-1 {
						program_state.Cursor = current + 1
					}
				}

			case path := <-program_state.PathChan:
				program_state.PathRWmutex.Lock();
				program_state.Paths = append(program_state.Paths, path)
				program_state.PathRWmutex.Unlock();

				program_state.QueryMutex.Lock();
				current_query := program_state.CurrentQuery
				program_state.QueryMutex.Unlock()

				program_state.PathRWmutex.RLock();
				if tmuxify.KMPSearch(current_query , path , program_state.Lps) > 0{
					program_state.Display = append(program_state.Display, i)
				}
				program_state.PathRWmutex.RUnlock();
				i++

			case <- program_state.Ticker.C:
				app.QueueUpdateDraw(func() {
					_, _, _, height := list.GetRect()
	
					start := max(program_state.Cursor - height/2 , 0)
					end := max(program_state.Cursor + height/2 , len(program_state.Display) - 1)
					if end - start < height {
						start = max(end - height , 0)
					}

					list.Clear()

					program_state.QueryMutex.Lock();
					current_query := program_state.CurrentQuery
					program_state.QueryMutex.Unlock()

					program_state.PathRWmutex.RLock();
					program_state.Display = tmuxify.GetMatch(&program_state.Paths , current_query)
					
					for i , path_idx := range program_state.Display {
						if(i >= start && i <= end){
							if (i <= program_state.Cursor + height) && (i >= program_state.Cursor - height/2){ 
								list.AddItem(program_state.Paths[path_idx], "", 0, nil)
							}
						}
					}

					program_state.PathRWmutex.RUnlock();
					list.SetCurrentItem(program_state.Cursor)
				})

			case current_input := <-program_state.InputChan:
				program_state.PathRWmutex.RLock();
				if len(current_input) > 0 {
					program_state.Display = tmuxify.GetMatch(&program_state.Paths, current_input)
				}
				program_state.PathRWmutex.RUnlock();
			}
		}
	}()



	var selected_dir string
	input_element.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp , tcell.KeyDown:
			go func(){
				program_state.KeyChan <- int(event.Key())
			}()

		case tcell.KeyESC:
			program_state.Done <- true;
			app.Stop()

		case tcell.KeyEnter:
			program_state.Done <- true;
			app.Stop()

			if list.GetItemCount() > 0 {
				program_state.PathRWmutex.RLock();
				selected_dir = program_state.Paths[program_state.Display[list.GetCurrentItem()]]
				program_state.PathRWmutex.RUnlock();

				fileInfo, err := os.Stat(selected_dir)
				if err == nil {
					if !fileInfo.IsDir() {
						selected_dir = filepath.Dir(selected_dir)
					}
				}

				config := tmuxify.GetConfig(selected_dir)
				tmuxify.SetupTmux(programContext, config, selected_dir)
			}
		}
		return event
	})

	input_element.SetChangedFunc(func(text string) {
		go func(){
			program_state.QueryMutex.Lock();
			program_state.CurrentQuery = text
			program_state.QueryMutex.Unlock();

			program_state.Lps = tmuxify.CreateLps(text)
		}()
	})

	root.AddItem(list, 0, 1, false)
	root.AddItem(input_element, 1, 0, true) 

	if err := app.SetRoot(root, true).SetFocus(input_element).Run(); err != nil {
		fmt.Print(err)
	}
}
