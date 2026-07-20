package tmuxify

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var default_ignores = []string{
	".git",
	"node_modules",
	".cache",
	".bun",
	".cargo",
}

func ReadDir(roots []string , ignore []string , maxDepth int) (chan string) {
	if len(ignore) == 0 {
		ignore = default_ignores
	}

	home , _ := os.UserHomeDir();
	if len(roots) == 0 {
		roots = append(roots, "")
	}

	path_chan := make(chan string)

	go func(){
		for _ , root := range roots{
			abs, err := filepath.Abs(filepath.Join(home , root))

			if err != nil {
				fmt.Print(err)
			}
			rootDepth := strings.Count(abs, string(filepath.Separator))

			err = filepath.WalkDir(abs, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				currentDepth := strings.Count(path, string(filepath.Separator)) - rootDepth
				if d.IsDir() && currentDepth >= maxDepth {
					return filepath.SkipDir
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
