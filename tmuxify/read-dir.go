package tmuxify

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

var default_ignores = []string{
	".git",
	"node_modules",
	".cache",
	".bun",
	".cargo",
	".wrangler",
}

func ReadDir(roots []string , ignore []string) (chan string) {
	if len(ignore) == 0 {
		ignore = default_ignores
	}

	if len(roots) == 0 {
		home , _ := os.UserHomeDir();
		roots = append(roots, home)
	}

	path_chan := make(chan string)

	go func(){
		for _ , root := range roots{
			abs, err := filepath.Abs(root)
			if err != nil {
				fmt.Print(err)
			}

			err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				base := filepath.Base(path)

				if slices.Contains(ignore , base) {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				
				path_chan <- path 
				return nil
			})
		}
		close(path_chan)
	}()

	return path_chan 
}
